import * as React from 'react'
import { MessageSquare, Plus, History, Settings } from 'lucide-react'
import { cn } from '@/utils'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import type { ChatSession } from '@/types'

interface SidebarProps {
  sessions: ChatSession[]
  currentSession: ChatSession | null
  onNewSession: () => void
  onSelectSession: (session: ChatSession) => void
  className?: string
}

export function Sidebar({
  sessions,
  currentSession,
  onNewSession,
  onSelectSession,
  className,
}: SidebarProps) {
  return (
    <div className={cn('flex h-full w-64 flex-col border-r bg-muted/50', className)}>
      <div className="p-4">
        <Button onClick={onNewSession} className="w-full">
          <Plus className="mr-2 h-4 w-4" />
          New Chat
        </Button>
      </div>

      <ScrollArea className="flex-1">
        <div className="space-y-1 px-3">
          <h3 className="mb-2 px-3 text-xs font-medium text-muted-foreground">
            Recent Chats
          </h3>
          {sessions.length === 0 ? (
            <p className="px-3 text-sm text-muted-foreground">No chats yet</p>
          ) : (
            sessions.map((session) => (
              <button
                key={session.id}
                onClick={() => onSelectSession(session)}
                className={cn(
                  'flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors',
                  currentSession?.id === session.id
                    ? 'bg-primary text-primary-foreground'
                    : 'hover:bg-muted'
                )}
              >
                <MessageSquare className="h-4 w-4" />
                <span className="flex-1 truncate text-left">{session.title}</span>
              </button>
            ))
          )}
        </div>
      </ScrollArea>

      <div className="border-t p-4">
        <button className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm text-muted-foreground hover:bg-muted">
          <History className="h-4 w-4" />
          <span>History</span>
        </button>
        <button className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm text-muted-foreground hover:bg-muted">
          <Settings className="h-4 w-4" />
          <span>Settings</span>
        </button>
      </div>
    </div>
  )
}
