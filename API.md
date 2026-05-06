# 职达招聘信息API接口文档

## 基础信息
- **Base URL**: `http://localhost:8080/api`
- **Content-Type**: `application/json`

## 接口列表

### 1. 获取招聘信息列表
```http
GET /api/jobs
```

**查询参数**:
- `page` (int): 页码，默认1
- `size` (int): 每页大小，默认20，最大100
- `company_name` (string): 公司名称，模糊查询
- `company_type` (string): 公司类型，精确匹配
- `industry` (string): 行业，精确匹配
- `recruitment_type` (string): 招聘类型
- `degree_requirement` (string): 学历要求
- `work_city` (string): 工作城市
- `position_name` (string): 岗位名称，模糊查询
- `keyword` (string): 关键词搜索，在多个字段中搜索
- `sort_by` (string): 排序字段，默认created_at
- `sort_order` (string): 排序方向，asc/desc，默认desc

**响应**:
```json
{
  "jobs": [
    {
      "id": 1,
      "company_name": "广东世纪晓教育科技有限公司",
      "company_scale": "5000-10000人",
      "company_type": "其他企业（含民营企业等）",
      "industry": "教育",
      "recruitment_type": "校园招聘",
      "degree_requirement": "本科、硕士",
      "position_count": "100",
      "position_name": "教育/培训",
      "major_requirement": "计算机相关专业",
      "job_description": "岗位职责...",
      "official_link": "https://job.hitwh.edu.cn/...",
      "work_city": "杭州|温州|台州",
      "deadline": "2026-06-30",
      "application_method": "邮箱投递",
      "salary": "10-18万",
      "created_at": "2024-11-30T12:00:00Z",
      "updated_at": "2024-11-30T12:00:00Z"
    }
  ],
  "total": 330,
  "page": 1,
  "size": 20,
  "total_pages": 17
}
```

### 2. 获取单个招聘信息
```http
GET /api/jobs/{id}
```

**响应**: 返回单个招聘信息对象

### 3. 搜索招聘信息
```http
GET /api/jobs/search
```

**查询参数**: 同获取列表接口，主要使用`keyword`参数进行全文搜索

### 4. 获取统计信息
```http
GET /api/jobs/stats
```

**响应**:
```json
{
  "total": 330,
  "company_type_stats": {
    "国有企业": 45,
    "其他企业（含民营企业等）": 285
  },
  "industry_stats": {
    "制造业": 89,
    "信息传输、软件和信息技术服务业": 67,
    "教育": 34
  },
  "degree_stats": {
    "本科": 156,
    "硕士": 98,
    "博士": 23
  },
  "city_stats": {
    "北京": 45,
    "上海": 38,
    "深圳": 32
  },
  "recruitment_stats": {
    "校园招聘": 280,
    "社会招聘": 50
  }
}
```

### 5. 手动触发同步
```http
POST /api/sync
```

**查询参数**:
- `csv_path` (string): CSV文件路径，默认 `../job-scheduler/auto_job_list.csv`

**响应**:
```json
{
  "success": true,
  "message": "同步任务已启动"
}
```

### 6. 获取同步状态
```http
GET /api/sync/status
```

**响应**:
```json
{
  "sync_in_progress": false,
  "last_sync_time": "2024-11-30T12:00:00Z",
  "csv_info": {
    "exists": true,
    "size": 1024000,
    "modified_time": "2024-11-30T11:55:00Z",
    "record_count": 330
  },
  "db_record_count": 330
}
```

### 7. 健康检查
```http
GET /health
```

**响应**:
```json
{
  "status": "ok",
  "message": "Job Backend Service is running"
}
```

## 错误响应格式
```json
{
  "error": "错误类型",
  "message": "具体错误信息"
}
```

## 常用查询示例

**按公司类型筛选**:
```
GET /api/jobs?company_type=国有企业&page=1&size=10
```

**按城市和学历筛选**:
```
GET /api/jobs?work_city=北京&degree_requirement=硕士
```

**关键词搜索**:
```
GET /api/jobs/search?keyword=软件开发&size=20
```

**按薪资排序**:
```
GET /api/jobs?sort_by=salary&sort_order=desc
```
