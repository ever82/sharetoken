import { useState, useCallback, useEffect, useRef } from 'react'
import type {
  Message,
  Task,
  ChatSession,
  ServiceRecommendation,
  Intent,
} from '@/types'
import { apiService } from '@/services/api'
import { a2aService } from '@/services/a2a'
import { walletService } from '@/services/wallet'
import { authService } from '@/services/auth'

export function useChat() {
  const [messages, setMessages] = useState<Message[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [currentIntent, setCurrentIntent] = useState<Intent | null>(null)
  const [recommendations, setRecommendations] = useState<ServiceRecommendation[]>([])

  const sendMessage = useCallback(async (content: string, signed = false) => {
    setIsLoading(true)

    // Add user message
    const userMessageId = Date.now().toString()
    const userMessage: Message = {
      id: userMessageId,
      role: 'user',
      content,
      timestamp: Date.now(),
    }
    setMessages((prev) => [...prev, userMessage])

    try {
      // Detect intent
      const intent = await apiService.detectIntent(content)
      setCurrentIntent(intent)

      // Get service recommendations if applicable
      if (intent.type !== 'unknown' && intent.confidence > 0.5) {
        const recs = await apiService.getServiceRecommendations(content)
        setRecommendations(recs)
      }

      // Send message with signature if requested
      let response: Message
      if (signed && authService.isAuthenticated()) {
        response = await apiService.signAndSendMessage(content)
      } else {
        response = await apiService.sendMessage(content)
      }

      // Add assistant message
      const assistantMessage: Message = {
        ...response,
        intent,
        services: recommendations.length > 0 ? recommendations : undefined,
      }
      setMessages((prev) => [...prev, assistantMessage])
    } catch (error) {
      // Add error message
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: 'Sorry, I encountered an error. Please try again.',
        timestamp: Date.now(),
      }
      setMessages((prev) => [...prev, errorMessage])
    } finally {
      setIsLoading(false)
    }
  }, [recommendations])

  const clearChat = useCallback(() => {
    setMessages([])
    setCurrentIntent(null)
    setRecommendations([])
  }, [])

  return {
    messages,
    isLoading,
    currentIntent,
    recommendations,
    sendMessage,
    clearChat,
  }
}

export function useTasks() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [isLoading, setIsLoading] = useState(false)

  const fetchTasks = useCallback(async () => {
    setIsLoading(true)
    try {
      const data = await apiService.getTasks()
      // Merge with A2A tasks
      const a2aTasks = a2aService.getCachedTasks().map((task) => ({
        id: task.id,
        name: task.agentId,
        description: `A2A task for skill ${task.skillId}`,
        status: task.status === 'completed' ? 'completed' : 'running',
        progress: task.status === 'completed' ? 100 : 0,
        createdAt: task.createdAt,
        result: task.result
          ? {
              content: task.result.content,
              format: task.result.format as 'text' | 'json' | 'markdown' | 'code',
            }
          : undefined,
      }))
      setTasks([...a2aTasks, ...data])
    } finally {
      setIsLoading(false)
    }
  }, [])

  const invokeService = useCallback(
    async (serviceId: string, params: Record<string, unknown>) => {
      setIsLoading(true)
      try {
        // Check if this is an A2A agent
        const agent = await a2aService.getAgent(serviceId)
        if (agent) {
          const task = await a2aService.submitTask(
            serviceId,
            'default',
            params,
            undefined // Will sign if authenticated
          )
          const newTask: Task = {
            id: task.id,
            name: agent.name,
            description: `A2A task from ${agent.name}`,
            status: 'pending',
            progress: 0,
            createdAt: task.createdAt,
          }
          setTasks((prev) => [newTask, ...prev])
          return newTask
        }

        const task = await apiService.invokeService(serviceId, params)
        setTasks((prev) => [task, ...prev])
        return task
      } finally {
        setIsLoading(false)
      }
    },
    []
  )

  const updateTaskProgress = useCallback((taskId: string, progress: number) => {
    setTasks((prev) =>
      prev.map((task) =>
        task.id === taskId ? { ...task, progress } : task
      )
    )
  }, [])

  const updateTaskStatus = useCallback((taskId: string, status: Task['status']) => {
    setTasks((prev) =>
      prev.map((task) =>
        task.id === taskId ? { ...task, status } : task
      )
    )
  }, [])

  const pollTaskResult = useCallback(
    async (taskId: string, onUpdate: (task: Task) => void) => {
      try {
        const result = await apiService.pollTaskResult(taskId, onUpdate)
        return result
      } catch (error) {
        console.error('Failed to poll task result:', error)
        throw error
      }
    },
    []
  )

  useEffect(() => {
    fetchTasks()
  }, [fetchTasks])

  return {
    tasks,
    isLoading,
    invokeService,
    updateTaskProgress,
    updateTaskStatus,
    pollTaskResult,
    refreshTasks: fetchTasks,
  }
}

