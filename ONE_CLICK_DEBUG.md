# 一键联调方案

## 目标

通过一次命令完成以下动作：

1. 设置后端运行环境变量
2. 启动 `job-backend` 服务
3. 等待健康检查通过
4. 自动触发一次 CSV 同步
5. 打印联调结果（岗位总数、统计总数）

## 脚本位置

- PowerShell脚本：`scripts/one-click-debug.ps1`
- Windows双击入口：`scripts/one-click-debug.cmd`

## 使用方式

在项目根目录执行：

```powershell
.\scripts\one-click-debug.ps1
```

或直接双击：

```text
scripts\one-click-debug.cmd
```

## 可选参数

- `-RunCrawler`：先执行一次 `职达工作流/codes/main.py` 再启动后端
- `-Port 8080`：指定后端端口（默认 8080）

示例：

```powershell
.\scripts\one-click-debug.ps1 -RunCrawler -Port 8080
```

## 默认路径规则

脚本默认使用：

- Python脚本：`职达工作流/codes/main.py`
- CSV文件：`职达工作流/auto_job_list.csv`

若默认CSV不存在，自动回退到：

- `职达工作流/for_end_ui/codes/auto_job_list.csv`

## 成功标志

脚本最后出现：

- `=== 一键联调完成 ===`
- 打印后端地址、CSV路径、健康检查和接口验证地址

## 常见问题

1. **后端未启动成功**
   - 查看脚本自动打开的新 PowerShell 窗口日志。
2. **CSV不存在**
   - 先在 `职达工作流` 目录生成 `auto_job_list.csv`。
3. **Python命令不可用（仅在 `-RunCrawler` 时）**
   - 先确认 `python --version` 可执行。
