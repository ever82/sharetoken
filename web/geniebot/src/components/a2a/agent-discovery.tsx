import * as React from 'react'
import { Bot, Brain, Workflow, Search, RefreshCw, CheckCircle } from 'lucide-react'
import { cn, truncate } from '@/utils'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'
import { a2aService } from '@/services/a2a'
import type { Agent, Skill } from '@/types'

const serviceIcons = {
  llm: Brain,
  agent: Bot,
  workflow: Workflow,
}

const serviceColors = {
  llm: 'text-blue-500 bg-blue-500/10',
  agent: 'text-green-500 bg-green-500/10',
  workflow: 'text-purple-500 bg-purple-500/10',
}

export function A2AAgentDiscovery() {
  const [agents, setAgents] = React.useState<Agent[]>([])
  const [searchQuery, setSearchQuery] = React.useState('')
  const [isLoading, setIsLoading] = React.useState(false)
  const [selectedAgent, setSelectedAgent] = React.useState<Agent | null>(null)
  const [error, setError] = React.useState<string | null>(null)

  const loadAgents = async () => {
    setIsLoading(true)
    setError(null)
    try {
      const discovered = await a2aService.discoverAgents()
      setAgents(discovered)
    } catch (err) {
      setError('Failed to discover agents')
    } finally {
      setIsLoading(false)
    }
  }

  React.useEffect(() => {
    // Load cached agents first
    const cached = a2aService.getCachedAgents()
    if (cached.length > 0) {
      setAgents(cached)
    }
    // Then refresh from network
    loadAgents()

    // Connect to A2A WebSocket
    a2aService.connectWebSocket()

    // Subscribe to A2A messages
    const unsubscribe = a2aService.onMessage((message) => {
      if (message.type === 'task_update') {
        // Handle task updates
        console.log('Task update received:', message)
      }
    })

    return () => {
      unsubscribe()
    }
  }, [])

  const handleSearch = async () => {
    if (!searchQuery.trim()) {
      loadAgents()
      return
    }

    setIsLoading(true)
    setError(null)
    try {
      const results = await a2aService.searchAgents(searchQuery)
      setAgents(results)
    } catch (err) {
      setError('Failed to search agents')
    } finally {
      setIsLoading(false)
    }
  }

  const filteredAgents = React.useMemo(() => {
    if (!searchQuery) return agents
    return agents.filter(
      (agent) =>
        agent.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        agent.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
        agent.capabilities.some((cap) =>
          cap.toLowerCase().includes(searchQuery.toLowerCase())
        )
    )
  }, [agents, searchQuery])

  return (
    <Card className="w-full">
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">A2A Agent Discovery</CardTitle>
          <Button
            variant="ghost"
            size="icon"
            onClick={loadAgents}
            disabled={isLoading}
          >
            <RefreshCw
              className={cn('h-4 w-4', isLoading && 'animate-spin')}
            />
          </Button>
        </div>
        <div className="flex gap-2 mt-2">
          <Input
            placeholder="Search agents by name or capability..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
          />
          <Button size="icon" onClick={handleSearch} disabled={isLoading}>
            <Search className="h-4 w-4" />
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        {error && (
          <div className="mb-4 p-3 bg-destructive/10 text-destructive rounded-md text-sm">
            {error}
          </div>
        )}

        <ScrollArea className="h-[300px]">
          {filteredAgents.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-muted-foreground">
              <Bot className="h-12 w-12 mb-2 opacity-50" />
              <p>No agents discovered yet</p>
              <p className="text-sm">Click refresh to discover agents</p>
            </div>
          ) : (
            <div className="space-y-3">
              {filteredAgents.map((agent) => (
                <AgentCard
                  key={agent.id}
                  agent={agent}
                  isSelected={selectedAgent?.id === agent.id}
                  onSelect={() => setSelectedAgent(agent)}
                />
              ))}
            </div>
          )}
        </ScrollArea>
      </CardContent>
    </Card>
  )
}

interface AgentCardProps {
  agent: Agent
  isSelected: boolean
  onSelect: () => void
}

