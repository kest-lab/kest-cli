import { useState, useEffect } from 'react'
import { Copy, ChevronDown, ChevronRight, Pencil, Save, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { toast } from 'sonner'
import { useUpdateAPISpec } from '@/hooks/use-kest-api'
import type { APISpec } from '@/types/kest-api'

const METHOD_COLORS: Record<string, string> = {
  GET: 'bg-blue-500 text-white',
  POST: 'bg-green-500 text-white',
  PUT: 'bg-amber-500 text-white',
  PATCH: 'bg-orange-500 text-white',
  DELETE: 'bg-red-500 text-white',
  HEAD: 'bg-purple-500 text-white',
  OPTIONS: 'bg-gray-500 text-white',
}

interface APIDetailPanelProps {
  spec: APISpec
  projectId: number
}

export function APIDetailPanel({ spec, projectId }: APIDetailPanelProps) {
  const [activeTab, setActiveTab] = useState('params')
  const [editing, setEditing] = useState(false)
  const [editSummary, setEditSummary] = useState(spec.summary || '')
  const [editDescription, setEditDescription] = useState(spec.description || '')
  const [editTags, setEditTags] = useState(spec.tags?.join(', ') || '')

  const updateMutation = useUpdateAPISpec()

  // Sync edit fields when spec changes
  useEffect(() => {
    setEditSummary(spec.summary || '')
    setEditDescription(spec.description || '')
    setEditTags(spec.tags?.join(', ') || '')
    setEditing(false)
  }, [spec.id])

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    toast.success('Copied to clipboard')
  }

  const handleSave = () => {
    const tags = editTags.split(',').map(t => t.trim()).filter(Boolean)
    updateMutation.mutate(
      {
        projectId,
        id: spec.id,
        data: {
          summary: editSummary,
          description: editDescription,
          tags,
        },
      },
      {
        onSuccess: () => {
          toast.success('API spec updated')
          setEditing(false)
        },
        onError: () => {
          toast.error('Failed to update API spec')
        },
      }
    )
  }

  const handleCancel = () => {
    setEditSummary(spec.summary || '')
    setEditDescription(spec.description || '')
    setEditTags(spec.tags?.join(', ') || '')
    setEditing(false)
  }

  const pathParams = spec.parameters?.filter(p => p.in === 'path') || []
  const queryParams = spec.parameters?.filter(p => p.in === 'query') || []
  const headerParams = spec.parameters?.filter(p => p.in === 'header') || []
  const hasParams = pathParams.length > 0 || queryParams.length > 0
  const hasHeaders = headerParams.length > 0
  const hasBody = !!spec.request_body
  const hasResponses = !!(spec.responses && Object.keys(spec.responses).length > 0)

  return (
    <div className="flex flex-col h-full">
      {/* URL Bar - Postman style */}
      <div className="flex items-center gap-2 p-3 border-b bg-card">
        <Badge className={`${METHOD_COLORS[spec.method] || 'bg-gray-500 text-white'} text-xs font-bold px-2.5 py-1 rounded-md shrink-0`}>
          {spec.method}
        </Badge>
        <div className="flex-1 flex items-center bg-muted rounded-md px-3 py-2 font-mono text-sm overflow-hidden">
          <span className="truncate">{spec.path}</span>
        </div>
        <Button variant="ghost" size="icon" className="shrink-0 h-8 w-8" onClick={() => copyToClipboard(spec.path)}>
          <Copy className="h-3.5 w-3.5" />
        </Button>
      </div>

      {/* Title + Description — editable */}
      <div className="px-4 py-3 border-b">
        {editing ? (
          <div className="space-y-2">
            <Input
              value={editSummary}
              onChange={e => setEditSummary(e.target.value)}
              placeholder="Summary"
              className="text-lg font-semibold h-9"
            />
            <Textarea
              value={editDescription}
              onChange={e => setEditDescription(e.target.value)}
              placeholder="Description"
              className="text-sm min-h-[60px] resize-none"
              rows={2}
            />
            <Input
              value={editTags}
              onChange={e => setEditTags(e.target.value)}
              placeholder="Tags (comma separated)"
              className="text-xs h-7"
            />
            <div className="flex gap-1.5 pt-1">
              <Button size="sm" className="h-7 text-xs" onClick={handleSave} disabled={updateMutation.isPending}>
                <Save className="h-3 w-3 mr-1" /> {updateMutation.isPending ? 'Saving...' : 'Save'}
              </Button>
              <Button size="sm" variant="ghost" className="h-7 text-xs" onClick={handleCancel}>
                <X className="h-3 w-3 mr-1" /> Cancel
              </Button>
            </div>
          </div>
        ) : (
          <div className="group relative">
            <div className="flex items-start justify-between">
              <h2 className="text-lg font-semibold">{spec.summary || spec.path}</h2>
              <Button
                variant="ghost"
                size="icon"
                className="h-7 w-7 opacity-0 group-hover:opacity-100 transition-opacity shrink-0"
                onClick={() => setEditing(true)}
              >
                <Pencil className="h-3.5 w-3.5" />
              </Button>
            </div>
            {spec.description && (
              <p className="text-sm text-muted-foreground mt-1">{spec.description}</p>
            )}
            {spec.tags && spec.tags.length > 0 && (
              <div className="flex gap-1.5 mt-2">
                {spec.tags.map((tag, i) => (
                  <Badge key={i} variant="outline" className="text-[10px]">{tag}</Badge>
                ))}
              </div>
            )}
          </div>
        )}
      </div>

      {/* Tabs: Params / Body / Headers / Responses */}
      <div className="flex-1 overflow-auto">
        <Tabs value={activeTab} onValueChange={setActiveTab} className="h-full flex flex-col">
          <TabsList className="w-full justify-start rounded-none border-b bg-transparent px-4 h-auto py-0">
            <TabsTrigger value="params" className="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent text-xs py-2">
              Params {hasParams && <Badge variant="secondary" className="ml-1 text-[9px] px-1 py-0">{pathParams.length + queryParams.length}</Badge>}
            </TabsTrigger>
            <TabsTrigger value="body" className="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent text-xs py-2">
              Body {hasBody && <span className="ml-1 w-1.5 h-1.5 rounded-full bg-green-500 inline-block" />}
            </TabsTrigger>
            <TabsTrigger value="headers" className="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent text-xs py-2">
              Headers {hasHeaders && <Badge variant="secondary" className="ml-1 text-[9px] px-1 py-0">{headerParams.length}</Badge>}
            </TabsTrigger>
            <TabsTrigger value="responses" className="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent text-xs py-2">
              Responses {hasResponses && <Badge variant="secondary" className="ml-1 text-[9px] px-1 py-0">{Object.keys(spec.responses!).length}</Badge>}
            </TabsTrigger>
          </TabsList>

          <div className="flex-1 overflow-auto p-4">
            <TabsContent value="params" className="mt-0">
              {pathParams.length > 0 && (
                <div className="mb-4">
                  <h4 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">Path Parameters</h4>
                  <ParamTable params={pathParams} />
                </div>
              )}
              {queryParams.length > 0 && (
                <div className="mb-4">
                  <h4 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">Query Parameters</h4>
                  <ParamTable params={queryParams} />
                </div>
              )}
              {!hasParams && (
                <div className="text-center py-8 text-sm text-muted-foreground">
                  No parameters defined for this endpoint
                </div>
              )}
            </TabsContent>

            <TabsContent value="body" className="mt-0">
              {spec.request_body ? (
                <div>
                  <div className="flex items-center justify-between mb-3">
                    <Badge variant="outline" className="text-xs">{spec.request_body.content_type || 'application/json'}</Badge>
                    {spec.request_body.required && <Badge className="bg-orange-100 text-orange-700 text-[10px]">Required</Badge>}
                  </div>
                  {spec.request_body.description && (
                    <p className="text-sm text-muted-foreground mb-3">{spec.request_body.description}</p>
                  )}
                  {spec.request_body.schema && (
                    <div>
                      <h4 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">Schema</h4>
                      <SchemaView schema={spec.request_body.schema} />
                    </div>
                  )}
                </div>
              ) : (
                <div className="text-center py-8 text-sm text-muted-foreground">
                  No request body for this endpoint
                </div>
              )}
            </TabsContent>

            <TabsContent value="headers" className="mt-0">
              {headerParams.length > 0 ? (
                <ParamTable params={headerParams} />
              ) : (
                <div className="text-center py-8 text-sm text-muted-foreground">
                  No custom headers defined
                </div>
              )}
            </TabsContent>

            <TabsContent value="responses" className="mt-0">
              {spec.responses && Object.keys(spec.responses).length > 0 ? (
                <div className="space-y-3">
                  {Object.entries(spec.responses).map(([code, resp]) => (
                    <ResponseItem key={code} code={code} response={resp} />
                  ))}
                </div>
              ) : (
                <div className="text-center py-8 text-sm text-muted-foreground">
                  No responses defined
                </div>
              )}
            </TabsContent>
          </div>
        </Tabs>
      </div>
    </div>
  )
}

