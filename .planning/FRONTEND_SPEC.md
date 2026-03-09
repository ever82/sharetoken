# Frontend UI Specification: GenieBot Client Plugin

**Version:** 3.0.0
**Updated:** 2026-03-02
**Module Type:** User Plugin
**Tech Stack:** React 18 + TypeScript + Zustand + Tailwind CSS + @cosmjs/stargate

---

## Overview

GenieBot是 ShareTokens 服务市场的用户端插件，提供 AI 对话界面，帮助用户便捷地发现和调用三层服务（LLM/Agent/Workflow）。

**重要定位:**
- GenieBot是**用户插件**，不是核心模块
- 核心模块中不包含GenieBot，节点可以不安装GenieBot而正常运行
- GenieBot调用服务市场的 API，不直接执行 AI 任务

```
+-----------------------------------------------------------------------------+
|                              User Layer                                      |
|  +-----------------------------------------------------------------------+  |
|  |  GenieBot Interface (Client Plugin) - This Document                   |  |
|  |  - Natural language chat interface                                    |  |
|  |  - Service browsing and selection                                     |  |
|  |  - Task management and tracking                                       |  |
|  +-----------------------------------------------------------------------+  |
+-----------------------------------------------------------------------------+
                              |
                              v
+-----------------------------------------------------------------------------+
|                          Service Marketplace (Core Module)                   |
|  Three-tier Service: LLM (Level 1) | Agent (Level 2) | Workflow (Level 3)   |
+-----------------------------------------------------------------------------+
```

---

## 1. Tech Stack Confirmation

### 1.1 Core Technology Stack

| Category | Technology | Version | Purpose |
|----------|------------|---------|---------|
| Framework | React | 18.x | UI framework |
| Language | TypeScript | 5.x | Type safety |
| State Management | Zustand | 4.x | Global state |
| Styling | Tailwind CSS | 3.x | Utility-first CSS |
| Build Tool | Vite | 5.x | Fast development |
| Chain Integration | @cosmjs/stargate | 0.32.x | Cosmos SDK client |
| Wallet | Keplr | - | Browser extension |
| HTTP Client | Axios | 1.x | API requests |
| Real-time | native EventSource | - | SSE streaming |
| WebSocket | native WebSocket | - | Real-time updates |

### 1.2 Development Dependencies

| Package | Purpose |
|---------|---------|
| @types/react | TypeScript types |
| @types/react-dom | TypeScript types |
| eslint + prettier | Code quality |
| vitest | Unit testing |
| @testing-library/react | Component testing |
| playwright | E2E testing |
| react-i18next | Internationalization |

### 1.3 Package Installation

```bash
# Core dependencies
npm install react react-dom zustand @cosmjs/stargate @cosmjs/tendermint-rpc axios react-router-dom react-i18next

# Dev dependencies
npm install -D typescript @types/react @types/react-dom vite @vitejs/plugin-react tailwindcss postcss autoprefixer eslint prettier vitest @testing-library/react playwright
```

---

## 2. Component Architecture

### 2.1 Directory Structure

