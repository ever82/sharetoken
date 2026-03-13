import type { Agent, A2ATask, A2AResult, A2AMessage, Skill } from '@/types'

const A2A_API_BASE = '/api/a2a'

/**
 * A2A Protocol Service
 * Implements Agent-to-Agent communication protocol for GenieBot
 */
export class A2AService {
  private static instance: A2AService
  private ws: WebSocket | null = null
  private messageHandlers: ((message: A2AMessage) => void)[] = []
  private agents: Map<string, Agent> = new Map()
  private tasks: Map<string, A2ATask> = new Map()

  static getInstance(): A2AService {
    if (!A2AService.instance) {
      A2AService.instance = new A2AService()
    }
    return A2AService.instance
  }

  /**
   * Discover available agents from the network
   */
  async discoverAgents(): Promise<Agent[]> {
    try {
      const response = await fetch(`${A2A_API_BASE}/agents`)

      if (!response.ok) {
        throw new Error('Failed to discover agents')
      }

      const agents: Agent[] = await response.json()

      // Cache agents
      agents.forEach((agent) => {
        this.agents.set(agent.id, agent)
      })

      return agents
    } catch (error) {
      console.error('Agent discovery failed:', error)
      // Return cached agents if available
      return Array.from(this.agents.values())
    }
  }

  /**
   * Search agents by capability
   */
  async searchAgents(capability: string): Promise<Agent[]> {
    try {
      const response = await fetch(
        `${A2A_API_BASE}/agents/search?capability=${encodeURIComponent(capability)}`
      )

      if (!response.ok) {
        throw new Error('Failed to search agents')
      }

      return response.json()
    } catch (error) {
      console.error('Agent search failed:', error)
      // Filter cached agents
      return Array.from(this.agents.values()).filter((agent) =>
        agent.capabilities.some((cap) =>
          cap.toLowerCase().includes(capability.toLowerCase())
        )
      )
    }
  }

  /**
   * Get agent details
   */
  async getAgent(agentId: string): Promise<Agent | null> {
    // Check cache first
    if (this.agents.has(agentId)) {
      return this.agents.get(agentId)!
    }

    try {
      const response = await fetch(`${A2A_API_BASE}/agents/${agentId}`)

      if (!response.ok) {
        throw new Error('Failed to get agent')
      }

      const agent: Agent = await response.json()
      this.agents.set(agent.id, agent)
      return agent
    } catch (error) {
      console.error('Failed to get agent:', error)
      return null
    }
  }

  /**
   * Submit a task to an agent
   */
  async submitTask(
    agentId: string,
    skillId: string,
    parameters: Record<string, unknown>,
    signature?: string
  ): Promise<A2ATask> {
    const taskData = {
      agentId,
      skillId,
      parameters,
      signature,
    }

    const response = await fetch(`${A2A_API_BASE}/tasks`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(taskData),
    })

    if (!response.ok) {
      throw new Error('Failed to submit task')
    }

    const task: A2ATask = await response.json()
    this.tasks.set(task.id, task)

