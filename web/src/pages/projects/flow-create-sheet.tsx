import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription, SheetFooter,
} from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { cn } from '@/lib/utils'
import { Loader2, Workflow, LogIn, Database, Shield, FileCode } from 'lucide-react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { flowService } from '@/services/flow.service'
import type { CreateStepRequest, CreateEdgeRequest } from '@/types/kest-api'
import { toast } from 'sonner'

// ========== Flow Templates ==========

interface FlowTemplate {
  id: string
  name: string
  description: string
  icon: React.ReactNode
  tags: string[]
  steps: CreateStepRequest[]
  edges: CreateEdgeRequest[]
}

const TEMPLATES: FlowTemplate[] = [
  {
    id: 'empty',
    name: 'Empty Flow',
    description: 'Start from scratch with a blank canvas.',
    icon: <FileCode className="h-5 w-5" />,
    tags: ['blank'],
    steps: [],
    edges: [],
  },
  {
    id: 'login',
    name: 'User Login Flow',
    description: 'Register a user, login, then fetch profile with the returned token.',
    icon: <LogIn className="h-5 w-5" />,
    tags: ['auth', 'jwt'],
    steps: [
      {
        name: 'Register User',
        sort_order: 0,
        method: 'POST',
        url: '{{base_url}}/v1/register',
        headers: '{"Content-Type": "application/json"}',
        body: '{"username": "testuser_{{$timestamp}}", "email": "test_{{$timestamp}}@example.com", "password": "SecurePass123!"}',
        captures: 'user_id=data.id',
        asserts: 'status=201',
        position_x: 250,
        position_y: 0,
      },
      {
        name: 'Login',
        sort_order: 1,
        method: 'POST',
        url: '{{base_url}}/v1/login',
        headers: '{"Content-Type": "application/json"}',
        body: '{"username": "testuser_{{$timestamp}}", "password": "SecurePass123!"}',
        captures: 'token=data.access_token',
        asserts: 'status=200\nbody.data.access_token=exists',
        position_x: 250,
        position_y: 180,
      },
      {
        name: 'Get Profile',
        sort_order: 2,
        method: 'GET',
        url: '{{base_url}}/v1/users/profile',
        headers: '{"Authorization": "Bearer {{token}}"}',
        body: '',
        captures: '',
        asserts: 'status=200\nbody.data.username=exists',
        position_x: 250,
        position_y: 360,
      },
    ],
    edges: [
      { source_step_id: 0, target_step_id: 1 },
      { source_step_id: 1, target_step_id: 2 },
    ],
  },
  {
    id: 'crud',
    name: 'CRUD Resource Flow',
    description: 'Login, then Create → Read → Update → Delete a project resource.',
    icon: <Database className="h-5 w-5" />,
    tags: ['crud', 'rest'],
    steps: [
      {
        name: 'Login',
        sort_order: 0,
        method: 'POST',
        url: '{{base_url}}/v1/login',
        headers: '{"Content-Type": "application/json"}',
        body: '{"username": "{{admin_user}}", "password": "{{admin_pass}}"}',
        captures: 'token=data.access_token',
        asserts: 'status=200',
        position_x: 250,
        position_y: 0,
      },
      {
        name: 'Create Project',
        sort_order: 1,
        method: 'POST',
        url: '{{base_url}}/v1/projects',
        headers: '{"Content-Type": "application/json", "Authorization": "Bearer {{token}}"}',
        body: '{"name": "Test Project {{$timestamp}}", "description": "Created by flow test"}',
        captures: 'project_id=data.id',
        asserts: 'status=201\nbody.data.name=exists',
        position_x: 250,
        position_y: 180,
      },
      {
        name: 'Get Project',
        sort_order: 2,
        method: 'GET',
        url: '{{base_url}}/v1/projects/{{project_id}}',
        headers: '{"Authorization": "Bearer {{token}}"}',
        body: '',
        captures: '',
        asserts: 'status=200\nbody.data.id={{project_id}}',
        position_x: 250,
        position_y: 360,
      },
      {
        name: 'Update Project',
        sort_order: 3,
        method: 'PUT',
        url: '{{base_url}}/v1/projects/{{project_id}}',
        headers: '{"Content-Type": "application/json", "Authorization": "Bearer {{token}}"}',
        body: '{"name": "Updated Project", "description": "Updated by flow test"}',
        captures: '',
        asserts: 'status=200',
        position_x: 250,
        position_y: 540,
      },
      {
        name: 'Delete Project',
        sort_order: 4,
        method: 'DELETE',
        url: '{{base_url}}/v1/projects/{{project_id}}',
        headers: '{"Authorization": "Bearer {{token}}"}',
        body: '',
        captures: '',
        asserts: 'status=204',
        position_x: 250,
        position_y: 720,
      },
    ],
    edges: [
      { source_step_id: 0, target_step_id: 1 },
      { source_step_id: 1, target_step_id: 2 },
      { source_step_id: 2, target_step_id: 3 },
      { source_step_id: 3, target_step_id: 4 },
    ],
  },
  {
    id: 'auth-resource',
    name: 'Auth + API Spec Flow',
    description: 'Login, create a project, add a category, then create an API spec under it.',
    icon: <Shield className="h-5 w-5" />,
    tags: ['auth', 'api-spec'],
    steps: [
      {
        name: 'Login',
        sort_order: 0,
        method: 'POST',
        url: '{{base_url}}/v1/login',
        headers: '{"Content-Type": "application/json"}',
        body: '{"username": "{{admin_user}}", "password": "{{admin_pass}}"}',
        captures: 'token=data.access_token',
        asserts: 'status=200',
        position_x: 250,
        position_y: 0,
      },
      {
        name: 'Create Project',
        sort_order: 1,
        method: 'POST',
        url: '{{base_url}}/v1/projects',
        headers: '{"Content-Type": "application/json", "Authorization": "Bearer {{token}}"}',
        body: '{"name": "Flow Test Project {{$timestamp}}"}',
        captures: 'project_id=data.id',
        asserts: 'status=201',
        position_x: 250,
        position_y: 180,
      },
      {
        name: 'Create Category',
        sort_order: 2,
        method: 'POST',
        url: '{{base_url}}/v1/projects/{{project_id}}/categories',
        headers: '{"Content-Type": "application/json", "Authorization": "Bearer {{token}}"}',
        body: '{"name": "Users", "description": "User management APIs"}',
        captures: 'category_id=data.id',
        asserts: 'status=201',
        position_x: 250,
        position_y: 360,
      },
      {
        name: 'Create API Spec',
        sort_order: 3,
        method: 'POST',
        url: '{{base_url}}/v1/projects/{{project_id}}/api-specs',
        headers: '{"Content-Type": "application/json", "Authorization": "Bearer {{token}}"}',
        body: '{"method": "GET", "path": "/v1/users/profile", "summary": "Get User Profile", "version": "1.0.0", "category_id": {{category_id}}}',
        captures: 'spec_id=data.id',
        asserts: 'status=201\nbody.data.method=GET',
        position_x: 250,
        position_y: 540,
      },
    ],
    edges: [
      { source_step_id: 0, target_step_id: 1 },
      { source_step_id: 1, target_step_id: 2 },
      { source_step_id: 2, target_step_id: 3 },
    ],
  },
]

