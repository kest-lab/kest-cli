import { useMemo, useState } from 'react'
import { toast } from 'sonner'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Textarea } from '@/components/ui/textarea'
import {
  useAllAPISpecs,
  useCategories,
  useCreateTestCase,
  useDeleteTestCase,
  useDuplicateTestCase,
  useGenerateTestCasesFromSpec,
  useRunTestCase,
  useTestCase,
  useTestCases,
  useUpdateTestCase,
} from '@/hooks/use-kest-api'
import type {
  APICategory,
  APISpec,
  CreateTestCaseRequest,
  ListTestCasesParams,
  TestCase,
} from '@/types/kest-api'

interface TestCasesPanelProps {
  projectId: number
}

const METHOD_OPTIONS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'] as const

const parseMaybeJson = (value: string) => {
  const trimmed = value.trim()
  if (!trimmed) return undefined
  return JSON.parse(trimmed)
}

const toPretty = (value: any) => (value == null ? '' : JSON.stringify(value, null, 2))

export function TestCasesPanel({ projectId }: TestCasesPanelProps) {
  const [filters, setFilters] = useState<ListTestCasesParams>({ page: 1, per_page: 20 })
  const [createOpen, setCreateOpen] = useState(false)
  const [generateOpen, setGenerateOpen] = useState(false)
  const [duplicateTarget, setDuplicateTarget] = useState<TestCase | null>(null)
  const [runTarget, setRunTarget] = useState<TestCase | null>(null)
  const [detailId, setDetailId] = useState<number | null>(null)
  const [editing, setEditing] = useState(false)
  const [runResult, setRunResult] = useState<any>(null)

  const [form, setForm] = useState({
    name: '',
    description: '',
    method: 'GET',
    path: '',
    api_spec_id: '',
    environment: '',
    category_id: '',
    request_headers: '{}',
    request_body: '',
    expected_status: '200',
    expected_response: '',
    variables: '',
    setup_script: '',
    teardown_script: '',
  })

  const [duplicateName, setDuplicateName] = useState('')
  const [duplicateEnv, setDuplicateEnv] = useState('')
  const [duplicateDescription, setDuplicateDescription] = useState('')
  const [runEnvironment, setRunEnvironment] = useState('')
  const [runVariables, setRunVariables] = useState('')
  const [runAsync, setRunAsync] = useState(false)
  const [generateForm, setGenerateForm] = useState({
    api_spec_id: '',
    environment: '',
    category_id: '',
    generate_positive_tests: true,
    generate_negative_tests: true,
    include_auth_tests: true,
    max_tests_per_endpoint: 5,
  })

  const { data: listData, isLoading } = useTestCases(projectId, filters)
  const { data: categoriesData } = useCategories(projectId, { per_page: 1000 })
  const { data: specsData } = useAllAPISpecs(projectId)
  const detailQuery = useTestCase(projectId, detailId || 0)

  const createMutation = useCreateTestCase(projectId)
  const updateMutation = useUpdateTestCase(projectId)
  const deleteMutation = useDeleteTestCase(projectId)
  const duplicateMutation = useDuplicateTestCase(projectId)
  const runMutation = useRunTestCase(projectId)
  const generateMutation = useGenerateTestCasesFromSpec(projectId)

  const testCases: TestCase[] = (() => {
    if (!listData) return []
    if (Array.isArray(listData)) return listData as any
    if (Array.isArray((listData as any).items)) return (listData as any).items
    if (Array.isArray((listData as any).data?.items)) return (listData as any).data.items
    if (Array.isArray((listData as any).data)) return (listData as any).data
    return []
  })()

  const categories = useMemo(() => {
    if (!categoriesData) return [] as APICategory[]
    if (Array.isArray(categoriesData)) return categoriesData as APICategory[]
    if (Array.isArray((categoriesData as any)?.items)) return (categoriesData as any).items as APICategory[]
    if (Array.isArray((categoriesData as any)?.data?.items)) return (categoriesData as any).data.items as APICategory[]
    if (Array.isArray((categoriesData as any)?.data)) return (categoriesData as any).data as APICategory[]
    return [] as APICategory[]
  }, [categoriesData])

  const apiSpecs = useMemo(() => {
    if (!specsData) return [] as APISpec[]
    if (Array.isArray(specsData)) return specsData as APISpec[]
    if (Array.isArray((specsData as any)?.items)) return (specsData as any).items as APISpec[]
    if (Array.isArray((specsData as any)?.data?.items)) return (specsData as any).data.items as APISpec[]
    if (Array.isArray((specsData as any)?.data)) return (specsData as any).data as APISpec[]
    return [] as APISpec[]
  }, [specsData])

  const resetForm = () => {
    setForm({
      name: '',
      description: '',
      method: 'GET',
      path: '',
      api_spec_id: '',
      environment: '',
      category_id: '',
      request_headers: '{}',
      request_body: '',
      expected_status: '200',
      expected_response: '',
      variables: '',
      setup_script: '',
      teardown_script: '',
    })
  }

  const openCreate = () => {
    setEditing(false)
    setDetailId(null)
    resetForm()
    setCreateOpen(true)
  }

  const openEdit = (tc: TestCase) => {
    setEditing(true)
    setDetailId(tc.id)
    setForm({
      name: tc.name || '',
      description: tc.description || '',
      method: tc.method || 'GET',
      path: tc.path || '',
      api_spec_id: tc.api_spec_id ? String(tc.api_spec_id) : '',
      environment: tc.environment || tc.env || '',
      category_id: tc.category_id ? String(tc.category_id) : '',
      request_headers: toPretty(tc.request_headers || {}),
      request_body: toPretty(tc.request_body),
      expected_status: String(tc.expected_status || 200),
      expected_response: toPretty(tc.expected_response),
      variables: toPretty(tc.variables),
      setup_script: tc.setup_script || '',
      teardown_script: tc.teardown_script || '',
    })
    setCreateOpen(true)
  }

  const submitForm = async () => {
    try {
      const payload: CreateTestCaseRequest = {
        name: form.name,
        description: form.description || undefined,
        method: form.method as any,
        path: form.path,
        api_spec_id: form.api_spec_id ? Number(form.api_spec_id) : undefined,
        environment: form.environment || undefined,
        category_id: form.category_id ? Number(form.category_id) : undefined,
        request_headers: parseMaybeJson(form.request_headers || '{}'),
        request_body: parseMaybeJson(form.request_body),
        expected_status: Number(form.expected_status || 200),
        expected_response: parseMaybeJson(form.expected_response),
        variables: parseMaybeJson(form.variables),
        setup_script: form.setup_script || undefined,
        teardown_script: form.teardown_script || undefined,
      } as any

      if (!payload.name || !payload.method || !payload.path) {
        toast.error('Name / Method / Path are required')
        return
      }

      if (editing && detailId) {
        await updateMutation.mutateAsync({ id: detailId, data: payload })
        toast.success('Test case updated')
      } else {
        await createMutation.mutateAsync(payload)
        toast.success('Test case created')
      }

      setCreateOpen(false)
      resetForm()
      setEditing(false)
    } catch (err: any) {
      toast.error(err?.message || 'Failed to save test case')
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this test case?')) return
    try {
      await deleteMutation.mutateAsync(id)
      toast.success('Deleted')
      if (detailId === id) setDetailId(null)
    } catch (err: any) {
      toast.error(err?.message || 'Failed to delete')
    }
  }

  const handleDuplicate = async () => {
    if (!duplicateTarget || !duplicateName.trim()) return
    try {
      await duplicateMutation.mutateAsync({
        id: duplicateTarget.id,
        data: {
          name: duplicateName.trim(),
          environment: duplicateEnv || undefined,
          description: duplicateDescription || undefined,
        },
      })
      toast.success('Duplicated')
      setDuplicateTarget(null)
      setDuplicateName('')
      setDuplicateEnv('')
      setDuplicateDescription('')
    } catch (err: any) {
      toast.error(err?.message || 'Failed to duplicate')
    }
  }

  const handleRun = async () => {
    if (!runTarget) return
    try {
      const result = await runMutation.mutateAsync({
        id: runTarget.id,
        data: {
          environment: runEnvironment || undefined,
          variables: parseMaybeJson(runVariables),
          async: runAsync,
        },
      })
      setRunResult(result)
      toast.success('Run request submitted')
    } catch (err: any) {
      toast.error(err?.message || 'Failed to run')
    }
  }

  const handleGenerate = async () => {
    if (!generateForm.api_spec_id) {
      toast.error('API Spec is required')
      return
    }
    try {
      const result = await generateMutation.mutateAsync({
        api_spec_id: Number(generateForm.api_spec_id),
        environment: generateForm.environment || undefined,
        category_id: generateForm.category_id ? Number(generateForm.category_id) : undefined,
        options: {
          generate_positive_tests: generateForm.generate_positive_tests,
          generate_negative_tests: generateForm.generate_negative_tests,
          include_auth_tests: generateForm.include_auth_tests,
          max_tests_per_endpoint: Number(generateForm.max_tests_per_endpoint || 5),
        },
      })
      toast.success(`Generated ${result.generated}, updated ${result.updated}`)
      setGenerateOpen(false)
    } catch (err: any) {
      toast.error(err?.message || 'Failed to generate')
    }
  }

  return (
    <Card className="border-none shadow-none h-full">
      <CardHeader className="px-0 pt-0">
        <div className="flex items-start justify-between gap-3">
          <div>
            <CardTitle className="text-xl">Test Cases</CardTitle>
            <CardDescription>Manage, generate and run API test cases</CardDescription>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => setGenerateOpen(true)}>Generate From Spec</Button>
            <Button onClick={openCreate}>Create Test Case</Button>
          </div>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-6 gap-2 mt-3">
          <Input
            placeholder="Keyword"
            value={filters.keyword || ''}
            onChange={(e) => setFilters((s) => ({ ...s, page: 1, keyword: e.target.value || undefined }))}
          />
          <Select
            value={filters.status || 'all'}
            onValueChange={(value) => setFilters((s) => ({ ...s, page: 1, status: value === 'all' ? undefined : value as any }))}
          >
            <SelectTrigger><SelectValue placeholder="Status" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Status</SelectItem>
              <SelectItem value="active">active</SelectItem>
              <SelectItem value="inactive">inactive</SelectItem>
              <SelectItem value="archived">archived</SelectItem>
            </SelectContent>
          </Select>
          <Select
            value={filters.api_spec_id ? String(filters.api_spec_id) : 'all'}
            onValueChange={(value) => setFilters((s) => ({ ...s, page: 1, api_spec_id: value === 'all' ? undefined : Number(value) }))}
          >
            <SelectTrigger><SelectValue placeholder="API Spec" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All API Specs</SelectItem>
              {apiSpecs.map((spec) => (
                <SelectItem key={spec.id} value={String(spec.id)}>{spec.method} {spec.path}</SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Select
            value={filters.category_id ? String(filters.category_id) : 'all'}
            onValueChange={(value) => setFilters((s) => ({ ...s, page: 1, category_id: value === 'all' ? undefined : Number(value) }))}
          >
            <SelectTrigger><SelectValue placeholder="Category" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Categories</SelectItem>
              {categories.map((cat) => (
                <SelectItem key={cat.id} value={String(cat.id)}>{cat.name}</SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Input
            placeholder="Environment"
            value={filters.env || ''}
            onChange={(e) => setFilters((s) => ({ ...s, page: 1, env: e.target.value || undefined }))}
          />
          <Select
            value={String(filters.per_page || 20)}
            onValueChange={(value) => setFilters((s) => ({ ...s, page: 1, per_page: Number(value) }))}
          >
            <SelectTrigger><SelectValue placeholder="Page Size" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="10">10 / page</SelectItem>
              <SelectItem value="20">20 / page</SelectItem>
              <SelectItem value="50">50 / page</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </CardHeader>
      <CardContent className="px-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Method</TableHead>
              <TableHead>Path</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Last Run</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {!isLoading && testCases.length === 0 && (
              <TableRow>
                <TableCell colSpan={6} className="text-center text-muted-foreground py-8">
                  No test cases found
                </TableCell>
              </TableRow>
            )}
            {testCases.map((tc) => (
              <TableRow key={tc.id}>
                <TableCell className="font-medium">{tc.name}</TableCell>
                <TableCell>
                  <Badge variant="outline">{tc.method || '-'}</Badge>
                </TableCell>
                <TableCell className="font-mono text-xs">{tc.path || '-'}</TableCell>
                <TableCell>{tc.status || '-'}</TableCell>
                <TableCell>{tc.last_run_status || '-'}</TableCell>
                <TableCell>
                  <div className="flex gap-2">
                    <Button size="sm" variant="outline" onClick={() => setDetailId(tc.id)}>View</Button>
                    <Button size="sm" variant="outline" onClick={() => openEdit(tc)}>Edit</Button>
                    <Button size="sm" variant="outline" onClick={() => { setDuplicateTarget(tc); setDuplicateName(`${tc.name} Copy`) }}>Duplicate</Button>
                    <Button size="sm" variant="outline" onClick={() => { setRunTarget(tc); setRunResult(null) }}>Run</Button>
                    <Button size="sm" variant="destructive" onClick={() => handleDelete(tc.id)}>Delete</Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        <div className="flex justify-end gap-2 mt-3">
          <Button
            variant="outline"
            onClick={() => setFilters((s) => ({ ...s, page: Math.max(1, (s.page || 1) - 1) }))}
            disabled={(filters.page || 1) <= 1}
          >
            Prev
          </Button>
          <Button
            variant="outline"
            onClick={() => setFilters((s) => ({ ...s, page: (s.page || 1) + 1 }))}
          >
            Next
          </Button>
        </div>
      </CardContent>

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="max-w-3xl w-[95vw] max-h-[90vh] overflow-hidden flex flex-col pb-2">
          <DialogHeader>
            <DialogTitle>{editing ? 'Edit Test Case' : 'Create Test Case'}</DialogTitle>
            <DialogDescription>Define request, expectation and runtime settings.</DialogDescription>
          </DialogHeader>
          <div className="overflow-y-auto pr-1 flex-1">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              <div className="space-y-1">
                <Label>Name</Label>
                <Input value={form.name} onChange={(e) => setForm((s) => ({ ...s, name: e.target.value }))} />
              </div>
              <div className="space-y-1">
                <Label>Method</Label>
                <Select value={form.method} onValueChange={(value) => setForm((s) => ({ ...s, method: value }))}>
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    {METHOD_OPTIONS.map((m) => <SelectItem key={m} value={m}>{m}</SelectItem>)}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1 md:col-span-2">
                <Label>Path</Label>
                <Input placeholder="/users" value={form.path} onChange={(e) => setForm((s) => ({ ...s, path: e.target.value }))} />
              </div>
              <div className="space-y-1 md:col-span-2">
                <Label>Description</Label>
                <Input value={form.description} onChange={(e) => setForm((s) => ({ ...s, description: e.target.value }))} />
              </div>
              <div className="space-y-1">
                <Label>API Spec</Label>
                <Select value={form.api_spec_id || 'none'} onValueChange={(value) => setForm((s) => ({ ...s, api_spec_id: value === 'none' ? '' : value }))}>
                  <SelectTrigger><SelectValue placeholder="Optional" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">None</SelectItem>
                    {apiSpecs.map((spec) => (
                      <SelectItem key={spec.id} value={String(spec.id)}>{spec.method} {spec.path}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1">
                <Label>Category</Label>
                <Select value={form.category_id || 'none'} onValueChange={(value) => setForm((s) => ({ ...s, category_id: value === 'none' ? '' : value }))}>
                  <SelectTrigger><SelectValue placeholder="Optional" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">None</SelectItem>
                    {categories.map((cat) => <SelectItem key={cat.id} value={String(cat.id)}>{cat.name}</SelectItem>)}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1">
                <Label>Environment</Label>
                <Input value={form.environment} onChange={(e) => setForm((s) => ({ ...s, environment: e.target.value }))} />
              </div>
              <div className="space-y-1">
                <Label>Expected Status</Label>
                <Input value={form.expected_status} onChange={(e) => setForm((s) => ({ ...s, expected_status: e.target.value }))} />
              </div>
              <div className="space-y-1 md:col-span-2">
                <Label>Request Headers (JSON)</Label>
                <Textarea rows={4} value={form.request_headers} onChange={(e) => setForm((s) => ({ ...s, request_headers: e.target.value }))} />
              </div>
              <div className="space-y-1 md:col-span-2">
                <Label>Request Body (JSON)</Label>
                <Textarea rows={4} value={form.request_body} onChange={(e) => setForm((s) => ({ ...s, request_body: e.target.value }))} />
              </div>
              <div className="space-y-1 md:col-span-2">
                <Label>Expected Response (JSON)</Label>
                <Textarea rows={4} value={form.expected_response} onChange={(e) => setForm((s) => ({ ...s, expected_response: e.target.value }))} />
              </div>
              <div className="space-y-1 md:col-span-2">
                <Label>Variables (JSON)</Label>
                <Textarea rows={3} value={form.variables} onChange={(e) => setForm((s) => ({ ...s, variables: e.target.value }))} />
              </div>
              <div className="space-y-1 md:col-span-2">
                <Label>Setup Script</Label>
                <Textarea rows={3} value={form.setup_script} onChange={(e) => setForm((s) => ({ ...s, setup_script: e.target.value }))} />
              </div>
              <div className="space-y-1 md:col-span-2">
                <Label>Teardown Script</Label>
                <Textarea rows={3} value={form.teardown_script} onChange={(e) => setForm((s) => ({ ...s, teardown_script: e.target.value }))} />
              </div>
            </div>
          </div>
          <DialogFooter className="pt-3 border-t bg-background sticky bottom-0">
            <Button variant="outline" onClick={() => setCreateOpen(false)}>Cancel</Button>
            <Button onClick={submitForm} disabled={createMutation.isPending || updateMutation.isPending}>
              {editing ? 'Save' : 'Create'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!detailId} onOpenChange={(open) => { if (!open) setDetailId(null) }}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Test Case Detail</DialogTitle>
          </DialogHeader>
          <pre className="text-xs bg-muted p-3 rounded-md overflow-auto max-h-[60vh]">
            {JSON.stringify(detailQuery.data, null, 2)}
          </pre>
        </DialogContent>
      </Dialog>

      <Dialog open={!!duplicateTarget} onOpenChange={(open) => { if (!open) setDuplicateTarget(null) }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Duplicate Test Case</DialogTitle>
          </DialogHeader>
          <div className="space-y-3">
            <div className="space-y-1">
              <Label>New Name</Label>
              <Input value={duplicateName} onChange={(e) => setDuplicateName(e.target.value)} />
            </div>
            <div className="space-y-1">
              <Label>Environment</Label>
              <Input value={duplicateEnv} onChange={(e) => setDuplicateEnv(e.target.value)} />
            </div>
            <div className="space-y-1">
              <Label>Description</Label>
              <Input value={duplicateDescription} onChange={(e) => setDuplicateDescription(e.target.value)} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDuplicateTarget(null)}>Cancel</Button>
            <Button onClick={handleDuplicate} disabled={duplicateMutation.isPending}>Duplicate</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!runTarget} onOpenChange={(open) => { if (!open) { setRunTarget(null); setRunResult(null) } }}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Run Test Case</DialogTitle>
          </DialogHeader>
          <div className="space-y-3">
            <div className="space-y-1">
              <Label>Environment (optional)</Label>
              <Input value={runEnvironment} onChange={(e) => setRunEnvironment(e.target.value)} />
            </div>
            <div className="space-y-1">
              <Label>Variables (JSON, optional)</Label>
              <Textarea rows={4} value={runVariables} onChange={(e) => setRunVariables(e.target.value)} />
            </div>
            <label className="flex items-center gap-2 text-sm">
              <input type="checkbox" checked={runAsync} onChange={(e) => setRunAsync(e.target.checked)} />
              Run async
            </label>
            {runResult && (
              <pre className="text-xs bg-muted p-3 rounded-md overflow-auto max-h-[35vh]">
                {JSON.stringify(runResult, null, 2)}
              </pre>
            )}
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setRunTarget(null)}>Close</Button>
            <Button onClick={handleRun} disabled={runMutation.isPending}>Run</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={generateOpen} onOpenChange={setGenerateOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Generate From API Spec</DialogTitle>
            <DialogDescription>Auto-generate test cases for selected API spec.</DialogDescription>
          </DialogHeader>
          <div className="space-y-3">
            <div className="space-y-1">
              <Label>API Spec</Label>
              <Select value={generateForm.api_spec_id || 'none'} onValueChange={(value) => setGenerateForm((s) => ({ ...s, api_spec_id: value === 'none' ? '' : value }))}>
                <SelectTrigger><SelectValue placeholder="Select spec" /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">Select</SelectItem>
                  {apiSpecs.map((spec) => <SelectItem key={spec.id} value={String(spec.id)}>{spec.method} {spec.path}</SelectItem>)}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-1">
              <Label>Environment</Label>
              <Input value={generateForm.environment} onChange={(e) => setGenerateForm((s) => ({ ...s, environment: e.target.value }))} />
            </div>
            <div className="space-y-1">
              <Label>Category</Label>
              <Select value={generateForm.category_id || 'none'} onValueChange={(value) => setGenerateForm((s) => ({ ...s, category_id: value === 'none' ? '' : value }))}>
                <SelectTrigger><SelectValue placeholder="Optional" /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">None</SelectItem>
                  {categories.map((cat) => <SelectItem key={cat.id} value={String(cat.id)}>{cat.name}</SelectItem>)}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-1">
              <Label>Max Tests Per Endpoint</Label>
              <Input
                value={String(generateForm.max_tests_per_endpoint)}
                onChange={(e) => setGenerateForm((s) => ({ ...s, max_tests_per_endpoint: Number(e.target.value || 5) }))}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setGenerateOpen(false)}>Cancel</Button>
            <Button onClick={handleGenerate} disabled={generateMutation.isPending}>Generate</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
  )
}
