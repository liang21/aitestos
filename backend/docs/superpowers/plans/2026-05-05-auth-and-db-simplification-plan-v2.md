# 认证与数据库架构简化实施计划 v2

> **技术评审意见：**
> - 增加风险评估和回滚计划
> - 聚焦验收标准而非实现细节
> - 补充安全性和兼容性考虑
> - 优化任务粒度和依赖关系

---

## 执行摘要

**目标：** 简化认证机制和数据库架构，降低系统复杂度

**变更范围：**
- 认证：access token 有效期 15分钟 → 7 天
- 数据库：删除 soft delete、project prefix、module abbreviation、test case number
- 表名：project → projects, module → modules

**风险评估：**
- 🔴 高风险：数据库表结构变更（需要回滚计划）
- 🟡 中风险：API 破坏性变更（影响客户端兼容性）
- 🟢 低风险：认证 token 延长（需要安全评估）

**预计时间：** 2-3 天（包含测试和验证）

---

## 阶段 0：前置准备（新增）

### Goal: 确保变更安全可控，有完善的回滚机制

### Task 0.1: 风险评估与影响分析

**验收标准：**
- [ ] 完成下游系统影响分析报告
- [ ] 评估 token 延长的安全风险与缓解措施
- [ ] 确认所有依赖 prefix/abbreviation/number 的场景已识别

**执行要点：**
- 使用 `grep -r "Prefix\|Abbreviation\|CaseNumber"` 代码扫描
- 检查所有导出功能、报表、API 响应
- 评估 token 泄露后的影响窗口（7天 vs 15分钟）

### Task 0.2: 数据库回滚方案

**验收标准：**
- [ ] 创建回滚脚本 `004_rollback.sql`
- [ ] 在测试环境验证回滚脚本可用
- [ ] 文档化回滚触发条件

**回滚脚本模板：**
```sql
BEGIN;
-- 004_rollback.sql
-- 1. 恢复表名
ALTER TABLE projects RENAME TO project;
ALTER TABLE modules RENAME TO module;

-- 2. 恢复列
ALTER TABLE project ADD COLUMN prefix varchar(4);
ALTER TABLE module ADD COLUMN abbreviation varchar(4);
ALTER TABLE test_case ADD COLUMN number varchar(50);

-- 3. 恢复索引
-- 根据需要恢复索引

COMMIT;
```

### Task 0.3: 兼容性策略

**验收标准：**
- [ ] 确定 API 版本策略（v1 → v2？或保持 v1 但废弃字段）
- [ ] 制定客户端迁移指南
- [ ] 设置兼容期时间表

**建议：**
- 保持 API v1，标记字段为 deprecated
- 新字段使用 optional，提供迁移过渡期
- 考虑使用 GraphQL 逐步迁移

---

## 阶段 1：文档更新（简化）

### Task 1.1: 更新数据库 Schema

**验收标准：**
- [ ] `specs/aitemos_optimized.sql` 反映新表结构
- [ ] `specs/openapi.yaml` 移除废弃字段
- [ ] 功能规范更新用户故事

### Task 1.2: 创建迁移脚本

**验收标准：**
- [ ] `004_remove_prefix_deleted_at_rename_tables.sql` 语法正确
- [ ] `004_rollback.sql` 回滚脚本可用
- [ ] 在测试环境验证迁移无错误

---

## 阶段 2：领域层重构

### Task 2.1: 清理值对象

**验收标准：**
- [ ] 删除 `prefix.go`, `abbreviation.go`, `case_number.go`
- [ ] 更新 Project/Module/TestCase 聚合根
- [ ] 领域层测试全部通过

### Task 2.2: 更新错误定义

**验收标准：**
- [ ] 移除前缀/缩写/编号相关错误
- [ ] 更新错误码映射
- [ ] 验证错误处理逻辑完整

---

## 阶段 3：服务层适配

### Task 3.1: 更新认证服务

