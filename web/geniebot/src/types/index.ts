export interface Message {
  id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  timestamp: number
  intent?: Intent
  services?: ServiceRecommendation[]
  signature?: string
  verified?: boolean
}

export interface Intent {
  type: 'llm' | 'agent' | 'workflow' | 'task' | 'query' | 'unknown'
  confidence: number
  entities: Entity[]
}

export interface Entity {
  type: string
  value: string
  start: number
  end: number
}

export interface ServiceRecommendation {
  id: string
  name: string
  type: 'llm' | 'agent' | 'workflow'
  description: string
  confidence: number
  estimatedCost: string
  estimatedTime: string
}

export interface Task {
  id: string
  name: string
  description: string
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled'
  progress: number
  createdAt: number
  startedAt?: number
  completedAt?: number
  result?: TaskResult
  txHash?: string
}

export interface TaskResult {
  content: string
  format: 'text' | 'json' | 'markdown' | 'code'
  downloadUrl?: string
  attachments?: Attachment[]
}

export interface Attachment {
  name: string
  type: string
  size: number
  url: string
}

export interface ChatSession {
  id: string
  title: string
  messages: Message[]
  tasks: Task[]
  createdAt: number
  updatedAt: number
}

// Wallet Types
export interface KeplrWindow extends Window {
  keplr?: {
    enable: (chainId: string) => Promise<void>
    getKey: (chainId: string) => Promise<{
      name: string
      pubKey: Uint8Array
      address: string
      algo: string
    }>
    signArbitrary: (
      chainId: string,
      signer: string,
      data: string
    ) => Promise<{
      signature: string
      pub_key: {
        type: string
        value: string
      }
    }>
    experimentalSuggestChain: (chainInfo: ChainInfo) => Promise<void>
    getOfflineSigner: (chainId: string) => OfflineSigner
  }
  getOfflineSigner?: (chainId: string) => OfflineSigner
}

export interface OfflineSigner {
  getAccounts: () => Promise<Account[]>
  signDirect: (
    signerAddress: string,
    signDoc: SignDoc
  ) => Promise<DirectSignResponse>
}

export interface Account {
  address: string
  algo: string
  pubkey: Uint8Array
}

export interface SignDoc {
  bodyBytes: Uint8Array
  authInfoBytes: Uint8Array
  chainId: string
  accountNumber: bigint
}

export interface DirectSignResponse {
  signed: SignDoc
  signature: {
    pub_key: {
      type: string
      value: string
    }
    signature: string
  }
}

export interface ChainInfo {
  chainId: string
  chainName: string
  rpc: string
  rest: string
  bip44: {
    coinType: number
  }
  bech32Config: {
    bech32PrefixAccAddr: string
    bech32PrefixAccPub: string
    bech32PrefixValAddr: string
    bech32PrefixValPub: string
    bech32PrefixConsAddr: string
    bech32PrefixConsPub: string
  }
  currencies: Currency[]
  feeCurrencies: Currency[]
  stakeCurrency: Currency
}

export interface Currency {
  coinDenom: string
  coinMinimalDenom: string
  coinDecimals: number
}

export interface WalletState {
  address: string | null
  balance: string
  isConnected: boolean
  isConnecting: boolean
  chainId: string
}

// A2A Protocol Types
export interface Agent {
  id: string
  name: string
  description: string
  capabilities: string[]
  endpoint: string
  skills: Skill[]
  reputation: number
  pricePerUnit: string
}

export interface Skill {
  id: string
  name: string
  description: string
  parameters: Parameter[]
}

export interface Parameter {
  name: string
  type: string
  required: boolean
  description: string
}

export interface A2ATask {
  id: string
  agentId: string
  skillId: string
  parameters: Record<string, unknown>
  status: 'pending' | 'submitted' | 'running' | 'completed' | 'failed'
  createdAt: number
  updatedAt: number
  result?: A2AResult
}

export interface A2AResult {
  content: string
  format: string
  artifacts?: Artifact[]
}

export interface Artifact {
  name: string
  type: string
  content: string
}

export interface A2AMessage {
  type: 'task_request' | 'task_response' | 'task_update' | 'heartbeat'
  taskId: string
  payload: unknown
  timestamp: number
  signature?: string
}

// DID Types
export interface DIDDocument {
  id: string
  verificationMethod: VerificationMethod[]
  authentication: string[]
  assertionMethod: string[]
  service: ServiceEndpoint[]
}

export interface VerificationMethod {
  id: string
  type: string
  controller: string
  publicKeyHex?: string
  publicKeyBase64?: string
}

export interface ServiceEndpoint {
  id: string
  type: string
  serviceEndpoint: string
}

export interface AuthState {
  isAuthenticated: boolean
  isLoading: boolean
  did: string | null
  user: UserProfile | null
}

export interface UserProfile {
  address: string
  did: string
  name?: string
  avatar?: string
  reputation: number
  createdAt: number
}

export interface SignedMessage {
  message: string
  signature: string
  publicKey: string
  algorithm: string
}
