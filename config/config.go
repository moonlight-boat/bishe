package config

import (
	"os"
	"strconv"
)

type Config struct {
	// 数据库配置
	DatabasePath string
	
	// 服务器配置
	ServerPort string
	
	// Python脚本配置
	PythonScriptPath string
	CSVFilePath      string
	
	// 定时任务配置
	CronSchedule string
}

func LoadConfig() *Config {
	return &Config{
		DatabasePath:     getEnv("DATABASE_PATH", "./jobs.db"),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		PythonScriptPath: getEnv("PYTHON_SCRIPT_PATH", "/opt/job/job-scheduler/codes/main.py"),
		CSVFilePath:      getEnv("CSV_FILE_PATH", "/opt/job/job-scheduler/auto_job_list.csv"),
		CronSchedule:     getEnv("CRON_SCHEDULE", "0 0 12 * * *"), // 每天12:00
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
