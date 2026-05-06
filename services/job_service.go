package services

import (
	"fmt"
	"job-backend/database"
	"job-backend/models"
	"math"
	"strings"

	"gorm.io/gorm"
)

// JobService 招聘信息业务服务
type JobService struct {
	db *gorm.DB
}

// NewJobService 创建招聘信息服务
func NewJobService() *JobService {
	return &JobService{
		db: database.GetDB(),
	}
}

// GetJobs 获取招聘信息列表
func (s *JobService) GetJobs(params models.JobQueryParams) (*models.JobListResponse, error) {
	var jobs []models.JobPosting
	var total int64
	
	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 20
	}
	if params.Size > 100 {
		params.Size = 100
	}
	
	// 确保数据库连接可用
	if s.db == nil {
		s.db = database.GetDB()
	}
	
	// 构建查询
	query := s.db.Model(&models.JobPosting{})
	
	// 应用过滤条件
	query = s.applyFilters(query, params)
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取总数失败: %v", err)
	}
	
	// 应用排序
	query = s.applySorting(query, params)
	
	// 应用分页
	offset := (params.Page - 1) * params.Size
	if err := query.Offset(offset).Limit(params.Size).Find(&jobs).Error; err != nil {
		return nil, fmt.Errorf("查询招聘信息失败: %v", err)
	}
	
	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(params.Size)))
	
	return &models.JobListResponse{
		Jobs:       jobs,
		Total:      total,
		Page:       params.Page,
		Size:       params.Size,
		TotalPages: totalPages,
	}, nil
}

// GetJobByID 根据ID获取招聘信息
func (s *JobService) GetJobByID(id uint) (*models.JobPosting, error) {
	// 确保数据库连接可用
	if s.db == nil {
		s.db = database.GetDB()
	}
	
	var job models.JobPosting
	
	if err := s.db.First(&job, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("招聘信息不存在")
		}
		return nil, fmt.Errorf("查询招聘信息失败: %v", err)
	}
	
	return &job, nil
}

// SearchJobs 搜索招聘信息
func (s *JobService) SearchJobs(params models.JobQueryParams) (*models.JobListResponse, error) {
	// 搜索逻辑与GetJobs相同，但会应用关键词搜索
	return s.GetJobs(params)
}

// GetJobStats 获取招聘信息统计
func (s *JobService) GetJobStats() (*models.JobStatsResponse, error) {
	var total int64
	
	// 获取总数
	if err := s.db.Model(&models.JobPosting{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取总数失败: %v", err)
	}
	
	stats := &models.JobStatsResponse{
		Total:            total,
		CompanyTypeStats: make(map[string]int64),
		IndustryStats:    make(map[string]int64),
		DegreeStats:      make(map[string]int64),
		CityStats:        make(map[string]int64),
		RecruitmentStats: make(map[string]int64),
	}
	
	// 公司类型统计
	if err := s.getFieldStats("company_type", stats.CompanyTypeStats); err != nil {
		return nil, err
	}
	
	// 行业统计
	if err := s.getFieldStats("industry", stats.IndustryStats); err != nil {
		return nil, err
	}
	
	// 学历要求统计
	if err := s.getFieldStats("degree_requirement", stats.DegreeStats); err != nil {
		return nil, err
	}
	
	// 工作城市统计（处理多城市情况）
	if err := s.getCityStats(stats.CityStats); err != nil {
		return nil, err
	}
	
	// 招聘类型统计
	if err := s.getFieldStats("recruitment_type", stats.RecruitmentStats); err != nil {
		return nil, err
	}
	
	return stats, nil
}

// applyFilters 应用过滤条件
func (s *JobService) applyFilters(query *gorm.DB, params models.JobQueryParams) *gorm.DB {
	// 公司名称模糊查询
	if params.CompanyName != "" {
		query = query.Where("company_name LIKE ?", "%"+params.CompanyName+"%")
	}
	
	// 公司类型精确匹配
	if params.CompanyType != "" {
		query = query.Where("company_type = ?", params.CompanyType)
	}
	
	// 行业精确匹配
	if params.Industry != "" {
		query = query.Where("industry = ?", params.Industry)
	}
	
	// 招聘类型精确匹配
	if params.RecruitmentType != "" {
		query = query.Where("recruitment_type = ?", params.RecruitmentType)
	}
	
	// 学历要求匹配
	if params.DegreeRequirement != "" {
		query = query.Where("degree_requirement LIKE ?", "%"+params.DegreeRequirement+"%")
	}
	
	// 工作城市匹配
	if params.WorkCity != "" {
		query = query.Where("work_city LIKE ?", "%"+params.WorkCity+"%")
	}
	
	// 岗位名称模糊查询
	if params.PositionName != "" {
		query = query.Where("position_name LIKE ?", "%"+params.PositionName+"%")
	}
	
	// 关键词搜索（在多个字段中搜索）
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(
			"company_name LIKE ? OR position_name LIKE ? OR major_requirement LIKE ? OR job_description LIKE ?",
			keyword, keyword, keyword, keyword,
		)
	}
	
	return query
}

// applySorting 应用排序
func (s *JobService) applySorting(query *gorm.DB, params models.JobQueryParams) *gorm.DB {
	sortBy := params.SortBy
	sortOrder := params.SortOrder
	
	// 设置默认排序
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}
	
	// 验证排序字段
	validSortFields := map[string]bool{
		"id":                 true,
		"company_name":       true,
		"company_scale":      true,
		"company_type":       true,
		"industry":           true,
		"recruitment_type":   true,
		"degree_requirement": true,
		"position_name":      true,
		"work_city":          true,
		"deadline":           true,
		"salary":             true,
		"created_at":         true,
		"updated_at":         true,
	}
	
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}
	
	return query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
}

// getFieldStats 获取字段统计
func (s *JobService) getFieldStats(field string, stats map[string]int64) error {
	var results []struct {
		Value string
		Count int64
	}
	
	err := s.db.Model(&models.JobPosting{}).
		Select(fmt.Sprintf("%s as value, COUNT(*) as count", field)).
		Where(fmt.Sprintf("%s != ''", field)).
		Group(field).
		Order("count DESC").
		Limit(20).
		Scan(&results).Error
	
	if err != nil {
		return fmt.Errorf("获取%s统计失败: %v", field, err)
	}
	
	for _, result := range results {
		stats[result.Value] = result.Count
	}
	
	return nil
}

// getCityStats 获取城市统计（处理多城市情况）
func (s *JobService) getCityStats(stats map[string]int64) error {
	var jobs []models.JobPosting
	
	if err := s.db.Select("work_city").Where("work_city != ''").Find(&jobs).Error; err != nil {
		return fmt.Errorf("获取城市数据失败: %v", err)
	}
	
	cityCount := make(map[string]int64)
	
	for _, job := range jobs {
		// 处理多城市情况，用|或、分割
		cities := strings.FieldsFunc(job.WorkCity, func(r rune) bool {
			return r == '|' || r == '、' || r == '，' || r == ','
		})
		
		for _, city := range cities {
			city = strings.TrimSpace(city)
			if city != "" {
				cityCount[city]++
			}
		}
	}
	
	// 复制到结果map
	for city, count := range cityCount {
		stats[city] = count
	}
	
	return nil
}
