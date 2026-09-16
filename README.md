# 流量回放服务

纯 Go 标准库实现的流量录制与回放管理系统。

## 运行方式

```bash
cd origin
/Users/fengyin/.local/go/bin/go build ./cmd/server
./server
```

默认监听 `:8080`，可通过环境变量配置：
- `PORT` 或 `ADDR`：服务地址
- `MAX_PAGE_SIZE`：最大分页大小（默认 100）
- `API_KEY`：API 鉴权密钥（空则不鉴权）
- `LOG_LEVEL`：日志级别（debug/info/warn/error）

## 前端

`origin/web/` 下包含 `index.html`、`app.js`、`style.css`。
启动服务后访问 `http://localhost:8080/` 即可打开管理界面。

## API 列表

| 实体 | 方法 | 路径 | 说明 |
|------|------|------|------|
| RecordingSession | POST | /api/recording-sessions | 创建录制会话 |
| RecordingSession | GET | /api/recording-sessions | 列表查询 |
| RecordingSession | GET | /api/recording-sessions/{id} | 详情 |
| RecordingSession | PUT | /api/recording-sessions/{id} | 更新 |
| RecordingSession | DELETE | /api/recording-sessions/{id} | 删除 |
| RecordingSession | POST | /api/recording-sessions/{id}/start | 开始录制 |
| RecordingSession | POST | /api/recording-sessions/{id}/pause | 暂停录制 |
| RecordingSession | POST | /api/recording-sessions/{id}/stop | 停止录制 |
| TrafficRecord | POST | /api/traffic-records | 创建流量记录 |
| TrafficRecord | GET | /api/traffic-records | 列表查询 |
| TrafficRecord | GET | /api/traffic-records/{id} | 详情 |
| TrafficRecord | PUT | /api/traffic-records/{id} | 更新 |
| TrafficRecord | DELETE | /api/traffic-records/{id} | 删除 |
| TrafficRecord | POST | /api/traffic-records/import | 批量导入 |
| ReplayTask | POST | /api/replay-tasks | 创建回放任务 |
| ReplayTask | GET | /api/replay-tasks | 列表查询 |
| ReplayTask | GET | /api/replay-tasks/{id} | 详情 |
| ReplayTask | PUT | /api/replay-tasks/{id} | 更新 |
| ReplayTask | DELETE | /api/replay-tasks/{id} | 删除 |
| ReplayTask | POST | /api/replay-tasks/{id}/run | 运行回放 |
| ReplayTask | POST | /api/replay-tasks/{id}/complete | 完成任务 |
| ReplayTask | POST | /api/replay-tasks/{id}/fail | 标记失败 |
| ReplayTask | POST | /api/replay-tasks/batch | 批量创建 |
| TargetEnv | POST | /api/target-envs | 创建目标环境 |
| TargetEnv | GET | /api/target-envs | 列表查询 |
| TargetEnv | GET | /api/target-envs/{id} | 详情 |
| TargetEnv | PUT | /api/target-envs/{id} | 更新 |
| TargetEnv | DELETE | /api/target-envs/{id} | 删除 |
| FilterRule | POST | /api/filter-rules | 创建过滤规则 |
| FilterRule | GET | /api/filter-rules | 列表查询 |
| FilterRule | GET | /api/filter-rules/{id} | 详情 |
| FilterRule | PUT | /api/filter-rules/{id} | 更新 |
| FilterRule | DELETE | /api/filter-rules/{id} | 删除 |
| CompareReport | POST | /api/compare-reports | 创建对比报告 |
| CompareReport | GET | /api/compare-reports | 列表查询 |
| CompareReport | GET | /api/compare-reports/{id} | 详情 |
| CompareReport | PUT | /api/compare-reports/{id} | 更新 |
| CompareReport | DELETE | /api/compare-reports/{id} | 删除 |
| CompareReport | POST | /api/compare-reports/{id}/run | 运行对比 |
| DiffRecord | POST | /api/diff-records | 创建差异记录 |
| DiffRecord | GET | /api/diff-records | 列表查询 |
| DiffRecord | GET | /api/diff-records/{id} | 详情 |
| DiffRecord | PUT | /api/diff-records/{id} | 更新 |
| DiffRecord | DELETE | /api/diff-records/{id} | 删除 |
| Schedule | POST | /api/schedules | 创建调度配置 |
| Schedule | GET | /api/schedules | 列表查询 |
| Schedule | GET | /api/schedules/{id} | 详情 |
| Schedule | PUT | /api/schedules/{id} | 更新 |
| Schedule | DELETE | /api/schedules/{id} | 删除 |
| ReplayResult | POST | /api/replay-results | 创建回放结果 |
| ReplayResult | GET | /api/replay-results | 列表查询 |
| ReplayResult | GET | /api/replay-results/{id} | 详情 |
| ReplayResult | PUT | /api/replay-results/{id} | 更新 |
| ReplayResult | DELETE | /api/replay-results/{id} | 删除 |
| AuditLog | POST | /api/audit-logs | 创建审计日志 |
| AuditLog | GET | /api/audit-logs | 列表查询 |
| AuditLog | GET | /api/audit-logs/{id} | 详情 |
| AuditLog | PUT | /api/audit-logs/{id} | 更新 |
| AuditLog | DELETE | /api/audit-logs/{id} | 删除 |
| WebhookConfig | POST | /api/webhook-configs | 创建 Webhook |
| WebhookConfig | GET | /api/webhook-configs | 列表查询 |
| WebhookConfig | GET | /api/webhook-configs/{id} | 详情 |
| WebhookConfig | PUT | /api/webhook-configs/{id} | 更新 |
| WebhookConfig | DELETE | /api/webhook-configs/{id} | 删除 |
| ScriptTemplate | POST | /api/script-templates | 创建脚本模板 |
| ScriptTemplate | GET | /api/script-templates | 列表查询 |
| ScriptTemplate | GET | /api/script-templates/{id} | 详情 |
| ScriptTemplate | PUT | /api/script-templates/{id} | 更新 |
| ScriptTemplate | DELETE | /api/script-templates/{id} | 删除 |
| EnvVariable | POST | /api/env-variables | 创建环境变量 |
| EnvVariable | GET | /api/env-variables | 列表查询 |
| EnvVariable | GET | /api/env-variables/{id} | 详情 |
| EnvVariable | PUT | /api/env-variables/{id} | 更新 |
| EnvVariable | DELETE | /api/env-variables/{id} | 删除 |
| NotificationRule | POST | /api/notification-rules | 创建通知规则 |
| NotificationRule | GET | /api/notification-rules | 列表查询 |
| NotificationRule | GET | /api/notification-rules/{id} | 详情 |
| NotificationRule | PUT | /api/notification-rules/{id} | 更新 |
| NotificationRule | DELETE | /api/notification-rules/{id} | 删除 |
| ProjectConfig | POST | /api/project-configs | 创建项目配置 |
| ProjectConfig | GET | /api/project-configs | 列表查询 |
| ProjectConfig | GET | /api/project-configs/{id} | 详情 |
| ProjectConfig | PUT | /api/project-configs/{id} | 更新 |
| ProjectConfig | DELETE | /api/project-configs/{id} | 删除 |
| Tag | POST | /api/tags | 创建标签 |
| Tag | GET | /api/tags | 列表查询 |
| Tag | GET | /api/tags/{id} | 详情 |
| Tag | PUT | /api/tags/{id} | 更新 |
| Tag | DELETE | /api/tags/{id} | 删除 |
| Certificate | POST | /api/certificates | 创建证书 |
| Certificate | GET | /api/certificates | 列表查询 |
| Certificate | GET | /api/certificates/{id} | 详情 |
| Certificate | PUT | /api/certificates/{id} | 更新 |
| Certificate | DELETE | /api/certificates/{id} | 删除 |
| DataSource | POST | /api/data-sources | 创建数据源 |
| DataSource | GET | /api/data-sources | 列表查询 |
| DataSource | GET | /api/data-sources/{id} | 详情 |
| DataSource | PUT | /api/data-sources/{id} | 更新 |
| DataSource | DELETE | /api/data-sources/{id} | 删除 |
| Stats | GET | /api/stats/overview | 概览统计 |
| Stats | GET | /api/stats/method-distribution | 方法分布 |
| Stats | GET | /api/stats/env-distribution | 环境分布 |
| Stats | GET | /api/stats/endpoint-distribution | 接口分布 |
| Stats | GET | /api/stats/success-rate-trend | 成功率趋势 |
| Export | GET | /api/export/snapshot | 全量快照导出 |

## 中间件

- 请求日志
- Panic recovery
- API Key 鉴权（可选）
- 每 IP 限流（120 次/分钟）