export function useWebSocket() {
  const [isConnected, setIsConnected] = useState(false)
  const [lastMessage, setLastMessage] = useState<unknown>(null)

  useEffect(() => {
    apiService.connectWebSocket()
    setIsConnected(true)

    const handleMessage = (message: unknown) => {
      setLastMessage(message)
    }

    apiService.onMessage(handleMessage)

    return () => {
      apiService.offMessage(handleMessage)
      apiService.disconnectWebSocket()
      setIsConnected(false)
    }
  }, [])

  return { isConnected, lastMessage }
}

export function useSessions() {
  const [sessions, setSessions] = useState<ChatSession[]>([])
  const [currentSession, setCurrentSession] = useState<ChatSession | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  const fetchSessions = useCallback(async () => {
    setIsLoading(true)
    try {
      const data = await apiService.getSessions()
      setSessions(data)
    } finally {
      setIsLoading(false)
    }
  }, [])

  const createSession = useCallback(async () => {
    setIsLoading(true)
    try {
      const session = await apiService.createSession()
      setCurrentSession(session)
      setSessions((prev) => [session, ...prev])
      return session
    } finally {
      setIsLoading(false)
    }
  }, [])

  const selectSession = useCallback((session: ChatSession) => {
    setCurrentSession(session)
  }, [])

  useEffect(() => {
    fetchSessions()
  }, [fetchSessions])

  return {
    sessions,
    currentSession,
    isLoading,
    createSession,
    selectSession,
    refreshSessions: fetchSessions,
  }
}

export function useA2A() {
  const [agents, setAgents] = useState<Agent[]>([])
  const [isLoading, setIsLoading] = useState(false)

  const discoverAgents = useCallback(async () => {
    setIsLoading(true)
    try {
      const discovered = await a2aService.discoverAgents()
      setAgents(discovered)
      return discovered
    } finally {
      setIsLoading(false)
    }
  }, [])

  const submitTask = useCallback(
    async (
      agentId: string,
      skillId: string,
      parameters: Record<string, unknown>
    ) => {
      return a2aService.submitTask(agentId, skillId, parameters)
    },
    []
  )

  const pollForResult = useCallback(
    async (
      taskId: string,
      onUpdate: (task: { id: string; status: string; progress: number }) => void
    ) => {
      return a2aService.pollForResult(taskId, (a2aTask) => {
        onUpdate({
          id: a2aTask.id,
          status: a2aTask.status,
          progress: a2aTask.status === 'completed' ? 100 : 0,
        })
      })
    },
    []
  )

  useEffect(() => {
    // Get cached agents
    const cached = a2aService.getCachedAgents()
    if (cached.length > 0) {
      setAgents(cached)
    }

    // Discover agents
    discoverAgents()

    // Connect WebSocket
    a2aService.connectWebSocket()

    return () => {
      a2aService.disconnectWebSocket()
    }
  }, [discoverAgents])

  return {
    agents,
    isLoading,
    discoverAgents,
    submitTask,
    pollForResult,
  }
}
