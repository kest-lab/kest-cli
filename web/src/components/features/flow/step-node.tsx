import { memo } from 'react'
import { Handle, Position, type NodeProps } from '@xyflow/react'
import { Badge } from '@/components/ui/badge'
import { CheckCircle2, XCircle, Loader2, Clock } from 'lucide-react'

const METHOD_COLORS: Record<string, string> = {
  GET: 'bg-blue-500 text-white',
  POST: 'bg-green-500 text-white',
  PUT: 'bg-yellow-500 text-white',
  PATCH: 'bg-orange-500 text-white',
  DELETE: 'bg-red-500 text-white',
}

export interface StepNodeData {
  stepId: number
  name: string
  method: string
  url: string
  headers: string
  body: string
  captures: string
  asserts: string
  sort_order?: number
  runStatus?: 'pending' | 'running' | 'passed' | 'failed'
  durationMs?: number
  selected?: boolean
  [key: string]: unknown
}

function StepNodeComponent({ data, selected }: NodeProps) {
  const nodeData = data as unknown as StepNodeData
  const statusIcon = () => {
    switch (nodeData.runStatus) {
      case 'running': return <Loader2 className="h-4 w-4 text-blue-500 animate-spin" />
      case 'passed': return <CheckCircle2 className="h-4 w-4 text-green-500" />
      case 'failed': return <XCircle className="h-4 w-4 text-red-500" />
      default: return <Clock className="h-4 w-4 text-muted-foreground" />
    }
  }

  const borderColor = () => {
    switch (nodeData.runStatus) {
      case 'running': return 'border-blue-500 shadow-blue-100'
      case 'passed': return 'border-green-500 shadow-green-100'
      case 'failed': return 'border-red-500 shadow-red-100'
      default: return selected ? 'border-primary shadow-primary/10' : 'border-border'
    }
  }

  const captureCount = nodeData.captures ? nodeData.captures.split('\n').filter((l: string) => l.trim()).length : 0
  const assertCount = nodeData.asserts ? nodeData.asserts.split('\n').filter((l: string) => l.trim()).length : 0

  return (
    <div className={`bg-card rounded-lg border-2 shadow-md min-w-[220px] max-w-[280px] transition-all ${borderColor()}`}>
      <Handle type="target" position={Position.Top} className="!bg-primary !w-3 !h-3 !border-2 !border-background" />

      <div className="px-3 py-2 border-b border-border/50 flex items-center justify-between gap-2">
        <div className="flex items-center gap-2 min-w-0">
          <Badge className={`${METHOD_COLORS[nodeData.method] || 'bg-gray-500 text-white'} text-[10px] px-1.5 py-0 font-mono shrink-0`}>
            {nodeData.method}
          </Badge>
          <span className="text-xs font-medium truncate">{nodeData.name}</span>
        </div>
        {statusIcon()}
      </div>

      <div className="px-3 py-2">
        <p className="text-[11px] font-mono text-muted-foreground truncate">{nodeData.url}</p>
      </div>

      {(captureCount > 0 || assertCount > 0 || nodeData.durationMs !== undefined) && (
        <div className="px-3 py-1.5 border-t border-border/50 flex items-center gap-2 text-[10px] text-muted-foreground">
          {captureCount > 0 && <span>📤 {captureCount} capture{captureCount > 1 ? 's' : ''}</span>}
          {assertCount > 0 && <span>✓ {assertCount} assert{assertCount > 1 ? 's' : ''}</span>}
          {nodeData.durationMs !== undefined && nodeData.durationMs > 0 && (
            <span className="ml-auto">{nodeData.durationMs}ms</span>
          )}
        </div>
      )}

      <Handle type="source" position={Position.Bottom} className="!bg-primary !w-3 !h-3 !border-2 !border-background" />
    </div>
  )
}

export const StepNode = memo(StepNodeComponent)
