import {
  SigningStargateClient,
  GasPrice,
  StdFee,
} from '@cosmjs/stargate'
import type { OfflineSigner } from '@cosmjs/proto-signing'
import type {
  KeplrWindow,
  Account,
  WalletState,
  ChainInfo,
  SignedMessage,
} from '@/types'

// ShareTokens chain configuration
const SHARETOKENS_CHAIN_ID = 'sharetokens-1'
const SHARETOKENS_RPC = 'https://rpc.sharetokens.com'
const SHARETOKENS_REST = 'https://api.sharetokens.com'

const SHARETOKENS_CHAIN_INFO: ChainInfo = {
  chainId: SHARETOKENS_CHAIN_ID,
  chainName: 'ShareTokens',
  rpc: SHARETOKENS_RPC,
  rest: SHARETOKENS_REST,
  bip44: {
    coinType: 118,
  },
  bech32Config: {
    bech32PrefixAccAddr: 'share',
    bech32PrefixAccPub: 'sharepub',
    bech32PrefixValAddr: 'sharevaloper',
    bech32PrefixValPub: 'sharevaloperpub',
    bech32PrefixConsAddr: 'sharevalcons',
    bech32PrefixConsPub: 'sharevalconspub',
  },
  currencies: [
    {
      coinDenom: 'STT',
      coinMinimalDenom: 'ustt',
      coinDecimals: 6,
    },
  ],
  feeCurrencies: [
    {
      coinDenom: 'STT',
      coinMinimalDenom: 'ustt',
      coinDecimals: 6,
    },
  ],
  stakeCurrency: {
    coinDenom: 'STT',
    coinMinimalDenom: 'ustt',
    coinDecimals: 6,
  },
}

/**
 * Wallet Service
 * Manages Keplr wallet connection and blockchain interactions
 */
export class WalletService {
  private static instance: WalletService
  private client: SigningStargateClient | null = null
  private accounts: Account[] = []
  private chainId: string = SHARETOKENS_CHAIN_ID
  private listeners: Set<(state: WalletState) => void> = new Set()
  private state: WalletState = {
    address: null,
    balance: '0',
    isConnected: false,
    isConnecting: false,
    chainId: SHARETOKENS_CHAIN_ID,
  }

  static getInstance(): WalletService {
    if (!WalletService.instance) {
      WalletService.instance = new WalletService()
    }
    return WalletService.instance
  }

  private constructor() {
    this.setupKeplrListeners()
  }

  /**
   * Get current wallet state
   */
  getState(): WalletState {
    return { ...this.state }
  }

  /**
   * Subscribe to state changes
   */
  subscribe(listener: (state: WalletState) => void): () => void {
    this.listeners.add(listener)
    return () => this.listeners.delete(listener)
  }

  private notifyListeners(): void {
    this.listeners.forEach((listener) => listener(this.getState()))
  }

  private updateState(updates: Partial<WalletState>): void {
    this.state = { ...this.state, ...updates }
    this.notifyListeners()
  }

  /**
   * Check if Keplr is installed
   */
  async isKeplrInstalled(): Promise<boolean> {
    const keplrWindow = window as KeplrWindow

    if (keplrWindow.keplr) {
      return true
    }

    // Wait for Keplr to inject (max 1 second)
    if (document.readyState === 'complete') {
      return false
    }

    return new Promise((resolve) => {
      const timer = setTimeout(() => resolve(false), 1000)
      window.addEventListener('load', () => {
        clearTimeout(timer)
        resolve(!!keplrWindow.keplr)
      })
    })
  }

  /**
   * Connect to Keplr wallet
   */
  async connect(): Promise<string> {
    const keplrWindow = window as KeplrWindow

    try {
      this.updateState({ isConnecting: true })

      // Check Keplr installation
      const isInstalled = await this.isKeplrInstalled()
      if (!isInstalled) {
        throw new Error(
          'Please install Keplr extension. Visit https://www.keplr.app/'
        )
      }

      // Try to enable the chain
      try {
        await keplrWindow.keplr!.enable(this.chainId)
      } catch (enableError) {
        // Chain not registered, suggest it
        await keplrWindow.keplr!.experimentalSuggestChain(SHARETOKENS_CHAIN_INFO)
        await keplrWindow.keplr!.enable(this.chainId)
      }

      // Get offline signer
      const offlineSigner = keplrWindow.keplr!.getOfflineSigner(this.chainId)

      // Get accounts
      this.accounts = await offlineSigner.getAccounts()

      if (this.accounts.length === 0) {
        throw new Error('No accounts found in Keplr')
      }

      // Create signing client
      this.client = await SigningStargateClient.connectWithSigner(
        SHARETOKENS_RPC,
        offlineSigner,
        {
          gasPrice: GasPrice.fromString('0.025ustt'),
        }
      )

      const address = this.accounts[0].address

      // Get initial balance
      const balance = await this.getBalance(address)

      this.updateState({
        address,
        balance,
        isConnected: true,
        isConnecting: false,
      })

      return address
    } catch (error) {
      this.updateState({ isConnecting: false })
      throw error
    }
  }

