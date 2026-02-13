import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Plus, Edit, Trash2, FileText } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { toast } from 'sonner'
import { request } from '@/http'
import { CreateAPISpecDialog } from './create-api-spec-dialog'

interface ApiSpec {
  id: number
  project_id: number
  method: string
  path: string
  summary: string
  description?: string
  created_at: string
  updated_at: string
}

export function ApiSpecsPage() {
  const { projectId } = useParams<{ projectId: string }>()
  const queryClient = useQueryClient()
  const [createDialogOpen, setCreateDialogOpen] = useState(false)

  const { data: specs = [], isLoading } = useQuery({
    queryKey: ['api-specs', projectId],
    queryFn: async () => {
      const response = await request.get<{ code: number; data: ApiSpec[] }>(
        `/v1/projects/${projectId}/api-specs`
      )
      return response.data || []
    },
    enabled: !!projectId,
  })

  const deleteMutation = useMutation({
    mutationFn: async (specId: number) => {
      await request.delete(`/v1/projects/${projectId}/api-specs/${specId}`)
    },
    onSuccess: () => {
      toast.success('API spec deleted successfully')
      queryClient.invalidateQueries({ queryKey: ['api-specs', projectId] })
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to delete API spec')
    },
  })

  const getMethodColor = (method: string) => {
    const colors: Record<string, string> = {
      GET: 'bg-blue-500',
      POST: 'bg-green-500',
      PUT: 'bg-yellow-500',
      PATCH: 'bg-orange-500',
      DELETE: 'bg-red-500',
    }
    return colors[method] || 'bg-gray-500'
  }

  if (isLoading) {
    return (
      <div className="container mx-auto p-8">
        <div className="flex items-center justify-center h-64">
          <p className="text-muted-foreground">Loading API specs...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto p-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold">API Specifications</h1>
          <p className="text-muted-foreground mt-1">
            Manage your API endpoints and documentation
          </p>
        </div>
        <Button onClick={() => setCreateDialogOpen(true)}>
          <Plus className="w-4 h-4 mr-2" />
          Add API Spec
        </Button>
      </div>

      {specs.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-16">
            <FileText className="w-16 h-16 text-muted-foreground mb-4" />
            <h3 className="text-lg font-semibold mb-2">No API specs found</h3>
            <p className="text-muted-foreground mb-4">
              Create your first API specification to get started
            </p>
            <Button onClick={() => setCreateDialogOpen(true)}>
              <Plus className="w-4 h-4 mr-2" />
              Create API Spec
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4">
          {specs.map((spec) => (
            <Card key={spec.id}>
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-3">
                    <Badge className={`${getMethodColor(spec.method)} text-white`}>
                      {spec.method}
                    </Badge>
                    <div>
                      <CardTitle className="text-lg">{spec.path}</CardTitle>
                      <p className="text-sm text-muted-foreground mt-1">
                        {spec.summary}
                      </p>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <Button variant="ghost" size="sm">
                      <Edit className="w-4 h-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => deleteMutation.mutate(spec.id)}
                      disabled={deleteMutation.isPending}
                    >
                      <Trash2 className="w-4 h-4 text-red-600" />
                    </Button>
                  </div>
                </div>
              </CardHeader>
              {spec.description && (
                <CardContent>
                  <p className="text-sm text-muted-foreground">{spec.description}</p>
                </CardContent>
              )}
            </Card>
          ))}
        </div>
      )}

      <CreateAPISpecDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
        projectId={Number(projectId)}
      />
    </div>
  )
}
