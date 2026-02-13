import { useState, useEffect } from 'react'
import { X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import type { StepNodeData } from './step-node'

interface StepPanelProps {
  data: StepNodeData
  onChange: (data: Partial<StepNodeData>) => void
  onClose: () => void
  onDelete: () => void
}

export function StepPanel({ data, onChange, onClose, onDelete }: StepPanelProps) {
  const [name, setName] = useState(data.name)
  const [method, setMethod] = useState(data.method)
  const [url, setUrl] = useState(data.url)
  const [headers, setHeaders] = useState(data.headers || '')
  const [body, setBody] = useState(data.body || '')
  const [captures, setCaptures] = useState(data.captures || '')
  const [asserts, setAsserts] = useState(data.asserts || '')

  useEffect(() => {
    setName(data.name)
    setMethod(data.method)
    setUrl(data.url)
    setHeaders(data.headers || '')
    setBody(data.body || '')
    setCaptures(data.captures || '')
    setAsserts(data.asserts || '')
  }, [data.stepId])

  const handleSave = () => {
    onChange({ name, method, url, headers, body, captures, asserts })
  }

  return (
    <div className="w-[380px] border-l border-border bg-card h-full flex flex-col">
      <div className="flex items-center justify-between px-4 py-3 border-b border-border">
        <h3 className="font-semibold text-sm">Step Configuration</h3>
        <Button variant="ghost" size="icon" className="h-7 w-7" onClick={onClose}>
          <X className="h-4 w-4" />
        </Button>
      </div>

      <div className="flex-1 overflow-auto p-4 space-y-4">
        <div className="space-y-2">
          <Label className="text-xs">Name</Label>
          <Input
            value={name}
            onChange={e => setName(e.target.value)}
            onBlur={handleSave}
            placeholder="Step name"
            className="h-8 text-sm"
          />
        </div>

        <div className="grid grid-cols-[100px_1fr] gap-2">
          <div className="space-y-2">
            <Label className="text-xs">Method</Label>
            <Select value={method} onValueChange={v => { setMethod(v); setTimeout(handleSave, 0) }}>
              <SelectTrigger className="h-8 text-sm">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {['GET', 'POST', 'PUT', 'PATCH', 'DELETE'].map(m => (
                  <SelectItem key={m} value={m}>{m}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label className="text-xs">URL</Label>
            <Input
              value={url}
              onChange={e => setUrl(e.target.value)}
              onBlur={handleSave}
              placeholder="/v1/endpoint"
              className="h-8 text-sm font-mono"
            />
          </div>
        </div>

        <Tabs defaultValue="headers" className="w-full">
          <TabsList className="w-full h-8">
            <TabsTrigger value="headers" className="text-xs flex-1">Headers</TabsTrigger>
            <TabsTrigger value="body" className="text-xs flex-1">Body</TabsTrigger>
            <TabsTrigger value="captures" className="text-xs flex-1">Captures</TabsTrigger>
            <TabsTrigger value="asserts" className="text-xs flex-1">Asserts</TabsTrigger>
          </TabsList>

          <TabsContent value="headers" className="mt-3">
            <Textarea
              value={headers}
              onChange={e => setHeaders(e.target.value)}
              onBlur={handleSave}
              placeholder='{"Content-Type": "application/json", "Authorization": "Bearer {{token}}"}'
              className="font-mono text-xs min-h-[120px] resize-none"
            />
            <p className="text-[10px] text-muted-foreground mt-1">JSON format. Use {"{{variable}}"} for dynamic values.</p>
          </TabsContent>

          <TabsContent value="body" className="mt-3">
            <Textarea
              value={body}
              onChange={e => setBody(e.target.value)}
              onBlur={handleSave}
              placeholder='{"username": "test", "password": "pass"}'
              className="font-mono text-xs min-h-[120px] resize-none"
            />
            <p className="text-[10px] text-muted-foreground mt-1">JSON request body. Use {"{{variable}}"} for captured values.</p>
          </TabsContent>

          <TabsContent value="captures" className="mt-3">
            <Textarea
              value={captures}
              onChange={e => setCaptures(e.target.value)}
              onBlur={handleSave}
              placeholder={'access_token: data.access_token\nuser_id: data.user.id'}
              className="font-mono text-xs min-h-[120px] resize-none"
            />
            <p className="text-[10px] text-muted-foreground mt-1">One per line: variable_name: json.path</p>
          </TabsContent>

          <TabsContent value="asserts" className="mt-3">
            <Textarea
              value={asserts}
              onChange={e => setAsserts(e.target.value)}
              onBlur={handleSave}
              placeholder={'status == 200\nbody.code == 0\nbody.data.id exists\nduration < 1000ms'}
              className="font-mono text-xs min-h-[120px] resize-none"
            />
            <p className="text-[10px] text-muted-foreground mt-1">One per line: status ==, body.x ==, exists, duration {"<"}</p>
          </TabsContent>
        </Tabs>
      </div>

      <div className="p-4 border-t border-border flex justify-between">
        <Button variant="destructive" size="sm" onClick={onDelete}>
          Delete Step
        </Button>
        <Button size="sm" onClick={handleSave}>
          Apply
        </Button>
      </div>
    </div>
  )
}
