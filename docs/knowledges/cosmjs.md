# CosmJS 知识文档

CosmJS 是 Cosmos 生态系统中用于 JavaScript/TypeScript 客户端开发的核心库，支持 Web 应用、浏览器扩展和服务器端客户端（如水龙头、爬虫）等多种场景。

## 1. CosmJS 概述

### 什么是 CosmJS？

CosmJS 是 Cosmos 生态系统的 "瑞士军刀"，提供与 Cosmos SDK 区块链交互的完整 TypeScript/JavaScript 解决方案。它支持 Cosmos SDK 0.40+ (Stargate) 版本。

### 核心包结构

| 包名 | 描述 |
|------|------|
| `@cosmjs/stargate` | Cosmos SDK 0.40+ 客户端库 |
| `@cosmjs/proto-signing` | 基于 Protobuf 的签名实现 |
| `@cosmjs/cosmwasm` | CosmWasm 模块客户端 |
| `@cosmjs/crypto` | 密码学工具（哈希、签名、HD密钥派生） |
| `@cosmjs/encoding` | 编码辅助工具 |
| `@cosmjs/math` | 安全整数和金融数值处理 |
| `@cosmjs/tendermint-rpc` | Tendermint/CometBFT RPC 客户端 |
| `@cosmjs/launchpad` | Cosmos SDK 0.39 (Launchpad) 支持 |

### 安装

```bash
npm install @cosmjs/stargate @cosmjs/proto-signing
```

## 2. StargateClient - 只读查询客户端

`StargateClient` 提供与区块链的只读交互，无需钱包即可查询链上数据。

### 连接到区块链

```typescript
import { StargateClient } from "@cosmjs/stargate";

// 连接到 RPC 端点
const rpcEndpoint = "https://rpc.cosmos.network:443";
const client = await StargateClient.connect(rpcEndpoint);

// 验证连接
const chainId = await client.getChainId();
console.log("Connected to chain:", chainId);
```

### 查询余额

```typescript
const address = "cosmos1..."; // 你的地址
const denom = "uatom"; // 代币 denom

// 查询单一代币余额
const balance = await client.getBalance(address, denom);
console.log(`Balance: ${balance.amount} ${balance.denom}`);

// 查询所有余额
const allBalances = await client.getAllBalances(address);
allBalances.forEach((coin) => {
  console.log(`${coin.denom}: ${coin.amount}`);
});
```

### 查询账户信息

```typescript
// 获取账户信息（包含 sequence 和 accountNumber）
const account = await client.getAccount(address);
if (account) {
  console.log("Account number:", account.accountNumber);
  console.log("Sequence:", account.sequence);
  console.log("Pubkey:", account.pubkey);
}
```

### 查询区块和交易

```typescript
// 获取最新区块高度
const height = await client.getHeight();
console.log("Current height:", height);

// 查询交易
const txHash = "ABC123..."; // 交易哈希
const tx = await client.getTx(txHash);
if (tx) {
  console.log("Transaction found:");
  console.log("  Height:", tx.height);
  console.log("  Code:", tx.code); // 0 表示成功
  console.log("  Gas used:", tx.gasUsed);
  console.log("  Gas wanted:", tx.gasWanted);
}

// 搜索交易
const sentTxs = await client.searchTx([
  { key: "message.sender", value: address },
]);
```

## 3. SigningStargateClient - 交易签名客户端

`SigningStargateClient` 扩展了 `StargateClient`，增加了交易签名和广播功能。

### 创建钱包

```typescript
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

// 生成新钱包（12词助记词）
const wallet = await DirectSecp256k1HdWallet.generate(12);

// 生成 24 词助记词（更安全）
const wallet24 = await DirectSecp256k1HdWallet.generate(24);

// 从助记词恢复钱包
const mnemonic = "word1 word2 word3 ... word12";
const recoveredWallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
  prefix: "cosmos", // 地址前缀
});

// 获取账户地址
const accounts = await wallet.getAccounts();
const [firstAccount] = accounts;
console.log("Address:", firstAccount.address);
console.log("Pubkey:", firstAccount.pubkey);

// 获取助记词（仅用于测试，生产环境勿用）
const myMnemonic = wallet.mnemonic;
```

### 连接签名客户端

