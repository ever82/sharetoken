import * as React from 'react'
import { Header } from '@/components/layout/header'
import { Sidebar } from '@/components/layout/sidebar'
import { ChatContainer } from '@/components/chat/chat-container'
import { ServiceList } from '@/components/services/service-list'
import { TaskList } from '@/components/tasks/task-list'
import { ResultViewer } from '@/components/results/result-viewer'
import { ScrollArea } from '@/components/ui/scroll-area'
import { useChat, useTasks, useSessions } from '@/hooks'
import type { Task, TaskResult } from '@/types'

function App() {
  const [selectedResult, setSelectedResult] = React.useState<TaskResult | null>(null)
  const [isResultOpen, setIsResultOpen] = React.useState(false)

  const { messages, isLoading, sendMessage, clearChat } = useChat()
  const { tasks, invokeService, refreshTasks } = useTasks()
  const {
    sessions,
    currentSession,
    createSession,
    selectSession,
  } = useSessions()

  const handleViewResult = React.useCallback((taskId: string) => {
    const task = tasks.find((t) => t.id === taskId)
    if (task?.result) {
      setSelectedResult(task.result)
      setIsResultOpen(true)
    }
  }, [tasks])

  const handleInvokeService = React.useCallback(async (serviceId: string) => {
    await invokeService(serviceId, {})
    refreshTasks()
  }, [invokeService, refreshTasks])

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
              onSendMessage={sendMessage}
            />
          </div>

          <aside className="hidden w-80 border-l bg-background lg:block">
            <ScrollArea className="h-full">
              <div className="space-y-4 p-4">
                <ServiceList
                  services={[]}
                  onInvokeService={handleInvokeService}
                  isLoading={isLoading}
                />

                <TaskList
                  tasks={tasks}
                  onViewResult={handleViewResult}
                />
              </div>
            </ScrollArea>
          </aside>
        </main>
      </div>

      <ResultViewer
        result={selectedResult}
        isOpen={isResultOpen}
        onClose={() => setIsResultOpen(false)}
        onDownload={() => {
          // Trigger download
          if (selectedResult?.downloadUrl) {
            window.open(selectedResult.downloadUrl, '_blank')
          }
        }}
      />
    </div>
  )
}

export default App
