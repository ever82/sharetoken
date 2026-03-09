import type { Message, Task, ChatSession, ServiceRecommendation, Intent } from '@/types'

const API_BASE = '/api'

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

  async sendMessage(content: string, sessionId?: string): Promise<Message> {
    const response = await fetch(`${API_BASE}/chat/message`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
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
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content }),
    })

    if (!response.ok) {
      throw new Error('Failed to detect intent')
    }

    return response.json()
  }

  async getServiceRecommendations(content: string): Promise<ServiceRecommendation[]> {
    const response = await fetch(`${API_BASE}/services/recommend`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content }),
    })

    if (!response.ok) {
      throw new Error('Failed to get recommendations')
    }

    return response.json()
  }

  async invokeService(serviceId: string, params: Record<string, unknown>): Promise<Task> {
    const response = await fetch(`${API_BASE}/services/invoke`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ serviceId, params }),
    })

    if (!response.ok) {
      throw new Error('Failed to invoke service')
    }

    return response.json()
  }

  async getTasks(): Promise<Task[]> {
    const response = await fetch(`${API_BASE}/tasks`)

    if (!response.ok) {
      throw new Error('Failed to get tasks')
    }

    return response.json()
  }

  async getTask(taskId: string): Promise<Task> {
    const response = await fetch(`${API_BASE}/tasks/${taskId}`)

    if (!response.ok) {
      throw new Error('Failed to get task')
    }

    return response.json()
  }

  async downloadResult(taskId: string): Promise<Blob> {
    const response = await fetch(`${API_BASE}/tasks/${taskId}/download`)

    if (!response.ok) {
      throw new Error('Failed to download result')
    }

    return response.blob()
  }

  async getSessions(): Promise<ChatSession[]> {
    const response = await fetch(`${API_BASE}/sessions`)

    if (!response.ok) {
      throw new Error('Failed to get sessions')
    }

    return response.json()
  }

  async createSession(): Promise<ChatSession> {
    const response = await fetch(`${API_BASE}/sessions`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    })

    if (!response.ok) {
      throw new Error('Failed to create session')
    }

    return response.json()
  }

  connectWebSocket(): void {
    if (this.ws?.readyState === WebSocket.OPEN) return

    this.ws = new WebSocket('/websocket')

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      this.messageHandlers.forEach(handler => handler(data))
    }

    this.ws.onclose = () => {
      setTimeout(() => this.connectWebSocket(), 3000)
    }
  }

  disconnectWebSocket(): void {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  onMessage(handler: (message: unknown) => void): void {
    this.messageHandlers.push(handler)
  }

  offMessage(handler: (message: unknown) => void): void {
    this.messageHandlers = this.messageHandlers.filter(h => h !== handler)
  }
}

export const apiService = ApiService.getInstance()
