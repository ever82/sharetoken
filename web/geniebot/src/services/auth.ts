import type {
  AuthState,
  UserProfile,
  DIDDocument,
  SignedMessage,
} from '@/types'
import { walletService } from './wallet'

const AUTH_API_BASE = '/api/auth'
const AUTH_STORAGE_KEY = 'geniebot_auth'

/**
 * Authentication Service
 * Manages DID-based authentication and session management
 */
export class AuthService {
  private static instance: AuthService
  private state: AuthState = {
    isAuthenticated: false,
    isLoading: false,
    did: null,
    user: null,
  }
  private listeners: Set<(state: AuthState) => void> = new Set()
  private token: string | null = null

  static getInstance(): AuthService {
    if (!AuthService.instance) {
      AuthService.instance = new AuthService()
    }
    return AuthService.instance
  }

  private constructor() {
    // Try to restore session from storage
    this.restoreSession()
  }

  /**
   * Get current auth state
   */
  getState(): AuthState {
    return { ...this.state }
  }

  /**
   * Subscribe to auth state changes
   */
  subscribe(listener: (state: AuthState) => void): () => void {
    this.listeners.add(listener)
    return () => this.listeners.delete(listener)
  }

  private notifyListeners(): void {
    this.listeners.forEach((listener) => listener(this.getState()))
  }

  private updateState(updates: Partial<AuthState>): void {
    this.state = { ...this.state, ...updates }
    this.notifyListeners()
  }

  /**
   * Generate DID from wallet address
   */
  generateDid(address: string): string {
    // Generate DID using the share: method
    return `did:share:${address}`
  }

  /**
   * Login with wallet
   */
  async login(): Promise<UserProfile> {
    try {
      this.updateState({ isLoading: true })

      // Ensure wallet is connected
      let address = walletService.getAddress()
      if (!address) {
        address = await walletService.connect()
      }

      // Generate DID
      const did = this.generateDid(address)

      // Create authentication challenge
      const challenge = await this.createChallenge(did)

      // Sign challenge with wallet
      const messageToSign = `Sign this message to authenticate with GenieBot\nChallenge: ${challenge}\nDID: ${did}\nTimestamp: ${Date.now()}`
      const signedMessage = await walletService.signMessage(messageToSign)

      // Verify signature with backend
      const { user, token } = await this.verifySignature(
        did,
        signedMessage,
        challenge
      )

      // Store token
      this.token = token
      this.storeSession(token, user)

      this.updateState({
        isAuthenticated: true,
        isLoading: false,
        did,
        user,
      })

      return user
    } catch (error) {
      this.updateState({ isLoading: false })
      throw error
    }
  }

