import * as React from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { cn, formatRelativeTime } from '@/utils'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import type { Message } from '@/types'

interface ChatMessageProps {
  message: Message
}

export function ChatMessage({ message }: ChatMessageProps) {
  const isUser = message.role === 'user'
  const isAssistant = message.role === 'assistant'

  return (
    <div
      className={cn(
        'flex w-full gap-3 p-4',
        isUser ? 'flex-row-reverse' : 'flex-row',
        isAssistant && 'bg-muted/50'
      )}
    >
      <Avatar className="h-8 w-8">
        <AvatarFallback className={cn(
          isUser ? 'bg-primary text-primary-foreground' : 'bg-secondary'
        )}>
          {isUser ? 'U' : 'AI'}
        </AvatarFallback>
      </Avatar>

      <div className={cn('flex max-w-[80%] flex-col gap-2', isUser && 'items-end')}>
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium">
            {isUser ? 'You' : 'GenieBot'}
          </span>
          <span className="text-xs text-muted-foreground">
            {formatRelativeTime(message.timestamp)}
          </span>
        </div>

        <div
          className={cn(
            'rounded-lg px-4 py-2',
            isUser
              ? 'bg-primary text-primary-foreground'
              : 'bg-muted'
          )}
        >
          {isAssistant ? (
            <div className="prose prose-sm dark:prose-invert max-w-none">
              <ReactMarkdown remarkPlugins={[remarkGfm]}>
                {message.content}
              </ReactMarkdown>
            </div>
          ) : (
            <p className="text-sm">{message.content}</p>
          )}
        </div>

        {message.intent && message.intent.confidence > 0.5 && (
          <div className="flex items-center gap-2">
            <Badge variant="info" className="text-xs">
              Intent: {message.intent.type} ({Math.round(message.intent.confidence * 100)}%)
            </Badge>
          </div>
        )}

        {message.services && message.services.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {message.services.map((service) => (
              <Badge
                key={service.id}
                variant="outline"
                className="cursor-pointer hover:bg-accent"
              >
                {service.name}
              </Badge>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
