## WORLDSVERRE

本文件用于记录本仓库的约定与操作注意事项，便于多人协作与自动化工具稳定运行。

### 项目类型

- **后端**：Go 1.21（Gin + GORM）
- **数据库**：SQLite（默认 `./jobs.db`）
- **数据来源**：Python 脚本生成 CSV，再由后端导入

### 运行配置（环境变量）

- **DATABASE_PATH**：SQLite 路径，默认 `./jobs.db`
- **SERVER_PORT**：服务端口，默认 `8080`
- **PYTHON_SCRIPT_PATH**：Python 脚本路径（定时任务使用）
- **CSV_FILE_PATH**：CSV 文件路径（定时任务/手动同步使用）
- **CRON_SCHEDULE**：cron 表达式（带秒），默认 `0 0 12 * * *`

### API 约定

- **健康检查**：`GET /health`
- **统一前缀**：`/api`

### 操作注意事项

- 定时任务会执行 Python 脚本并导入 CSV；部署环境需保证 `python` 可用且路径正确。
- CSV 导入当前策略为“清空重建”（会先删除 `job_postings` 全表数据）。

