# Spec 追溯规则（强制）

**Last Updated**: 2026-03-03

---

## 核心原则

> **所有 AI 生成的文档和代码必须有 spec 依据**
> **区分"人定的"和"AI填的"，追溯依赖链直到源头 spec**

---

## Spec 目录结构

```
.spec/
├── product.md       # @product/* 产品规范
├── concepts.md      # @concepts/* 核心概念
├── architecture.md  # @architecture/* 架构规范
└── (future)         # 其他领域规范
```

### Spec ID 格式

```
@<domain>/<topic>-<sequence>

示例：
@product/vision-001
@concepts/entity-002
@architecture/tech-001
```

---

## 生成任何内容前必做

1. **扫描 `.spec/` 目录**，找到相关的 spec
2. **确认 spec 约束是什么**
3. **在生成的内容中标注 spec 来源**

---

## 生成任何内容后必做

1. **检查是否标注了 spec**
2. **检查 derived_from 是否完整**
3. **检查内容是否违反 spec 约束**

---

## Spec 标注格式

### MD 文件（YAML frontmatter）

```markdown
---
spec:
  source: @product/vision-001
  derived_from: []
  generated_by: claude-opus-4
  date: 2026-03-03
---

# 文档标题
...
```

### TS 文件（文件头注释）

```typescript
/**
 * @spec @architecture/tech-001
 * @derived_from docs/design/auth-flow.md
 * @generated_by claude-opus-4
 * @date 2026-03-03
 *
 * AI填充说明：
 * - 具体的错误码设计是AI根据最佳实践填充的
 * - 日志格式遵循 @concepts/workflow-001
 */
```

---

## 标注字段说明

| 字段 | 必填 | 说明 |
|------|------|------|
| `source` | ✅ | spec ID，来自 `.spec/` 目录 |
| `derived_from` | ✅ | 依赖的上游 AI 生成文件（空数组 `[]` 表示无依赖） |
| `generated_by` | ✅ | 生成内容的模型名 |
| `date` | ✅ | 生成日期 |

---

## 禁止事项

- ❌ 无 spec 依据生成核心逻辑
- ❌ spec 标注不完整就提交
- ❌ 内容违反 spec 约束

---

## 例外情况

以下内容**可以不标注 spec**：
- 纯工具函数
- 测试代码
- 配置文件

但如果有对应 spec，仍建议标注。

---

## 依赖链追溯示例

```
人类意图（.spec/）
└── @product/vision-001 "用户需要安全登录"
    │
    ├── docs/prd/auth-system.md
    │   spec: @product/vision-001
    │   derived_from: []
    │
    ├── docs/design/auth-flow.md
    │   spec: @product/vision-001
    │   derived_from: [docs/prd/auth-system.md]
    │
    └── src/services/AuthService.ts
        spec: @architecture/tech-001
        derived_from: [docs/design/auth-flow.md]
```

**追溯能力**：
- 修改 spec → AI 知道要更新整个链条
- 发现代码问题 → 追溯到设计 → PRD → 需求
- AI "幻觉"检测 → 检查内容是否有 spec 依据

---

## AI Agent 行为规则

### 生成阶段

```
1. 接收任务
2. 扫描 .spec/ 找到相关规范，根据任务修改规范文档  ← 强制
3. 生成内容
4. 给内容添加用到的 spec 标注           ← 强制
5. 自检：spec 完整？约束满足？ ← 强制
6. 输出
```

### 验证阶段

```
1. 检查 derived_from 文件是否存在
2. 检查 spec source 是否有效
3. 检查是否有循环依赖
4. 检查内容是否违反 spec 约束
```