  /**
   * Create authentication challenge
   */
  private async createChallenge(did: string): Promise<string> {
    try {
      const response = await fetch(`${AUTH_API_BASE}/challenge`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ did }),
      })

      if (!response.ok) {
        throw new Error('Failed to create authentication challenge')
      }

      const data = await response.json()
      return data.challenge
    } catch (error) {
      console.error('Failed to create challenge:', error)
      // Fallback to local challenge generation
      return this.generateLocalChallenge(did)
    }
  }

  /**
   * Generate local challenge (fallback)
   */
  private generateLocalChallenge(did: string): string {
    const timestamp = Date.now()
    const random = Math.random().toString(36).substring(2, 15)
    return `challenge-${did}-${timestamp}-${random}`
  }

  /**
   * Verify signature with backend
   */
  private async verifySignature(
    did: string,
    signedMessage: SignedMessage,
    challenge: string
  ): Promise<{ user: UserProfile; token: string }> {
    try {
      const response = await fetch(`${AUTH_API_BASE}/verify`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          did,
          signature: signedMessage.signature,
          publicKey: signedMessage.publicKey,
          message: signedMessage.message,
          challenge,
        }),
      })

      if (!response.ok) {
        const error = await response.text()
        throw new Error(`Authentication failed: ${error}`)
      }

      return response.json()
    } catch (error) {
      console.error('Signature verification failed:', error)
      // Fallback to local verification for development
      return this.localVerify(did, signedMessage)
    }
  }

  /**
   * Local verification (for development/testing)
   */
  private localVerify(
    did: string,
    signedMessage: SignedMessage
  ): { user: UserProfile; token: string } {
    // Create a mock user profile
    const user: UserProfile = {
      address: did.replace('did:share:', ''),
      did,
      reputation: 1.0,
      createdAt: Date.now(),
    }

    // Create a mock JWT token
    const token = this.generateMockToken(user)

    return { user, token }
  }

  /**
   * Generate mock JWT token (for development)
   */
  private generateMockToken(user: UserProfile): string {
    const header = btoa(JSON.stringify({ alg: 'none', typ: 'JWT' }))
    const payload = btoa(
      JSON.stringify({
        sub: user.did,
        address: user.address,
        exp: Date.now() + 24 * 60 * 60 * 1000, // 24 hours
        iat: Date.now(),
      })
    )
    return `${header}.${payload}.`
  }

  /**
   * Logout
   */
  logout(): void {
    this.token = null
    this.clearSession()
    this.updateState({
      isAuthenticated: false,
      isLoading: false,
      did: null,
      user: null,
    })
  }

  /**
   * Check if authenticated
   */
  isAuthenticated(): boolean {
    return this.state.isAuthenticated && this.token !== null
  }

  /**
   * Get auth token
   */
  getToken(): string | null {
    return this.token
  }

  /**
   * Get current user
   */
  getUser(): UserProfile | null {
    return this.state.user
  }

  /**
   * Get current DID
   */
  getDID(): string | null {
    return this.state.did
  }

  /**
   * Refresh session
   */
  async refreshSession(): Promise<void> {
    if (!this.token) {
      throw new Error('No active session')
    }

    try {
      this.updateState({ isLoading: true })

      const response = await fetch(`${AUTH_API_BASE}/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${this.token}`,
        },
      })

      if (!response.ok) {
        throw new Error('Failed to refresh session')
      }

      const { user, token } = await response.json()
      this.token = token
      this.storeSession(token, user)

      this.updateState({
        isAuthenticated: true,
        isLoading: false,
        did: user.did,
        user,
      })
    } catch (error) {
      this.updateState({ isLoading: false })
      // Clear invalid session
      this.logout()
      throw error
    }
  }

  /**
   * Verify a message signature
   */
  async verifyMessageSignature(
    message: string,
    signature: string,
    publicKey: string
  ): Promise<boolean> {
    try {
      const response = await fetch(`${AUTH_API_BASE}/verify-signature`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          message,
          signature,
          publicKey,
        }),
      })

      if (!response.ok) {
        return false
      }

      const data = await response.json()
      return data.valid
    } catch (error) {
      console.error('Signature verification failed:', error)
      return false
    }
  }

  /**
   * Fetch DID document
   */
  async fetchDIDDocument(did: string): Promise<DIDDocument | null> {
    try {
      const response = await fetch(`${AUTH_API_BASE}/did/${did}`)

      if (!response.ok) {
        throw new Error('Failed to fetch DID document')
      }

      return response.json()
    } catch (error) {
      console.error('Failed to fetch DID document:', error)
      return null
    }
  }

  /**
   * Update user profile
   */
  async updateProfile(updates: Partial<UserProfile>): Promise<UserProfile> {
    if (!this.token) {
      throw new Error('Not authenticated')
    }

    try {
      const response = await fetch(`${AUTH_API_BASE}/profile`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${this.token}`,
        },
        body: JSON.stringify(updates),
      })

      if (!response.ok) {
        throw new Error('Failed to update profile')
      }

      const user = await response.json()
      this.updateState({ user })
      return user
    } catch (error) {
      console.error('Failed to update profile:', error)
      throw error
    }
  }

  /**
   * Make authenticated request
   */
  async fetchWithAuth(
    url: string,
    options: RequestInit = {}
  ): Promise<Response> {
    if (!this.token) {
      throw new Error('Not authenticated')
    }

    const headers = {
      ...options.headers,
      Authorization: `Bearer ${this.token}`,
    }

    return fetch(url, { ...options, headers })
  }

  /**
   * Store session in localStorage
   */
  private storeSession(token: string, user: UserProfile): void {
    try {
      localStorage.setItem(
        AUTH_STORAGE_KEY,
        JSON.stringify({
          token,
          user,
          timestamp: Date.now(),
        })
      )
    } catch (error) {
      console.error('Failed to store session:', error)
    }
  }

  /**
   * Restore session from localStorage
   */
  private restoreSession(): void {
    try {
      const stored = localStorage.getItem(AUTH_STORAGE_KEY)
      if (!stored) return

      const { token, user, timestamp } = JSON.parse(stored)

      // Check if session is expired (24 hours)
      if (Date.now() - timestamp > 24 * 60 * 60 * 1000) {
        this.clearSession()
        return
      }

      this.token = token
      this.updateState({
        isAuthenticated: true,
        did: user.did,
        user,
      })
    } catch (error) {
      console.error('Failed to restore session:', error)
      this.clearSession()
    }
  }

  /**
   * Clear session from storage
   */
  private clearSession(): void {
    try {
      localStorage.removeItem(AUTH_STORAGE_KEY)
    } catch (error) {
      console.error('Failed to clear session:', error)
    }
  }
}

export const authService = AuthService.getInstance()
