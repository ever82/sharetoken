import type {
  Message,
  Task,
  ChatSession,
  ServiceRecommendation,
  Intent,
} from '@/types'
import { a2aService } from './a2a'
import { walletService } from './wallet'
import { authService } from './auth'

const API_BASE = '/api'

/**
 * Enhanced API Service
 * Integrates with A2A protocol, wallet, and authentication
 */
export class ApiService {
  private static instance: ApiService
  private ws: WebSocket | null = null
  private messageHandlers: ((message: unknown) => void)[] = []

  static getInstance(): ApiService {
    if (!ApiService.instance) {
      ApiService.instance = new ApiService()
    }
    return ApiService.instance
  }

  /**
   * Get authentication headers
   */
  private getHeaders(): Record<string, string> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }

    const token = authService.getToken()
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    return headers
  }

  async sendMessage(content: string, sessionId?: string): Promise<Message> {
    const response = await fetch(`${API_BASE}/chat/message`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ content, sessionId }),
    })

    if (!response.ok) {
      throw new Error('Failed to send message')
    }

    return response.json()
  }

  async detectIntent(content: string): Promise<Intent> {
    const response = await fetch(`${API_BASE}/chat/intent`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ content }),
    })

    if (!response.ok) {
      throw new Error('Failed to detect intent')
    }

    return response.json()
  }

  async getServiceRecommendations(
    content: string
  ): Promise<ServiceRecommendation[]> {
    // First try A2A service discovery
    const intent = await this.detectIntent(content)

    if (intent.type !== 'unknown' && intent.confidence > 0.5) {
      const agents = await a2aService.matchIntentToAgent(
        intent.type,
        intent.entities
      )

      if (agents) {
        // Convert agents to service recommendations
        return [
          {
            id: agents.id,
            name: agents.name,
            type:
              intent.type === 'llm'
                ? 'llm'
                : intent.type === 'agent'
                  ? 'agent'
                  : 'workflow',
            description: agents.description,
            confidence: intent.confidence,
            estimatedCost: agents.pricePerUnit,
            estimatedTime: '~2 min',
          },
        ]
      }
    }

    // Fallback to API
    const response = await fetch(`${API_BASE}/services/recommend`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ content }),
    })

    if (!response.ok) {
      throw new Error('Failed to get recommendations')
    }

    return response.json()
  }

  async invokeService(
    serviceId: string,
    params: Record<string, unknown>
  ): Promise<Task> {
    // Check if this is an A2A agent
    const agent = await a2aService.getAgent(serviceId)

    if (agent) {
      // Use A2A protocol
      const task = await a2aService.submitTask(serviceId, 'default', params)

      return {
        id: task.id,
        name: agent.name,
        description: `A2A task from ${agent.name}`,
        status: task.status === 'completed' ? 'completed' : 'running',
        progress: task.status === 'completed' ? 100 : 0,
        createdAt: task.createdAt,
        result: task.result
          ? {
              content: task.result.content,
              format: task.result.format as 'text' | 'json' | 'markdown' | 'code',
            }
          : undefined,
      }
    }

    // Fallback to regular API
    const response = await fetch(`${API_BASE}/services/invoke`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ serviceId, params }),
    })

    if (!response.ok) {
      throw new Error('Failed to invoke service')
    }

    return response.json()
  }

  async getTasks(): Promise<Task[]> {
    const response = await fetch(`${API_BASE}/tasks`, {
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      throw new Error('Failed to get tasks')
    }

    return response.json()
  }

  async getTask(taskId: string): Promise<Task> {
    // Check A2A tasks first
    const a2aTask = await a2aService.getTaskStatus(taskId)
    if (a2aTask) {
      return {
        id: a2aTask.id,
        name: 'A2A Task',
        description: 'Task from A2A protocol',
        status: a2aTask.status === 'completed' ? 'completed' : 'running',
        progress: a2aTask.status === 'completed' ? 100 : 0,
        createdAt: a2aTask.createdAt,
        result: a2aTask.result
          ? {
              content: a2aTask.result.content,
              format: a2aTask.result.format as
                | 'text'
                | 'json'
                | 'markdown'
                | 'code',
            }
          : undefined,
      }
    }

    const response = await fetch(`${API_BASE}/tasks/${taskId}`, {
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      throw new Error('Failed to get task')
    }

    return response.json()
  }

  async downloadResult(taskId: string): Promise<Blob> {
    const response = await fetch(`${API_BASE}/tasks/${taskId}/download`, {
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      throw new Error('Failed to download result')
    }

    return response.blob()
  }

  async getSessions(): Promise<ChatSession[]> {
    const response = await fetch(`${API_BASE}/sessions`, {
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      throw new Error('Failed to get sessions')
    }

    return response.json()
  }

  async createSession(): Promise<ChatSession> {
    const response = await fetch(`${API_BASE}/sessions`, {
      method: 'POST',
      headers: this.getHeaders(),
    })

    if (!response.ok) {
      throw new Error('Failed to create session')
    }

    return response.json()
  }

  async signAndSendMessage(content: string): Promise<Message> {
    // Sign message with wallet for verification
    const signedMessage = await walletService.signMessage(content)

    const response = await fetch(`${API_BASE}/chat/message`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({
        content,
        signature: signedMessage.signature,
        publicKey: signedMessage.publicKey,
      }),
    })

    if (!response.ok) {
      throw new Error('Failed to send signed message')
    }

    return response.json()
  }

  connectWebSocket(): void {
    if (this.ws?.readyState === WebSocket.OPEN) return

    this.ws = new WebSocket('/websocket')

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      this.messageHandlers.forEach((handler) => handler(data))
    }

    this.ws.onclose = () => {
      setTimeout(() => this.connectWebSocket(), 3000)
    }

    // Also connect A2A WebSocket
    a2aService.connectWebSocket()
  }

  disconnectWebSocket(): void {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    a2aService.disconnectWebSocket()
  }

  onMessage(handler: (message: unknown) => void): void {
    this.messageHandlers.push(handler)
  }

  offMessage(handler: (message: unknown) => void): void {
    this.messageHandlers = this.messageHandlers.filter((h) => h !== handler)
  }

  /**
   * Poll task for result using A2A protocol
   */
  async pollTaskResult(
    taskId: string,
    onUpdate: (task: Task) => void
  ): Promise<Task['result']> {
    const a2aResult = await a2aService.pollForResult(
      taskId,
      (a2aTask) => {
        onUpdate({
          id: a2aTask.id,
          name: 'A2A Task',
          description: 'Task from A2A protocol',
          status: a2aTask.status === 'completed' ? 'completed' : 'running',
          progress: a2aTask.status === 'completed' ? 100 : a2aTask.progress || 0,
          createdAt: a2aTask.createdAt,
        })
      }
    )

    return {
      content: a2aResult.content,
      format: a2aResult.format as 'text' | 'json' | 'markdown' | 'code',
    }
  }

  /**
   * Execute blockchain transaction
   */
  async executeTransaction(
    recipient: string,
    amount: string,
    memo: string = ''
  ): Promise<{ transactionHash: string }> {
    return walletService.sendTokens(recipient, amount, memo)
  }
}

export const apiService = ApiService.getInstance()
