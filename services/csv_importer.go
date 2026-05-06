package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"job-backend/database"
	"job-backend/models"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)

// CSVImporter CSV导入服务
type CSVImporter struct {
	db *gorm.DB
}

// NewCSVImporter 创建CSV导入服务
func NewCSVImporter() *CSVImporter {
	return &CSVImporter{
		db: database.GetDB(),
	}
}

// ImportFromCSV 从CSV文件导入数据（方案A：清空重建）
func (c *CSVImporter) ImportFromCSV(csvPath string) (int, error) {
	log.Printf("开始导入CSV文件: %s", csvPath)
	
	// 检查文件是否存在
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		return 0, fmt.Errorf("CSV文件不存在: %s", csvPath)
	}
	
	// 打开CSV文件
	file, err := os.Open(csvPath)
	if err != nil {
		return 0, fmt.Errorf("打开CSV文件失败: %v", err)
	}
	defer file.Close()
	
	// 创建CSV读取器
	reader := csv.NewReader(file)
	
	// 读取表头
	headers, err := reader.Read()
	if err != nil {
		return 0, fmt.Errorf("读取CSV表头失败: %v", err)
	}
	
	// 创建字段索引映射
	fieldIndexMap := make(map[string]int)
	for i, header := range headers {
		fieldIndexMap[strings.TrimSpace(header)] = i
	}
	
	log.Printf("CSV表头: %v", headers)
	
	// 确保数据库连接可用
	if c.db == nil {
		c.db = database.GetDB()
	}
	
	// 开启事务
	tx := c.db.Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("开启事务失败: %v", tx.Error)
	}
	
	// 清空现有数据
	if err := tx.Exec("DELETE FROM job_postings").Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("清空数据表失败: %v", err)
	}
	
	log.Println("已清空现有数据")
	
	// 读取并导入数据
	recordCount := 0
	batchSize := 100
	jobs := make([]models.JobPosting, 0, batchSize)
	
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("读取CSV记录失败: %v", err)
		}
		
		// 转换为JobPosting结构体
		job := models.JobPosting{
			CompanyName:       c.getFieldValue(record, fieldIndexMap, "单位名称"),
			CompanyScale:      c.getFieldValue(record, fieldIndexMap, "单位规模"),
			CompanyType:       c.getFieldValue(record, fieldIndexMap, "单位性质"),
			Industry:          c.getFieldValue(record, fieldIndexMap, "单位行业"),
			RecruitmentType:   c.getFieldValue(record, fieldIndexMap, "招聘类型"),
			DegreeRequirement: c.getFieldValue(record, fieldIndexMap, "学位要求"),
			PositionCount:     c.getFieldValue(record, fieldIndexMap, "招聘岗位数"),
			PositionName:      c.getFieldValue(record, fieldIndexMap, "岗位"),
			MajorRequirement:  c.getFieldValue(record, fieldIndexMap, "专业要求"),
			JobDescription:    c.getFieldValue(record, fieldIndexMap, "岗位要求"),
			OfficialLink:      c.getFieldValue(record, fieldIndexMap, "学校官网链接"),
			WorkCity:          c.getFieldValue(record, fieldIndexMap, "工作城市"),
			Deadline:          c.getFieldValue(record, fieldIndexMap, "截止日期"),
			ApplicationMethod: c.getFieldValue(record, fieldIndexMap, "投递口"),
			Salary:            c.getFieldValue(record, fieldIndexMap, "年薪"),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		
		jobs = append(jobs, job)
		recordCount++
		
		// 批量插入
		if len(jobs) >= batchSize {
			if err := tx.Create(&jobs).Error; err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("批量插入失败: %v", err)
			}
			jobs = jobs[:0] // 清空切片
			log.Printf("已导入 %d 条记录", recordCount)
		}
	}
	
	// 插入剩余记录
	if len(jobs) > 0 {
		if err := tx.Create(&jobs).Error; err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("插入剩余记录失败: %v", err)
		}
	}
	
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("提交事务失败: %v", err)
	}
	
	log.Printf("CSV导入完成，共导入 %d 条记录", recordCount)
	return recordCount, nil
}

// getFieldValue 根据字段名获取记录中的值
func (c *CSVImporter) getFieldValue(record []string, fieldMap map[string]int, fieldName string) string {
	if index, exists := fieldMap[fieldName]; exists && index < len(record) {
		return strings.TrimSpace(record[index])
	}
	return ""
}

// GetCSVInfo 获取CSV文件信息
func (c *CSVImporter) GetCSVInfo(csvPath string) (map[string]interface{}, error) {
	info := make(map[string]interface{})
	
	// 检查文件是否存在
	fileInfo, err := os.Stat(csvPath)
	if os.IsNotExist(err) {
		info["exists"] = false
		return info, nil
	}
	if err != nil {
		return nil, err
	}
	
	info["exists"] = true
	info["size"] = fileInfo.Size()
	info["modified_time"] = fileInfo.ModTime()
	
	// 读取记录数量
	file, err := os.Open(csvPath)
	if err != nil {
		return info, nil
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	recordCount := 0
	
	// 跳过表头
	if _, err := reader.Read(); err != nil {
		return info, nil
	}
	
	// 计算记录数
	for {
		if _, err := reader.Read(); err == io.EOF {
			break
		} else if err != nil {
			break
		}
		recordCount++
	}
	
	info["record_count"] = recordCount
	return info, nil
}