```
frontend/
+-- src/
|   +-- app/                          # App configuration
|   |   +-- App.tsx                   # Root component
|   |   +-- Router.tsx                # React Router setup
|   |   +-- providers/                # Context providers
|   |   |   +-- QueryProvider.tsx
|   |   |   +-- ThemeProvider.tsx
|   |   |   +-- WalletProvider.tsx
|   |   +-- config/
|   |       +-- chains.ts             # Chain configuration
|   |       +-- constants.ts          # App constants
|   |
|   +-- core/                         # Core infrastructure
|   |   +-- wallet/                   # Keplr wallet integration
|   |   |   +-- KeplrProvider.tsx     # Wallet context provider
|   |   |   +-- useWallet.ts          # Wallet hook
|   |   |   +-- useSigningClient.ts   # Signing client hook
|   |   |   +-- keplr.ts              # Keplr utilities
|   |   |   +-- types.ts              # Wallet types
|   |   |
|   |   +-- chain/                    # Chain interaction
|   |   |   +-- ChainClient.ts        # Cosmos SDK client
|   |   |   +-- QueryClient.ts        # Query client
|   |   |   +-- TxClient.ts           # Transaction client
|   |   |   +-- types.ts              # Chain types
|   |   |
|   |   +-- api/                      # API layer
|   |   |   +-- base.ts               # Base API client
|   |   |   +-- marketplace.ts        # Marketplace API
|   |   |   +-- service.ts            # Service API
|   |   |   +-- escrow.ts             # Escrow API
|   |   |   +-- dispute.ts            # Dispute API
|   |   |   +-- streaming.ts          # SSE handling
|   |   |
|   |   +-- layout/                   # Layout components
|   |   |   +-- AppLayout.tsx
|   |   |   +-- Sidebar.tsx
|   |   |   +-- Header.tsx
|   |   |   +-- Footer.tsx
|   |   |   +-- MobileNav.tsx
|   |   |
|   |   +-- theme/                    # Theming
|   |       +-- ThemeProvider.tsx
|   |       +-- tokens.css            # Design tokens
|   |
|   +-- stores/                       # Zustand stores
|   |   +-- index.ts                  # Store exports
|   |   +-- walletStore.ts            # Wallet state
|   |   +-- userStore.ts              # User/MQ state
|   |   +-- chatStore.ts              # Chat state
|   |   +-- serviceStore.ts           # Service browsing state
|   |   +-- uiStore.ts                # UI state (modals, etc.)
|   |
|   +-- components/                   # Reusable components
|   |   +-- common/                   # Common UI components
|   |   |   +-- Button/
|   |   |   +-- Input/
|   |   |   +-- Card/
|   |   |   +-- Modal/
|   |   |   +-- Toast/
|   |   |   +-- Spinner/
|   |   |   +-- Dropdown/
|   |   |   +-- Badge/
|   |   |   +-- Avatar/
|   |   |   +-- Tooltip/
|   |   |
|   |   +-- wallet/                   # Wallet components
|   |   |   +-- ConnectButton.tsx
|   |   |   +-- WalletModal.tsx
|   |   |   +-- BalanceDisplay.tsx
|   |   |   +-- AddressDisplay.tsx
|   |   |
|   |   +-- chat/                     # Chat components
|   |   |   +-- ChatContainer.tsx
|   |   |   +-- MessageList.tsx
|   |   |   +-- MessageBubble.tsx
|   |   |   +-- ChatInput.tsx
|   |   |   +-- TypingIndicator.tsx
|   |   |   +-- StreamingMessage.tsx
|   |   |   +-- WelcomeCard.tsx
|   |   |
|   |   +-- service/                  # Service components
|   |   |   +-- ServiceCard.tsx
|   |   |   +-- ServiceLevelBadge.tsx
|   |   |   +-- ServiceFilter.tsx
|   |   |   +-- ServiceDetail.tsx
|   |   |   +-- ServiceRequestForm.tsx
|   |   |
|   |   +-- mq/                       # MQ components
|   |   |   +-- MQBadge.tsx
|   |   |   +-- MQDisplay.tsx
|   |   |   +-- MQProgressBar.tsx
|   |   |
|   |   +-- dispute/                  # Dispute components
|   |   |   +-- DisputeCard.tsx
|   |   |   +-- DisputeForm.tsx
|   |   |   +-- EvidenceUpload.tsx
|   |   |   +-- VotingPanel.tsx
|   |   |
|   |   +-- idea/                     # Idea components
|   |   |   +-- IdeaCard.tsx
|   |   |   +-- EvalReport.tsx
|   |   |   +-- TokenEstimate.tsx
|   |   |
|   |   +-- token/                    # Token components
|   |       +-- TokenAmount.tsx
|   |       +-- TokenInput.tsx
|   |
|   +-- pages/                        # Page components
|   |   +-- ChatPage.tsx              # Main chat interface
|   |   +-- MarketplacePage.tsx       # Service marketplace
|   |   +-- ServiceDetailPage.tsx     # Service details
|   |   +-- MyServicesPage.tsx        # User's services
|   |   +-- DisputesPage.tsx          # Dispute list
|   |   +-- DisputeDetailPage.tsx     # Dispute details
|   |   +-- ProfilePage.tsx           # User profile
|   |   +-- IdeasPage.tsx             # Idea incubation
|   |   +-- SettingsPage.tsx          # Settings
|   |
|   +-- hooks/                        # Custom hooks
|   |   +-- useChat.ts
|   |   +-- useService.ts
|   |   +-- useDispute.ts
|   |   +-- useEscrow.ts
|   |   +-- useToast.ts
|   |   +-- useMediaQuery.ts
|   |   +-- useStreaming.ts
|   |
|   +-- utils/                        # Utility functions
|   |   +-- format.ts                 # Formatting utilities
|   |   +-- validation.ts             # Validation helpers
|   |   +-- constants.ts              # App constants
|   |   +-- helpers.ts                # General helpers
|   |
|   +-- types/                        # TypeScript types
|   |   +-- wallet.ts
|   |   +-- service.ts
|   |   +-- chat.ts
|   |   |   +-- dispute.ts
|   |   +-- mq.ts
|   |   +-- api.ts
|   |
|   +-- i18n/                         # Internationalization
|   |   +-- index.ts
|   |   +-- locales/
|   |       +-- en.json
|   |       +-- zh.json
|   |
|   +-- styles/
|       +-- globals.css
|       +-- variables.css
|
+-- public/
|   +-- favicon.ico
|   +-- logo.svg
|
+-- index.html
+-- package.json
+-- tsconfig.json
+-- vite.config.ts
+-- tailwind.config.js
+-- .env.example
```

### 2.2 Component Architecture Diagram

