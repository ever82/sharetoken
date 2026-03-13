import * as React from 'react'
import { Header } from '@/components/layout/header'
import { Sidebar } from '@/components/layout/sidebar'
import { ChatContainer } from '@/components/chat/chat-container'
import { ServiceList } from '@/components/services/service-list'
import { TaskList } from '@/components/tasks/task-list'
import { ResultViewer } from '@/components/results/result-viewer'
import { A2AAgentDiscovery } from '@/components/a2a/agent-discovery'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useChat, useTasks, useSessions } from '@/hooks'
import { useWallet } from '@/hooks/useWallet'
import { useAuth } from '@/hooks/useAuth'
import { a2aService } from '@/services/a2a'
import type { Task, TaskResult, Agent } from '@/types'

function App() {
  const [selectedResult, setSelectedResult] = React.useState<TaskResult | null>(null)
  const [isResultOpen, setIsResultOpen] = React.useState(false)
  const [selectedAgent, setSelectedAgent] = React.useState<Agent | null>(null)
  const [activeTab, setActiveTab] = React.useState('services')

  const { messages, isLoading, sendMessage, clearChat } = useChat()
  const { tasks, invokeService, refreshTasks } = useTasks()
  const {
    sessions,
    currentSession,
    createSession,
    selectSession,
  } = useSessions()
  const { state: walletState } = useWallet()
  const { state: authState } = useAuth()

  const handleViewResult = React.useCallback((taskId: string) => {
    const task = tasks.find((t) => t.id === taskId)
    if (task?.result) {
      setSelectedResult(task.result)
      setIsResultOpen(true)
    }
  }, [tasks])

  const handleInvokeService = React.useCallback(async (serviceId: string) => {
    try {
      await invokeService(serviceId, {})
      refreshTasks()
    } catch (error) {
      console.error('Failed to invoke service:', error)
    }
  }, [invokeService, refreshTasks])

  const handleSendMessage = React.useCallback(async (content: string) => {
    // If authenticated, sign the message
    if (authState.isAuthenticated) {
      try {
        await sendMessage(content, true) // signed message
      } catch (error) {
        console.error('Failed to send signed message:', error)
        // Fallback to unsigned message
        await sendMessage(content)
      }
    } else {
      await sendMessage(content)
    }
  }, [sendMessage, authState.isAuthenticated])

  // Initialize A2A service on mount
  React.useEffect(() => {
    a2aService.discoverAgents()

    // Subscribe to A2A messages
    const unsubscribe = a2aService.onMessage((message) => {
      if (message.type === 'task_update') {
        // Refresh tasks when we get updates
        refreshTasks()
      }
    })

    return () => {
      unsubscribe()
    }
  }, [refreshTasks])

  return (
    <div className="flex h-screen flex-col bg-background">
      <Header />

      <div className="flex flex-1 overflow-hidden">
        <Sidebar
          sessions={sessions}
          currentSession={currentSession}
          onNewSession={createSession}
          onSelectSession={selectSession}
        />

        <main className="flex flex-1 overflow-hidden">
          <div className="flex flex-1 flex-col">
            <ChatContainer
              messages={messages}
              isLoading={isLoading}
              onSendMessage={handleSendMessage}
            />
          </div>

          <aside className="hidden w-96 border-l bg-background lg:block">
            <Tabs value={activeTab} onValueChange={setActiveTab} className="h-full flex flex-col">
              <TabsList className="w-full justify-start rounded-none border-b px-4 py-2">
                <TabsTrigger value="services">Services</TabsTrigger>
                <TabsTrigger value="tasks">Tasks</TabsTrigger>
                <TabsTrigger value="agents">Agents</TabsTrigger>
              </TabsList>

              <div className="flex-1 overflow-hidden">
                <ScrollArea className="h-full">
                  <TabsContent value="services" className="m-0">
                    <div className="p-4 space-y-4">
                      <ServiceList
                        services={[]}
                        onInvokeService={handleInvokeService}
                        isLoading={isLoading}
                      />
                    </div>
                  </TabsContent>

                  <TabsContent value="tasks" className="m-0">
                    <div className="p-4">
                      <TaskList
                        tasks={tasks}
                        onViewResult={handleViewResult}
                      />
                    </div>
                  </TabsContent>

                  <TabsContent value="agents" className="m-0">
                    <div className="p-4">
                      <A2AAgentDiscovery />
                    </div>
                  </TabsContent>
                </ScrollArea>
              </div>
            </Tabs>
          </aside>
        </main>
      </div>

      <ResultViewer
        result={selectedResult}
        isOpen={isResultOpen}
        onClose={() => setIsResultOpen(false)}
        onDownload={() => {
          if (selectedResult?.downloadUrl) {
            window.open(selectedResult.downloadUrl, '_blank')
          }
        }}
      />
    </div>
  )
}

export default App
