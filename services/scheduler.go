package services

import (
	"job-backend/config"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	cron        *cron.Cron
	config      *config.Config
	csvImporter *CSVImporter
}

// NewScheduler 创建定时任务调度器
func NewScheduler(cfg *config.Config) *Scheduler {
	return &Scheduler{
		cron:        cron.New(cron.WithSeconds()),
		config:      cfg,
		csvImporter: NewCSVImporter(),
	}
}

// Start 启动定时任务
func (s *Scheduler) Start() {
	log.Println("启动定时任务调度器")
	
	// 添加定时任务：每天12:00执行
	_, err := s.cron.AddFunc(s.config.CronSchedule, s.executeJob)
	if err != nil {
		log.Printf("添加定时任务失败: %v", err)
		return
	}
	
	// 启动调度器
	s.cron.Start()
	log.Printf("定时任务已启动，执行计划: %s", s.config.CronSchedule)
	
	// 可选：启动时执行一次（用于测试）
	// go s.executeJob()
}

// Stop 停止定时任务
func (s *Scheduler) Stop() {
	if s.cron != nil {
		s.cron.Stop()
		log.Println("定时任务调度器已停止")
	}
}

// executeJob 执行定时任务
func (s *Scheduler) executeJob() {
	log.Println("开始执行定时任务：爬取招聘信息")
	startTime := time.Now()
	
	// 1. 执行Python爬虫脚本
	if err := s.runPythonScript(); err != nil {
		log.Printf("执行Python脚本失败: %v", err)
		return
	}
	
	// 2. 等待一段时间确保CSV文件生成完成
	time.Sleep(5 * time.Second)
	
	// 3. 导入CSV数据到数据库
	recordCount, err := s.csvImporter.ImportFromCSV(s.config.CSVFilePath)
	if err != nil {
		log.Printf("导入CSV失败: %v", err)
		return
	}
	
	executionTime := time.Since(startTime)
	log.Printf("定时任务执行完成，导入 %d 条记录，耗时: %v", recordCount, executionTime)
}

// runPythonScript 执行Python爬虫脚本
func (s *Scheduler) runPythonScript() error {
	log.Printf("执行Python脚本: %s", s.config.PythonScriptPath)
	
	// 获取脚本所在目录
	scriptDir := filepath.Dir(s.config.PythonScriptPath)
	scriptName := filepath.Base(s.config.PythonScriptPath)
	
	// 创建命令
	cmd := exec.Command("python", scriptName)
	cmd.Dir = scriptDir
	
	// 设置环境变量
	cmd.Env = os.Environ()
	
	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Python脚本执行失败: %v", err)
		log.Printf("脚本输出: %s", string(output))
		return err
	}
	
	log.Printf("Python脚本执行成功")
	log.Printf("脚本输出: %s", string(output))
	
	return nil
}

// ExecuteJobManually 手动执行任务（用于测试）
func (s *Scheduler) ExecuteJobManually() {
	log.Println("手动触发定时任务")
	go s.executeJob()
}

// GetNextRunTime 获取下次执行时间
func (s *Scheduler) GetNextRunTime() time.Time {
	entries := s.cron.Entries()
	if len(entries) > 0 {
		return entries[0].Next
	}
	return time.Time{}
}

// GetJobStatus 获取任务状态
func (s *Scheduler) GetJobStatus() map[string]interface{} {
	status := make(map[string]interface{})
	
	status["cron_schedule"] = s.config.CronSchedule
	status["next_run_time"] = s.GetNextRunTime()
	status["python_script_path"] = s.config.PythonScriptPath
	status["csv_file_path"] = s.config.CSVFilePath
	
	// 检查Python脚本是否存在
	if _, err := os.Stat(s.config.PythonScriptPath); err == nil {
		status["python_script_exists"] = true
	} else {
		status["python_script_exists"] = false
	}
	
	// 检查CSV文件是否存在
	if csvInfo, err := s.csvImporter.GetCSVInfo(s.config.CSVFilePath); err == nil {
		status["csv_file_info"] = csvInfo
	}
	
	return status
}