```
+-----------------------------------------------------------------------------+
|                           GenieBot Frontend Architecture                    |
+-----------------------------------------------------------------------------+

+-----------------------------------------------------------------------------+
|  Pages Layer (Route Components)                                             |
|  +-----------+  +-----------+  +-----------+  +-----------+                 |
|  | ChatPage  |  | MarketPage|  | Disputes  |  | Profile   |                 |
|  +-----------+  +-----------+  +-----------+  +-----------+                 |
+-----------------------------------------------------------------------------+
                              |
                              v
+-----------------------------------------------------------------------------+
|  Feature Components                                                         |
|  +----------------+  +----------------+  +----------------+                 |
|  | ChatContainer  |  | ServiceBrowser |  | DisputeList    |                 |
|  | MessageList    |  | ServiceCard    |  | VotingPanel    |                 |
|  | ChatInput      |  | ServiceFilter  |  | EvidenceUpload |                 |
|  +----------------+  +----------------+  +----------------+                 |
+-----------------------------------------------------------------------------+
                              |
                              v
+-----------------------------------------------------------------------------+
|  Common Components                                                          |
|  +-------+  +-------+  +-------+  +-------+  +-------+  +-------+          |
|  | Button|  | Input |  | Card  |  | Modal |  | Toast |  | Badge |          |
|  +-------+  +-------+  +-------+  +-------+  +-------+  +-------+          |
+-----------------------------------------------------------------------------+
                              |
                              v
+-----------------------------------------------------------------------------+
|  Core Infrastructure                                                        |
|  +------------------+  +------------------+  +------------------+           |
|  | WalletProvider   |  | ChainClient      |  | API Client       |           |
|  | (Keplr)          |  | (@cosmjs)        |  | (Axios + SSE)    |           |
|  +------------------+  +------------------+  +------------------+           |
+-----------------------------------------------------------------------------+
                              |
                              v
+-----------------------------------------------------------------------------+
|  State Management (Zustand)                                                 |
|  +------------+  +------------+  +------------+  +------------+             |
|  | walletStore|  | chatStore  |  | userStore  |  | serviceStore|            |
|  +------------+  +------------+  +------------+  +------------+             |
+-----------------------------------------------------------------------------+
```

---

## 3. Keplr Wallet Integration

### 3.1 Wallet Provider Architecture

```typescript
// src/core/wallet/KeplrProvider.tsx
import React, { createContext, useCallback, useEffect, useState } from 'react';
import { useWalletStore } from '@stores/walletStore';

// Chain configuration for ShareTokens
const SHARETOKENS_CHAIN_CONFIG = {
  chainId: 'sharetokens-1',
  chainName: 'ShareTokens',
  rpc: 'https://rpc.sharetokens.io',
  rest: 'https://api.sharetokens.io',
  stakeCurrency: {
    coinDenom: 'STT',
    coinMinimalDenom: 'ustt',
    coinDecimals: 6,
  },
  bip44: { coinType: 118 },
  bech32Config: {
    bech32PrefixAccAddr: 'cosmos',
    bech32PrefixAccPub: 'cosmospub',
    bech32PrefixValAddr: 'cosmosvaloper',
    bech32PrefixValPub: 'cosmosvaloperpub',
    bech32PrefixConsAddr: 'cosmosvalcons',
    bech32PrefixConsPub: 'cosmosvalconspub',
  },
  currencies: [{
    coinDenom: 'STT',
    coinMinimalDenom: 'ustt',
    coinDecimals: 6,
  }],
  feeCurrencies: [{
    coinDenom: 'STT',
    coinMinimalDenom: 'ustt',
    coinDecimals: 6,
    gasPriceStep: { low: 0.01, average: 0.025, high: 0.04 },
  }],
  features: ['stargate', 'ibc-transfer'],
};

interface WalletContextValue {
  isKeplrInstalled: boolean;
  connect: () => Promise<void>;
  disconnect: () => void;
  getSigningClient: () => Promise<SigningStargateClient | null>;
  signAndBroadcast: (msgs: EncodeObject[]) => Promise<DeliverTxResponse>;
}

export const WalletContext = createContext<WalletContextValue | null>(null);

export function KeplrProvider({ children }: { children: React.ReactNode }) {
  const {
    address, balance, isConnected, isConnecting,
    setAddress, setBalance, setConnected, setConnecting
  } = useWalletStore();

  const [isKeplrInstalled, setIsKeplrInstalled] = useState(false);
  const [signingClient, setSigningClient] = useState<SigningStargateClient | null>(null);

  // Check Keplr installation on mount
  useEffect(() => {
    const checkKeplr = async () => {
      if (typeof window !== 'undefined' && (window as any).keplr) {
        setIsKeplrInstalled(true);
      }
    };
    checkKeplr();

    // Listen for Keplr installation
    const handleKeplrReady = () => setIsKeplrInstalled(true);
    window.addEventListener('keplr_keystorechange', handleKeplrReady);
    return () => window.removeEventListener('keplr_keystorechange', handleKeplrReady);
  }, []);

  // Connect wallet
  const connect = useCallback(async () => {
    if (!isKeplrInstalled) {
      // Open Keplr install page
      window.open('https://www.keplr.app/get', '_blank');
      return;
    }

    try {
      setConnecting(true);

      // Enable the chain
      await (window as any).keplr.enable(SHARETOKENS_CHAIN_CONFIG.chainId);

      // Get offline signer
      const offlineSigner = (window as any).getOfflineSigner?.(SHARETOKENS_CHAIN_CONFIG.chainId);

      // Get accounts
      const accounts = await offlineSigner.getAccounts();
      const address = accounts[0].address;

      setAddress(address);
      setConnected(true);

      // Create signing client
      const client = await SigningStargateClient.connectWithSigner(
        SHARETOKENS_CHAIN_CONFIG.rpc,
        offlineSigner
      );
      setSigningClient(client);

      // Fetch balance
      const balance = await client.getAllBalances(address);
      setBalance(balance);

    } catch (error) {
      console.error('Failed to connect wallet:', error);
      setConnected(false);
    } finally {
      setConnecting(false);
    }
  }, [isKeplrInstalled, setAddress, setBalance, setConnected, setConnecting]);

  // Disconnect wallet
  const disconnect = useCallback(() => {
    setAddress(null);
    setBalance([]);
    setConnected(false);
    setSigningClient(null);
  }, [setAddress, setBalance, setConnected]);

  // Sign and broadcast transaction
  const signAndBroadcast = useCallback(async (msgs: EncodeObject[]) => {
    if (!signingClient || !address) {
      throw new Error('Wallet not connected');
    }

    const fee = {
      amount: [{ denom: 'ustt', amount: '5000' }],
      gas: '200000',
    };

    return await signingClient.signAndBroadcast(address, msgs, fee);
  }, [signingClient, address]);

  // Listen for account changes
  useEffect(() => {
    const handleKeyStoreChange = () => {
      // Reconnect to get new account
      if (isConnected) {
        connect();
      }
    };

    window.addEventListener('keplr_keystorechange', handleKeyStoreChange);
    return () => window.removeEventListener('keplr_keystorechange', handleKeyStoreChange);
  }, [isConnected, connect]);

  return (
    <WalletContext.Provider value={{
      isKeplrInstalled,
      connect,
      disconnect,
      getSigningClient: () => Promise.resolve(signingClient),
      signAndBroadcast,
    }}>
      {children}
    </WalletContext.Provider>
  );
}
```

