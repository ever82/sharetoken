# ShareToken 标准开发流程

> 基于 Issue #1/2/3 开发经验总结的标准流程，供后续 AI 开发参考

---

## 流程概览

```
┌─────────────────────────────────────────────────────────────────┐
│  Phase 1: 任务准备                                               │
│  ├── 1.1 创建 GitHub Issue                                       │
│  ├── 1.2 创建验收文档 (docs/achievements/tasks/issue-XXX.md)    │
│  └── 1.3 技术方案确认                                            │
├─────────────────────────────────────────────────────────────────┤
│  Phase 2: 开发实现                                               │
│  ├── 2.1 代码实现                                                │
│  ├── 2.2 单元测试 (TDD)                                         │
│  └── 2.3 本地验证                                                │
├─────────────────────────────────────────────────────────────────┤
│  Phase 3: 验收测试                                               │
│  ├── 3.1 自动化测试覆盖                                          │
│  ├── 3.2 人工验收测试                                            │
│  └── 3.3 问题修复迭代                                            │
├─────────────────────────────────────────────────────────────────┤
│  Phase 4: 任务完结                                               │
│  ├── 4.1 更新验收文档                                            │
│  ├── 4.2 移动到 done 目录                                        │
│  └── 4.3 关闭 GitHub Issue                                       │
└─────────────────────────────────────────────────────────────────┘
```

---

## Phase 1: 任务准备

### 1.1 创建 GitHub Issue

```bash
gh issue create \
  --title "[PX] ACH-DEV-XXX: Title" \
  --label "PX" \
  --body "$(cat <<'EOF'
## 验收标准
- [ ] 标准1
- [ ] 标准2

## 关联 Spec
SPEC-XXX
EOF
)"
```

**经验要点**:
- Issue 标题格式: `[优先级] ACH-DEV-编号: 简短描述`
- 必须包含验收标准清单
- 必须标注关联的 Spec 文档

### 1.2 创建验收文档

在 `docs/achievements/tasks/issue-XXX.md` 创建文档，模板：

```markdown
# Issue #X: ACH-DEV-XXX 标题

## 验收标准
1. 标准1
2. 标准2

## 自动化测试覆盖

### ✅ 已覆盖
- [x] 测试项1

### ⚠️ 部分覆盖
- [~] 测试项2

### ❌ 未覆盖（需人工验收）
- [ ] 测试项3

## 实际文件清单
| 文件 | 状态 | 说明 |
|------|------|------|
| path/to/file | ✅ 存在 | 说明 |

## 技术实现细节
...

## 验收总结
| 验收项 | 结果 | 说明 |
|--------|------|------|
| 项1 | ✅ 通过 | 说明 |

## 最终验收结论
**ACH-DEV-XXX** 验收完成度：**XX%**

**关键成果**:
1. ✅ 成果1
2. ✅ 成果2

**遗留项**:
- ⏭️ 延后项（记录到 postponed.md）
```