// ========== Component ==========

interface FlowCreateSheetProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  projectId: number
}

const METHOD_COLORS: Record<string, string> = {
  GET: 'bg-blue-100 text-blue-700',
  POST: 'bg-green-100 text-green-700',
  PUT: 'bg-amber-100 text-amber-700',
  DELETE: 'bg-red-100 text-red-700',
  PATCH: 'bg-orange-100 text-orange-700',
}

export function FlowCreateSheet({ open, onOpenChange, projectId }: FlowCreateSheetProps) {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [selectedTemplate, setSelectedTemplate] = useState<string>('empty')

  const template = TEMPLATES.find(t => t.id === selectedTemplate) || TEMPLATES[0]

  const createMutation = useMutation({
    mutationFn: async () => {
      // 1. Create the flow
      const flow = await flowService.create(projectId, {
        name: name || template.name,
        description: description || template.description,
      })

      // 2. If template has steps, save them
      if (template.steps.length > 0) {
        await flowService.save(projectId, flow.id, {
          name: flow.name,
          description: flow.description,
          steps: template.steps,
          edges: template.edges,
        })
      }

      return flow
    },
    onSuccess: (flow) => {
      queryClient.invalidateQueries({ queryKey: ['flows', projectId] })
      toast.success(`Flow "${flow.name}" created`)
      onOpenChange(false)
      resetForm()
      navigate(`/projects/${projectId}/flows/${flow.id}`)
    },
    onError: () => toast.error('Failed to create flow'),
  })

  const resetForm = () => {
    setName('')
    setDescription('')
    setSelectedTemplate('empty')
  }

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full flex flex-col p-0 overflow-hidden" style={{ maxWidth: '32rem' }}>
        <SheetHeader className="px-6 pt-6 pb-4 border-b">
          <SheetTitle className="flex items-center gap-2">
            <Workflow className="h-5 w-5 text-primary" />
            Create New Flow
          </SheetTitle>
          <SheetDescription>
            Choose a template to get started quickly, or create an empty flow.
          </SheetDescription>
        </SheetHeader>

        <ScrollArea className="flex-1">
          <div className="px-6 py-4 space-y-6">
            {/* Form Fields */}
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="flow-name">Flow Name</Label>
                <Input
                  id="flow-name"
                  placeholder={template.name || 'e.g. User Login Flow'}
                  value={name}
                  onChange={e => setName(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="flow-desc">Description</Label>
                <Textarea
                  id="flow-desc"
                  placeholder={template.description || 'Optional description'}
                  value={description}
                  onChange={e => setDescription(e.target.value)}
                  className="resize-none min-h-[60px]"
                  rows={2}
                />
              </div>
            </div>

            {/* Template Selection */}
            <div className="space-y-3">
              <Label className="text-sm font-semibold">Choose a Template</Label>
              <div className="grid gap-3">
                {TEMPLATES.map(tmpl => (
                  <button
                    key={tmpl.id}
                    type="button"
                    className={cn(
                      'w-full text-left border rounded-lg p-4 transition-all',
                      selectedTemplate === tmpl.id
                        ? 'border-primary bg-primary/5 ring-1 ring-primary'
                        : 'border-border hover:border-primary/40 hover:bg-accent/30'
                    )}
                    onClick={() => setSelectedTemplate(tmpl.id)}
                  >
                    <div className="flex items-start gap-3">
                      <div className={cn(
                        'mt-0.5 p-2 rounded-md',
                        selectedTemplate === tmpl.id ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'
                      )}>
                        {tmpl.icon}
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <span className="font-medium text-sm">{tmpl.name}</span>
                          {tmpl.steps.length > 0 && (
                            <Badge variant="secondary" className="text-[10px] px-1.5 py-0">
                              {tmpl.steps.length} steps
                            </Badge>
                          )}
                        </div>
                        <p className="text-xs text-muted-foreground mt-1">{tmpl.description}</p>
                        {tmpl.tags.length > 0 && (
                          <div className="flex gap-1 mt-2">
                            {tmpl.tags.map(tag => (
                              <Badge key={tag} variant="outline" className="text-[9px] px-1 py-0">{tag}</Badge>
                            ))}
                          </div>
                        )}
                      </div>
                    </div>
                  </button>
                ))}
              </div>
            </div>

            {/* Template Preview */}
            {template.steps.length > 0 && (
              <div className="space-y-3">
                <Label className="text-sm font-semibold">Steps Preview</Label>
                <div className="border rounded-lg overflow-hidden">
                  {template.steps.map((step, i) => (
                    <div
                      key={i}
                      className={cn(
                        'flex items-center gap-3 px-4 py-2.5 text-sm',
                        i > 0 && 'border-t'
                      )}
                    >
                      <span className="text-xs text-muted-foreground w-5 text-center font-mono">{i + 1}</span>
                      <Badge className={cn('text-[10px] font-bold px-1.5 py-0 shrink-0', METHOD_COLORS[step.method] || 'bg-gray-100 text-gray-700')}>
                        {step.method}
                      </Badge>
                      <div className="flex-1 min-w-0">
                        <span className="font-medium text-xs">{step.name}</span>
                        <p className="font-mono text-[10px] text-muted-foreground truncate">{step.url}</p>
                      </div>
                      {step.captures && (
                        <Badge variant="outline" className="text-[9px] px-1 py-0 shrink-0">capture</Badge>
                      )}
                    </div>
                  ))}
                </div>
                {template.edges.length > 0 && (
                  <p className="text-[10px] text-muted-foreground">
                    {template.edges.length} edge{template.edges.length > 1 ? 's' : ''} connecting steps sequentially. Variables are passed via captures.
                  </p>
                )}
              </div>
            )}
          </div>
        </ScrollArea>

        <SheetFooter className="px-6 py-4 border-t">
          <div className="flex gap-2 w-full">
            <Button variant="outline" className="flex-1" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button
              className="flex-1"
              onClick={() => createMutation.mutate()}
              disabled={createMutation.isPending}
            >
              {createMutation.isPending ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  Creating...
                </>
              ) : (
                <>
                  <Workflow className="h-4 w-4 mr-2" />
                  Create Flow
                </>
              )}
            </Button>
          </div>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  )
}