  /**
   * Disconnect wallet
   */
  disconnect(): void {
    this.client = null
    this.accounts = []
    this.updateState({
      address: null,
      balance: '0',
      isConnected: false,
      isConnecting: false,
    })
  }

  /**
   * Get wallet address
   */
  getAddress(): string | null {
    return this.state.address
  }

  /**
   * Get account balance
   */
  async getBalance(address?: string): Promise<string> {
    if (!this.client) {
      throw new Error('Wallet not connected')
    }

    const addr = address || this.state.address
    if (!addr) {
      throw new Error('No address available')
    }

    try {
      const balance = await this.client.getBalance(addr, 'ustt')
      // Convert from uSTT to STT
      const sttAmount = parseFloat(balance.amount) / 1_000_000
      return sttAmount.toFixed(6)
    } catch (error) {
      console.error('Failed to get balance:', error)
      return '0'
    }
  }

  /**
   * Refresh balance
   */
  async refreshBalance(): Promise<void> {
    if (this.state.address) {
      const balance = await this.getBalance(this.state.address)
      this.updateState({ balance })
    }
  }

  /**
   * Send tokens
   */
  async sendTokens(
    recipientAddress: string,
    amount: string,
    memo: string = ''
  ): Promise<{ transactionHash: string }> {
    if (!this.client || !this.state.address) {
      throw new Error('Wallet not connected')
    }

    // Convert STT to uSTT
    const uamount = Math.floor(parseFloat(amount) * 1_000_000).toString()

    const result = await this.client.sendTokens(
      this.state.address,
      recipientAddress,
      [{ denom: 'ustt', amount: uamount }],
      'auto',
      memo
    )

    // Refresh balance after sending
    await this.refreshBalance()

    return result
  }

  /**
   * Sign a message for authentication
   */
  async signMessage(message: string): Promise<SignedMessage> {
    const keplrWindow = window as KeplrWindow

    if (!keplrWindow.keplr || !this.state.address) {
      throw new Error('Keplr not connected')
    }

    try {
      const result = await keplrWindow.keplr.signArbitrary(
        this.chainId,
        this.state.address,
        message
      )

      return {
        message,
        signature: result.signature,
        publicKey: result.pub_key.value,
        algorithm: result.pub_key.type,
      }
    } catch (error) {
      console.error('Message signing failed:', error)
      throw new Error('Failed to sign message')
    }
  }

  /**
   * Broadcast a signed transaction
   */
  async broadcastTransaction(
    signedTx: Uint8Array
  ): Promise<{ transactionHash: string }> {
    if (!this.client) {
      throw new Error('Wallet not connected')
    }

    const result = await this.client.broadcastTx(signedTx)
    await this.refreshBalance()
    return result
  }

  /**
   * Estimate transaction fee
   */
  async estimateFee(
    messages: unknown[],
    memo: string = ''
  ): Promise<StdFee> {
    if (!this.client || !this.state.address) {
      throw new Error('Wallet not connected')
    }

    const gasEstimation = await this.client.simulate(
      this.state.address,
      messages,
      memo
    )

    const gasLimit = Math.ceil(gasEstimation * 1.3) // Add 30% buffer

    return {
      amount: [{ denom: 'ustt', amount: (gasLimit * 0.025).toString() }],
      gas: gasLimit.toString(),
    }
  }

  /**
   * Check if connected
   */
  isConnected(): boolean {
    return this.state.isConnected
  }

  /**
   * Get signing client (for advanced usage)
   */
  getClient(): SigningStargateClient | null {
    return this.client
  }

  /**
   * Setup Keplr event listeners
   */
  private setupKeplrListeners(): void {
    // Listen for account changes
    window.addEventListener('keplr_keystorechange', () => {
      console.log('Keplr account changed')
      // Reconnect with new account
      if (this.state.isConnected) {
        this.connect().catch((error) => {
          console.error('Failed to reconnect after account change:', error)
        })
      }
    })
  }

  /**
   * Get chain info
   */
  getChainInfo(): ChainInfo {
    return SHARETOKENS_CHAIN_INFO
  }

  /**
   * Get chain ID
   */
  getChainId(): string {
    return this.chainId
  }

  /**
   * Switch to a different chain
   */
  async switchChain(chainId: string): Promise<void> {
    if (this.chainId === chainId) return

    this.chainId = chainId

    if (this.state.isConnected) {
      // Reconnect to new chain
      await this.connect()
    }
  }

  /**
   * Verify if address is valid
   */
  isValidAddress(address: string): boolean {
    // Check if address starts with the bech32 prefix
    return address.startsWith('share') && address.length === 43
  }
}

export const walletService = WalletService.getInstance()