**经验要点** (Issue #1/2/3):
- Issue #1: 详细列出 CI/CD 配置文件清单
- Issue #2: 分阶段验证（单节点→多节点→UPnP）
- Issue #3: 区分代码实现和运行时测试

### 1.3 技术方案确认

开发前必须确认：
1. **技术选型** - 使用什么工具/库（如 goupnp for UPnP）
2. **实现范围** - 哪些功能必须现在做，哪些可以延后
3. **验收标准** - 明确的通过/失败标准

**决策记录示例** (Issue #2 UPnP):
- 方案A: 简单端口映射 → 选择方案B: 完整 UPnP 实现
- 延后项: 区块浏览器（需部署后测试）

---

## Phase 2: 开发实现

### 2.1 代码实现

**开发顺序** (基于 Issue #1/2/3 经验):

1. **先写测试** (TDD)
   ```bash
   # 创建测试脚本
   scripts/test_xxx.sh

   # 运行测试
   go test ./...
   ```

2. **核心功能实现**
   - 遵循项目已有代码风格
   - 添加必要的错误处理
   - 关键代码添加注释

3. **配置文件更新**
   - 新增配置项要同步更新示例配置
   - 文档中要说明配置含义

**经验要点**:
- Issue #1: 先创建本地测试脚本，验证通过后再配置 GitHub Actions
- Issue #2: 多节点网络先修复启动脚本，再验证网络通信
- Issue #3: 前端代码与 CLI 功能并行开发，CLI 先验证

### 2.2 单元测试

**测试覆盖率要求**:
- P0 功能：核心逻辑必须有单元测试
- P1 功能：关键路径有测试
- P2/P3 功能：可选

**测试文件位置**:
```
x/module/
├── keeper/
│   ├── keeper.go
│   └── keeper_test.go  # 对应测试
```

### 2.3 本地验证

**验证清单** (根据 Issue #1/2/3):

```bash
# 1. 代码编译
make build

# 2. 单元测试
make test

# 3. Lint 检查
make lint

# 4. 功能验证（根据issue类型）
# Issue #1 类型: CI/CD
./scripts/test_cicd.sh

# Issue #2 类型: 网络功能
./scripts/devnet_multi.sh
# 验证: 节点启动、区块产出、P2P连接

# Issue #3 类型: 交易功能
./bin/sharetokend tx bank send ...
# 验证: 交易成功、余额更新
```

---

## Phase 3: 验收测试

### 3.1 自动化测试覆盖

在验收文档中记录：

```markdown
### ✅ 已覆盖
- [x] 配置文件存在性检查
- [x] 单元测试通过
- [x] 构建成功

### ⚠️ 部分覆盖
- [~] 需要特定环境的测试

### ❌ 未覆盖（需人工验收）
- [ ] 需要浏览器/移动设备的测试
- [ ] 需要真实网络环境的测试（UPnP）
```

**经验要点** (Issue #2 UPnP):
- 代码实现可以自动化验证
- 但实际 UPnP 功能需要真实路由器环境
- 决策：代码完成即标记进度，实际测试通过后更新状态

### 3.2 人工验收测试
所谓人工验收也是需要AI claude code自己先想尽办法代表人去验收，实在不行的先设置成延后处理（不至于卡在那里等人）

**人工验收步骤**:
1. 按照验收文档中的"测试步骤"执行
2. 记录实际输出
3. 对比预期结果
4. 更新验收状态

### 3.3 问题修复迭代

**发现问题的处理流程**:

```
发现问题 → 记录到验收文档 → 分析原因 → 修复 → 重新验证
```

**经验要点** (Issue #2):
- 原始问题: UPnP 实现复杂
- 修复过程: 使用 goupnp 库简化
- 验证: 实际测试通过后才标记完成

**修复记录格式** (验收文档中):
```markdown
### 已修复问题 ✅
1. ~~问题描述~~ - **已修复**：
   - 修复措施
   - 验证结果
```

---

## Phase 4: 任务完结

### 4.1 更新验收文档

完成时必须更新：

```markdown
## 最终验收结论
**ACH-DEV-XXX** 验收完成度：**100%**

| 验收项 | 状态 | 完成度 |
|--------|------|--------|
| 项1 | ✅ 通过 | 100% |
| 项2 | ⏭️ 延后 | postponed.md |

**遗留项**:
- ⏭️ 延后项说明（记录到 postponed.md）
```

### 4.2 移动到 done 目录

```bash
# 任务完成后执行
mv docs/achievements/tasks/issue-XXX.md docs/achievements/done/
```

**目录结构**:
```
docs/achievements/
├── done/           # 已完成的 issues
│   ├── issue-001.md
│   ├── issue-002.md
│   └── issue-003.md
├── tasks/          # 进行中的 issues（空表示无进行中任务）
├── postponed.md    # 延后项记录
└── for-dev.md      # 总开发计划
```

### 4.3 关闭 GitHub Issue

```bash
gh issue close XXX --comment "$(cat <<'EOF'
✅ **验收完成**

所有验收项已通过：
- ✅ 标准1
- ✅ 标准2

详见本地文档：`docs/achievements/done/issue-XXX.md`

备注：XXX（如有延后项）
EOF
)"
```

---

## 关键决策点

### 何时延后实现？

参考 Issue #2/3 的经验，以下情况应延后：

| 情况 | 处理方式 | 示例 |
|------|---------|------|
| 需要部署后环境 | 记录到 postponed.md | 区块浏览器 |
| 需要特定硬件/软件 | 记录到 postponed.md | WalletConnect 移动端测试 |
| 不影响核心功能 | 记录到 postponed.md | 前端运行时测试 |

**延后项记录格式** (postponed.md):
```markdown
## Issue #X: ACH-DEV-XXX

### 延后项

#### 1. 延后项名称

**原始验收项**: XXX

**当前状态**: ⏭️ 延后实现

**延后原因**: XXX

**前置条件**:
- [ ] 条件1

**计划实现阶段**: ACH-DEV-XXX

**验收标准（待实现）**:
- [ ] 标准1
```

---

## 工具与命令速查

### 常用命令

```bash
# GitHub Issue 管理
gh issue create --title "[PX] ACH-DEV-XXX: Title" --label "PX"
gh issue list --state open
gh issue close XXX --comment "完成评论"

# 本地测试
make build
make test
make lint

# 开发网络
./scripts/devnet_multi.sh
./scripts/devnet_multi.sh status
./scripts/devnet_multi.sh stop

# 交易测试
./bin/sharetokend query bank balances <address>
./bin/sharetokend tx bank send <from> <to> <amount>
```

### 通知命令

```bash
# 完成任务
Bash(echo "CC-NOTIFY: [完成] ACH-DEV-XXX - 描述")

# 阻塞问题
Bash(echo "CC-NOTIFY: [阻塞] 问题描述")

# 需要确认
Bash(echo "CC-NOTIFY: [确认] 需要确认的操作")
```

---

## 经验总结

### Issue #1 (CI/CD Pipeline) 经验
- **关键成功因素**: 先本地脚本验证，再配置 GitHub Actions
- **教训**: 初始脚本要考虑路径和依赖问题
- **最佳实践**: 创建 `test_cicd.sh` 本地测试脚本

### Issue #2 (Blockchain Network) 经验
- **关键成功因素**: 分阶段验证（单节点→多节点→UPnP）
- **教训**: UPnP 需要真实环境测试，不能仅依赖代码
- **最佳实践**: 延后无法本地验证的项到部署阶段

### Issue #3 (Wallet & Token) 经验
- **关键成功因素**: CLI 先验证核心功能，前端代码并行完成
- **教训**: 浏览器/移动端测试需要特定环境，应延后
- **最佳实践**: 代码实现与运行时测试分开验收

---

## 附录：验收文档模板

见 `docs/knowledges/issue-template.md`

---

*最后更新: 2026-03-10*
*基于 Issue #1/2/3 开发经验总结*
