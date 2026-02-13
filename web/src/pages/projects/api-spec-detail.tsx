import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import {
    ArrowLeft,
    Settings,
    Code,
    Database,
    Clock,
    Tag,
    Globe,
    Copy,
    Plus
} from 'lucide-react'
import { useAPISpecWithExamples } from '@/hooks/use-kest-api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { toast } from 'sonner'

export function APISpecDetailPage() {
    const { id, sid } = useParams<{ id: string, sid: string }>()
    const projectId = parseInt(id || '0')
    const specId = parseInt(sid || '0')

    const { data: spec, isLoading } = useAPISpecWithExamples(projectId, specId)

    const [activeTab, setActiveTab] = useState('definition')

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

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text)
        toast.success('Copied to clipboard')
    }

    if (isLoading) {
        return (
            <div className="container mx-auto p-8">
                <Skeleton className="h-8 w-64 mb-4" />
                <Skeleton className="h-4 w-96 mb-8" />
                <Skeleton className="h-[400px]" />
            </div>
        )
    }

    if (!spec) {
        return (
            <div className="container mx-auto p-8 text-center">
                <h2 className="text-2xl font-bold mb-2">API Specification not found</h2>
                <Link to={`/projects/${projectId}`}>
                    <Button>Back to Project</Button>
                </Link>
            </div>
        )
    }

    return (
        <div className="container mx-auto p-8">
            <div className="mb-6">
                <Link to={`/projects/${projectId}`}>
                    <Button variant="ghost" className="mb-4">
                        <ArrowLeft className="w-4 h-4 mr-2" />
                        Back to Project
                    </Button>
                </Link>

                <div className="flex items-start justify-between">
                    <div className="flex-1">
                        <div className="flex items-center gap-3 mb-2">
                            <Badge className={`${getMethodColor(spec.method)} text-sm px-3 py-1 font-bold`}>
                                {spec.method}
                            </Badge>
                            <code className="text-xl font-mono bg-muted px-3 py-1 rounded-md">
                                {spec.path}
                            </code>
                            <Button variant="ghost" size="icon" onClick={() => copyToClipboard(spec.path)}>
                                <Copy className="w-4 h-4" />
                            </Button>
                        </div>
                        <h1 className="text-4xl font-bold mb-2">{spec.summary || spec.path}</h1>
                        {spec.description && <p className="text-gray-600 mb-4">{spec.description}</p>}

                        <div className="flex items-center gap-4 text-sm text-muted-foreground">
                            <span className="flex items-center gap-1">
                                <Clock className="w-4 h-4" />
                                Updated {new Date(spec.updated_at).toLocaleDateString()}
                            </span>
                            <span className="flex items-center gap-1">
                                <Tag className="w-4 h-4" />
                                {spec.tags?.join(', ') || 'No tags'}
                            </span>
                            <span className="flex items-center gap-1">
                                <Globe className="w-4 h-4" />
                                {spec.is_public ? 'Public' : 'Private'}
                            </span>
                        </div>
                    </div>

                    <div className="flex gap-2">
                        <Button variant="outline">
                            <Code className="w-4 h-4 mr-2" />
                            Test API
                        </Button>
                        <Button variant="outline" size="icon">
                            <Settings className="w-4 h-4" />
                        </Button>
                    </div>
                </div>
            </div>

            <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
                <TabsList className="bg-muted p-1 rounded-lg">
                    <TabsTrigger value="definition">Definition</TabsTrigger>
                    <TabsTrigger value="examples">Examples ({spec.examples?.length || 0})</TabsTrigger>
                    <TabsTrigger value="tests">Test Cases</TabsTrigger>
                    <TabsTrigger value="history">History</TabsTrigger>
                </TabsList>

                <TabsContent value="definition" className="space-y-6 animate-in fade-in duration-300">
                    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                        <div className="lg:col-span-2 space-y-6">
                            {/* Description */}
                            <Card>
                                <CardHeader>
                                    <CardTitle className="text-lg">Description</CardTitle>
                                </CardHeader>
                                <CardContent>
                                    <div className="prose prose-sm max-w-none">
                                        {spec.description || 'No detailed description provided.'}
                                    </div>
                                </CardContent>
                            </Card>

                            {/* Request Parameters */}
                            <Card>
                                <CardHeader>
                                    <CardTitle className="text-lg">Request Parameters</CardTitle>
                                </CardHeader>
                                <CardContent className="space-y-6">
                                    {spec.parameters?.filter(p => p.in === 'path').length ? (
                                        <div className="space-y-3">
                                            <h4 className="text-sm font-bold flex items-center gap-2">
                                                <Badge variant="outline">Path Params</Badge>
                                            </h4>
                                            <ParameterTable parameters={spec.parameters.filter(p => p.in === 'path')} />
                                        </div>
                                    ) : null}

                                    {spec.parameters?.filter(p => p.in === 'query').length ? (
                                        <div className="space-y-3">
                                            <h4 className="text-sm font-bold flex items-center gap-2">
                                                <Badge variant="outline">Query Params</Badge>
                                            </h4>
                                            <ParameterTable parameters={spec.parameters.filter(p => p.in === 'query')} />
                                        </div>
                                    ) : null}

                                    {spec.parameters?.filter(p => p.in === 'header').length ? (
                                        <div className="space-y-3">
                                            <h4 className="text-sm font-bold flex items-center gap-2">
                                                <Badge variant="outline">Headers</Badge>
                                            </h4>
                                            <ParameterTable parameters={spec.parameters.filter(p => p.in === 'header')} />
                                        </div>
                                    ) : null}

                                    {(!spec.parameters || spec.parameters.length === 0) && (
                                        <p className="text-sm text-center text-muted-foreground py-4">No parameters defined</p>
                                    )}
                                </CardContent>
                            </Card>

                            {/* Request Body */}
                            {spec.request_body && (
                                <Card>
                                    <CardHeader>
                                        <div className="flex items-center justify-between">
                                            <CardTitle className="text-lg">Request Body</CardTitle>
                                            <Badge variant="secondary">{spec.request_body.content_type}</Badge>
                                        </div>
                                        {spec.request_body.description && (
                                            <CardDescription>{spec.request_body.description}</CardDescription>
                                        )}
                                    </CardHeader>
                                    <CardContent>
                                        <div className="bg-muted p-4 rounded-md overflow-x-auto">
                                            <pre className="text-xs font-mono">
                                                {JSON.stringify(spec.request_body.schema, null, 2)}
                                            </pre>
                                        </div>
                                    </CardContent>
                                </Card>
                            )}
                        </div>

                        <div className="space-y-6">
                            {/* Status & Metadata */}
                            <Card>
                                <CardHeader>
                                    <CardTitle className="text-lg">Metadata</CardTitle>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <div className="flex items-center justify-between">
                                        <span className="text-sm text-muted-foreground">Status</span>
                                        <Badge variant={spec.status === 'done' ? 'default' : 'outline'}>
                                            {(spec.status || 'undone').toUpperCase()}
                                        </Badge>
                                    </div>
                                    <div className="flex items-center justify-between">
                                        <span className="text-sm text-muted-foreground">Version</span>
                                        <span className="text-sm font-medium">{spec.version}</span>
                                    </div>
                                    <div className="flex items-center justify-between">
                                        <span className="text-sm text-muted-foreground">Category</span>
                                        <span className="text-sm font-medium">Core API</span>
                                    </div>
                                </CardContent>
                            </Card>

                            {/* Responses */}
                            <Card>
                                <CardHeader>
                                    <CardTitle className="text-lg">Responses</CardTitle>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    {spec.responses && Object.entries(spec.responses).map(([code, resp]) => (
                                        <div key={code} className="p-3 border rounded-md hover:bg-accent cursor-pointer transition-colors">
                                            <div className="flex items-center justify-between mb-2">
                                                <Badge className={parseInt(code) < 300 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}>
                                                    {code}
                                                </Badge>
                                                <span className="text-xs text-muted-foreground">{resp.content_type}</span>
                                            </div>
                                            <p className="text-sm font-medium line-clamp-1">{resp.description}</p>
                                        </div>
                                    ))}
                                    {!spec.responses && (
                                        <p className="text-sm text-center text-muted-foreground py-2">No responses defined</p>
                                    )}
                                </CardContent>
                            </Card>
                        </div>
                    </div>
                </TabsContent>

                <TabsContent value="examples" className="mt-6">
                    <div className="grid grid-cols-1 gap-6">
                        {spec.examples && spec.examples.length > 0 ? (
                            spec.examples.map((example) => (
                                <Card key={example.id}>
                                    <CardHeader>
                                        <div className="flex items-center justify-between">
                                            <CardTitle className="text-lg">{example.name}</CardTitle>
                                            <Badge className={example.response_status < 400 ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}>
                                                {example.response_status} - {example.duration_ms}ms
                                            </Badge>
                                        </div>
                                    </CardHeader>
                                    <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                        <div className="space-y-2">
                                            <h4 className="text-sm font-bold">Request Body</h4>
                                            <pre className="p-3 bg-muted rounded-md text-xs font-mono overflow-x-auto max-h-[300px]">
                                                {JSON.stringify(example.request_body, null, 2)}
                                            </pre>
                                        </div>
                                        <div className="space-y-2">
                                            <h4 className="text-sm font-bold">Response Body</h4>
                                            <pre className="p-3 bg-muted rounded-md text-xs font-mono overflow-x-auto max-h-[300px]">
                                                {JSON.stringify(example.response_body, null, 2)}
                                            </pre>
                                        </div>
                                    </CardContent>
                                </Card>
                            ))
                        ) : (
                            <div className="text-center py-20 border-2 border-dashed rounded-lg">
                                <Database className="w-12 h-12 text-muted-foreground mx-auto mb-4 opacity-20" />
                                <h3 className="text-lg font-medium">No examples yet</h3>
                                <p className="text-muted-foreground mb-4">Run tests or record traffic to see real-world examples</p>
                                <Button>
                                    <Plus className="w-4 h-4 mr-2" />
                                    Add Manual Example
                                </Button>
                            </div>
                        )}
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    )
}

function ParameterTable({ parameters }: { parameters: any[] }) {
    return (
        <div className="border rounded-md overflow-hidden">
            <table className="w-full text-sm">
                <thead className="bg-muted">
                    <tr>
                        <th className="px-4 py-2 text-left font-medium">Name</th>
                        <th className="px-4 py-2 text-left font-medium">Type</th>
                        <th className="px-4 py-2 text-left font-medium">Status</th>
                        <th className="px-4 py-2 text-left font-medium">Description</th>
                    </tr>
                </thead>
                <tbody className="divide-y">
                    {parameters.map((p, idx) => (
                        <tr key={idx} className="hover:bg-accent/50 transition-colors">
                            <td className="px-4 py-3 font-mono text-xs">{p.name}</td>
                            <td className="px-4 py-3">
                                <Badge variant="outline" className="text-[10px] uppercase">{p.schema?.type || 'string'}</Badge>
                            </td>
                            <td className="px-4 py-3">
                                {p.required ? <Badge className="bg-orange-100 text-orange-700 text-[10px]">Required</Badge> : <span className="text-xs text-muted-foreground">Optional</span>}
                            </td>
                            <td className="px-4 py-3 text-muted-foreground">{p.description || '-'}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}
