import { useState, useCallback, useRef, useEffect, useMemo } from 'react'
import { useParams, Link } from 'react-router-dom'
import {
  ReactFlow,
  Controls,
  Background,
  BackgroundVariant,
  addEdge,
  useNodesState,
  useEdgesState,
  type Connection,
  type Edge,
  type Node,
  MarkerType,
  Panel,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { ArrowLeft, Plus, Save, Play, Loader2, CheckCircle2, XCircle, History } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/components/ui/sheet'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { flowService } from '@/services/flow.service'
import { StepNode, type StepNodeData } from '@/components/features/flow/step-node'
import { StepPanel } from '@/components/features/flow/step-panel'
import { env } from '@/config/env'
import { getAuthTokens } from '@/store/auth-store'
import type { FlowStep, FlowEdge as FlowEdgeType, FlowStepEvent, FlowRun } from '@/types/kest-api'
import { toast } from 'sonner'

const nodeTypes = { step: StepNode }

let stepCounter = 0

function flowToNodes(steps: FlowStep[]): Node[] {
  return steps.map((step, i) => ({
    id: `step-${step.id}`,
    type: 'step',
    position: { x: step.position_x || 250, y: step.position_y || i * 180 },
    data: {
      stepId: step.id,
      name: step.name,
      method: step.method,
      url: step.url,
      headers: step.headers || '',
      body: step.body || '',
      captures: step.captures || '',
      asserts: step.asserts || '',
      sort_order: step.sort_order,
    } satisfies StepNodeData,
  }))
}

function flowToEdges(edges: FlowEdgeType[]): Edge[] {
  return edges.map(edge => ({
    id: `edge-${edge.id}`,
    source: `step-${edge.source_step_id}`,
    target: `step-${edge.target_step_id}`,
    markerEnd: { type: MarkerType.ArrowClosed, width: 16, height: 16 },
    style: { strokeWidth: 2 },
    data: { variable_mapping: edge.variable_mapping || '' },
  }))
}

export function FlowEditorPage() {
  const { id, fid } = useParams<{ id: string; fid: string }>()
  const projectId = parseInt(id || '0')
  const flowId = parseInt(fid || '0')
  const queryClient = useQueryClient()

  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([])
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([])
  const [selectedNode, setSelectedNode] = useState<Node | null>(null)
  const [isRunning, setIsRunning] = useState(false)
  const [runStatus, setRunStatus] = useState<Record<number, { status: string; durationMs?: number }>>({})
  const [showHistory, setShowHistory] = useState(false)
  const [dirty, setDirty] = useState(false)
  const initialized = useRef(false)

  const { data: flow, isLoading } = useQuery({
    queryKey: ['flow', projectId, flowId],
    queryFn: () => flowService.get(projectId, flowId),
  })

  const { data: runsData } = useQuery({
    queryKey: ['flow-runs', projectId, flowId],
    queryFn: () => flowService.listRuns(projectId, flowId),
  })

  // Initialize nodes/edges from flow data
  useEffect(() => {
    if (flow && !initialized.current) {
      setNodes(flowToNodes(flow.steps || []))
      setEdges(flowToEdges(flow.edges || []))
      stepCounter = Math.max(0, ...(flow.steps || []).map(s => s.id)) + 1
      initialized.current = true
    }
  }, [flow])

  const nodeTypes_ = useMemo(() => nodeTypes, [])

  const onConnect = useCallback((params: Connection) => {
    setEdges(eds => addEdge({
      ...params,
      markerEnd: { type: MarkerType.ArrowClosed, width: 16, height: 16 },
      style: { strokeWidth: 2 },
    }, eds))
    setDirty(true)
  }, [])

  const onNodeClick = useCallback((_: React.MouseEvent, node: Node) => {
    setSelectedNode(node)
  }, [])

  const onPaneClick = useCallback(() => {
    setSelectedNode(null)
  }, [])

  // Add new step
  const addStep = useCallback(() => {
    const newId = `step-new-${++stepCounter}`
    const yOffset = nodes.length * 180
    const newNode: Node = {
      id: newId,
      type: 'step',
      position: { x: 250, y: yOffset + 50 },
      data: {
        stepId: 0,
        name: `Step ${nodes.length + 1}`,
        method: 'GET',
        url: '/v1/',
        headers: '{"Content-Type": "application/json"}',
        body: '',
        captures: '',
        asserts: 'status == 200',
        sort_order: nodes.length,
      } satisfies StepNodeData,
    }
    setNodes(nds => [...nds, newNode])
    setSelectedNode(newNode)
    setDirty(true)
  }, [nodes])

  // Update step data from panel
  const updateNodeData = useCallback((nodeId: string, updates: Partial<StepNodeData>) => {
    setNodes(nds => nds.map(n =>
      n.id === nodeId ? { ...n, data: { ...n.data, ...updates } } : n
    ))
    setSelectedNode(prev => prev && prev.id === nodeId ? { ...prev, data: { ...prev.data, ...updates } } : prev)
    setDirty(true)
  }, [])

  // Delete step
  const deleteNode = useCallback((nodeId: string) => {
    setNodes(nds => nds.filter(n => n.id !== nodeId))
    setEdges(eds => eds.filter(e => e.source !== nodeId && e.target !== nodeId))
    setSelectedNode(null)
    setDirty(true)
  }, [])

  // Save flow
  const saveMutation = useMutation({
    mutationFn: async () => {
      const stepMap = new Map<string, number>()
      const steps = nodes.map((node, i) => {
        const d = node.data as StepNodeData
        stepMap.set(node.id, i)
        return {
          name: String(d.name),
          sort_order: i,
          method: String(d.method),
          url: String(d.url),
          headers: String(d.headers || ''),
          body: String(d.body || ''),
          captures: String(d.captures || ''),
          asserts: String(d.asserts || ''),
          position_x: Math.round(node.position.x),
          position_y: Math.round(node.position.y),
        }
      })

      // We need to save first to get step IDs, then create edges
      // Use the save endpoint which replaces all steps and edges
      const edgeData = edges.map(e => {
        const sourceIdx = stepMap.get(e.source) ?? 0
        const targetIdx = stepMap.get(e.target) ?? 0
        return {
          source_step_id: sourceIdx,
          target_step_id: targetIdx,
          variable_mapping: (e.data as any)?.variable_mapping || '',
        }
      })

      return flowService.save(projectId, flowId, {
        name: flow?.name,
        description: flow?.description,
        steps,
        edges: edgeData,
      })
    },
    onSuccess: (savedFlow) => {
      // Re-initialize with saved data (now has real IDs)
      initialized.current = false
      queryClient.invalidateQueries({ queryKey: ['flow', projectId, flowId] })
      setDirty(false)
      toast.success('Flow saved')
    },
    onError: () => toast.error('Failed to save flow'),
  })

  // Run flow
  const runFlow = useCallback(async () => {
    // Save first if dirty
    if (dirty) {
      await saveMutation.mutateAsync()
      // Wait for query to refetch
      await queryClient.refetchQueries({ queryKey: ['flow', projectId, flowId] })
    }

    setIsRunning(true)
    setRunStatus({})

    try {
      // Create run
      const run = await flowService.run(projectId, flowId)

      // Connect to SSE
      const { accessToken } = getAuthTokens()
      const baseUrl = env.VITE_API_URL || ''
      const sseUrl = `${baseUrl}/v1/projects/${projectId}/flows/${flowId}/runs/${run.id}/events`

      // EventSource doesn't support auth headers, so we use fetch with ReadableStream
      const response = await fetch(sseUrl, {
        headers: { Authorization: `Bearer ${accessToken}` },
      })

      if (!response.body) throw new Error('No response body')

      const reader = response.body.getReader()
      const decoder = new TextDecoder()
      let buffer = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() || ''

        let eventType = ''
        for (const line of lines) {
          if (line.startsWith('event: ')) {
            eventType = line.slice(7).trim()
          } else if (line.startsWith('data: ') && eventType) {
            try {
              const event: FlowStepEvent = JSON.parse(line.slice(6))
              // Update node status
              setRunStatus(prev => ({
                ...prev,
                [event.step_id]: {
                  status: event.status,
                  durationMs: event.data?.duration_ms,
                },
              }))

              // Update node visual
              setNodes(nds => nds.map(n => {
                const d = n.data as StepNodeData
                if (d.stepId === event.step_id) {
                  return {
                    ...n,
                    data: {
                      ...n.data,
                      runStatus: event.status as StepNodeData['runStatus'],
                      durationMs: event.data?.duration_ms,
                    },
                  }
                }
                return n
              }))
            } catch {}
            eventType = ''
          }
        }
      }

      queryClient.invalidateQueries({ queryKey: ['flow-runs', projectId, flowId] })
      toast.success('Flow execution completed')
    } catch (err: any) {
      toast.error(err?.message || 'Flow execution failed')
    } finally {
      setIsRunning(false)
    }
  }, [dirty, projectId, flowId, flow])

  // Clear run status
  const clearRunStatus = useCallback(() => {
    setRunStatus({})
    setNodes(nds => nds.map(n => ({
      ...n,
      data: { ...n.data, runStatus: undefined, durationMs: undefined },
    })))
  }, [])

  if (isLoading) {
    return (
      <div className="h-screen flex items-center justify-center">
        <Skeleton className="h-[600px] w-full max-w-4xl" />
      </div>
    )
  }

  const runs = runsData?.items || []

  return (
    <div className="h-[calc(100vh-64px)] flex flex-col">
      {/* Toolbar */}
      <div className="flex items-center justify-between px-4 py-2 border-b border-border bg-background">
        <div className="flex items-center gap-3">
          <Link to={`/projects/${projectId}`}>
            <Button variant="ghost" size="sm">
              <ArrowLeft className="h-4 w-4 mr-1" />
              Back
            </Button>
          </Link>
          <div>
            <h1 className="text-sm font-semibold">{flow?.name}</h1>
            <p className="text-xs text-muted-foreground">{flow?.description || 'No description'}</p>
          </div>
          {dirty && <Badge variant="outline" className="text-orange-500 border-orange-300">Unsaved</Badge>}
        </div>

        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={() => setShowHistory(true)}>
            <History className="h-4 w-4 mr-1" />
            Runs ({runs.length})
          </Button>
          <Button variant="outline" size="sm" onClick={addStep}>
            <Plus className="h-4 w-4 mr-1" />
            Add Step
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => saveMutation.mutate()}
            disabled={saveMutation.isPending || !dirty}
          >
            <Save className="h-4 w-4 mr-1" />
            {saveMutation.isPending ? 'Saving...' : 'Save'}
          </Button>
          <Button
            size="sm"
            onClick={runFlow}
            disabled={isRunning || nodes.length === 0}
            className="bg-green-600 hover:bg-green-700"
          >
            {isRunning ? (
              <><Loader2 className="h-4 w-4 mr-1 animate-spin" />Running...</>
            ) : (
              <><Play className="h-4 w-4 mr-1" />Run</>
            )}
          </Button>
        </div>
      </div>

      {/* Canvas + Panel */}
      <div className="flex-1 flex">
        <div className="flex-1">
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={changes => {
              onNodesChange(changes)
              if (!isRunning) {
                const hasPositionChange = changes.some(c => c.type === 'position' && c.dragging)
                const hasRemove = changes.some(c => c.type === 'remove')
                if (hasPositionChange || hasRemove) setDirty(true)
              }
            }}
            onEdgesChange={changes => {
              onEdgesChange(changes)
              if (!isRunning) {
                const hasRemove = changes.some(c => c.type === 'remove')
                if (hasRemove) setDirty(true)
              }
            }}
            onConnect={onConnect}
            onNodeClick={onNodeClick}
            onPaneClick={onPaneClick}
            nodeTypes={nodeTypes_}
            fitView
            fitViewOptions={{ padding: 0.3 }}
            deleteKeyCode={['Backspace', 'Delete']}
            className="bg-muted/30"
          >
            <Controls />
            <Background variant={BackgroundVariant.Dots} gap={20} size={1} />
            {Object.keys(runStatus).length > 0 && (
              <Panel position="top-left">
                <Button variant="outline" size="sm" onClick={clearRunStatus}>
                  Clear Results
                </Button>
              </Panel>
            )}
          </ReactFlow>
        </div>

        {selectedNode && (
          <StepPanel
            data={selectedNode.data as StepNodeData}
            onChange={(updates) => updateNodeData(selectedNode.id, updates)}
            onClose={() => setSelectedNode(null)}
            onDelete={() => deleteNode(selectedNode.id)}
          />
        )}
      </div>

      {/* Run History Sheet */}
      <Sheet open={showHistory} onOpenChange={setShowHistory}>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>Run History</SheetTitle>
          </SheetHeader>
          <div className="mt-4 space-y-3">
            {runs.length === 0 ? (
              <p className="text-sm text-muted-foreground">No runs yet</p>
            ) : (
              runs.map((run: FlowRun) => (
                <div key={run.id} className="flex items-center justify-between p-3 rounded-lg border">
                  <div className="flex items-center gap-2">
                    {run.status === 'passed' ? (
                      <CheckCircle2 className="h-4 w-4 text-green-500" />
                    ) : run.status === 'failed' ? (
                      <XCircle className="h-4 w-4 text-red-500" />
                    ) : (
                      <Loader2 className="h-4 w-4 text-blue-500 animate-spin" />
                    )}
                    <div>
                      <p className="text-sm font-medium">Run #{run.id}</p>
                      <p className="text-xs text-muted-foreground">
                        {new Date(run.created_at).toLocaleString()}
                      </p>
                    </div>
                  </div>
                  <Badge variant={run.status === 'passed' ? 'default' : run.status === 'failed' ? 'destructive' : 'outline'}>
                    {run.status}
                  </Badge>
                </div>
              ))
            )}
          </div>
        </SheetContent>
      </Sheet>
    </div>
  )
}