### 3.2 Wallet Hook

```typescript
// src/core/wallet/useWallet.ts
import { useContext } from 'react';
import { WalletContext } from './KeplrProvider';

export function useWallet() {
  const context = useContext(WalletContext);
  const store = useWalletStore();

  if (!context) {
    throw new Error('useWallet must be used within KeplrProvider');
  }

  return {
    // State from store
    address: store.address,
    balance: store.balance,
    isConnected: store.isConnected,
    isConnecting: store.isConnecting,

    // Methods from context
    isKeplrInstalled: context.isKeplrInstalled,
    connect: context.connect,
    disconnect: context.disconnect,
    signAndBroadcast: context.signAndBroadcast,

    // Helpers
    truncateAddress: (start = 6, end = 4) => {
      if (!store.address) return '';
      return `${store.address.slice(0, start)}...${store.address.slice(-end)}`;
    },

    getSttBalance: () => {
      const stt = store.balance.find(b => b.denom === 'ustt');
      return stt ? Number(stt.amount) / 1_000_000 : 0;
    },
  };
}
```

### 3.3 Connect Button Component

```typescript
// src/components/wallet/ConnectButton.tsx
import { useWallet } from '@core/wallet/useWallet';
import { Button } from '@components/common/Button';
import { WalletModal } from './WalletModal';
import { useState } from 'react';

export function ConnectButton() {
  const { isConnected, isConnecting, address, connect, disconnect, truncateAddress, getSttBalance } = useWallet();
  const [showModal, setShowModal] = useState(false);

  if (isConnecting) {
    return (
      <Button variant="secondary" disabled>
        <Spinner size="sm" />
        Connecting...
      </Button>
    );
  }

  if (!isConnected) {
    return (
      <>
        <Button variant="primary" onClick={connect}>
          <WalletIcon />
          Connect Wallet
        </Button>
      </>
    );
  }

  return (
    <>
      <Button variant="ghost" onClick={() => setShowModal(true)}>
        <div className="flex items-center gap-2">
          <div className="w-2 h-2 rounded-full bg-green-500" />
          <span>{truncateAddress()}</span>
          <Badge variant="success">{getSttBalance().toFixed(2)} STT</Badge>
        </div>
      </Button>

      <WalletModal isOpen={showModal} onClose={() => setShowModal(false)} />
    </>
  );
}
```

### 3.4 Transaction Flow

```
+-----------------------------------------------------------------------------+
|                        Transaction Flow (Keplr)                              |
+-----------------------------------------------------------------------------+

User Action (e.g., "Start Service")
              |
              v
+---------------------------+
| 1. Build Message          |
|   - Create MsgExecuteContract|
|   - Or custom module msg  |
+---------------------------+
              |
              v
+---------------------------+
| 2. Estimate Gas           |
|   - Simulate tx           |
|   - Calculate fees        |
+---------------------------+
              |
              v
+---------------------------+
| 3. Request Signature      |
|   - Keplr popup opens     |
|   - User reviews tx       |
|   - User approves/rejects |
+---------------------------+
              |
        +-----+-----+
        |           |
    Approved    Rejected
        |           |
        v           v
+---------------+  +----------------+
| 4. Broadcast  |  | Handle Error   |
|    to chain   |  | Show toast     |
+---------------+  +----------------+
        |
        v
+---------------------------+
| 5. Wait for Confirmation  |
|   - Poll tx status        |
|   - Show progress         |
+---------------------------+
        |
        v
+---------------------------+
| 6. Update UI State        |
|   - Refresh balances      |
|   - Show success toast    |
|   - Navigate if needed    |
+---------------------------+
```

