# Phase 03: Service Provider Plugins - Research

**Researched:** 2026-03-03
**Domain:** AI Agent OS, API Key Security, Agent Orchestration, WASM Sandboxing
**Confidence:** HIGH

## Summary

本研究针对 ShareTokens 服务提供者插件的三个核心组件进行技术调研：LLM API Key 托管插件、Agent 执行器插件（基于 OpenFang）、以及 Workflow 执行器插件。

**Primary recommendation:**
- **OpenFang 已可直接使用** - 它是一个成熟的 Rust Agent OS，具备 16 层安全防护、7 个内置 Hands、40+ 渠道适配器，可作为 Agent 执行器的核心基础
- **API Key 安全存储** - 使用 Rust `ring` 或 `aes-gcm` crate 实现 AEAD 加密，配合 HashiCorp Vault 或本地加密存储
- **Workflow 编排** - 结合 OpenFang 的 Hands 系统与 CrewAI/LangGraph 的多 Agent 协作模式

---

## Standard Stack

### Core (Rust 生态)

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| **OpenFang** | v0.1.0 | Agent Operating System | 成熟的 Rust Agent OS，16 层安全，7 个内置 Hands，MIT 许可证 |
| **ring** | 0.17+ | AEAD 加密 (AES-256-GCM, ChaCha20-Poly1305) | FIPS 验证，高性能，广泛使用 |
| **aes-gcm** | 0.10+ | AES-GCM 加密 | RustCrypto 官方，纯 Rust 实现 |
| **wasmtime** | 35.0+ | WASM 运行时沙箱 | Bytecode Alliance 官方，安全隔离 |
| **octocrab** | 0.43+ | GitHub API 客户端 | 现代、可扩展的 Rust GitHub API 库 |
| **vaultrs** | 0.7+ | HashiCorp Vault 客户端 | 企业级密钥管理 |

### Supporting (加密与安全)

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| **orion** | 0.17+ | 密码学原语库 | 需要更多加密算法时 |
| **zeroize** | 1.8+ | 内存安全擦除 | API Key 使用后自动清零 |
| **subtle** | 2.6+ | 常量时间操作 | 防止时序攻击 |
| **rand** | 0.9+ | 安全随机数生成 | 密钥和 IV 生成 |
| **base64** | 0.22+ | Base64 编解码 | API Key 存储 |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| OpenFang | OpenClaw | OpenClaw 功能较少，无自主调度，安全层级仅 3 层 |
| OpenFang | LangGraph/AutoGen | Python 生态，性能较低，无内置沙箱 |
| ring | aws-lc-rs | AWS 实现，功能类似但更偏向 AWS 生态 |
| wasmtime | wasmer | Wasmer 功能类似，但 Wasmtime 由 Bytecode Alliance 维护 |
| vaultrs | 本地加密存储 | 本地存储适合小型部署，Vault 适合企业级 |

**Installation:**

```toml
# Cargo.toml for LLM Provider Plugin
[dependencies]
ring = "0.17"
aes-gcm = "0.10"
zeroize = "1.8"
rand = "0.9"
base64 = "0.22"
serde = { version = "1.0", features = ["derive"] }
tokio = { version = "1.0", features = ["full"] }

# Cargo.toml for Agent/Workflow Executor Plugin
[dependencies]
# OpenFang integration (需要作为子模块或通过 git 引入)
# wasmtime for sandboxing
wasmtime = "35.0"
wasmtime-wasi = "35.0"

# GitHub integration
octocrab = "0.43"

# Secret management
vaultrs = "0.7"
```

---

## 类似项目分析

### 1. LLM API Key 托管插件

#### 类似项目