    return task
  }

  /**
   * Get task status and result
   */
  async getTaskStatus(taskId: string): Promise<A2ATask | null> {
    // Check cache first
    const cached = this.tasks.get(taskId)
    if (cached && ['completed', 'failed'].includes(cached.status)) {
      return cached
    }

    try {
      const response = await fetch(`${A2A_API_BASE}/tasks/${taskId}`)

      if (!response.ok) {
        throw new Error('Failed to get task status')
      }

      const task: A2ATask = await response.json()
      this.tasks.set(task.id, task)
      return task
    } catch (error) {
      console.error('Failed to get task status:', error)
      return cached || null
    }
  }

  /**
   * Poll for task result
   */
  async pollForResult(
    taskId: string,
    onUpdate: (task: A2ATask) => void,
    interval: number = 2000,
    timeout: number = 300000
  ): Promise<A2AResult> {
    const startTime = Date.now()

    return new Promise((resolve, reject) => {
      const poll = async () => {
        try {
          const task = await this.getTaskStatus(taskId)

          if (!task) {
            reject(new Error('Task not found'))
            return
          }

          onUpdate(task)

          if (task.status === 'completed' && task.result) {
            resolve(task.result)
            return
          }

          if (task.status === 'failed') {
            reject(new Error('Task failed'))
            return
          }

          if (Date.now() - startTime > timeout) {
            reject(new Error('Task polling timeout'))
            return
          }

          setTimeout(poll, interval)
        } catch (error) {
          reject(error)
        }
      }

      poll()
    })
  }

  /**
   * Cancel a running task
   */
  async cancelTask(taskId: string): Promise<boolean> {
    try {
      const response = await fetch(`${A2A_API_BASE}/tasks/${taskId}/cancel`, {
        method: 'POST',
      })

      return response.ok
    } catch (error) {
      console.error('Failed to cancel task:', error)
      return false
    }
  }

  /**
   * Connect to A2A WebSocket for real-time updates
   */
  connectWebSocket(): void {
    if (this.ws?.readyState === WebSocket.OPEN) return

    this.ws = new WebSocket('/websocket/a2a')

    this.ws.onmessage = (event) => {
      try {
        const message: A2AMessage = JSON.parse(event.data)
        this.handleMessage(message)
      } catch (error) {
        console.error('Failed to parse A2A message:', error)
      }
    }

    this.ws.onclose = () => {
      console.log('A2A WebSocket closed, reconnecting...')
      setTimeout(() => this.connectWebSocket(), 3000)
    }

    this.ws.onerror = (error) => {
      console.error('A2A WebSocket error:', error)
    }
  }

  /**
   * Disconnect WebSocket
   */
  disconnectWebSocket(): void {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  /**
   * Subscribe to A2A messages
   */
  onMessage(handler: (message: A2AMessage) => void): () => void {
    this.messageHandlers.push(handler)
    return () => {
      this.messageHandlers = this.messageHandlers.filter((h) => h !== handler)
    }
  }

  private handleMessage(message: A2AMessage): void {
    // Update task cache if it's a task update
    if (message.type === 'task_update') {
      const task = message.payload as A2ATask
      this.tasks.set(task.id, task)
    }

    // Notify all handlers
    this.messageHandlers.forEach((handler) => handler(message))
  }

  /**
   * Get cached agents
   */
  getCachedAgents(): Agent[] {
    return Array.from(this.agents.values())
  }

  /**
   * Get cached tasks
   */
  getCachedTasks(): A2ATask[] {
    return Array.from(this.tasks.values())
  }

  /**
   * Find best agent for a task based on capabilities and reputation
   */
  findBestAgent(
    requiredCapabilities: string[],
    minReputation: number = 0.5
  ): Agent | null {
    const agents = Array.from(this.agents.values()).filter(
      (agent) =>
        agent.reputation >= minReputation &&
        requiredCapabilities.every((cap) =>
          agent.capabilities.some((c) => c.toLowerCase() === cap.toLowerCase())
        )
    )

    if (agents.length === 0) return null

    // Sort by reputation and price
    return agents.sort((a, b) => {
      const scoreA = a.reputation * 100 - parseFloat(a.pricePerUnit)
      const scoreB = b.reputation * 100 - parseFloat(b.pricePerUnit)
      return scoreB - scoreA
    })[0]
  }

  /**
   * Match intent to agent capabilities
   */
  async matchIntentToAgent(
    intentType: string,
    entities: { type: string; value: string }[]
  ): Promise<Agent | null> {
    // Map intent types to required capabilities
    const capabilityMap: Record<string, string[]> = {
      llm: ['text-generation', 'completion'],
      agent: ['agent-execution', 'task-running'],
      workflow: ['workflow-orchestration', 'multi-step'],
      task: ['task-execution', 'computation'],
      query: ['information-retrieval', 'query-processing'],
    }

    const requiredCapabilities = capabilityMap[intentType] || ['general']

    // Try to find from cached agents first
    let bestAgent = this.findBestAgent(requiredCapabilities)

    if (!bestAgent) {
      // Try to discover new agents
      await this.discoverAgents()
      bestAgent = this.findBestAgent(requiredCapabilities)
    }

    return bestAgent
  }
}

export const a2aService = A2AService.getInstance()
