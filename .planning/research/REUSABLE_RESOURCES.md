# ShareTokens 可复用资源研究报告

> **研究目的**: 逐个模块分析可参考项目和可用代码库，避免重复造轮子
> **研究时间**: 2026-03-03
> **有效期**: 6个月（至2026-09-03）

---

## 目录

1. [核心模块 (Core Modules)](#一核心模块-core-modules)
2. [Provider 插件模块](#二provider-插件模块)
3. [前端与基础设施](#三前端与基础设施)
4. [快速启动技术组合](#四快速启动技术组合)
5. [开发量节省估算](#五开发量节省估算)
6. [下一步行动建议](#六下一步行动建议)

---

## 一、核心模块 (Core Modules)

### 1.1 Identity 模块 (x/identity) - 身份认证

**职责**: 实名认证、身份哈希存储、女巫攻击防护

| 类别 | 项目/库 | 说明 | 链接 |
|------|---------|------|------|
| **可参考项目** | Hypersign | Cosmos SDK构建的SSI身份网络，W3C认可 | [hypersign.id](https://hypersign.id) |
| | Microsoft Entra Verified ID | 开源SDK实现可验证凭证 | [Microsoft](https://azure.microsoft.com/solutions/blockchain/) |
| **可用代码库** | Cosmos SDK x/auth | 账户管理，直接使用 | [cosmos-sdk](https://github.com/cosmos/cosmos-sdk) |
| | DID Methods | W3C DID标准实现 | [w3c.github.io/did](https://w3c.github.io/did/) |

**快速实现建议**:
- Cosmos SDK Auth模块 + 自定义验证逻辑
- 预估代码量: ~200行Go代码

---

### 1.2 Escrow 模块 (x/escrow) - 托管支付

**职责**: 支付托管、释放、争议锁定、裁决分配

| 类别 | 项目/库 | 说明 |
|------|---------|------|
| **可参考项目** | Coral Protocol | Solana上的托管合约，支持条件释放、自动退款 |
| | FISCO BCOS | 企业级金融区块链，有成熟的存证合约模板 |
| **可用代码库** | Cosmos SDK x/bank | 代币转账，直接使用 |
| | Cosmos SDK Authz | 授权模块，可用于托管授权 |

**参考代码模式** (Solidity风格，可移植到Cosmos):
```solidity
contract Escrow {
    address public buyer;
    address public seller;
    bool public delivered;

    function confirmDelivery() external {
        require(msg.sender == buyer, "Only buyer can confirm");
        delivered = true;
        payable(seller).transfer(address(this).balance);
    }
}
```

**快速实现建议**:
- 参考Coral的Escrow设计模式
- 预估代码量: ~300行Go代码

---

### 1.3 Service Market 模块 (x/compute) - 服务市场

**职责**: 服务注册、发现、定价、路由（三层服务）

| 类别 | 项目/库 | 说明 |
|------|---------|------|
| **可参考项目** | Akash Network | Cosmos上的去中心化计算市场 |
| | Osmosis | Cosmos上的DEX，有成熟的交易撮合逻辑 |
| **可用代码库** | Cosmos SDK | 基础框架 |
| | IBC Module | 跨链服务调用 |

**这是核心业务模块**，需要定制开发
- 预估代码量: ~2000行Go代码

---

### 1.4 Trust System 模块 (x/trust) - 信誉+仲裁

**职责**: MQ评分、争议仲裁、陪审团机制、零和重分配

| 类别 | 项目/库 | 说明 |
|------|---------|------|
| **可参考项目** | Kleros | 去中心化仲裁服务，陪审团机制 |
| | Aragon Court | DAO争议解决 |
| | Shentu Network | Cosmos上的安全评分系统 |
| **可用代码库** | Cosmos SDK x/staking | 可参考验证人选择逻辑 |

**快速实现建议**:
- MQ评分算法需定制
- 仲裁可参考Kleros设计
- 预估代码量: ~1600行Go代码

---

## 二、Provider 插件模块

### 2.1 LLM Provider (P02) - API聚合网关

**职责**: 托管API密钥、请求路由、Token计量

#### ⭐ 强烈推荐: LiteLLM

| 属性 | 值 |
|------|-----|
| GitHub Stars | 20k+ |
| 支持模型 | 100+ LLM APIs |
| 协议 | OpenAI格式兼容 |
| 功能 | 统一API、负载均衡、成本追踪、限流 |

**其他开源网关**:

| 项目 | 说明 |
|------|------|
| OneAPI | Docker一键部署，支持所有主流模型 |
| Bella OpenAPI | 类OpenRouter，支持chat/embedding/TTS/ASR |
| NyaProxy | 智能API管理器，负载均衡 |

**技术选型**:
```
方案A（最快）: 直接使用 LiteLLM Proxy Server - 零代码
方案B（自托管）: OneAPI + Docker Compose - 1天部署
方案C（企业级）: 自建网关 + LiteLLM SDK
```

**快速启动** (LiteLLM):
```bash
# 安装
pip install litellm

# 启动代理服务器
litellm --model gpt-3.5-turbo --api_key sk-xxx

# 或使用Docker
docker run -p 4000:4000 ghcr.io/berriai/litellm:main-latest
```

---

### 2.2 Agent Provider (P03) - AI Agent框架

**职责**: Agent注册、任务执行、沙箱隔离

#### 主流Agent框架对比

| 框架 | Stars | 特点 | 适用场景 |
|------|-------|------|----------|
| **LangChain** | 100k+ | 生态最大，工具链丰富，Python/JS双支持 | 通用Agent开发 |
| **LangGraph** | - | 状态图工作流，适合复杂Agent | 复杂状态流转 |
| **AutoGen (微软)** | 40k+ | 多Agent协作框架 | 多Agent协作 |
| **CrewAI** | 活跃 | 角色扮演型Agent | 团队协作场景 |
| **MetaGPT** | 活跃 | 软件开发（PM+架构师+工程师） | 软件自动化 |
| **OpenHands** | 活跃 | 软件开发自动化 | 代码生成 |
| **Semantic Kernel (微软)** | - | 轻量级，Azure集成好 | 企业级应用 |

**技术选型建议**:
```
推荐组合: LangChain + LangGraph + OpenFang
- LangChain: 基础Agent能力（工具调用、记忆）
- LangGraph: 复杂状态流转
- OpenFang: 按项目需求集成
```

---

### 2.3 Workflow Provider (P04) - 工作流引擎

**职责**: 工作流注册、执行、状态管理、失败恢复

#### 主流工作流引擎对比

| 引擎 | Stars | 特点 | 适用场景 |
|------|-------|------|----------|
| **n8n** | 79k+ | 可视化工作流，400+集成，AI原生支持 | 快速上线、AI工作流 |
| **Temporal** | 12k+ | 长时任务编排，强一致性，"Workflow-as-Code" | 高可靠性、长时任务 |
| **Prefect** | 活跃 | Airflow现代化替代，函数式API | 数据管道 |
| **Apache Airflow** | 38k+ | 数据管道编排，成熟稳定 | ETL场景 |
| **Dify** | 国内活跃 | AI应用开发平台，可视化工作流 | AI应用快速开发 |
| **Windmill** | 活跃 | 开发者友好，多语言支持 | 内部工具 |

**技术选型建议**:
```
场景推荐：
- 快速上线 + AI工作流 → n8n（低代码，可视化）
- 高可靠性 + 长时任务 → Temporal
- 复杂状态管理 → LangGraph
- 企业级BPM → Camunda
```

**n8n 快速启动**:
```bash
docker run -it --rm --name n8n -p 5678:5678 -v ~/.n8n:/home/node/.n8n n8nio/n8n
```

---

## 三、前端与基础设施

### 3.1 GenieBot Interface - React前端

#### AI聊天界面参考项目

| 项目 | Stars | 技术栈 | 说明 |
|------|-------|--------|------|
| **Chatbot-UI** | 30k+ | Next.js 14, TypeScript, Supabase | 最流行的开源ChatGPT克隆 |
| **OpenWebUI** | 94k+ | Python, Docker | 功能最完整的自托管AI聊天界面 |
| **assistant-ui** | - | React Components | 专门的AI聊天React组件库 |
| **LobeChat** | 活跃 | Next.js | 现代化AI聊天界面，插件化 |

#### Cosmos生态前端参考

| 项目 | 技术栈 | 说明 |
|------|--------|------|
| Keplr Wallet | React, TypeScript | Cosmos生态最大钱包 |
| Ping.pub | Vue.js | 轻量级Cosmos浏览器 |
| Osmosis Frontend | React, TypeScript | 最大的Cosmos DEX前端 |

---

### 3.2 可用UI组件库

#### React聊天组件

| 库名称 | 特点 | 适用场景 |
|--------|------|----------|
| **Vercel AI SDK useChat** | 自动处理流式、状态、错误重试 | 生产级AI聊天首选 |
| **@chatscope/chat-ui-kit-react** | 30+独立组件，虚拟滚动 | 通用聊天应用 |
| **Streamdown** | 专为AI流式输出设计 | AI助手界面 |
| **Ant Design X 2.0** | 即插即用智能聊天组件 | 企业级应用 |

#### Web3钱包连接组件

| 库名称 | 说明 | 支持钱包 |
|--------|------|----------|
| **@cosmos-kit/react** | 统一的钱包连接接口 | Keplr, Cosmostation, Leap, WalletConnect, MetaMask |
| **@cosmos-kit/keplr** | Keplr专用集成 | Keplr Extension, Keplr Mobile |

**推荐**: 使用 Cosmos Kit，比手动集成Keplr节省70%代码量

#### 通用UI组件库

| 库名称 | Stars | 特点 |
|--------|-------|------|
| **shadcn/ui** | 100k+ | 非传统npm包，CLI复制到项目，100%可定制 |
| **Radix UI** | - | 无样式可访问性组件，shadcn/ui的底层 |
| **Tailwind CSS** | - | 原子化CSS |

---

### 3.3 Cosmos SDK 基础设施

#### Ignite CLI 常用命令

| 功能 | 命令 |
|------|------|
| 创建新链 | `ignite scaffold chain github.com/user/chain` |
| 添加模块 | `ignite scaffold module token --dep bank` |
| 添加消息 | `ignite scaffold message create-post title body` |
| 启动测试网 | `ignite chain serve` |
| 构建二进制 | `ignite chain build` |

#### CosmJS 核心包

| 包名 | 功能 |
|------|------|
| `@cosmjs/stargate` | 主客户端（SigningStargateClient） |
| `@cosmjs/tendermint-rpc` | Tendermint RPC客户端 |
| `@cosmjs/proto-signing` | Protobuf签名 |
| `@cosmjs/cosmwasm-stargate` | CosmWasm交互 |

#### Keplr 集成代码示例

```typescript
// 方式1: 通过 Cosmos Kit（推荐）
import { WalletProvider } from '@cosmos-kit/react';
import { wallets as keplrWallets } from '@cosmos-kit/keplr';

const sharetokensChain = {
  chainName: 'ShareTokens',
  chainId: 'sharetokens-1',
  rpc: 'https://rpc.sharetokens.io',
  rest: 'https://api.sharetokens.io',
  stakeCurrency: {
    coinDenom: 'STT',
    coinMinimalDenom: 'ustt',
    coinDecimals: 6,
  },
};

<WalletProvider chains={[sharetokensChain]} wallets={keplrWallets}>
  <App />
</WalletProvider>

// 方式2: 直接集成
if (!window.keplr) {
  window.open('https://www.keplr.app/get');
  return;
}
await window.keplr.enable(chainId);
const offlineSigner = window.getOfflineSigner(chainId);
```

---

## 四、快速启动技术组合

### 推荐技术栈

```
┌─────────────────────────────────────────────────────────────────┐
│                    ShareTokens 快速启动方案                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  核心链 (Go)                                                    │
│  ├── 框架: Cosmos SDK + Ignite CLI                             │
│  ├── P2P: CometBFT (内置)                                      │
│  ├── 身份: x/auth + 自定义验证                                  │
│  ├── 托管: x/bank + 自定义Escrow                                │
│  └── 信誉: 参考Kleros设计                                       │
│                                                                 │
│  LLM Provider                                                   │
│  └── 方案: LiteLLM Proxy (零开发，配置即用)                      │
│                                                                 │
│  Agent Provider                                                 │
│  └── 框架: LangChain + LangGraph                                │
│                                                                 │
│  Workflow Provider                                              │
│  └── 方案A: n8n (快速上线)                                      │
│  └── 方案B: Temporal (高可靠)                                   │
│                                                                 │
│  前端 (React + TypeScript)                                      │
│  ├── 框架: React 18 + Vite                                     │
│  ├── 状态管理: Zustand                                         │
│  ├── UI组件: shadcn/ui                                         │
│  ├── 钱包: Cosmos Kit + Keplr                                  │
│  ├── AI聊天: Vercel AI SDK                                     │
│  └── 链交互: CosmJS                                             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 前端快速启动命令

```bash
# 1. 创建项目
npm create vite@latest geniebot -- --template react-ts
cd geniebot

# 2. 安装核心依赖
npm install zustand @cosmjs/stargate @cosmjs/tendermint-rpc \
  @cosmos-kit/react @cosmos-kit/keplr axios react-router-dom \
  @ai-sdk/react lucide-react

# 3. 安装开发依赖
npm install -D typescript @types/react @types/react-dom vite \
  @vitejs/plugin-react tailwindcss postcss autoprefixer \
  eslint prettier vitest @testing-library/react

# 4. 初始化 shadcn/ui
npx shadcn-ui@latest init

# 5. 添加常用组件
npx shadcn-ui@latest add button card input dialog dropdown-menu toast badge avatar

# 6. 启动开发服务器
npm run dev
```

---

## 五、开发量节省估算

| 模块 | 原始估算 | 使用开源方案后 | 节省 |
|------|----------|---------------|------|
| P2P/Wallet | ~2000行 | 0 (CometBFT+Keplr) | **100%** |
| LLM Provider | ~1500行 | ~100行配置 | **93%** |
| Agent Provider | ~2500行 | ~500行集成 | **80%** |
| Workflow Provider | ~1500行 | ~200行配置 | **87%** |
| 聊天UI | ~1500行 | ~300行定制 | **80%** |
| Identity | ~600行 | ~200行定制 | **67%** |
| Escrow | ~600行 | ~300行定制 | **50%** |
| Service Market | ~2000行 | ~2000行（需开发） | **0%** |
| Trust System | ~1600行 | ~800行定制 | **50%** |
| **总计** | **~13800行** | **~4400行** | **68%** |

---

## 六、下一步行动建议

### 立即可用（零开发）

| 组件 | 方案 | 行动 |
|------|------|------|
| LLM网关 | LiteLLM | `pip install litellm && litellm --model gpt-3.5-turbo` |
| 工作流 | n8n | `docker run -p 5678:5678 n8nio/n8n` |
| 钱包 | Keplr | 安装浏览器扩展 |

### 需要集成

| 组件 | 方案 | 工作量 |
|------|------|--------|
| Agent框架 | LangChain + LangGraph | 1-2周 |
| 前端钱包 | Cosmos Kit | 2-3天 |
| 链交互 | CosmJS | 1周 |

### 需要开发

| 模块 | 说明 | 优先级 |
|------|------|--------|
| Identity验证逻辑 | 实名认证、防女巫 | P1 |
| Escrow合约 | 托管、释放、争议锁定 | P1 |
| Service Market | 服务注册、发现、路由 | P0（核心） |
| MQ评分算法 | 信誉评分、零和重分配 | P1 |

---

## 参考资源

### 官方文档

- [Vercel AI SDK](https://sdk.vercel.ai/docs)
- [Cosmos Kit](https://cosmoskit.com)
- [CosmJS](https://github.com/cosmos/cosmjs)
- [shadcn/ui](https://ui.shadcn.com)
- [Zustand](https://zustand-demo.pmnd.rs/)
- [Ignite CLI](https://ignite.com)
- [LiteLLM](https://github.com/BerriAI/litellm)
- [n8n](https://n8n.io)
- [Temporal](https://temporal.io)

### GitHub参考项目

- [LangChain](https://github.com/langchain-ai/langchain)
- [Chatbot-UI](https://github.com/mckaywrigley/chatbot-ui)
- [Keplr Wallet](https://github.com/chainapsis/keplr-wallet)
- [Ping.pub](https://github.com/ping-pub/explorer)
- [Kleros](https://github.com/kleros)

---

**研究完成时间**: 2026-03-03
**置信度**: HIGH
**下次更新**: 2026-06-03