function ParamTable({ params }: { params: any[] }) {
  return (
    <div className="border rounded-md overflow-hidden text-sm">
      <table className="w-full">
        <thead className="bg-muted/50">
          <tr>
            <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground">Name</th>
            <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground">Type</th>
            <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground">Required</th>
            <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground">Description</th>
          </tr>
        </thead>
        <tbody className="divide-y">
          {params.map((p, idx) => (
            <tr key={idx} className="hover:bg-accent/30">
              <td className="px-3 py-2 font-mono text-xs font-medium">{p.name}</td>
              <td className="px-3 py-2">
                <Badge variant="outline" className="text-[10px]">{p.schema?.type || 'string'}</Badge>
              </td>
              <td className="px-3 py-2">
                {p.required
                  ? <Badge className="bg-orange-100 text-orange-700 text-[10px]">Required</Badge>
                  : <span className="text-xs text-muted-foreground">Optional</span>
                }
              </td>
              <td className="px-3 py-2 text-muted-foreground text-xs">{p.description || '—'}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function SchemaView({ schema }: { schema: any }) {
  if (!schema) return null

  if (schema.type === 'object' && schema.properties) {
    const required = schema.required || []
    return (
      <div className="border rounded-md overflow-hidden text-sm">
        <table className="w-full">
          <thead className="bg-muted/50">
            <tr>
              <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground">Field</th>
              <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground">Type</th>
              <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground">Required</th>
              <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground">Details</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            {Object.entries(schema.properties).map(([name, prop]: [string, any]) => (
              <tr key={name} className="hover:bg-accent/30">
                <td className="px-3 py-2 font-mono text-xs font-medium">{name}</td>
                <td className="px-3 py-2">
                  <Badge variant="outline" className="text-[10px]">{prop.type || 'any'}</Badge>
                  {prop.enum && <span className="text-[10px] text-muted-foreground ml-1">enum</span>}
                </td>
                <td className="px-3 py-2">
                  {required.includes(name)
                    ? <Badge className="bg-orange-100 text-orange-700 text-[10px]">Required</Badge>
                    : <span className="text-xs text-muted-foreground">Optional</span>
                  }
                </td>
                <td className="px-3 py-2 text-xs text-muted-foreground">
                  {prop.format && <Badge variant="outline" className="text-[9px] mr-1">{prop.format}</Badge>}
                  {prop.enum && <span className="font-mono">[{prop.enum.join(', ')}]</span>}
                  {prop.minLength && <span>min: {prop.minLength}</span>}
                  {prop.default !== undefined && <span>default: {String(prop.default)}</span>}
                  {prop.items?.type && <span>items: {prop.items.type}</span>}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    )
  }

  // Fallback: show raw JSON
  return (
    <pre className="bg-muted p-3 rounded-md text-xs font-mono overflow-x-auto">
      {JSON.stringify(schema, null, 2)}
    </pre>
  )
}

function ResponseItem({ code, response }: { code: string; response: any }) {
  const [open, setOpen] = useState(parseInt(code) < 300)
  const isSuccess = parseInt(code) < 400

  return (
    <div className="border rounded-md overflow-hidden">
      <button
        className="w-full flex items-center gap-2 px-3 py-2.5 hover:bg-accent/30 transition-colors text-left"
        onClick={() => setOpen(!open)}
      >
        {open ? <ChevronDown className="h-3.5 w-3.5 shrink-0" /> : <ChevronRight className="h-3.5 w-3.5 shrink-0" />}
        <Badge className={`${isSuccess ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'} text-xs`}>
          {code}
        </Badge>
        <span className="text-sm flex-1 truncate">{response.description}</span>
        {response.content_type && (
          <span className="text-[10px] text-muted-foreground">{response.content_type}</span>
        )}
      </button>
      {open && response.schema && (
        <div className="px-3 pb-3 border-t">
          <pre className="bg-muted p-3 rounded-md text-xs font-mono overflow-x-auto mt-2">
            {JSON.stringify(response.schema, null, 2)}
          </pre>
        </div>
      )}
    </div>
  )
}