---

## 4. State Management (Zustand)

### 4.1 Store Architecture

```typescript
// src/stores/walletStore.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface Coin {
  denom: string;
  amount: string;
}

interface WalletState {
  address: string | null;
  balance: Coin[];
  isConnected: boolean;
  isConnecting: boolean;

  // Actions
  setAddress: (address: string | null) => void;
  setBalance: (balance: Coin[]) => void;
  setConnected: (connected: boolean) => void;
  setConnecting: (connecting: boolean) => void;
  reset: () => void;
}

const initialState = {
  address: null,
  balance: [],
  isConnected: false,
  isConnecting: false,
};

export const useWalletStore = create<WalletState>()(
  persist(
    (set) => ({
      ...initialState,

      setAddress: (address) => set({ address }),
      setBalance: (balance) => set({ balance }),
      setConnected: (isConnected) => set({ isConnected }),
      setConnecting: (isConnecting) => set({ isConnecting }),
      reset: () => set(initialState),
    }),
    {
      name: 'sharetokens-wallet',
      partialize: (state) => ({
        address: state.address,
        isConnected: state.isConnected,
      }),
    }
  )
);
```

```typescript
// src/stores/chatStore.ts
import { create } from 'zustand';

interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: number;
  isStreaming?: boolean;
  serviceLevel?: 1 | 2 | 3;
  serviceId?: string;
  serviceCard?: ServiceRecommendation;
}

interface ChatState {
  messages: Message[];
  currentConversationId: string | null;
  isTyping: boolean;
  streamingMessageId: string | null;

  // Actions
  addMessage: (message: Omit<Message, 'id' | 'timestamp'>) => string;
  updateMessage: (id: string, content: string) => void;
  setStreaming: (id: string | null) => void;
  setTyping: (typing: boolean) => void;
  clearChat: () => void;
  loadConversation: (conversationId: string) => Promise<void>;
}

export const useChatStore = create<ChatState>((set, get) => ({
  messages: [],
  currentConversationId: null,
  isTyping: false,
  streamingMessageId: null,

  addMessage: (message) => {
    const id = crypto.randomUUID();
    set((state) => ({
      messages: [
        ...state.messages,
        { ...message, id, timestamp: Date.now() },
      ],
    }));
    return id;
  },

  updateMessage: (id, content) => {
    set((state) => ({
      messages: state.messages.map((m) =>
        m.id === id ? { ...m, content } : m
      ),
    }));
  },

  setStreaming: (id) => set({ streamingMessageId: id }),
  setTyping: (isTyping) => set({ isTyping }),
  clearChat: () => set({ messages: [], currentConversationId: null }),

  loadConversation: async (conversationId) => {
    // Load from API
    set({ currentConversationId: conversationId });
  },
}));
```

```typescript
// src/stores/userStore.ts
import { create } from 'zustand';

interface Identity {
  verified: boolean;
  level: 'none' | 'basic' | 'verified' | 'premium';
  registeredAt: number;
}

interface MQInfo {
  score: number;
  level: 1 | 2 | 3 | 4 | 5;
  tierName: string;
  lastUpdated: number;
}

interface UserState {
  identity: Identity | null;
  mqInfo: MQInfo | null;
  preferences: {
    language: string;
    theme: 'light' | 'dark' | 'system';
    currency: string;
  };

  // Actions
  setIdentity: (identity: Identity | null) => void;
  setMQInfo: (mq: MQInfo | null) => void;
  updatePreferences: (prefs: Partial<UserState['preferences']>) => void;
  fetchUserData: (address: string) => Promise<void>;
}

export const useUserStore = create<UserState>((set) => ({
  identity: null,
  mqInfo: null,
  preferences: {
    language: 'en',
    theme: 'system',
    currency: 'USD',
  },

  setIdentity: (identity) => set({ identity }),
  setMQInfo: (mqInfo) => set({ mqInfo }),
  updatePreferences: (prefs) => set((state) => ({
    preferences: { ...state.preferences, ...prefs },
  })),

  fetchUserData: async (address) => {
    // Fetch identity and MQ info from chain/API
    // This would call the API layer
  },
}));
```

### 4.2 State Flow Diagram

