package handlers

import (
	"job-backend/config"
	"job-backend/database"
	"job-backend/models"
	"job-backend/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	jobService    = services.NewJobService()
	csvImporter   = services.NewCSVImporter()
	lastSyncTime  time.Time
	syncInProgress bool
)

// GetJobs 获取招聘信息列表
func GetJobs(c *gin.Context) {
	var params models.JobQueryParams
	
	// 绑定查询参数
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": err.Error(),
		})
		return
	}
	
	// 调用服务层
	result, err := jobService.GetJobs(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "查询失败",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, result)
}

// GetJobByID 根据ID获取招聘信息
func GetJobByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的ID",
			"message": "ID必须是数字",
		})
		return
	}
	
	job, err := jobService.GetJobByID(uint(id))
	if err != nil {
		if err.Error() == "招聘信息不存在" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "未找到",
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "查询失败",
				"message": err.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, job)
}

// SearchJobs 搜索招聘信息
func SearchJobs(c *gin.Context) {
	var params models.JobQueryParams
	
	// 绑定查询参数
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": err.Error(),
		})
		return
	}
	
	// 调用服务层搜索
	result, err := jobService.SearchJobs(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "搜索失败",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, result)
}

// GetJobStats 获取招聘信息统计
func GetJobStats(c *gin.Context) {
	stats, err := jobService.GetJobStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取统计失败",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

// ManualSync 手动触发数据同步
func ManualSync(c *gin.Context) {
	if syncInProgress {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "同步进行中",
			"message": "已有同步任务在进行中，请稍后再试",
		})
		return
	}
	
	// 获取CSV文件路径
	cfg := config.LoadConfig()
	csvPath := c.DefaultQuery("csv_path", cfg.CSVFilePath)
	
	// 异步执行同步
	go func() {
		syncInProgress = true
		defer func() {
			syncInProgress = false
		}()
		
		startTime := time.Now()
		
		// 导入CSV
		recordCount, err := csvImporter.ImportFromCSV(csvPath)
		if err != nil {
			// 这里可以记录到日志系统
			return
		}
		
		lastSyncTime = time.Now()
		executionTime := time.Since(startTime)
		
		// 这里可以发送通知或记录日志
		_ = recordCount
		_ = executionTime
	}()
	
	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "同步任务已启动",
	})
}

// SyncStatus 获取同步状态
func SyncStatus(c *gin.Context) {
	// 获取CSV文件信息
	cfg := config.LoadConfig()
	csvPath := c.DefaultQuery("csv_path", cfg.CSVFilePath)
	csvInfo, err := csvImporter.GetCSVInfo(csvPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取CSV信息失败",
			"message": err.Error(),
		})
		return
	}
	
	// 获取数据库记录数
	var dbRecordCount int64
	database.GetDB().Model(&models.JobPosting{}).Count(&dbRecordCount)
	
	status := gin.H{
		"sync_in_progress": syncInProgress,
		"last_sync_time":   lastSyncTime,
		"csv_info":         csvInfo,
		"db_record_count":  dbRecordCount,
	}
	
	c.JSON(http.StatusOK, status)
}
