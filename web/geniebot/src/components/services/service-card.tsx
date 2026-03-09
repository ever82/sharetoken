import * as React from 'react'
import { Bot, Workflow, Brain, Clock, Coins, Play } from 'lucide-react'
import { cn } from '@/utils'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import type { ServiceRecommendation } from '@/types'

interface ServiceCardProps {
  service: ServiceRecommendation
  onInvoke: (serviceId: string) => void
  isLoading?: boolean
}

const serviceIcons = {
  llm: Brain,
  agent: Bot,
  workflow: Workflow,
}

const serviceColors = {
  llm: 'text-blue-500',
  agent: 'text-green-500',
  workflow: 'text-purple-500',
}

export function ServiceCard({ service, onInvoke, isLoading = false }: ServiceCardProps) {
  const Icon = serviceIcons[service.type]

  return (
    <Card className="w-full">
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-2">
            <Icon className={cn('h-5 w-5', serviceColors[service.type])} />
            <CardTitle className="text-base">{service.name}</CardTitle>
          </div>
          <Badge variant="secondary">
            {Math.round(service.confidence * 100)}% match
          </Badge>
        </div>
        <CardDescription>{service.description}</CardDescription>
      </CardHeader>
      <CardContent className="pb-3">
        <div className="flex gap-4 text-sm text-muted-foreground">
          <div className="flex items-center gap-1">
            <Coins className="h-4 w-4" />
            <span>{service.estimatedCost}</span>
          </div>
          <div className="flex items-center gap-1">
            <Clock className="h-4 w-4" />
            <span>{service.estimatedTime}</span>
          </div>
        </div>
      </CardContent>
      <CardFooter>
        <Button
          onClick={() => onInvoke(service.id)}
          disabled={isLoading}
          className="w-full"
          size="sm"
        >
          {isLoading ? (
            'Running...'
          ) : (
            <>
              <Play className="mr-2 h-4 w-4" />
              Run Service
            </>
          )}
        </Button>
      </CardFooter>
    </Card>
  )
}
