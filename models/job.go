package models

import (
	"time"
)

// JobPosting 招聘信息模型
type JobPosting struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	CompanyName       string    `gorm:"size:255;index" json:"company_name"`        // 单位名称
	CompanyScale      string    `gorm:"size:100" json:"company_scale"`             // 单位规模
	CompanyType       string    `gorm:"size:100;index" json:"company_type"`        // 单位性质
	Industry          string    `gorm:"size:100;index" json:"industry"`            // 单位行业
	RecruitmentType   string    `gorm:"size:50;index" json:"recruitment_type"`     // 招聘类型
	DegreeRequirement string    `gorm:"size:100;index" json:"degree_requirement"`  // 学位要求
	PositionCount     string    `gorm:"size:20" json:"position_count"`             // 招聘岗位数
	PositionName      string    `gorm:"type:text;index" json:"position_name"`      // 岗位
	MajorRequirement  string    `gorm:"type:text" json:"major_requirement"`        // 专业要求
	JobDescription    string    `gorm:"type:text" json:"job_description"`          // 岗位要求
	OfficialLink      string    `gorm:"size:500" json:"official_link"`             // 学校官网链接
	WorkCity          string    `gorm:"size:200;index" json:"work_city"`           // 工作城市
	Deadline          string    `gorm:"size:50" json:"deadline"`                   // 截止日期
	ApplicationMethod string    `gorm:"type:text" json:"application_method"`       // 投递口
	Salary            string    `gorm:"size:100;index" json:"salary"`              // 年薪
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName 指定表名
func (JobPosting) TableName() string {
	return "job_postings"
}

// JobQueryParams 查询参数
type JobQueryParams struct {
	Page              int    `form:"page"`                                        // 页码
	Size              int    `form:"size"`                                        // 每页大小
	CompanyName       string `form:"company_name"`                                // 公司名称（模糊查询）
	CompanyType       string `form:"company_type"`                                // 公司类型
	Industry          string `form:"industry"`                                    // 行业
	RecruitmentType   string `form:"recruitment_type"`                            // 招聘类型
	DegreeRequirement string `form:"degree_requirement"`                          // 学历要求
	WorkCity          string `form:"work_city"`                                   // 工作城市
	PositionName      string `form:"position_name"`                               // 岗位名称（模糊查询）
	Keyword           string `form:"keyword"`                                     // 关键词搜索
	SortBy            string `form:"sort_by"`                                     // 排序字段
	SortOrder         string `form:"sort_order"`                                  // 排序方向
}

// JobListResponse 列表响应
type JobListResponse struct {
	Jobs       []JobPosting `json:"jobs"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	Size       int          `json:"size"`
	TotalPages int          `json:"total_pages"`
}

// JobStatsResponse 统计响应
type JobStatsResponse struct {
	Total              int64                    `json:"total"`
	CompanyTypeStats   map[string]int64         `json:"company_type_stats"`
	IndustryStats      map[string]int64         `json:"industry_stats"`
	DegreeStats        map[string]int64         `json:"degree_stats"`
	CityStats          map[string]int64         `json:"city_stats"`
	RecruitmentStats   map[string]int64         `json:"recruitment_stats"`
}

// SyncResponse 同步响应
type SyncResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	RecordsCount int   `json:"records_count"`
	ExecutionTime string `json:"execution_time"`
}
