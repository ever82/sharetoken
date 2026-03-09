import * as React from 'react'
import { Download, X, FileText, Code, FileJson } from 'lucide-react'
import { cn } from '@/utils'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { ScrollArea } from '@/components/ui/scroll-area'
import type { TaskResult } from '@/types'

interface ResultViewerProps {
  result: TaskResult | null
  isOpen: boolean
  onClose: () => void
  onDownload: () => void
}

const formatIcons = {
  text: FileText,
  markdown: FileText,
  code: Code,
  json: FileJson,
}

export function ResultViewer({
  result,
  isOpen,
  onClose,
  onDownload,
}: ResultViewerProps) {
  if (!result) return null

  const FormatIcon = formatIcons[result.format]

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-4xl max-h-[80vh]">
        <DialogHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <FormatIcon className="h-5 w-5" />
              <DialogTitle>Task Result</DialogTitle>
            </div>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" onClick={onDownload}>
                <Download className="mr-2 h-4 w-4" />
                Download
              </Button>
              <Button variant="ghost" size="icon" onClick={onClose}>
                <X className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </DialogHeader>

        <ScrollArea className="h-[500px] rounded-md border bg-muted p-4">
          {result.format === 'json' ? (
            <pre className="text-sm">
              <code>{JSON.stringify(JSON.parse(result.content), null, 2)}</code>
            </pre>
          ) : result.format === 'code' ? (
            <pre className="text-sm">
              <code>{result.content}</code>
            </pre>
          ) : (
            <div className="prose prose-sm dark:prose-invert max-w-none">
              {result.content}
            </div>
          )}
        </ScrollArea>

        {result.attachments && result.attachments.length > 0 && (
          <div className="border-t pt-4">
            <h4 className="text-sm font-medium mb-2">Attachments</h4>
            <div className="flex flex-wrap gap-2">
              {result.attachments.map((attachment) => (
                <Button
                  key={attachment.name}
                  variant="outline"
                  size="sm"
                  onClick={() => window.open(attachment.url, '_blank')}
                >
                  <Download className="mr-2 h-4 w-4" />
                  {attachment.name} ({formatFileSize(attachment.size)})
                </Button>
              ))}
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  )
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
