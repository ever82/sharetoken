# Keplr Wallet 开发指南

## 概述

Keplr 是 Cosmos 生态系统最流行的钱包，支持所有 Cosmos SDK 链。ShareTokens 项目中，用户通过 Keplr 钱包连接 GenieBot 界面进行操作。

## 1. Keplr 简介

- **定位**：Cosmos 生态官方推荐钱包
- **支持链**：所有基于 Cosmos SDK 的链（Cosmos Hub, Osmosis, Juno, Secret 等）
- **核心功能**：
  - 多链资产管理
  - 跨链转账（IBC）
  - Stake/Delegate
  - dApp 连接与交易签名

## 2. 平台支持

### 浏览器扩展
- Chrome, Firefox, Brave
- 安装：[Chrome Web Store](https://chrome.google.com/webstore/detail/keplr/dmkamcknogkgcdfhhbddcghachkejeap)

### 移动端 App
- iOS App Store / Google Play
- 支持 WalletConnect 连接

## 3. 链集成与配置

### 注册链到 Keplr

```typescript
// 在 dApp 中注册自定义链
await window.keplr.experimentalSuggestChain({
  chainId: 'sharetokens-1',
  chainName: 'ShareTokens',
  rpc: 'https://rpc.sharetokens.com',
  rest: 'https://api.sharetokens.com',
  bip44: {
    coinType: 118, // Cosmos coin type
  },
  bech32Config: {
    bech32PrefixAccAddr: 'share',
    bech32PrefixAccPub: 'sharepub',
    bech32PrefixValAddr: 'sharevaloper',
    bech32PrefixValPub: 'sharevaloperpub',
    bech32PrefixConsAddr: 'sharevalcons',
    bech32PrefixConsPub: 'sharevalconspub',
  },
  currencies: [{
    coinDenom: 'SHARE',
    coinMinimalDenom: 'ushare',
    coinDecimals: 6,
  }],
  feeCurrencies: [{
    coinDenom: 'SHARE',
    coinMinimalDenom: 'ushare',
    coinDecimals: 6,
  }],
  stakeCurrency: {
    coinDenom: 'SHARE',
    coinMinimalDenom: 'ushare',
    coinDecimals: 6,
  },
});
```

### 检测链是否已启用

```typescript
// 检查链是否已在 Keplr 中启用
const chainId = 'sharetokens-1';

try {
  await window.keplr.enable(chainId);
  console.log('Chain enabled successfully');
} catch (error) {
  console.error('Chain not enabled or not installed');
  // 可尝试使用 experimentalSuggestChain 注册链
}
```

## 4. dApp 钱包连接

### 检测 Keplr 是否安装

```typescript
// 等待 Keplr 注入
const checkKeplr = async (): Promise<boolean> => {
  if (window.keplr) {
    return true;
  }

  // 等待注入完成（最多 1 秒）
  if (document.readyState === 'complete') {
    return window.keplr !== undefined;
  }

  return new Promise((resolve) => {
    const timer = setTimeout(() => resolve(false), 1000);
    window.addEventListener('keplr_keystorechange', () => {
      clearTimeout(timer);
      resolve(true);
    });
  });
};
```

### 连接并获取账户

```typescript
import { OfflineAminoSigner, OfflineDirectSigner } from '@cosmjs/proto-signing';

interface KeplrAccount {
  address: string;
  algo: string;
  pubkey: Uint8Array;
}

const connectKeplr = async (chainId: string): Promise<KeplrAccount[]> => {
  // 启用链
  await window.keplr.enable(chainId);

  // 获取离线签名者
  const offlineSigner = window.getOfflineSigner(chainId);

  // 获取账户信息
  const accounts = await offlineSigner.getAccounts();

  return accounts.map(acc => ({
    address: acc.address,
    algo: acc.algo,
    pubkey: acc.pubkey,
  }));
};
```

### 监听账户变化

```typescript
// 监听钱包账户切换
window.addEventListener('keplr_keystorechange', () => {
  console.log('Keplr account changed');
  // 重新获取账户信息
  reconnectWallet();
});
```

## 5. 交易签名流程

### 使用 Amino 签名（JSON 格式）

```typescript
import { SigningStargateClient } from '@cosmjs/stargate';

const signTransactionAmino = async (
  chainId: string,
  fromAddress: string,
  toAddress: string,
  amount: { denom: string; amount: string }[]
) => {
  // 获取离线签名者
  const offlineSigner = window.getOfflineSigner(chainId);

  // 创建 SigningStargateClient
  const client = await SigningStargateClient.connectWithSigner(
    'https://rpc.sharetokens.com',
    offlineSigner
  );

  // 发送交易
  const result = await client.sendTokens(
    fromAddress,
    toAddress,
    amount,
    'auto', // 自动计算 gas
    'Transfer via GenieBot'
  );

  console.log('Transaction hash:', result.transactionHash);
  return result;
};
```

### 使用 Direct 签名（Protobuf 格式）

```typescript
import { SigningStargateClient } from '@cosmjs/stargate';

const signTransactionDirect = async (chainId: string) => {
  // 获取 Direct 签名者
  const offlineSigner = window.getOfflineSigner(chainId);

  // 强制使用 Direct 模式
  const client = await SigningStargateClient.connectWithSigner(
    'https://rpc.sharetokens.com',
    offlineSigner,
    { broadcastPollIntervalMs: 300, broadcastTimeoutMs: 8000 }
  );

  return client;
};
```

### 自定义消息签名

```typescript
// 签名任意消息（用于身份验证）
const signArbitrary = async (
  chainId: string,
  signer: string,
  data: object
): Promise<{ signature: string; pub_key: { type: string; value: string } }> => {
  const result = await window.keplr.signArbitrary(chainId, signer, data);
  return result;
};

// 验证签名
import { verifyADR36Amino } from '@keplr-wallet/cosmos';

const verifySignature = (
  signer: string,
  data: object,
  signature: { pub_key: { value: string }; signature: string }
): boolean => {
  return verifyADR36Amino(
    'share', // bech32 前缀
    signer,
    Buffer.from(JSON.stringify(data)),
    signature.pub_key,
    signature.signature
  );
};
```

## 6. Keplr + CosmJS 集成模式

### 完整连接示例

```typescript
import { SigningStargateClient, GasPrice } from '@cosmjs/stargate';
import { Registry } from '@cosmjs/proto-signing';

// GenieBot 钱包服务
class KeplrWalletService {
  private client: SigningStargateClient | null = null;
  private accounts: KeplrAccount[] = [];
  private chainId: string;

  constructor(chainId: string = 'sharetokens-1') {
    this.chainId = chainId;
  }

  async connect(): Promise<string> {
    // 1. 检查 Keplr 安装
    if (!window.keplr) {
      throw new Error('Please install Keplr extension');
    }

    // 2. 启用链
    await window.keplr.enable(this.chainId);

    // 3. 获取签名者
    const offlineSigner = window.getOfflineSigner(this.chainId);

    // 4. 获取账户
    this.accounts = await offlineSigner.getAccounts();

    // 5. 创建客户端
    this.client = await SigningStargateClient.connectWithSigner(
      'https://rpc.sharetokens.com',
      offlineSigner,
      {
        gasPrice: GasPrice.fromString('0.025ushare'),
      }
    );

    return this.accounts[0].address;
  }

  async getBalance(address: string): Promise<string> {
    if (!this.client) throw new Error('Not connected');

    const balance = await this.client.getBalance(address, 'ushare');
    return balance.amount;
  }

  async sendTokens(
    toAddress: string,
    amount: string,
    memo: string = ''
  ): Promise<string> {
    if (!this.client || !this.accounts.length) {
      throw new Error('Not connected');
    }

    const result = await this.client.sendTokens(
      this.accounts[0].address,
      toAddress,
      [{ denom: 'ushare', amount }],
      'auto',
      memo
    );

    return result.transactionHash;
  }

  getAddress(): string | null {
    return this.accounts[0]?.address || null;
  }

  disconnect(): void {
    this.client = null;
    this.accounts = [];
  }
}

// 使用示例
const walletService = new KeplrWalletService();

// 在 React/Vue 组件中
async function handleConnect() {
  try {
    const address = await walletService.connect();
    console.log('Connected:', address);
  } catch (error) {
    console.error('Connection failed:', error);
  }
}
```

### 处理自定义消息类型

```typescript
import { MsgExecuteContract } from 'cosmjs-types/cosmwasm/wasm/v1/tx';

// 对于 CosmWasm 合约交互
const executeContract = async (
  chainId: string,
  contractAddress: string,
  msg: object,
  funds: { denom: string; amount: string }[] = []
) => {
  const offlineSigner = window.getOfflineSigner(chainId);
  const accounts = await offlineSigner.getAccounts();

  const client = await SigningStargateClient.connectWithSigner(
    'https://rpc.sharetokens.com',
    offlineSigner
  );

  const msgExecute = {
    typeUrl: '/cosmwasm.wasm.v1.MsgExecuteContract',
    value: MsgExecuteContract.fromPartial({
      sender: accounts[0].address,
      contract: contractAddress,
      msg: Buffer.from(JSON.stringify(msg)),
      funds: funds,
    }),
  };

  const result = await client.signAndBroadcast(
    accounts[0].address,
    [msgExecute],
    'auto'
  );

  return result.transactionHash;
};
```

## 常见问题

### Q: 如何处理用户拒绝授权？

```typescript
try {
  await window.keplr.enable(chainId);
} catch (error) {
  if (error.message?.includes('rejected')) {
    // 用户拒绝了连接请求
    console.log('User rejected the connection request');
  }
}
```

### Q: 如何处理链未注册？

```typescript
const ensureChainRegistered = async (chainConfig: ChainConfig) => {
  try {
    await window.keplr.enable(chainConfig.chainId);
  } catch (error) {
    if (error.message?.includes('no chain info')) {
      // 尝试注册链
      await window.keplr.experimentalSuggestChain(chainConfig);
    }
  }
};
```

### Q: 如何检测网络切换？

```typescript
// Keplr 不提供网络切换事件，需要定期检查
let currentChainId: string | null = null;

setInterval(async () => {
  const key = await window.keplr.getKey(chainId);
  if (key.address !== currentAddress) {
    currentAddress = key.address;
    // 触发账户变化处理
  }
}, 1000);
```

## 参考资源

- [Keplr 官方文档](https://docs.keplr.app/)
- [CosmJS 文档](https://cosmos.github.io/cosmjs/)
- [Keplr API 参考](https://docs.keplr.app/api/)
- [Cosmos SDK 文档](https://docs.cosmos.network/)

## ShareTokens 项目集成要点

1. **连接流程**：用户点击"连接钱包" -> 检测 Keplr -> 启用 sharetokens-1 链 -> 获取地址
2. **交易签名**：通过 CosmJS SigningStargateClient 签名广播交易
3. **状态管理**：监听 `keplr_keystorechange` 事件处理账户切换
4. **错误处理**：优雅处理未安装、拒绝授权、链未注册等情况
