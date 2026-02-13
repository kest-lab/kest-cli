import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
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
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { kestApi } from '@/services/kest-api.service'
import { useCategories } from '@/hooks/use-kest-api'
import { APICategory } from '@/types/kest-api'

const apiSpecSchema = z.object({
  method: z.enum(['GET', 'POST', 'PUT', 'PATCH', 'DELETE']),
  path: z.string().min(1, 'Path is required').startsWith('/', 'Path must start with /'),
  summary: z.string().min(1, 'Summary is required'),
  description: z.string().optional(),
  category_id: z.string().optional(),
})

type APISpecFormData = z.infer<typeof apiSpecSchema>

interface CreateAPISpecDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  projectId: number
}

export function CreateAPISpecDialog({ open, onOpenChange, projectId }: CreateAPISpecDialogProps) {
  const queryClient = useQueryClient()
  const { data: categoriesData } = useCategories(projectId)
  const categories: APICategory[] = (() => {
    if (!categoriesData) return []
    if (Array.isArray(categoriesData)) return categoriesData
    if (Array.isArray((categoriesData as any)?.items)) return (categoriesData as any).items
    if (Array.isArray((categoriesData as any)?.data)) return (categoriesData as any).data
    return []
  })()

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    setValue,
  } = useForm<APISpecFormData>({
    resolver: zodResolver(apiSpecSchema),
    defaultValues: {
      method: 'GET',
    },
  })

  const createMutation = useMutation({
    mutationFn: async (data: APISpecFormData) => {
      const categoryId = data.category_id ? parseInt(data.category_id, 10) : undefined
      return await kestApi.apiSpec.create(projectId, {
        ...data,
        category_id: categoryId && categoryId > 0 ? categoryId : undefined,
        version: '1.0.0',
      })
    },
    onSuccess: () => {
      toast.success('API specification created successfully!')
      queryClient.invalidateQueries({ queryKey: ['api-specs', projectId] })
      reset()
      onOpenChange(false)
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to create API specification')
    },
  })

  const onSubmit = (data: APISpecFormData) => {
    createMutation.mutate(data)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle>Create API Specification</DialogTitle>
          <DialogDescription>
            Define a new API endpoint for your project (手动创建)
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 gap-4">
              <div className="col-span-1">
                <Label htmlFor="method">Method *</Label>
                <Select
                  onValueChange={(value) => setValue('method', value as any)}
                  defaultValue="GET"
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Method" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="GET">GET</SelectItem>
                    <SelectItem value="POST">POST</SelectItem>
                    <SelectItem value="PUT">PUT</SelectItem>
                    <SelectItem value="PATCH">PATCH</SelectItem>
                    <SelectItem value="DELETE">DELETE</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="col-span-3">
                <Label htmlFor="path">Path *</Label>
                <Input
                  id="path"
                  placeholder="/api/users"
                  {...register('path')}
                />
                {errors.path && (
                  <p className="text-sm text-red-600 mt-1">{errors.path.message}</p>
                )}
              </div>
            </div>

            <div className="grid gap-2">
              <Label htmlFor="summary">Summary *</Label>
              <Input
                id="summary"
                placeholder="Get list of users"
                {...register('summary')}
              />
              {errors.summary && (
                <p className="text-sm text-red-600">{errors.summary.message}</p>
              )}
            </div>

            <div className="grid gap-2">
              <Label htmlFor="category">Category</Label>
              <Select onValueChange={(value) => setValue('category_id', value)}>
                <SelectTrigger>
                  <SelectValue placeholder="Select a category" />
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

            <div className="grid gap-2">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                placeholder="Detailed description of the API endpoint..."
                rows={4}
                {...register('description')}
              />
            </div>

            <div className="text-sm text-muted-foreground">
              💡 提示：后续可以通过 Kest CLI 自动导入 API 规范
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                reset()
                onOpenChange(false)
              }}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={createMutation.isPending}>
              {createMutation.isPending ? 'Creating...' : 'Create API Spec'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
