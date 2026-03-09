import * as React from 'react'
import { ChatMessage } from './chat-message'
import { ChatInput } from './chat-input'
import { ScrollArea } from '@/components/ui/scroll-area'
import type { Message } from '@/types'

interface ChatContainerProps {
  messages: Message[]
  isLoading: boolean
  onSendMessage: (message: string) => void
}

export function ChatContainer({
  messages,
  isLoading,
  onSendMessage,
}: ChatContainerProps) {
  const scrollRef = React.useRef<HTMLDivElement>(null)
  const messagesEndRef = React.useRef<HTMLDivElement>(null)

  // Scroll to bottom when messages change
  React.useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  return (
    <div className="flex h-full flex-col">
      <ScrollArea className="flex-1 p-4">
        <div className="flex flex-col gap-2">
          {messages.length === 0 ? (
            <div className="flex h-full items-center justify-center py-20">
              <div className="text-center">
                <h3 className="text-lg font-semibold">Welcome to GenieBot</h3>
                <p className="text-muted-foreground">
                  Start a conversation by typing a message below.
                </p>
                <div className="mt-4 text-sm text-muted-foreground">
                  <p>Try asking:</p>
                  <ul className="mt-2 space-y-1">
                    <li>"Write a code snippet for..."</li>
                    <li>"Research about blockchain"</li>
                    <li>"Create a workflow for..."</li>
                    <li>"Analyze this data"</li>
                  </ul>
                </div>
              </div>
            </div>
          ) : (
            messages.map((message) => (
              <ChatMessage key={message.id} message={message} />
            ))
          )}
          <div ref={messagesEndRef} />
        </div>
      </ScrollArea>

      <ChatInput
        onSend={onSendMessage}
        isLoading={isLoading}
        placeholder="Type a message... (Shift+Enter for new line)"
      />
    </div>
  )
}