```
+-----------------------------------------------------------------------------+
|                        State Management Flow                                 |
+-----------------------------------------------------------------------------+

                    +-------------------+
                    |  User Action      |
                    +-------------------+
                              |
                              v
+-----------------------------------------------------------------------------+
|  Zustand Store Actions                                                      |
|  +----------------+  +----------------+  +----------------+                 |
|  | useWalletStore |  | useChatStore   |  | useUserStore   |                 |
|  | - connect()    |  | - addMessage() |  | - fetchUserData|                 |
|  | - disconnect() |  | - setTyping()  |  | - setMQInfo()  |                 |
|  +----------------+  +----------------+  +----------------+                 |
+-----------------------------------------------------------------------------+
                              |
                              v
+-----------------------------------------------------------------------------+
|  API / Chain Layer                                                          |
|  +----------------+  +----------------+  +----------------+                 |
|  | ChainClient    |  | MarketplaceAPI |  | WebSocket       |                 |
|  | (CosmJS)       |  | (REST)         |  | (Real-time)     |                 |
|  +----------------+  +----------------+  +----------------+                 |
+-----------------------------------------------------------------------------+
                              |
                              v
+-----------------------------------------------------------------------------+
|  React Components (Subscribe to Store)                                      |
|  +----------------+  +----------------+  +----------------+                 |
|  | <ConnectButton>|  | <ChatContainer>|  | <MQBadge>      |                 |
|  | subscribes to  |  | subscribes to  |  | subscribes to  |                 |
|  | walletStore    |  | chatStore      |  | userStore      |                 |
|  +----------------+  +----------------+  +----------------+                 |
+-----------------------------------------------------------------------------+
```

---

## 5. Chain Interaction API Layer

### 5.1 Chain Client

```typescript
// src/core/chain/ChainClient.ts
import { QueryClient, setupBankExtension, setupAuthExtension } from '@cosmjs/stargate';
import { Tendermint34Client } from '@cosmjs/tendermint-rpc';

export class ChainClient {
  private tmClient: Tendermint34Client;
  private queryClient: QueryClient;

  constructor(rpcUrl: string) {
    this.tmClient = await Tendermint34Client.connect(rpcUrl);
    this.queryClient = QueryClient.withExtensions(
      this.tmClient,
      setupBankExtension,
      setupAuthExtension
    );
  }

  // Bank queries
  async getBalance(address: string, denom: string) {
    return this.queryClient.bank.balance(address, denom);
  }

  async getAllBalances(address: string) {
    return this.queryClient.bank.allBalances(address);
  }

  // Account queries
  async getAccount(address: string) {
    return this.queryClient.auth.account(address);
  }

  // Custom module queries (gRPC)
  async getMQInfo(address: string): Promise<MQInfo> {
    // Query x/trust module
    const response = await this.queryClient.queryUnverified(
      `/sharetokens/trust/mq/${address}`,
      {}
    );
    return MQInfo.fromProto(response);
  }

  async getService(serviceId: string): Promise<Service> {
    const response = await this.queryClient.queryUnverified(
      `/sharetokens/compute/service/${serviceId}`,
      {}
    );
    return Service.fromProto(response);
  }

  async getDispute(disputeId: string): Promise<Dispute> {
    const response = await this.queryClient.queryUnverified(
      `/sharetokens/trust/dispute/${disputeId}`,
      {}
    );
    return Dispute.fromProto(response);
  }
}
```

### 5.2 Transaction Client

```typescript
// src/core/chain/TxClient.ts
import { SigningStargateClient, EncodeObject } from '@cosmjs/stargate';
import { useWallet } from '@core/wallet/useWallet';

export class TxClient {
  constructor(private signingClient: SigningStargateClient) {}

  // Service marketplace transactions
  async registerService(registration: ServiceRegistration): Promise<TxResult> {
    const msg: EncodeObject = {
      typeUrl: '/sharetokens.compute.v1.MsgRegisterService',
      value: {
        owner: registration.owner,
        name: registration.name,
        level: registration.level,
        category: registration.category,
        description: registration.description,
        pricing: registration.pricing,
        endpoint: registration.endpoint,
      },
    };

    return this.signAndBroadcast([msg]);
  }

  async createServiceRequest(request: ServiceRequest): Promise<TxResult> {
    const msg: EncodeObject = {
      typeUrl: '/sharetokens.compute.v1.MsgSubmitRequest',
      value: {
        consumer: request.consumer,
        serviceId: request.serviceId,
        params: request.params,
        budget: request.budget,
      },
    };

    return this.signAndBroadcast([msg]);
  }

  // Escrow transactions
  async createEscrow(params: CreateEscrowParams): Promise<TxResult> {
    const msg: EncodeObject = {
      typeUrl: '/sharetokens.escrow.v1.MsgCreateEscrow',
      value: {
        creator: params.creator,
        beneficiary: params.beneficiary,
        amount: params.amount,
        duration: params.duration,
        metadata: params.metadata,
      },
    };

    return this.signAndBroadcast([msg]);
  }

  async releaseEscrow(escrowId: string, releaser: string): Promise<TxResult> {
    const msg: EncodeObject = {
      typeUrl: '/sharetokens.escrow.v1.MsgReleaseEscrow',
      value: { escrowId, releaser },
    };

    return this.signAndBroadcast([msg]);
  }

  // Dispute transactions
  async createDispute(params: DisputeParams): Promise<TxResult> {
    const msg: EncodeObject = {
      typeUrl: '/sharetokens.trust.v1.MsgCreateDispute',
      value: {
        plaintiff: params.plaintiff,
        defendant: params.defendant,
        orderId: params.orderId,
        title: params.title,
        description: params.description,
        evidence: params.evidence,
      },
    };

    return this.signAndBroadcast([msg]);
  }

  async castVote(disputeId: string, juror: string, verdict: string): Promise<TxResult> {
    const msg: EncodeObject = {
      typeUrl: '/sharetokens.trust.v1.MsgCastVote',
      value: { disputeId, juror, verdict },
    };

    return this.signAndBroadcast([msg]);
  }

  // Helper
  private async signAndBroadcast(msgs: EncodeObject[]): Promise<TxResult> {
    const address = this.getAddress();
    const fee = await this.estimateFee(msgs);

    const result = await this.signingClient.signAndBroadcast(address, msgs, fee);

    return {
      txHash: result.transactionHash,
      height: result.height,
      gasUsed: result.gasUsed,
      events: result.events,
    };
  }

  private async estimateFee(msgs: EncodeObject[]): Promise<StdFee> {
    const gasEstimation = await this.signingClient.simulate(
      this.getAddress(),
      msgs,
      ''
    );

    return {
      amount: [{ denom: 'ustt', amount: '5000' }],
      gas: Math.ceil(gasEstimation * 1.3).toString(), // 30% buffer
    };
  }
}
```

