import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Plus, Trash2, Clock, CheckCircle2, XCircle, Workflow } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { flowService } from '@/services/flow.service'
import { toast } from 'sonner'
import { FlowCreateSheet } from './flow-create-sheet'

interface FlowListProps {
  projectId: number
}

export function FlowList({ projectId }: FlowListProps) {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [showCreate, setShowCreate] = useState(false)

  const { data, isLoading } = useQuery({
    queryKey: ['flows', projectId],
    queryFn: () => flowService.list(projectId),
  })

  const deleteMutation = useMutation({
    mutationFn: (flowId: number) => flowService.delete(projectId, flowId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['flows', projectId] })
      toast.success('Flow deleted')
    },
  })

  const flows = data?.items || []

  const statusIcon = (status?: string) => {
    switch (status) {
      case 'passed': return <CheckCircle2 className="h-4 w-4 text-green-500" />
      case 'failed': return <XCircle className="h-4 w-4 text-red-500" />
      default: return <Clock className="h-4 w-4 text-muted-foreground" />
    }
  }

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[1, 2, 3].map(i => <Skeleton key={i} className="h-24 w-full" />)}
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold">Test Flows</h2>
          <p className="text-muted-foreground text-sm">Visual API test scenarios with variable passing</p>
        </div>
        <Button onClick={() => setShowCreate(true)}>
          <Plus className="h-4 w-4 mr-2" />
          New Flow
        </Button>
      </div>

      {flows.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-16">
            <Workflow className="h-12 w-12 text-muted-foreground mb-4" />
            <h3 className="text-lg font-semibold mb-2">No test flows yet</h3>
            <p className="text-muted-foreground text-sm mb-4">Create your first flow to start testing API scenarios</p>
            <Button onClick={() => setShowCreate(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Create Flow
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4">
          {flows.map(flow => (
            <Card
              key={flow.id}
              className="cursor-pointer hover:border-primary/50 transition-colors"
              onClick={() => navigate(`/projects/${projectId}/flows/${flow.id}`)}
            >
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <Workflow className="h-5 w-5 text-primary" />
                    <div>
                      <CardTitle className="text-base">{flow.name}</CardTitle>
                      {flow.description && (
                        <CardDescription className="mt-1">{flow.description}</CardDescription>
                      )}
                    </div>
                  </div>
                  <div className="flex items-center gap-2" onClick={e => e.stopPropagation()}>
                    <Badge variant="outline">{flow.step_count || 0} steps</Badge>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8 text-destructive hover:text-destructive"
                      onClick={() => deleteMutation.mutate(flow.id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </CardHeader>
            </Card>
          ))}
        </div>
      )}

      <FlowCreateSheet open={showCreate} onOpenChange={setShowCreate} projectId={projectId} />
    </div>
  )
}
