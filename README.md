# 职达招聘信息后端服务 (Job Backend)

> 基于 Go + Gin 框架的招聘信息数据管理与 API 服务，为"智能招聘数据分析平台"提供后端数据支持。

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-v1.9+-00ADD8?logo=go)](https://gin-gonic.com/)
[![SQLite](https://img.shields.io/badge/SQLite-3-003B57?logo=sqlite)](https://sqlite.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 📋 项目简介

本项目是"智能招聘数据分析平台"的后端服务，主要负责：

- **招聘数据采集管理**：接收并存储来自 Python 爬虫的招聘数据（CSV 格式）
- **数据查询服务**：提供丰富的筛选、搜索、分页查询能力
- **数据统计分析**：按公司类型、行业、学历要求、工作城市等维度统计
- **定时任务调度**：每日自动执行爬虫脚本并同步数据到数据库
- **一键联调支持**：内置 PowerShell 脚本实现环境一键搭建与验证

---

## ✨ 功能特性

| 模块 | 功能描述 |
|------|---------|
| 🗄️ **数据管理** | CSV 批量导入、SQLite 持久化存储、自动表结构迁移 |
| 🔍 **智能查询** | 多条件组合筛选、关键词全文搜索、灵活排序、分页返回 |
| 📊 **数据统计** | 公司类型/行业/学历/城市/招聘类型等多维度统计 |
| ⏰ **定时同步** | Cron 表达式配置，自动执行爬虫 + 数据导入 |
| 🔧 **手动同步** | 提供 HTTP 接口手动触发数据同步 |
| 🚀 **一键联调** | PowerShell 脚本实现环境检查、服务启动、自动验证 |
| 🌐 **跨域支持** | 内置 CORS 中间件，支持前端跨域调用 |

---

## 🏗️ 技术架构

```
┌─────────────────────────────────────────────────────────────┐
│                        客户端 (前端/爬虫)                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Job Backend (Go + Gin)                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐ │
│  │ handlers │  │ services │  │  models  │  │   config     │ │
│  │  HTTP层  │──│ 业务逻辑 │──│ 数据模型 │  │  环境配置    │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────┘ │
│       │              │                                      │
│       │       ┌──────┴──────┐     ┌─────────────────┐      │
│       │       │ CSVImporter │     │    Scheduler    │      │
│       │       │  CSV导入器  │     │  定时任务调度器  │      │
│       │       └─────────────┘     └─────────────────┘      │
│       │                                                     │
│  ┌────┴─────────────────────────────────────────────────┐   │
│  │                   database (SQLite + GORM)            │   │
│  │              jobs.db → job_postings 表               │   │
│  └───────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## 📁 项目目录

```
job-backend/
├── config/                    # 配置模块
│   └── config.go              # 环境变量与配置加载
├── database/                  # 数据库模块
│   └── db.go                  # SQLite 连接与自动迁移
├── handlers/                  # HTTP 请求处理器
│   └── job_handler.go         # 招聘相关的 REST API 接口
├── models/                    # 数据模型
│   └── job.go                 # JobPosting 实体与查询参数结构
├── services/                  # 业务服务层
│   ├── job_service.go         # 招聘信息查询与统计逻辑
│   ├── csv_importer.go        # CSV 数据导入服务
│   └── scheduler.go           # 定时任务调度器
├── scripts/                   # 运维脚本
│   ├── one-click-debug.ps1    # 一键联调 PowerShell 脚本
│   ├── one-click-debug.cmd    # Windows 双击入口
│   └── one-click-debug-safe.ps1 # 安全模式联调脚本
├── 职达工作流/                # Python 爬虫子模块（独立仓库）
│   └── codes/main.py          # 招聘数据爬虫入口
├── main.go                    # 应用入口：初始化+路由+启动服务
├── go.mod / go.sum            # Go 依赖管理
├── API.md                     # 详细 API 接口文档
├── ONE_CLICK_DEBUG.md         # 一键联调使用说明
├── 接口文档.md                 # 中文接口文档
└── README.md                  # 本文件
```

---

## 🚀 快速开始

### 前置要求

- [Go](https://golang.org/dl/) 1.21 或更高版本
- [Python](https://www.python.org/) 3.8+（如需运行爬虫）
- Windows / Linux / macOS

### 1. 克隆项目

```bash
git clone https://github.com/moonlight-boat/bishe.git
cd bishe/job-backend
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 运行服务

```bash
go run main.go
```

服务默认启动在 `http://localhost:8080`

### 4. 验证启动

```bash
curl http://localhost:8080/health
```

预期响应：
```json
{
  "status": "ok",
  "message": "Job Backend Service is running"
}
```

---

## ⚙️ 环境变量配置

所有配置均通过环境变量读取，未设置时使用默认值：

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `DATABASE_PATH` | SQLite 数据库文件路径 | `./jobs.db` |
| `SERVER_PORT` | HTTP 服务端口 | `8080` |
| `PYTHON_SCRIPT_PATH` | Python 爬虫脚本路径 | `/opt/job/job-scheduler/codes/main.py` |
| `CSV_FILE_PATH` | 爬虫输出的 CSV 文件路径 | `/opt/job/job-scheduler/auto_job_list.csv` |
| `CRON_SCHEDULE` | 定时同步 Cron 表达式 | `0 0 12 * * *`（每天 12:00） |

### 示例：自定义配置启动

```bash
# Linux/macOS
export SERVER_PORT=9090
export DATABASE_PATH=./data/my_jobs.db
export CSV_FILE_PATH=./data/auto_job_list.csv
go run main.go

# Windows PowerShell
$env:SERVER_PORT="9090"
$env:DATABASE_PATH="./data/my_jobs.db"
go run main.go
```

---

## 🔌 API 接口概览

Base URL: `http://localhost:8080/api`

### 招聘信息接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/jobs` | 获取招聘信息列表（支持分页、筛选、排序） |
| `GET` | `/api/jobs/:id` | 根据 ID 获取单个招聘信息 |
| `GET` | `/api/jobs/search` | 关键词全文搜索 |
| `GET` | `/api/jobs/stats` | 获取多维度统计信息 |

### 数据同步接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/sync` | 手动触发 CSV 数据同步 |
| `GET` | `/api/sync/status` | 获取同步任务状态 |

### 系统接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/health` | 健康检查 |

> 📖 **完整接口文档** 请参阅 [`API.md`](./API.md) 或 [`接口文档.md`](./接口文档.md)

---

## 📊 数据模型

### JobPosting（招聘信息）

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | uint | 主键 ID |
| `company_name` | string | 单位名称（索引） |
| `company_scale` | string | 单位规模 |
| `company_type` | string | 单位性质：国企/民企等（索引） |
| `industry` | string | 所属行业（索引） |
| `recruitment_type` | string | 招聘类型：校招/社招（索引） |
| `degree_requirement` | string | 学历要求：本科/硕士/博士（索引） |
| `position_count` | string | 招聘岗位数量 |
| `position_name` | string | 岗位名称（索引） |
| `major_requirement` | string | 专业要求 |
| `job_description` | string | 岗位要求/职责 |
| `official_link` | string | 学校官网链接 |
| `work_city` | string | 工作城市，支持多城市（索引） |
| `deadline` | string | 投递截止日期 |
| `application_method` | string | 投递方式 |
| `salary` | string | 年薪范围（索引） |
| `created_at` | time | 记录创建时间 |
| `updated_at` | time | 记录更新时间 |

---

## 🔄 数据同步流程

### 自动同步（定时任务）

```
每天 12:00
    │
    ▼
┌───────────────┐
│ 执行 Python   │  →  职达工作流/codes/main.py
│ 爬虫脚本      │     生成 auto_job_list.csv
└───────────────┘
    │
    ▼
┌───────────────┐
│ 等待 5 秒     │  →  确保 CSV 写入完成
└───────────────┘
    │
    ▼
┌───────────────┐
│ CSV 导入数据库 │  →  清空旧数据 → 批量插入新数据
│ (事务保证)    │     默认每批 100 条
└───────────────┘
    │
    ▼
  完成，记录日志
```

### 手动同步

```bash
# 通过 API 手动触发
curl -X POST "http://localhost:8080/api/sync?csv_path=/path/to/your.csv"

# 响应
{
  "success": true,
  "message": "同步任务已启动"
}
```

---

## 🧪 一键联调

项目提供 PowerShell 一键联调脚本，可自动完成：环境检查 → 服务启动 → 健康验证 → CSV 同步 → 结果输出。

### 使用方式

```powershell
# 在项目根目录执行
.\scripts\one-click-debug.ps1

# 或带参数：先执行爬虫再启动
.\scripts\one-click-debug.ps1 -RunCrawler -Port 8080

# 或直接双击
scripts\one-click-debug.cmd
```

详见 [`ONE_CLICK_DEBUG.md`](./ONE_CLICK_DEBUG.md)

---

## 🛠️ 开发指南

### 添加新的筛选条件

在 `models/job.go` 的 `JobQueryParams` 中添加字段，然后在 `services/job_service.go` 的 `applyFilters` 方法中实现查询逻辑：

```go
// models/job.go
type JobQueryParams struct {
    // ... 现有字段
    NewField string `form:"new_field"`
}

// services/job_service.go
func (s *JobService) applyFilters(query *gorm.DB, params models.JobQueryParams) *gorm.DB {
    // ... 现有条件
    if params.NewField != "" {
        query = query.Where("new_field = ?", params.NewField)
    }
    return query
}
```

### 添加新的 API 接口

1. 在 `handlers/job_handler.go` 中实现 Handler 函数
2. 在 `main.go` 的路由组中注册新路由
3. 更新 `API.md` 文档

---

## 📦 构建与部署

### 本地构建

```bash
# 编译可执行文件（Linux）
GOOS=linux GOARCH=amd64 go build -o job-backend main.go

# 编译可执行文件（Windows）
GOOS=windows GOARCH=amd64 go build -o job-backend.exe main.go
```

### 生产部署建议

1. **使用反向代理**：Nginx / Caddy 提供 HTTPS 和负载均衡
2. **进程守护**：systemd / supervisor / pm2 确保服务常驻
3. **数据库备份**：定期备份 `jobs.db` 文件
4. **日志收集**：将标准输出重定向到日志文件或使用日志服务

### systemd 示例

```ini
# /etc/systemd/system/job-backend.service
[Unit]
Description=Job Backend Service
After=network.target

[Service]
Type=simple
User=www
WorkingDirectory=/opt/job-backend
Environment="SERVER_PORT=8080"
Environment="DATABASE_PATH=/data/jobs.db"
ExecStart=/opt/job-backend/job-backend
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

---

## 🤝 项目关系

本项目是"智能招聘数据分析平台"的核心后端，与以下模块协同工作：

```
┌──────────────────────┐
│   前端可视化平台      │  ← Vue/React 数据大屏
│   (数据分析展示)      │
└──────────┬───────────┘
           │ REST API
           ▼
┌──────────────────────┐
│   Job Backend        │  ← 本项目：数据管理 + API 服务
│   (Go + SQLite)      │
└──────────┬───────────┘
           │ CSV / 定时调用
           ▼
┌──────────────────────┐
│   职达工作流          │  ← Python 爬虫子模块
│   (爬虫 + 数据清洗)   │
└──────────────────────┘
```

---

## 📄 开源协议

本项目采用 [MIT License](LICENSE) 开源协议。

---

## 🙋 常见问题

**Q: 数据库文件在哪里？**  
A: 默认在项目根目录生成 `jobs.db`，可通过 `DATABASE_PATH` 环境变量修改。

**Q: CSV 文件格式要求？**  
A: 必须包含以下表头（中文）：`单位名称`、`单位规模`、`单位性质`、`单位行业`、`招聘类型`、`学位要求`、`招聘岗位数`、`岗位`、`专业要求`、`岗位要求`、`学校官网链接`、`工作城市`、`截止日期`、`投递口`、`年薪`

**Q: 如何修改定时任务时间？**  
A: 设置 `CRON_SCHEDULE` 环境变量，使用标准 Cron 表达式。例如 `0 0 9 * * *` 表示每天 9:00。

**Q: 支持其他数据库吗？**  
A: 当前使用 SQLite 轻量部署，如需切换 MySQL/PostgreSQL，只需修改 `database/db.go` 中的 GORM 驱动配置即可。

---

> 💡 **提示**：如遇到其他问题，欢迎提交 [Issue](https://github.com/moonlight-boat/bishe/issues) 或查阅 [`API.md`](./API.md) 获取更详细的技术文档。