### 5.3 REST API Client

```typescript
// src/core/api/base.ts
import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

export class BaseApiClient {
  protected client: AxiosInstance;

  constructor(baseURL: string) {
    this.client = axios.create({
      baseURL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Add response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        const message = error.response?.data?.message || error.message;
        throw new Error(message);
      }
    );
  }

  protected async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.get<T>(url, config);
    return response.data;
  }

  protected async post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.post<T>(url, data, config);
    return response.data;
  }
}
```

```typescript
// src/core/api/marketplace.ts
import { BaseApiClient } from './base';

export class MarketplaceAPI extends BaseApiClient {
  // Service discovery
  async getServices(params: ServiceQuery): Promise<PaginatedResponse<Service>> {
    return this.get('/marketplace/services', { params });
  }

  async getService(serviceId: string): Promise<Service> {
    return this.get(`/marketplace/services/${serviceId}`);
  }

  async getServiceReviews(serviceId: string): Promise<Review[]> {
    return this.get(`/marketplace/services/${serviceId}/reviews`);
  }

  // Service requests
  async createRequest(request: ServiceRequestInput): Promise<RequestResult> {
    return this.post('/marketplace/requests', request);
  }

  async getRequest(requestId: string): Promise<ServiceRequest> {
    return this.get(`/marketplace/requests/${requestId}`);
  }

  // Streaming (SSE)
  streamRequest(requestId: string, onMessage: (delta: string) => void, onDone: () => void): () => void {
    const eventSource = new EventSource(
      `${this.baseURL}/marketplace/requests/${requestId}/stream`
    );

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.done) {
        eventSource.close();
        onDone();
      } else {
        onMessage(data.delta);
      }
    };

    eventSource.onerror = () => {
      eventSource.close();
    };

    return () => eventSource.close();
  }
}
```

---

## 6. Key Page Component Breakdown

### 6.1 Chat Page (Main Interface)

```typescript
// src/pages/ChatPage.tsx
import { ChatContainer } from '@components/chat/ChatContainer';
import { ChatInput } from '@components/chat/ChatInput';
import { WelcomeCard } from '@components/chat/WelcomeCard';
import { useChatStore } from '@stores/chatStore';
import { useWallet } from '@core/wallet/useWallet';

export function ChatPage() {
  const { messages, isTyping } = useChatStore();
  const { isConnected } = useWallet();

  return (
    <div className="flex flex-col h-full">
      {/* Main chat area */}
      <div className="flex-1 overflow-hidden">
        {messages.length === 0 ? (
          <WelcomeCard />
        ) : (
          <ChatContainer messages={messages} isTyping={isTyping} />
        )}
      </div>

      {/* Input area - fixed at bottom */}
      <div className="border-t border-gray-200 p-4">
        {!isConnected && (
          <div className="mb-4 p-3 bg-yellow-50 rounded-lg text-sm text-yellow-800">
            Please connect your wallet to start using GenieBot.
          </div>
        )}
        <ChatInput disabled={!isConnected} />
      </div>
    </div>
  );
}
```

### 6.2 Service Marketplace Page

```typescript
// src/pages/MarketplacePage.tsx
import { useState } from 'react';
import { ServiceFilter } from '@components/service/ServiceFilter';
import { ServiceCard } from '@components/service/ServiceCard';
import { ServiceLevelTabs } from '@components/service/ServiceLevelTabs';
import { useServices } from '@hooks/useService';

export function MarketplacePage() {
  const [level, setLevel] = useState<ServiceLevel | null>(null);
  const [filters, setFilters] = useState<ServiceFilters>({});

  const { services, isLoading, error } = useServices({ level, ...filters });

  return (
    <div className="space-y-6">
      <header>
        <h1>Service Marketplace</h1>
        <p>Browse and use AI services across three levels</p>
      </header>

      {/* Level tabs */}
      <ServiceLevelTabs value={level} onChange={setLevel} />

      {/* Filters */}
      <ServiceFilter values={filters} onChange={setFilters} />

      {/* Service grid */}
      {isLoading ? (
        <ServiceGridSkeleton />
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {services.map((service) => (
            <ServiceCard key={service.id} service={service} />
          ))}
        </div>
      )}
    </div>
  );
}
```

