import { useState, useCallback, useEffect, useRef } from 'react'
import type { Message, Task, ChatSession, ServiceRecommendation, Intent } from '@/types'
import { apiService } from '@/services/api'

export function useChat() {
  const [messages, setMessages] = useState<Message[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [currentIntent, setCurrentIntent] = useState<Intent | null>(null)
  const [recommendations, setRecommendations] = useState<ServiceRecommendation[]>([])

  const sendMessage = useCallback(async (content: string) => {
    setIsLoading(true)

    // Add user message
    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content,
      timestamp: Date.now(),
    }
    setMessages(prev => [...prev, userMessage])

    try {
      // Detect intent
      const intent = await apiService.detectIntent(content)
      setCurrentIntent(intent)

      // Get service recommendations if applicable
      if (intent.type !== 'unknown' && intent.confidence > 0.5) {
        const recs = await apiService.getServiceRecommendations(content)
        setRecommendations(recs)
      }

      // Send message to API
      const response = await apiService.sendMessage(content)

      // Add assistant message
      const assistantMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: response.content,
        timestamp: Date.now(),
        intent,
        services: recommendations,
      }
      setMessages(prev => [...prev, assistantMessage])
    } catch (error) {
      // Add error message
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: 'Sorry, I encountered an error. Please try again.',
        timestamp: Date.now(),
      }
      setMessages(prev => [...prev, errorMessage])
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
      setTasks(data)
    } finally {
      setIsLoading(false)
    }
  }, [])

  const invokeService = useCallback(async (serviceId: string, params: Record<string, unknown>) => {
    setIsLoading(true)
    try {
      const task = await apiService.invokeService(serviceId, params)
      setTasks(prev => [task, ...prev])
      return task
    } finally {
      setIsLoading(false)
    }
  }, [])

  const updateTaskProgress = useCallback((taskId: string, progress: number) => {
    setTasks(prev =>
      prev.map(task =>
        task.id === taskId ? { ...task, progress } : task
      )
    )
  }, [])

  const updateTaskStatus = useCallback((taskId: string, status: Task['status']) => {
    setTasks(prev =>
      prev.map(task =>
        task.id === taskId ? { ...task, status } : task
      )
    )
  }, [])

  useEffect(() => {
    fetchTasks()
  }, [fetchTasks])

  return {
    tasks,
    isLoading,
    invokeService,
    updateTaskProgress,
    updateTaskStatus,
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
      setSessions(prev => [session, ...prev])
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
