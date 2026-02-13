import { useState } from 'react'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'
import { toast } from 'sonner'
import { kestApi } from '@/services/kest-api.service'
import { CreateAPISpecRequest, Parameter, ResponseSpec, APICategory } from '@/types/kest-api'
import { FileUp, AlertCircle } from 'lucide-react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useCategories } from '@/hooks/use-kest-api'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

interface ImportAPISpecDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    projectId: number
}

export function ImportAPISpecDialog({ open, onOpenChange, projectId }: ImportAPISpecDialogProps) {
    const queryClient = useQueryClient()
    const [format, setFormat] = useState<'openapi' | 'swagger' | 'postman'>('openapi')
    const [content, setContent] = useState('')
    const [isParsing, setIsParsing] = useState(false)
    const [targetCategoryId, setTargetCategoryId] = useState<number | null>(null)

    const { data: categoriesData } = useCategories(projectId)
    const categories: APICategory[] = (() => {
        if (!categoriesData) return []
        if (Array.isArray(categoriesData)) return categoriesData
        if (Array.isArray((categoriesData as any)?.items)) return (categoriesData as any).items
        if (Array.isArray((categoriesData as any)?.data)) return (categoriesData as any).data
        return []
    })()

    const importMutation = useMutation({
        mutationFn: (specs: CreateAPISpecRequest[]) => kestApi.apiSpec.import(projectId, specs),
        onSuccess: () => {
            toast.success('APIs imported successfully!')
            queryClient.invalidateQueries({ queryKey: ['api-specs', projectId] })
            onOpenChange(false)
            setContent('')
        },
        onError: (error: any) => {
            toast.error(error.message || 'Failed to import APIs')
        },
    })

    const parseOpenAPI = (data: any): CreateAPISpecRequest[] => {
        const specs: CreateAPISpecRequest[] = []
        const paths = data.paths || {}

        Object.entries(paths).forEach(([path, methods]: [string, any]) => {
            Object.entries(methods).forEach(([method, detail]: [string, any]) => {
                // Skip non-HTTP methods and metadata
                if (!['get', 'post', 'put', 'delete', 'patch', 'head', 'options'].includes(method.toLowerCase())) return

                const spec: CreateAPISpecRequest = {
                    project_id: projectId,
                    method: method.toUpperCase(),
                    path: path,
                    summary: detail.summary || detail.operationId || `${method.toUpperCase()} ${path}`,
                    description: detail.description || '',
                    tags: detail.tags || [],
                    version: data.info?.version || '1.0.0',
                    category_id: targetCategoryId || undefined,
                    parameters: [] as Parameter[],
                    responses: {} as Record<string, ResponseSpec>,
                }

                // Parse Parameters (query, path, header, cookie all go into single array)
                const rawParams = [...(detail.parameters || []), ...(methods.parameters || [])]
                rawParams.forEach((p: any) => {
                    const param: Parameter = {
                        name: p.name,
                        in: p.in,
                        required: !!p.required,
                        description: p.description,
                        schema: p.schema || { type: 'string' },
                        example: p.example,
                    }
                    spec.parameters!.push(param)
                })

                // Parse Request Body (OpenAPI 3.0)
                if (detail.requestBody) {
                    const content = detail.requestBody.content || {}
                    const firstType = Object.keys(content)[0]
                    if (firstType) {
                        spec.request_body = {
                            description: detail.requestBody.description,
                            required: !!detail.requestBody.required,
                            content_type: firstType,
                            schema: content[firstType].schema || { type: 'object' },
                        }
                    }
                }

                // Parse Responses
                if (detail.responses) {
                    Object.entries(detail.responses).forEach(([code, r]: [string, any]) => {
                        const content = r.content || {}
                        const firstType = Object.keys(content)[0]
                        spec.responses![code] = {
                            description: r.description || '',
                            content_type: firstType || 'application/json',
                            schema: content[firstType]?.schema || { type: 'object' },
                        }
                    })
                }

                specs.push(spec)
            })
        })

        return specs
    }

    const handleImport = async () => {
        if (!content.trim()) return
        setIsParsing(true)
        try {
            let data: any
            try {
                data = JSON.parse(content)
            } catch (e) {
                throw new Error('Invalid JSON content. Currently only JSON format is supported.')
            }

            let specs: CreateAPISpecRequest[] = []
            if (format === 'openapi' || format === 'swagger') {
                specs = parseOpenAPI(data)
            } else {
                throw new Error('Format not supported yet')
            }

            if (specs.length === 0) {
                throw new Error('No valid API endpoints found in the content')
            }

            importMutation.mutate(specs)
        } catch (err: any) {
            toast.error(err.message)
        } finally {
            setIsParsing(false)
        }
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[700px]">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2">
                        <FileUp className="w-5 h-5" />
                        Import APIs
                    </DialogTitle>
                    <DialogDescription>
                        Import multiple API specifications from Swagger, OpenAPI or Postman collections.
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-4 py-4">
                    <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="format" className="text-right">Source Format</Label>
                        <Select value={format} onValueChange={(v: any) => setFormat(v)}>
                            <SelectTrigger className="col-span-3">
                                <SelectValue placeholder="Select format" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="openapi">OpenAPI 3.0 (JSON)</SelectItem>
                                <SelectItem value="swagger">Swagger 2.0 (JSON)</SelectItem>
                                <SelectItem value="postman" disabled>Postman Collection (Coming soon)</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="category" className="text-right">Assign Category</Label>
                        <Select onValueChange={(v) => setTargetCategoryId(v === '0' ? null : parseInt(v))}>
                            <SelectTrigger className="col-span-3">
                                <SelectValue placeholder="Default Category (Optional)" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="0">None</SelectItem>
                                {categories.map((cat: APICategory) => (
                                    <SelectItem key={cat.id} value={cat.id.toString()}>
                                        {cat.name}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="space-y-2">
                        <div className="flex items-center justify-between">
                            <Label htmlFor="content">JSON Content</Label>
                            <Button variant="ghost" size="sm" onClick={() => setContent('')} className="text-xs h-7">
                                Clear
                            </Button>
                        </div>
                        <Textarea
                            id="content"
                            placeholder='Paste your JSON content here...'
                            value={content}
                            onChange={(e) => setContent(e.target.value)}
                            className="font-mono text-xs min-h-[300px]"
                        />
                    </div>

                    <Alert variant="default" className="bg-blue-50 border-blue-200">
                        <AlertCircle className="h-4 w-4 text-blue-600" />
                        <AlertTitle className="text-blue-800 text-xs font-bold">Pro Tip</AlertTitle>
                        <AlertDescription className="text-blue-700 text-xs">
                            Existing APIs with the same method and path will be updated. Others will be created.
                        </AlertDescription>
                    </Alert>
                </div>

                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>Cancel</Button>
                    <Button
                        onClick={handleImport}
                        disabled={importMutation.isPending || isParsing || !content}
                    >
                        {importMutation.isPending ? 'Importing...' : 'Start Import'}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}
