export interface Message {
  id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  timestamp: number
  intent?: Intent
  services?: ServiceRecommendation[]
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
