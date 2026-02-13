import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus, Search, FileText, FileUp, FileDown } from 'lucide-react'
import { toast } from 'sonner'
import { kestApi } from '@/services/kest-api.service'
import type { APISpec } from '@/types/kest-api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { CreateAPISpecDialog } from './create-api-spec-dialog'
import { ImportAPISpecDialog } from './import-api-spec-dialog'

interface APISpecListProps {
  projectId: number
  apiSpecs: APISpec[]
  isLoading: boolean
}

export function APISpecList({ projectId, apiSpecs, isLoading }: APISpecListProps) {
  const [searchQuery, setSearchQuery] = useState('')
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [importDialogOpen, setImportDialogOpen] = useState(false)

  const filteredSpecs = apiSpecs.filter((spec) =>
    spec.summary?.toLowerCase().includes(searchQuery.toLowerCase()) ||
    spec.path?.toLowerCase().includes(searchQuery.toLowerCase()) ||
    spec.description?.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const getMethodColor = (method: string) => {
    const colors: Record<string, string> = {
      GET: 'bg-blue-100 text-blue-800',
      POST: 'bg-green-100 text-green-800',
      PUT: 'bg-yellow-100 text-yellow-800',
      PATCH: 'bg-orange-100 text-orange-800',
      DELETE: 'bg-red-100 text-red-800',
    }
    return colors[method] || 'bg-gray-100 text-gray-800'
  }

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      done: 'bg-green-100 text-green-800',
      undone: 'bg-gray-100 text-gray-800',
      deprecated: 'bg-red-100 text-red-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[1, 2, 3].map((i) => (
          <Skeleton key={i} className="h-24" />
        ))}
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="relative flex-1 max-w-md">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
          <Input
            placeholder="Search APIs..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10"
          />
        </div>
        <div className="flex items-center gap-2 shrink-0">
          <Button variant="outline" onClick={() => setImportDialogOpen(true)}>
            <FileUp className="w-4 h-4 mr-2" />
            Import
          </Button>
          <Button variant="outline" onClick={async () => {
            try {
              const res = await kestApi.apiSpec.export(projectId, 'json')
              const blob = new Blob([JSON.stringify(res, null, 2)], { type: 'application/json' })
              const url = URL.createObjectURL(blob)
              const a = document.createElement('a')
              a.href = url
              a.download = `api-specs-${projectId}.json`
              a.click()
              URL.revokeObjectURL(url)
              toast.success('Export successful')
            } catch {
              toast.error('Failed to export API specs')
            }
          }}>
            <FileDown className="w-4 h-4 mr-2" />
            Export
          </Button>
          <Button onClick={() => setCreateDialogOpen(true)}>
            <Plus className="w-4 h-4 mr-2" />
            Add API
          </Button>
        </div>
      </div>

      {filteredSpecs.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <FileText className="w-16 h-16 text-gray-300 mb-4" />
            <h3 className="text-lg font-semibold mb-2">No API specifications found</h3>
            <p className="text-gray-600 mb-4">
              {searchQuery ? 'Try a different search term' : 'Add your first API specification'}
            </p>
            {!searchQuery && (
              <Button onClick={() => setCreateDialogOpen(true)}>
                <Plus className="w-4 h-4 mr-2" />
                Add API Specification
              </Button>
            )}
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {filteredSpecs.map((spec) => (
            <Card key={spec.id} className="hover:shadow-md transition-shadow">
              <CardHeader className="pb-3">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                      <Badge className={getMethodColor(spec.method)}>
                        {spec.method}
                      </Badge>
                      <Link to={`/projects/${projectId}/api-specs/${spec.id}`} className="text-sm font-mono bg-gray-50 px-2 py-1 rounded hover:text-blue-600">
                        {spec.path}
                      </Link>
                    </div>
                    <CardTitle className="text-lg">{spec.summary || spec.path}</CardTitle>
                    {spec.description && (
                      <CardDescription className="mt-1 line-clamp-2">{spec.description}</CardDescription>
                    )}
                  </div>
                  {spec.status && (
                    <Badge className={getStatusColor(spec.status)}>
                      {spec.status}
                    </Badge>
                  )}
                </div>
              </CardHeader>
              {(spec.tags && spec.tags.length > 0) && (
                <CardContent className="pt-0">
                  <div className="flex gap-2 flex-wrap">
                    {spec.tags.map((tag, idx) => (
                      <Badge key={idx} variant="outline" className="text-xs">
                        {tag}
                      </Badge>
                    ))}
                  </div>
                </CardContent>
              )}
            </Card>
          ))}
        </div>
      )}

      <CreateAPISpecDialog
        projectId={projectId}
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
      />
      <ImportAPISpecDialog
        projectId={projectId}
        open={importDialogOpen}
        onOpenChange={setImportDialogOpen}
      />
    </div>
  )
}