```typescript
import { SigningStargateClient, GasPrice } from "@cosmjs/stargate";

const rpcEndpoint = "https://rpc.cosmos.network:443";
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
  prefix: "cosmos",
});

// 使用默认 Gas 价格
const client = await SigningStargateClient.connectWithSigner(
  rpcEndpoint,
  wallet
);

// 使用自定义 Gas 价格
const gasPrice = GasPrice.fromString("0.025uatom");
const clientWithGas = await SigningStargateClient.connectWithSigner(
  rpcEndpoint,
  wallet,
  { gasPrice }
);

// 获取账户地址
const [account] = await wallet.getAccounts();
```

## 4. 交易构建与签名

### 发送代币

```typescript
import { coins } from "@cosmjs/stargate";

const recipientAddress = "cosmos1recipient...";
const amount = coins(1000000, "uatom"); // 1 ATOM (1,000,000 uatom)
const memo = "Payment for services";

// 发送代币
const result = await client.sendTokens(
  account.address,
  recipientAddress,
  amount,
  "auto", // 自动估算 Gas
  memo
);

console.log("Transaction broadcast successful!");
console.log("Transaction hash:", result.transactionHash);
console.log("Height:", result.height);
console.log("Gas used:", result.gasUsed);
console.log("Code:", result.code); // 0 表示成功
```

### 自定义 Gas 费用

```typescript
import { calculateFee, GasPrice } from "@cosmjs/stargate";

// 方式1：使用 "auto" 自动估算
const result1 = await client.sendTokens(
  account.address,
  recipientAddress,
  amount,
  "auto",
  memo
);

// 方式2：手动指定 Gas 限制
const gasLimit = 200000;
const fee = calculateFee(gasLimit, "0.025uatom");
const result2 = await client.sendTokens(
  account.address,
  recipientAddress,
  amount,
  fee,
  memo
);

// 方式3：模拟 Gas 并添加余量
const gasEstimation = await client.simulate(
  account.address,
  [
    {
      typeUrl: "/cosmos.bank.v1beta1.MsgSend",
      value: {
        fromAddress: account.address,
        toAddress: recipientAddress,
        amount: amount,
      },
    },
  ],
  memo
);
const adjustedGas = Math.ceil(gasEstimation * 1.3); // 添加 30% 余量
const adjustedFee = calculateFee(adjustedGas, "0.025uatom");
```

### 构建自定义消息

```typescript
import { MsgSendEncodeObject } from "@cosmjs/stargate";

// 构建银行转账消息
const msgSend: MsgSendEncodeObject = {
  typeUrl: "/cosmos.bank.v1beta1.MsgSend",
  value: {
    fromAddress: account.address,
    toAddress: recipientAddress,
    amount: coins(1000000, "uatom"),
  },
};

// 签名并广播
const result = await client.signAndBroadcast(
  account.address,
  [msgSend],
  "auto",
  "Custom transaction"
);
```

### 多消息交易

```typescript
import { MsgDelegateEncodeObject, MsgSendEncodeObject } from "@cosmjs/stargate";

// 构建多个消息
const messages = [
  {
    typeUrl: "/cosmos.bank.v1beta1.MsgSend",
    value: {
      fromAddress: account.address,
      toAddress: recipientAddress,
      amount: coins(500000, "uatom"),
    },
  },
  {
    typeUrl: "/cosmos.staking.v1beta1.MsgDelegate",
    value: {
      delegatorAddress: account.address,
      validatorAddress: "cosmosvaloper1...",
      amount: { denom: "uatom", amount: "500000" },
    },
  },
];

// 一次交易中发送多个消息
const result = await client.signAndBroadcast(
  account.address,
  messages,
  "auto",
  "Multi-message transaction"
);
```

## 5. 查询客户端扩展

CosmJS 使用模块化扩展来查询不同模块的数据。

### 基本查询扩展