| 项目 | 链接 | 相似之处 | 可学习架构 |
|------|------|----------|------------|
| **LiteLLM** | [github.com/BerriAI/litellm](https://github.com/BerriAI/litellm) | 多提供商 API 代理、统一接口、负载均衡 | 统一 API 格式、成本追踪、速率限制 |
| **AI-PROXY** | [aceproxy.xyz](https://aceproxy.xyz) | 多供应商网关、密钥保护、管理后台 | 团队协作、管理界面设计 |
| **openai-http-proxy** | [pypi.org/project/openai-http-proxy](https://pypi.org/project/openai-http-proxy/) | 可配置路由、多 LLM 支持 | TOML 配置模式、路由规则设计 |
| **Claude-Proxy** | [github.com/yinxulai/claude-proxy](https://github.com/yinxulai/claude-proxy) | API 格式转换、流式响应 | 多提供商格式统一 |

#### 可复用代码库

| 库名 | 链接 | 许可证 | 成熟度 | 特别关注 |
|------|------|--------|--------|----------|
| **ring** | [briansmith/ring](https://github.com/briansmith/ring) | ISC/OpenSSL | 3.5k stars, 活跃 | AEAD 加密，FIPS 验证 |
| **aes-gcm** | [RustCrypto/AEADs](https://github.com/RustCrypto/AEADs) | MIT/Apache-2.0 | 活跃 | 纯 Rust AES-GCM 实现 |
| **vaultrs** | [jkelleyrtp/vaultrs](https://github.com/jkelleyrtp/vaultrs) | MIT | 活跃 | HashiCorp Vault Rust 客户端 |

### 2. Agent 执行器插件 (基于 OpenFang)

#### OpenFang 详细分析

| 属性 | 详情 |
|------|------|
| **GitHub** | [RightNow-AI/openfang](https://github.com/RightNow-AI/openfang) |
| **官网** | [openfang.sh](https://www.openfang.sh) |
| **许可证** | MIT (免费可商用) |
| **Stars** | 6,211+ (4 天内) |
| **语言** | 纯 Rust |
| **代码量** | 137,000+ 行，14 个 crates |
| **测试** | 1,767+ 测试，零 clippy 警告 |
| **版本** | v0.1.0 (2026年2月24日) |

#### OpenFang 核心能力

**7 个内置 Hands (自主能力包):**

| Hand | 功能 | ShareTokens 用例 |
|------|------|------------------|
| **Researcher** | 深度研究，CRAAP 评估，APA 引用 | 想法评估、市场分析 |
| **Lead** | 销售线索发现，ICP 匹配，评分 | 资源匹配服务 |
| **Collector** | OSINT 情报，知识图谱构建 | 竞品监控、数据收集 |
| **Clip** | 视频下载、剪辑、字幕、发布 | 内容创作工作流 |
| **Predictor** | 超级预测，置信区间，Brier 评分 | 趋势预测 |
| **Twitter** | 自主 Twitter 账号管理 | 社交媒体自动化 |
| **Browser** | Web 自动化，Playwright 桥接 | 自动化工作流 |

**16 层安全系统:**

| # | 系统 | 功能 |
|---|------|------|
| 1 | WASM 双计量沙箱 | 工具代码在 WASM 中运行，燃料计量 + 周期中断 |
| 2 | Merkle 哈希链审计 | 每个操作加密链接，篡改检测 |
| 3 | 信息流污点追踪 | 秘密从源头到汇的追踪 |
| 4 | Ed25519 签名代理清单 | 代理身份和能力集加密签名 |
| 5 | SSRF 防护 | 阻止私有 IP、云元数据端点 |
| 6 | 秘密归零 | `Zeroizing<String>` 自动擦除 API Key |
| 7 | OFP 互认证 | HMAC-SHA256 随机数验证 |
| 8 | 能力门控 | 基于角色的访问控制 |
| 9 | 安全头 | CSP, X-Frame-Options, HSTS |
| 10 | 健康端点编辑 | 公开健康检查最小化信息 |
| 11 | 子进程沙箱 | `env_clear()` + 选择性变量传递 |
| 12 | Prompt 注入扫描 | 检测覆盖尝试、数据泄露模式 |
| 13 | 循环防护 | SHA256 工具调用循环检测 |
| 14 | 会话修复 | 7 阶段消息历史验证 |
| 15 | 路径遍历防护 | 规范化 + 符号链接防护 |
| 16 | GCRA 速率限制 | 成本感知令牌桶 |

#### 与 OpenClaw 对比

| 特性 | OpenFang | OpenClaw |
|------|----------|----------|
| **冷启动** | 180ms | 5.98s |
| **空闲内存** | 40MB | 394MB |
| **安装大小** | 32MB | 500MB |
| **安全层级** | 16 | 3 |
| **渠道适配器** | 40 | 13 |
| **自主调度** | Yes | No |
| **WASM 沙箱** | Yes | No |

#### 可复用代码库

| 库名 | 链接 | 许可证 | 成熟度 | 特别关注 |
|------|------|--------|--------|----------|
| **OpenFang** | [RightNow-AI/openfang](https://github.com/RightNow-AI/openfang) | MIT | 6k+ stars | 完整 Agent OS，可直接集成 |
| **wasmtime** | [bytecodealliance/wasmtime](https://github.com/bytecodealliance/wasmtime) | Apache-2.0 | 15k+ stars | WASM 运行时沙箱 |

### 3. Workflow 执行器插件

#### 类似项目

| 项目 | 链接 | 相似之处 | 可学习架构 |
|------|------|----------|------------|
| **CrewAI** | [crewAIInc/crewAI](https://github.com/crewAIInc/crewAI) | 多 Agent 协作、角色-任务映射 | Manager-Worker-Reviewer 模式 |
| **LangGraph** | [langchain-ai/langgraph](https://github.com/langchain-ai/langgraph) | 图工作流、状态机 | 复杂工作流编排 |
| **Deer-Flow** | [bytedance/deer-flow](https://github.com/bytedance/deer-flow) | 多 Agent 协作框架 | 21k+ stars，字节跳动开源 |
| **Claude-Flow** | GitHub Trending | 企业级自主工作流编排 | 智能编排平台 |

#### Agent 编排框架对比

| 框架 | 语言 | Stars | 特点 | 适用场景 |
|------|------|-------|------|----------|
| **LangChain** | Python | 89k+ | 最成熟，丰富生态 | 单 Agent + 工具调用 |
| **CrewAI** | Python | 25k+ | 角色协作，低学习曲线 | 多 Agent 内容生成 |
| **LangGraph** | Python | 活跃 | 图工作流，状态检查点 | 复杂工作流 |
| **OpenFang Hands** | Rust | 6k+ | 自主调度，16 层安全 | 生产级自主任务 |

#### GitHub API 集成

| 库名 | 链接 | 许可证 | 成熟度 | 特别关注 |
|------|------|--------|--------|----------|
| **octocrab** | [XAMPPRocky/octocrab](https://github.com/XAMPPRocky/octocrab) | MIT/Apache-2.0 | 活跃 (v0.43) | 现代 GitHub API 客户端 |

---

## Architecture Patterns

### 推荐项目结构

```
plugins/
├── llm-provider/                    # P02: LLM Provider Plugin (Rust)
│   ├── Cargo.toml
│   ├── src/
│   │   ├── lib.rs
│   │   ├── key_management/
│   │   │   ├── mod.rs
│   │   │   ├── store.rs            # AEAD 加密存储
│   │   │   ├── rotate.rs           # 密钥轮换
│   │   │   └── audit.rs            # 使用审计
│   │   ├── providers/
│   │   │   ├── mod.rs
│   │   │   ├── base.rs             # Provider trait
│   │   │   ├── openai.rs
│   │   │   ├── anthropic.rs
│   │   │   └── google.rs
│   │   ├── routing/
│   │   │   ├── mod.rs
│   │   │   ├── router.rs           # 请求路由
│   │   │   └── load_balancer.rs
│   │   ├── billing/
│   │   │   ├── mod.rs
│   │   │   └── meter.rs            # Token 计量
│   │   └── security/
│   │       ├── mod.rs
│   │       └── encryption.rs       # AES-256-GCM
│   └── tests/
│
├── agent-provider/                  # P03: Agent Provider Plugin (Rust + OpenFang)
│   ├── Cargo.toml
│   ├── openfang/                    # OpenFang 子模块或引用
│   ├── src/
│   │   ├── lib.rs
│   │   ├── openfang_adapter/
│   │   │   ├── mod.rs
│   │   │   ├── client.rs           # OpenFang SDK 客户端
│   │   │   └── hands/
│   │   │       ├── researcher.rs
│   │   │       ├── coder.rs
│   │   │       └── writer.rs
│   │   ├── task_executor/
│   │   │   ├── mod.rs
│   │   │   ├── state_machine.rs    # 任务状态机
│   │   │   ├── scheduler.rs
│   │   │   └── monitor.rs
│   │   ├── security/
│   │   │   ├── mod.rs
│   │   │   ├── sandbox.rs          # WASM 沙箱
│   │   │   └── limits.rs           # 资源限制
│   │   └── reporting/
│   │       ├── mod.rs
│   │       └── metrics.rs
│   └── tests/
│
└── workflow-provider/               # P04: Workflow Provider Plugin (Rust)
    ├── Cargo.toml
    ├── src/
    │   ├── lib.rs
    │   ├── workflows/
    │   │   ├── mod.rs
    │   │   ├── base.rs             # 基础 Workflow trait
    │   │   ├── software_dev.rs     # 软件开发工作流
    │   │   ├── content_creation.rs # 内容创作工作流
    │   │   └── business.rs         # 商业规划工作流
    │   ├── executor/
    │   │   ├── mod.rs
    │   │   ├── runner.rs           # 工作流运行器
    │   │   ├── state_machine.rs    # 工作流状态机
    │   │   └── recovery.rs         # 故障恢复
    │   ├── nodes/
    │   │   ├── mod.rs
    │   │   ├── auto.rs             # 自动节点
    │   │   ├── human_gate.rs       # 人工审批节点
    │   │   └── external.rs         # 外部服务节点
    │   ├── github_integration/
    │   │   ├── mod.rs
    │   │   ├── client.rs           # octocrab 客户端
    │   │   ├── repo_manager.rs
    │   │   └── pr_handler.rs
    │   └── monitoring/
    │       ├── mod.rs
    │       └── progress.rs
    └── tests/
```

### Pattern 1: AEAD 加密存储 (API Key 托管)

**What:** 使用 AES-256-GCM 或 ChaCha20-Poly1305 加密 API Keys，支持认证加密

**When to use:** 存储敏感 API Keys、用户凭证

**Example:**

```rust
// Source: RustCrypto AEADs + ring documentation
use aes_gcm::{
    aead::{Aead, KeyInit, OsRng},
    Aes256Gcm, Nonce,
};
use zeroize::Zeroizing;
use rand::RngCore;

/// 加密的 API Key 存储
pub struct EncryptedKeyStore {
    cipher: Aes256Gcm,
}

#[derive(serde::Serialize, serde::Deserialize)]
pub struct StoredKey {
    pub provider: String,
    pub ciphertext: Vec<u8>,
    pub nonce: Vec<u8>,
    pub created_at: chrono::DateTime<chrono::Utc>,
    pub last_used: Option<chrono::DateTime<chrono::Utc>>,
}

impl EncryptedKeyStore {
    pub fn new(master_key: &[u8; 32]) -> Self {
        let cipher = Aes256Gcm::new_from_slice(master_key)
            .expect("Invalid key length");
        Self { cipher }
    }

    /// 加密并存储 API Key
    pub fn store_key(
        &self,
        provider: &str,
        api_key: &str,
    ) -> Result<StoredKey, KeyStoreError> {
        // 生成随机 nonce (96-bit)
        let mut nonce_bytes = [0u8; 12];
        OsRng.fill_bytes(&mut nonce_bytes);
        let nonce = Nonce::from_slice(&nonce_bytes);

        // 加密 (Zeroizing 确保密钥在使用后被擦除)
        let key_bytes = Zeroizing::new(api_key.as_bytes().to_vec());
        let ciphertext = self.cipher
            .encrypt(nonce, key_bytes.as_ref())
            .map_err(|_| KeyStoreError::EncryptionFailed)?;

        Ok(StoredKey {
            provider: provider.to_string(),
            ciphertext,
            nonce: nonce_bytes.to_vec(),
            created_at: chrono::Utc::now(),
            last_used: None,
        })
    }

    /// 解密 API Key
    pub fn retrieve_key(
        &self,
        stored: &StoredKey,
    ) -> Result<Zeroizing<String>, KeyStoreError> {
        let nonce = Nonce::from_slice(&stored.nonce);
        let plaintext = self.cipher
            .decrypt(nonce, stored.ciphertext.as_ref())
            .map_err(|_| KeyStoreError::DecryptionFailed)?;

        // 返回 Zeroizing 包装，确保使用后擦除
        let key_string = String::from_utf8(plaintext)
            .map_err(|_| KeyStoreError::InvalidUtf8)?;
        Ok(Zeroizing::new(key_string))
    }
}
```

### Pattern 2: OpenFang 集成 (Agent 执行器)

**What:** 通过 OpenFang SDK 调用内置 Hands 和自定义 Agents

**When to use:** 执行 Agent 任务、调度自主工作

**Example:**

```rust
// Source: OpenFang architecture documentation
use openfang_kernel::{AgentExecutor, TaskConfig, HandType};
use openfang_runtime::Runtime;

pub struct OpenFangAdapter {
    runtime: Runtime,
}

impl OpenFangAdapter {
    pub async fn new() -> Result<Self, OpenFangError> {
        let runtime = Runtime::new()
            .with_security_layers(16)  // 启用所有 16 层安全
            .with_wasm_sandbox(true)
            .build()
            .await?;

        Ok(Self { runtime })
    }

    /// 激活 Researcher Hand 进行想法评估
    pub async fn evaluate_idea(
        &self,
        idea: &IdeaInput,
    ) -> Result<EvaluationResult, OpenFangError> {
        let task = TaskConfig {
            hand: HandType::Researcher,
            input: serde_json::to_value(idea)?,
            constraints: TaskConstraints {
                max_tokens: 200_000,
                max_duration: Duration::from_secs(300),
                ..Default::default()
            },
        };

        let result = self.runtime.execute_task(task).await?;

        Ok(EvaluationResult {
            value_score: result.metrics.get("value_score").copied(),
            feasibility: result.metrics.get("feasibility").copied(),
            token_estimate: result.metrics.get("token_estimate").copied(),
            report: result.output,
        })
    }

    /// 使用 Coder Agent 执行代码生成
    pub async fn generate_code(
        &self,
        spec: &CodeSpec,
    ) -> Result<CodeResult, OpenFangError> {
        let task = TaskConfig {
            hand: HandType::Custom("coder".to_string()),
            input: serde_json::to_value(spec)?,
            ..Default::default()
        };

        let result = self.runtime.execute_task(task).await?;

        Ok(CodeResult {
            files: result.artifacts,
            tests_passed: result.metrics.get("tests_passed").copied(),
        })
    }
}
```

### Pattern 3: WASM 沙箱隔离 (安全执行)

**What:** 使用 Wasmtime 在隔离环境中执行不可信代码

**When to use:** 执行用户定义的 Agent 技能、工具代码

**Example:**

```rust
// Source: Wasmtime documentation + OpenFang security model
use wasmtime::*;
use wasmtime_wasi::preview1::{self, WasiP1Ctx};

pub struct SecureSandbox {
    engine: Engine,
    module: Module,
}

impl SecureSandbox {
    pub fn new(wasm_bytes: &[u8]) -> Result<Self, SandboxError> {
        // 配置安全引擎
        let mut config = Config::new();
        config
            .cranelift_opt_level(OptLevel::Speed)  // 性能优化
            .consume_fuel(true)                     // 启用燃料计量
            .epoch_interruption(true);              // 启用周期中断

        let engine = Engine::new(&config)?;
        let module = Module::new(&engine, wasm_bytes)?;

        Ok(Self { engine, module })
    }

    /// 在沙箱中执行函数
    pub fn execute(
        &self,
        func_name: &str,
        input: &[u8],
        max_fuel: u64,
    ) -> Result<Vec<u8>, SandboxError> {
        let mut store = Store::new(&self.engine, ());

        // 设置燃料限制 (双计量：燃料 + 周期)
        store.set_fuel(max_fuel)?;

        // 配置 WASI 文件系统隔离
        let wasi = preview1::add_to_store(&mut store, WasiP1Ctx::builder()
            .preopened_dir("/sandbox", "/tmp", Capabilities::empty())?
            .build())?;

        let instance = Instance::new(&mut store, &self.module, &[])?;
        let func = instance.get_typed_func::<(i32, i32), i32>(&mut store, func_name)?;

        // 执行并检查超时
        let result = func.call(&mut store, (input.as_ptr() as i32, input.len() as i32))?;

        // 检查剩余燃料
        let remaining = store.get_fuel()?;
        if remaining == 0 {
            return Err(SandboxError::FuelExhausted);
        }

        Ok(/* 结果 */)
    }
}
```

### Pattern 4: GitHub 集成工作流 (PR 创建)

**What:** 使用 octocrab 创建仓库、生成代码、提交 PR

**When to use:** 软件开发工作流、自动化部署

**Example:**

```rust
// Source: octocrab documentation
use octocrab::{Octocrab, models};
use octocrab::params::repos::Reference;

pub struct GitHubWorkflow {
    client: Octocrab,
}

impl GitHubWorkflow {
    pub async fn new(token: &str) -> Result<Self, GitHubError> {
        let client = Octocrab::builder()
            .personal_token(token.to_string())
            .build()?;

        Ok(Self { client })
    }

    /// 创建 PR 并提交代码
    pub async fn create_pr_with_code(
        &self,
        owner: &str,
        repo: &str,
        branch: &str,
        files: Vec<CodeFile>,
        pr_title: &str,
        pr_body: &str,
    ) -> Result<PullRequest, GitHubError> {
        // 1. 获取默认分支的最新 commit
        let repo_info = self.client.repos(owner, repo).get().await?;
        let default_branch = repo_info.default_branch.unwrap_or_else(|| "main".to_string());

        let ref_info = self.client
            .repos(owner, repo)
            .get_ref(&Reference::Branch(default_branch.clone()))
            .await?;

        // 2. 创建新分支
        self.client
            .repos(owner, repo)
            .create_reference(
                &format!("refs/heads/{}", branch),
                ref_info.object.sha,
            )
            .await?;

        // 3. 提交文件
        for file in &files {
            self.client
                .repos(owner, repo)
                .create_file(
                    file.path.as_str(),
                    &format!("Add {}", file.path),
                    file.content.as_bytes(),
                    &Reference::Branch(branch.to_string()),
                )
                .await?;
        }

        // 4. 创建 PR
        let pr = self.client
            .pulls(owner, repo)
            .create(pr_title, branch, &default_branch)
            .body(pr_body)
            .send()
            .await?;

        Ok(pr)
    }
}
```

### Anti-Patterns to Avoid

1. **硬编码 API Keys:** 不要将 API Keys 硬编码在代码或配置文件中，必须使用 AEAD 加密存储
2. **明文传输凭证:** 所有 API Key 传输必须通过 TLS，内存中使用 `Zeroizing` 包装
3. **无沙箱执行外部代码:** 用户定义的 Agent 技能必须在 WASM 沙箱中执行
4. **忽略 OpenFang 安全层:** 不要绕过 OpenFang 的 16 层安全机制
5. **单一密钥加密所有数据:** 每个用户/提供商应使用独立的密钥派生

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| API Key 加密 | 自定义加密算法 | `ring` / `aes-gcm` crate | 经过验证，防止侧信道攻击 |
| Agent 调度系统 | 自己写调度器 | OpenFang 内核 | 16 层安全、成熟稳定 |
| WASM 沙箱 | 自己实现隔离 | `wasmtime` | Bytecode Alliance 维护 |
| GitHub API 调用 | 手动 HTTP 请求 | `octocrab` | 处理分页、认证、错误 |
| 密钥管理 | 本地文件存储 | HashiCorp Vault (`vaultrs`) | 企业级审计、轮换、访问控制 |
| 多 Agent 协作 | 自己写编排逻辑 | CrewAI 模式 + OpenFang Hands | 角色协作模式已验证 |
| 状态机 | 手动状态管理 | OpenFang 状态机 / LangGraph | 检查点、恢复、审计 |

**Key insight:** 在 Agent OS 领域，OpenFang 已经提供了生产级的解决方案。与其从零开始构建 Agent 执行器，不如直接集成 OpenFang 并扩展其 Hands 系统。安全和可靠性是最难正确实现的部分，必须使用经过验证的库。

---

## Common Pitfalls

### Pitfall 1: API Key 泄露

**What goes wrong:** API Key 被意外提交到 Git、记录到日志、或通过错误消息泄露

**Why it happens:** 开发便利性考虑、缺乏安全意识、CI/CD 配置错误

**How to avoid:**
1. 使用 `Zeroizing<String>` 包装所有 API Key
2. 永远不要将 `.env` 文件提交到 Git
3. 在日志中过滤敏感字段
4. 使用 HashiCorp Vault 或类似的密钥管理系统

**Warning signs:** 日志中出现 `sk-` 或 `xoxb-` 前缀，Git diff 显示密钥变更

### Pitfall 2: WASM 沙箱逃逸

**What goes wrong:** 恶意 WASM 代码突破沙箱，访问主机资源

**Why it happens:** 配置错误、WASI 权限过大、未设置资源限制

**How to avoid:**
1. 使用 OpenFang 的双计量沙箱（燃料 + 周期中断）
2. 严格限制 WASI preopens 目录
3. 禁用危险的系统调用（`personality`, `madvise`）
4. 启用 eBPF 监控异常行为

**Warning signs:** WASM 模块执行时间异常长、尝试访问未授权路径

### Pitfall 3: Agent 无限循环

**What goes wrong:** Agent 在工具调用之间无限循环，消耗大量 Token 和时间

**Why it happens:** 缺乏循环检测、工具返回值设计问题

**How to avoid:**
1. 使用 OpenFang 的 SHA256 循环检测（第 13 层安全）
2. 设置最大工具调用次数限制（默认 50 次）
3. 实现断路器模式
4. 监控 ping-pong 调用模式

**Warning signs:** 任务执行时间超过 30 分钟、Token 使用量异常

### Pitfall 4: GitHub API 速率限制

**What goes wrong:** 频繁的 API 调用触发 GitHub 速率限制，工作流失败

**Why it happens:** 未实现重试逻辑、批量操作未优化

**How to avoid:**
1. 使用 `octocrab` 的内置速率限制处理
2. 实现指数退避重试
3. 使用 GraphQL API 批量获取数据
4. 缓存频繁访问的资源

**Warning signs:** HTTP 403 响应、`X-RateLimit-Remaining` 接近 0

### Pitfall 5: 跨插件状态不一致

**What goes wrong:** LLM Provider、Agent Provider、Workflow Provider 之间的状态不同步

**Why it happens:** 缺乏统一的状态管理、异步操作竞态

**How to avoid:**
1. 使用 OpenFang 的 Merkle 哈希链审计（第 2 层安全）
2. 实现事件溯源模式
3. 关键操作使用分布式锁
4. 定期状态同步检查

**Warning signs:** 计费金额与实际使用不符、任务状态显示异常

---

## Code Examples

### 完整的 LLM Provider Plugin 初始化

```rust
// plugins/llm-provider/src/lib.rs
use ring::aead::{Aad, BoundKey, Nonce, SealingKey, OpeningKey, UnboundKey, AES_256_GCM};
use zeroize::Zeroizing;
use serde::{Deserialize, Serialize};

pub struct LLMProviderPlugin {
    key_store: EncryptedKeyStore,
    router: RequestRouter,
    billing: BillingEngine,
}

impl LLMProviderPlugin {
    pub async fn new(config: PluginConfig) -> Result<Self, PluginError> {
        // 从安全存储加载主密钥
        let master_key = Self::load_master_key(&config.key_path)?;

        Ok(Self {
            key_store: EncryptedKeyStore::new(&master_key),
            router: RequestRouter::new(&config.providers),
            billing: BillingEngine::new(&config.billing),
        })
    }

    /// 处理 LLM 请求
    pub async fn process_request(
        &self,
        request: LLMRequest,
    ) -> Result<LLMResponse, PluginError> {
        // 1. 路由到合适的提供商
        let provider = self.router.select_provider(&request)?;

        // 2. 检索加密的 API Key
        let stored_key = self.key_store.get(&provider.name)?;
        let api_key = self.key_store.retrieve_key(&stored_key)?;

        // 3. 执行请求
        let response = provider.execute(&api_key, request.clone()).await?;

        // 4. 记录使用量
        self.billing.record_usage(UsageMetrics {
            tokens_used: response.usage.total_tokens,
            provider: provider.name.clone(),
            request_id: request.id,
        });

        Ok(response)
    }
}
```

### OpenFang Agent 任务执行

```rust
// plugins/agent-provider/src/openfang_adapter/client.rs
use openfang_kernel::{Kernel, Task, Hand, SecurityConfig};
use openfang_runtime::Runtime;

pub struct OpenFangClient {
    kernel: Kernel,
}

impl OpenFangClient {
    pub async fn execute_agent_task(
        &self,
        agent_type: AgentType,
        task_input: TaskInput,
    ) -> Result<TaskResult, AgentError> {
        // 配置安全沙箱
        let security = SecurityConfig {
            wasm_sandbox: true,
            fuel_limit: 100_000,
            network_restricted: true,
            file_access: FileAccess::SandboxOnly,
        };

        // 创建任务
        let task = Task::new()
            .with_hand(match agent_type {
                AgentType::Coder => Hand::Custom("coder"),
                AgentType::Researcher => Hand::Researcher,
                AgentType::Writer => Hand::Custom("writer"),
            })
            .with_input(task_input)
            .with_security(security)
            .with_timeout(Duration::from_secs(300));

        // 执行并监控
        let execution = self.kernel.submit_task(task).await?;

        // 等待完成或超时
        let result = execution.wait().await?;

        // 获取审计日志 (Merkle 链)
        let audit_trail = self.kernel.get_audit_trail(result.task_id)?;

        Ok(TaskResult {
            output: result.output,
            tokens_used: result.metrics.tokens,
            execution_time: result.metrics.duration,
            audit_hash: audit_trail.root_hash,
        })
    }
}
```

### Workflow 人工审批门

```rust
// plugins/workflow-provider/src/nodes/human_gate.rs
use std::sync::Arc;
use tokio::sync::mpsc;

pub struct HumanGate {
    gate_id: String,
    timeout: Duration,
    reminders: Vec<Reminder>,
    escalation: Option<EscalationPolicy>,
}

impl HumanGate {
    pub async fn wait_for_approval(
        &self,
        context: WorkflowContext,
    ) -> Result<GateResponse, WorkflowError> {
        let (tx, mut rx) = mpsc::channel(1);

        // 发送初始通知
        self.notify_user(&context).await?;

        // 设置超时
        tokio::select! {
            // 用户响应
            response = rx.recv() => {
                match response {
                    Some(GateResponse::Approved) => Ok(GateResponse::Approved),
                    Some(GateResponse::Rejected(reason)) => {
                        Err(WorkflowError::UserRejected(reason))
                    }
                    None => Err(WorkflowError::ChannelClosed),
                }
            }

            // 超时处理
            _ = tokio::time::sleep(self.timeout) => {
                // 触发升级
                if let Some(ref policy) = self.escalation {
                    self.escalate(policy, &context).await?;
                }
                Err(WorkflowError::Timeout)
            }
        }
    }
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Python Agent 框架 | Rust Agent OS (OpenFang) | 2026 | 15x 性能提升，16 层安全 |
| 明文 API Key 存储 | AEAD 加密 + Vault | 2024-2025 | 防止密钥泄露 |
| 单 Agent 执行 | 多 Agent 协作 (Hands) | 2025 | 自主调度，24/7 运行 |
| Docker 隔离 | WASM 沙箱 | 2024-2025 | 更轻量，更快启动 |
| 简单日志 | Merkle 审计链 | 2025 | 不可篡改的审计追踪 |

**Deprecated/outdated:**

- **OpenClaw:** 被 OpenFang 取代，无自主调度，安全层级仅 3 层
- **简单 Base64 编码:** 不是加密，已完全不安全
- **无沙箱 Agent 执行:** 存在严重安全风险

---

## Open Questions

1. **OpenFang 许可证兼容性**
   - What we know: OpenFang 使用 MIT 许可证，可商用
   - What's unclear: 是否需要开源衍生作品
   - Recommendation: MIT 允许闭源商用，可直接集成

2. **多租户隔离**
   - What we know: OpenFang 有 16 层安全
   - What's unclear: 多用户场景下的隔离粒度
   - Recommendation: 在 OpenFang 之上实现租户级别的密钥隔离

3. **跨链支付结算**
   - What we know: 需要 STT 代币结算
   - What's unclear: 与 Cosmos SDK 的集成方式
   - Recommendation: 通过 x/compute 模块触发链上结算事件

---

## Validation Architecture

> 注意: 此部分需要根据实际测试框架配置

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Rust built-in tests + tokio-test |
| Config file | `Cargo.toml` |
| Quick run command | `cargo test --lib` |
| Full suite command | `cargo test --workspace` |

### Phase Requirements -> Test Map
| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| LLM-PLUGIN-01 | 安全存储 API Keys | unit | `cargo test key_management::store::tests` | No - Wave 0 |
| LLM-PLUGIN-02 | 监控 API Key 使用量 | unit | `cargo test billing::tests` | No - Wave 0 |
| LLM-PLUGIN-03 | 多提供商支持 | integration | `cargo test providers::integration_tests` | No - Wave 0 |
| AGENT-PLUGIN-01 | OpenFang 安装配置 | integration | `cargo test openfang_adapter::tests` | No - Wave 0 |
| AGENT-PLUGIN-02~04 | Agent 任务执行 | integration | `cargo test task_executor::tests` | No - Wave 0 |
| AGENT-PLUGIN-09 | WASM 沙箱安全 | security | `cargo test security::sandbox::tests` | No - Wave 0 |
| WF-PLUGIN-01~06 | GitHub 集成工作流 | e2e | `cargo test github_integration::tests` | No - Wave 0 |

### Sampling Rate
- **Per task commit:** `cargo test --lib`
- **Per wave merge:** `cargo test --workspace`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps
- [ ] `plugins/llm-provider/tests/key_management_tests.rs`
- [ ] `plugins/agent-provider/tests/openfang_integration_tests.rs`
- [ ] `plugins/workflow-provider/tests/github_workflow_tests.rs`
- [ ] Framework install: `cargo install cargo-nextest` - if using nextest

---

## Sources

### Primary (HIGH confidence)
- [OpenFang GitHub](https://github.com/RightNow-AI/openfang) - Agent OS architecture, security model, Hands system
- [ring crate documentation](https://briansmith.github.io/ring/) - AEAD encryption APIs
- [Wasmtime documentation](https://docs.wasmtime.dev/) - WASM runtime and security
- [octocrab crate](https://docs.rs/octocrab) - GitHub API client

### Secondary (MEDIUM confidence)
- [LiteLLM GitHub](https://github.com/BerriAI/litellm) - API proxy patterns
- [CrewAI documentation](https://docs.crewai.com) - Multi-agent orchestration
- [WASI Security Analysis](https://blog.csdn.net/2501_91980039/article/details/148089256) - WASM sandbox security

### Tertiary (LOW confidence)
- [AI-PROXY](https://aceproxy.xyz) - Multi-vendor gateway (需要验证)
- [vaultrs](https://github.com/jkelleyrtp/vaultrs) - Vault client (需要验证 API 兼容性)

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - OpenFang, ring, wasmtime 都是成熟项目
- Architecture: HIGH - 基于 OpenFang 官方架构文档
- Pitfalls: MEDIUM - 部分来自社区经验，需要实际验证

**Research date:** 2026-03-03
**Valid until:** 30 days - OpenFang 发展迅速，需定期更新
