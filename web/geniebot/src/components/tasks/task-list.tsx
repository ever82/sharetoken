import * as React from 'react'
import { CheckCircle2, Clock, XCircle, Loader2, PauseCircle } from 'lucide-react'
import { cn, formatRelativeTime } from '@/utils'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { ScrollArea } from '@/components/ui/scroll-area'
import type { Task } from '@/types'

interface TaskListProps {
  tasks: Task[]
  onViewResult?: (taskId: string) => void
}

const statusIcons = {
  pending: Clock,
  running: Loader2,
  completed: CheckCircle2,
  failed: XCircle,
  cancelled: PauseCircle,
}

const statusColors = {
  pending: 'text-yellow-500',
  running: 'text-blue-500',
  completed: 'text-green-500',
  failed: 'text-red-500',
  cancelled: 'text-gray-500',
}

const statusVariants = {
  pending: 'warning',
  running: 'info',
  completed: 'success',
  failed: 'destructive',
  cancelled: 'secondary',
} as const

function TaskItem({ task, onViewResult }: { task: Task; onViewResult?: (taskId: string) => void }) {
  const StatusIcon = statusIcons[task.status]
  const isRunning = task.status === 'running'

  return (
    <div className="border-b last:border-b-0 py-4">
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-start gap-3">
          <StatusIcon
            className={cn('h-5 w-5', statusColors[task.status], isRunning && 'animate-spin')}
          />
          <div>
            <h4 className="font-medium">{task.name}</h4>
            <p className="text-sm text-muted-foreground">{task.description}</p>
            <p className="text-xs text-muted-foreground mt-1">
              {formatRelativeTime(task.createdAt)}
            </p>
          </div>
        </div>
        <Badge variant={statusVariants[task.status]}>
          {task.status}
        </Badge>
      </div>

      {(isRunning || task.progress > 0) && (
        <div className="mt-3">
          <Progress value={task.progress} className="h-2" />
          <p className="text-xs text-muted-foreground mt-1 text-right">
            {task.progress}%
          </p>
        </div>
      )}

      {task.status === 'completed' && task.result && onViewResult && (
        <div className="mt-3">
          <button
            onClick={() => onViewResult(task.id)}
            className="text-sm text-primary hover:underline"
          >
            View Result
          </button>
        </div>
      )}
    </div>
  )
}

export function TaskList({ tasks, onViewResult }: TaskListProps) {
  const sortedTasks = React.useMemo(() => {
    return [...tasks].sort((a, b) => b.createdAt - a.createdAt)
  }, [tasks])

  const runningTasks = sortedTasks.filter((t) => t.status === 'running')
  const completedTasks = sortedTasks.filter((t) => t.status === 'completed')
  const otherTasks = sortedTasks.filter(
    (t) => t.status !== 'running' && t.status !== 'completed'
  )

  return (
    <Card className="h-full">
      <CardHeader>
        <CardTitle>Tasks ({tasks.length})</CardTitle>
      </CardHeader>
      <CardContent>
        <ScrollArea className="h-[400px]">
          {tasks.length === 0 ? (
            <div className="flex h-32 items-center justify-center text-muted-foreground">
              No tasks yet
            </div>
          ) : (
            <div className="space-y-4">
              {runningTasks.length > 0 && (
                <div>
                  <h4 className="text-sm font-medium text-muted-foreground mb-2">
                    Running
                  </h4>
                  {runningTasks.map((task) => (
                    <TaskItem key={task.id} task={task} onViewResult={onViewResult} />
                  ))}
                </div>
              )}

              {completedTasks.length > 0 && (
                <div>
                  <h4 className="text-sm font-medium text-muted-foreground mb-2">
                    Completed
                  </h4>
                  {completedTasks.map((task) => (
                    <TaskItem key={task.id} task={task} onViewResult={onViewResult} />
                  ))}
                </div>
              )}

              {otherTasks.length > 0 && (
                <div>
                  <h4 className="text-sm font-medium text-muted-foreground mb-2">
                    Other
                  </h4>
                  {otherTasks.map((task) => (
                    <TaskItem key={task.id} task={task} onViewResult={onViewResult} />
                  ))}
                </div>
              )}
            </div>
          )}
        </ScrollArea>
      </CardContent>
    </Card>
  )
}