```typescript
import {
  StargateClient,
  QueryClient,
  setupBankExtension,
  setupAuthExtension,
  setupStakingExtension,
} from "@cosmjs/stargate";
import { Tendermint34Client } from "@cosmjs/tendermint-rpc";

// 创建底层 Tendermint 客户端
const tmClient = await Tendermint34Client.connect(rpcEndpoint);

// 创建带有扩展的 QueryClient
const queryClient = QueryClient.withExtensions(
  tmClient,
  setupBankExtension,
  setupAuthExtension,
  setupStakingExtension
);

// 使用 Bank 扩展查询
const totalSupply = await queryClient.bank.unverified.totalSupply();
console.log("Total supply:", totalSupply);

// 使用 Staking 扩展查询
const validators = await queryClient.staking.unverified.validators("BOND_STATUS_BONDED");
console.log("Active validators:", validators);
```

### IBC 查询扩展

```typescript
import { setupIbcExtension } from "@cosmjs/stargate";

const queryClient = QueryClient.withExtensions(
  tmClient,
  setupIbcExtension
);

// 查询 IBC 通道
const channels = await queryClient.ibc.unverified.channels();
console.log("IBC channels:", channels);

// 查询连接
const connections = await queryClient.ibc.unverified.connections();
console.log("IBC connections:", connections);
```

## 6. 完整示例：ShareTokens 前端集成

以下是一个适用于 ShareTokens 项目的完整示例：

```typescript
import {
  SigningStargateClient,
  StargateClient,
  GasPrice,
  coins,
} from "@cosmjs/stargate";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";

// ShareTokens 链配置
const SHARETOKENS_CONFIG = {
  rpcEndpoint: "https://rpc.sharetokens.io:443",
  chainId: "sharetokens-1",
  denom: "ushare",
  prefix: "share",
  gasPrice: "0.025ushare",
};

// 创建客户端类
class ShareTokensClient {
  private client: SigningStargateClient | null = null;
  private queryClient: StargateClient | null = null;
  private wallet: DirectSecp256k1HdWallet | null = null;

  // 初始化只读客户端
  async connectReadOnly() {
    this.queryClient = await StargateClient.connect(SHARETOKENS_CONFIG.rpcEndpoint);
    return this.queryClient;
  }

  // 初始化签名客户端
  async connectWithWallet(mnemonic: string) {
    this.wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
      prefix: SHARETOKENS_CONFIG.prefix,
    });

    this.client = await SigningStargateClient.connectWithSigner(
      SHARETOKENS_CONFIG.rpcEndpoint,
      this.wallet,
      { gasPrice: GasPrice.fromString(SHARETOKENS_CONFIG.gasPrice) }
    );

    return this.client;
  }

  // 获取账户地址
  getAddress(): string | null {
    if (!this.wallet) return null;
    const accounts = this.wallet.getAccounts();
    return accounts[0]?.address ?? null;
  }

  // 查询余额
  async getBalance(address: string) {
    if (!this.queryClient) {
      await this.connectReadOnly();
    }
    return this.queryClient!.getBalance(address, SHARETOKENS_CONFIG.denom);
  }

  // 发送代币
  async sendTokens(
    recipientAddress: string,
    amount: string,
    memo?: string
  ) {
    if (!this.client) {
      throw new Error("Client not initialized with wallet");
    }

    const [account] = await this.wallet!.getAccounts();
    const result = await this.client.sendTokens(
      account.address,
      recipientAddress,
      coins(parseInt(amount), SHARETOKENS_CONFIG.denom),
      "auto",
      memo ?? ""
    );

    return {
      txHash: result.transactionHash,
      height: result.height,
      code: result.code,
      gasUsed: result.gasUsed,
      rawLog: result.rawLog,
    };
  }

  // 广播自定义消息
  async broadcastMessage(msg: any, memo?: string) {
    if (!this.client) {
      throw new Error("Client not initialized with wallet");
    }

    const [account] = await this.wallet!.getAccounts();
    const result = await this.client.signAndBroadcast(
      account.address,
      [msg],
      "auto",
      memo ?? ""
    );

    return result;
  }

  // 断开连接
  disconnect() {
    this.client = null;
    this.queryClient = null;
    this.wallet = null;
  }
}

// 使用示例
async function main() {
  const sharetokensClient = new ShareTokensClient();

  // 只读查询
  const queryClient = await sharetokensClient.connectReadOnly();
  const balance = await sharetokensClient.getBalance("share1...");
  console.log("Balance:", balance);

  // 签名交易
  const mnemonic = "your mnemonic words here...";
  await sharetokensClient.connectWithWallet(mnemonic);

  // 发送代币
  const txResult = await sharetokensClient.sendTokens(
    "share1recipient...",
    "1000000",
    "Test transfer"
  );
  console.log("Transaction result:", txResult);
}
```