### 6.3 Idea Incubation Page

```typescript
// src/pages/IdeasPage.tsx
import { IdeaCard } from '@components/idea/IdeaCard';
import { EvalReport } from '@components/idea/EvalReport';
import { TokenEstimate } from '@components/idea/TokenEstimate';
import { useIdeas } from '@hooks/useIdeas';

export function IdeasPage() {
  const { ideas, selectedIdea, selectIdea, evaluateIdea } = useIdeas();

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      {/* Ideas list */}
      <div className="lg:col-span-2 space-y-4">
        <h2>My Ideas</h2>
        {ideas.map((idea) => (
          <IdeaCard
            key={idea.id}
            idea={idea}
            selected={selectedIdea?.id === idea.id}
            onClick={() => selectIdea(idea)}
          />
        ))}
      </div>

      {/* Idea detail / evaluation */}
      <div className="space-y-4">
        {selectedIdea && (
          <>
            <EvalReport evaluation={selectedIdea.evaluation} />
            <TokenEstimate estimate={selectedIdea.tokenEstimate} />
            <Button onClick={() => evaluateIdea(selectedIdea.id)}>
              Refresh Evaluation
            </Button>
          </>
        )}
      </div>
    </div>
  );
}
```

---

## 7. Responsive Design

### 7.1 Breakpoints (Tailwind)

```typescript
// tailwind.config.js
module.exports = {
  theme: {
    screens: {
      'sm': '640px',
      'md': '768px',
      'lg': '1024px',
      'xl': '1280px',
      '2xl': '1536px',
    },
  },
};
```

### 7.2 Layout Adaptations

| Component | Mobile (<768px) | Tablet (768-1024px) | Desktop (>1024px) |
|-----------|-----------------|---------------------|-------------------|
| Navigation | Hamburger drawer | Collapsed sidebar | Full sidebar |
| Chat | Full-screen | Full-width | 2/3 width + sidebar |
| Service Cards | Single column | 2 columns | 3 columns |
| Modals | Full-screen slide | Centered modal | Centered modal |
| Tables | Card list view | Scrollable table | Full table |

### 7.3 Mobile Navigation

```typescript
// src/core/layout/MobileNav.tsx
export function MobileNav() {
  const location = useLocation();

  const navItems = [
    { path: '/chat', icon: ChatIcon, label: 'Chat' },
    { path: '/market', icon: MarketIcon, label: 'Market' },
    { path: '/ideas', icon: IdeaIcon, label: 'Ideas' },
    { path: '/profile', icon: UserIcon, label: 'Profile' },
  ];

  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-white border-t md:hidden">
      <div className="flex justify-around py-2">
        {navItems.map((item) => (
          <Link
            key={item.path}
            to={item.path}
            className={cn(
              'flex flex-col items-center p-2 min-w-[44px] min-h-[44px]',
              location.pathname === item.path ? 'text-primary' : 'text-gray-500'
            )}
          >
            <item.icon />
            <span className="text-xs mt-1">{item.label}</span>
          </Link>
        ))}
      </div>
    </nav>
  );
}
```

---

## 8. Design Tokens

### 8.1 Colors

```css
:root {
  /* Brand */
  --color-primary: #6366F1;
  --color-primary-dark: #4F46E5;
  --color-primary-light: #818CF8;

  /* Service Levels */
  --color-level-1: #3B82F6;  /* Blue - LLM */
  --color-level-2: #8B5CF6;  /* Purple - Agent */
  --color-level-3: #F59E0B;  /* Gold - Workflow */

  /* Semantic */
  --color-success: #10B981;
  --color-warning: #F59E0B;
  --color-error: #EF4444;
  --color-info: #3B82F6;

  /* MQ Levels */
  --mq-lv1: #9CA3AF;  /* Gray - Newcomer */
  --mq-lv2: #3B82F6;  /* Blue - Member */
  --mq-lv3: #10B981;  /* Green - Trusted */
  --mq-lv4: #8B5CF6;  /* Purple - Expert */
  --mq-lv5: #F59E0B;  /* Gold - Guardian */
}
```

---

## 9. Installation and Setup

### 9.1 Development Setup

```bash
# Clone and install
git clone https://github.com/sharetokens/geniebot-frontend
cd geniebot-frontend
npm install

# Configure environment
cp .env.example .env.local
# Edit .env.local with your configuration

# Run development server
npm run dev
```

### 9.2 Environment Variables

```bash
# .env.example
VITE_CHAIN_ID=sharetokens-1
VITE_RPC_URL=https://rpc.sharetokens.io
VITE_API_URL=https://api.sharetokens.io
VITE_WS_URL=wss://ws.sharetokens.io
```

---

*Document Version: 3.0.0*
*Last Updated: 2026-03-02*
*Changes: Added detailed component architecture, Keplr integration flow, state management with Zustand, and chain interaction API layer*