**验收标准：**
- [ ] access token 有效期改为 7 天
- [ ] 更新相关文档说明安全考虑
- [ ] 认证测试通过

**安全考虑：**
- 评估是否需要增加设备指纹/IP 绑定
- 考虑 token 刷新策略的安全边界

### Task 3.2: 更新业务服务

**验收标准：**
- [ ] 项目/模块/测试用例服务移除编号生成逻辑
- [ ] 服务层测试全部通过
- [ ] API 响应不包含已删除字段

---

## 阶段 4：持久层迁移

### Task 4.1: Repository 层更新

**验收标准：**
- [ ] 所有 Repository 使用新表名
- [ ] SQL 查询移除已删除列
- [ ] Repository 测试通过

**注意：** 集成测试需要 Docker，如果环境不可用，可使用 mock 测试验证逻辑

### Task 4.2: 数据迁移执行

**验收标准：**
- [ ] 在测试环境执行迁移
- [ ] 验证数据完整性
- [ ] 执行回滚验证

---

## 阶段 5：传输层适配

### Task 5.1: 更新 HTTP Handler

**验收标准：**
- [ ] Request/Response 结构体移除废弃字段
- [ ] Handler 测试通过
- [ ] API 文档更新

### Task 5.2: 兼容性处理

**验收标准：**
- [ ] 旧版客户端仍可正常工作（兼容期内）
- [ ] 新客户端不依赖废弃字段
- [ ] 弃弃字段在响应中标记为 deprecated

---

## 阶段 6：验证与发布

### Task 6.1: 全面测试

**验收标准：**
- [ ] `go test ./...` 全部通过
- [ ] `go vet ./...` 无警告
- [ ] `go build ./cmd/...` 成功
- [ ] 集成测试通过（如有环境）

### Task 6.2: 代码审查与发布

**验收标准：**
- [ ] 代码审查通过
- [ ] 提交 PR 并合并
- [ ] 标注发布说明
- [ ] 部署到测试环境验证

---

## 附录

### A. 安全评估 Checklist

- [ ] Token 延长到 7 天的影响：
  - [ ] 如果 token 泄露，攻击者有 7 天访问权限 vs 15 分钟
  - [ ] 缓解措施：____（如 IP 限制、设备指纹、异常检测）
  - [ ] 是否需要增加 token 撤销功能？

### B. API 兼容性策略

**选项 1：保持字段但标记废弃**
```json
{
  "id": "...",
  "name": "...",
  "prefix": null,    // deprecated, will be removed in v2
  "abbreviation": null  // deprecated
}
```

**选项 2：API 版本升级**
- `/api/v1/...` 保留旧字段
- `/api/v2/...` 移除废弃字段

**推荐：** 选项 1（渐进式迁移）

### C. 回滚决策树

```
遇到问题时评估：
1. 数据丢失风险？ → 立即回滚
2. 功能缺失但数据完整？→ 评估修复时间 vs 回滚成本
3. 性能问题？→ 不回滚，优化
4. 兼容性问题？→ 修复客户端代码
```

### D. 提交策略优化

**原计划问题：** 23 个小提交，过于碎片

**优化后：**
1. 领域层重构：1-2 个 commits
2. 服务层适配：1-2 个 commits  
3. 持久层变更：1 commit（含迁移脚本）
4. 传输层适配：1-2 个 commits
5. 测试修复：1 commit

**总计：** 约 6-8 个有意义的 commits

---

## 变更总结

**原计划：** 23 个任务，4-6 小时，过度细节

**优化后：** 6 个阶段，约 2-3 天，聚焦验收标准

**关键改进：**
1. ✅ 增加风险评估和回滚机制
2. ✅ 简化任务粒度，聚焦目标
3. ✅ 补充安全性和兼容性考虑
4. ✅ 优化提交策略，减少碎片化
5. ✅ 明确验收标准，便于验证