function AgentCard({ agent, isSelected, onSelect }: AgentCardProps) {
  const primaryCapability = agent.capabilities[0] || 'general'
  const Icon = serviceIcons[primaryCapability as keyof typeof serviceIcons] || Bot
  const colorClass =
    serviceColors[primaryCapability as keyof typeof serviceColors] ||
    'text-gray-500 bg-gray-500/10'

  return (
    <button
      onClick={onSelect}
      className={cn(
        'w-full p-3 rounded-lg border transition-all',
        'hover:bg-muted/50 text-left',
        isSelected && 'border-primary bg-primary/5'
      )}
    >
      <div className="flex items-start gap-3">
        <div
          className={cn(
            'h-10 w-10 rounded-lg flex items-center justify-center shrink-0',
            colorClass
          )}
        >
          <Icon className="h-5 w-5" />
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between gap-2">
            <h4 className="font-medium truncate">{agent.name}</h4>
            <Badge variant="secondary" className="shrink-0">
              {agent.reputation.toFixed(1)} ★
            </Badge>
          </div>
          <p className="text-sm text-muted-foreground line-clamp-2">
            {agent.description}
          </p>
          <div className="flex items-center gap-2 mt-2">
            <Badge variant="outline" className="text-xs">
              {agent.pricePerUnit}
            </Badge>
            {agent.skills.slice(0, 2).map((skill) => (
              <span
                key={skill.id}
                className="text-xs text-muted-foreground"
              >
                {skill.name}
              </span>
            ))}
          </div>
        </div>
      </div>
    </button>
  )
}

interface A2ATaskSubmissionProps {
  agent: Agent
  onSubmit: (task: { agentId: string; parameters: Record<string, unknown> }) => void
  onCancel: () => void
}

export function A2ATaskSubmission({
  agent,
  onSubmit,
  onCancel,
}: A2ATaskSubmissionProps) {
  const [selectedSkill, setSelectedSkill] = React.useState<Skill | null>(
    agent.skills[0] || null
  )
  const [parameters, setParameters] = React.useState<Record<string, string>>({})
  const [isSubmitting, setIsSubmitting] = React.useState(false)

  const handleSubmit = async () => {
    setIsSubmitting(true)
    try {
      // Convert parameters to correct types
      const typedParameters: Record<string, unknown> = {}
      if (selectedSkill) {
        selectedSkill.parameters.forEach((param) => {
          const value = parameters[param.name]
          if (param.type === 'number') {
            typedParameters[param.name] = parseFloat(value) || 0
          } else if (param.type === 'boolean') {
            typedParameters[param.name] = value === 'true'
          } else {
            typedParameters[param.name] = value
          }
        })
      }

      onSubmit({ agentId: agent.id, parameters: typedParameters })
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Submit Task</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <label className="text-sm font-medium">Agent</label>
          <div className="p-3 bg-muted rounded-lg">
            <p className="font-medium">{agent.name}</p>
            <p className="text-sm text-muted-foreground">
              {agent.description}
            </p>
          </div>
        </div>

        {agent.skills.length > 0 && (
          <div className="space-y-2">
            <label className="text-sm font-medium">Skill</label>
            <select
              className="w-full p-2 border rounded-md bg-background"
              value={selectedSkill?.id || ''}
              onChange={(e) => {
                const skill = agent.skills.find((s) => s.id === e.target.value)
                setSelectedSkill(skill || null)
                setParameters({})
              }}
            >
              {agent.skills.map((skill) => (
                <option key={skill.id} value={skill.id}>
                  {skill.name}
                </option>
              ))}
            </select>
          </div>
        )}

        {selectedSkill?.parameters.map((param) => (
          <div key={param.name} className="space-y-2">
            <label className="text-sm font-medium">
              {param.name}
              {param.required && <span className="text-destructive">*</span>}
            </label>
            {param.type === 'boolean' ? (
              <select
                className="w-full p-2 border rounded-md bg-background"
                value={parameters[param.name] || 'false'}
                onChange={(e) =>
                  setParameters((prev) => ({
                    ...prev,
                    [param.name]: e.target.value,
                  }))
                }
              >
                <option value="true">Yes</option>
                <option value="false">No</option>
              </select>
            ) : (
              <Input
                type={param.type === 'number' ? 'number' : 'text'}
                placeholder={param.description}
                value={parameters[param.name] || ''}
                onChange={(e) =>
                  setParameters((prev) => ({
                    ...prev,
                    [param.name]: e.target.value,
                  }))
                }
                required={param.required}
              />
            )}
            <p className="text-xs text-muted-foreground">
              {param.description}
            </p>
          </div>
        ))}

        <div className="flex gap-2 pt-4">
          <Button
            className="flex-1"
            onClick={handleSubmit}
            disabled={isSubmitting}
          >
            {isSubmitting ? 'Submitting...' : 'Submit Task'}
          </Button>
          <Button variant="outline" onClick={onCancel}>
            Cancel
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