## 7. 错误处理

```typescript
import { StargateClient, isDeliverTxFailure } from "@cosmjs/stargate";

try {
  const result = await client.sendTokens(
    senderAddress,
    recipientAddress,
    amount,
    "auto"
  );

  // 检查交易是否失败
  if (isDeliverTxFailure(result)) {
    console.error("Transaction failed:", result.rawLog);
    // 处理失败逻辑
  } else {
    console.log("Transaction succeeded:", result.transactionHash);
  }
} catch (error) {
  if (error instanceof Error) {
    console.error("Error:", error.message);

    // 常见错误类型
    if (error.message.includes("insufficient funds")) {
      console.error("Insufficient balance");
    } else if (error.message.includes("sequence mismatch")) {
      console.error("Account sequence mismatch - retry transaction");
    } else if (error.message.includes("out of gas")) {
      console.error("Transaction ran out of gas");
    }
  }
}
```

## 8. 最佳实践

### 安全建议

1. **助记词安全**
   - 永远不要在代码中硬编码助记词
   - 不要将助记词写入日志
   - 使用环境变量或安全存储

2. **地址验证**
   ```typescript
   import { StargateClient } from "@cosmjs/stargate";

   // 验证地址有效性
   const client = await StargateClient.connect(rpcEndpoint);
   const account = await client.getAccount(address);
   if (!account) {
     throw new Error("Invalid address");
   }
   ```

3. **Gas 估算**
   ```typescript
   // 使用 "auto" 并添加 1.4 倍安全余量（CosmJS 默认）
   const result = await client.sendTokens(
     senderAddress,
     recipientAddress,
     amount,
     "auto"
   );
   ```

### 性能优化

```typescript
// 复用客户端连接
let clientInstance: StargateClient | null = null;

async function getClient(): Promise<StargateClient> {
  if (!clientInstance) {
    clientInstance = await StargateClient.connect(rpcEndpoint);
  }
  return clientInstance;
}
```

## 9. 常用消息类型

```typescript
// 银行转账
const msgSend = {
  typeUrl: "/cosmos.bank.v1beta1.MsgSend",
  value: {
    fromAddress: senderAddress,
    toAddress: recipientAddress,
    amount: coins(1000000, "uatom"),
  },
};

// 质押
const msgDelegate = {
  typeUrl: "/cosmos.staking.v1beta1.MsgDelegate",
  value: {
    delegatorAddress: delegatorAddress,
    validatorAddress: validatorAddress,
    amount: { denom: "uatom", amount: "1000000" },
  },
};

// 取消质押
const msgUndelegate = {
  typeUrl: "/cosmos.staking.v1beta1.MsgUndelegate",
  value: {
    delegatorAddress: delegatorAddress,
    validatorAddress: validatorAddress,
    amount: { denom: "uatom", amount: "1000000" },
  },
};

// 提取奖励
const msgWithdrawRewards = {
  typeUrl: "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
  value: {
    delegatorAddress: delegatorAddress,
    validatorAddress: validatorAddress,
  },
};

// IBC 转账
const msgTransfer = {
  typeUrl: "/ibc.applications.transfer.v1.MsgTransfer",
  value: {
    sourcePort: "transfer",
    sourceChannel: "channel-0",
    token: { denom: "uatom", amount: "1000000" },
    sender: senderAddress,
    receiver: recipientAddress,
    timeoutHeight: {
      revisionNumber: BigInt(1),
      revisionHeight: BigInt(1000000),
    },
    timeoutTimestamp: BigInt(0),
  },
};
```

## 10. 参考资源

- [CosmJS 官方文档](https://cosmos.github.io/cosmjs/)
- [CosmJS GitHub](https://github.com/cosmos/cosmjs)
- [Cosmos SDK 文档](https://docs.cosmos.network/)
- [Cosmos Tutorials](https://tutorials.cosmos.network/tutorials/7-cosmjs/)

---

*此文档为 ShareTokens 项目 GenieBot UI 和客户端 SDK 开发提供参考。*
