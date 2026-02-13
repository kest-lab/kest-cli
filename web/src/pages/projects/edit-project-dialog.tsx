import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { toast } from 'sonner'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { kestApi } from '@/services/kest-api.service'
import { Project } from '@/types/kest-api'

const projectSchema = z.object({
  name: z.string().min(1, 'Project name is required').max(100),
  description: z.string().optional(),
  platform: z.string().optional(),
})

type ProjectFormData = z.infer<typeof projectSchema>

interface EditProjectDialogProps {
  project: Project | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function EditProjectDialog({ project, open, onOpenChange }: EditProjectDialogProps) {
  const queryClient = useQueryClient()
  
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    setValue,
  } = useForm<ProjectFormData>({
    resolver: zodResolver(projectSchema),
  })

  // Set form values when project changes
  useEffect(() => {
    if (project) {
      setValue('name', project.name)
      setValue('description', project.description || '')
      setValue('platform', project.platform || '')
    }
  }, [project, setValue])

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: ProjectFormData }) => 
      kestApi.project.update(id, data),
    onSuccess: () => {
      toast.success('Project updated successfully!')
      queryClient.invalidateQueries({ queryKey: ['projects'] })
      onOpenChange(false)
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to update project')
    },
  })

  const onSubmit = (data: ProjectFormData) => {
    if (!project) return
    updateMutation.mutate({ id: project.id, data })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Edit Project</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="name">Project Name *</Label>
              <Input
                id="name"
                placeholder="My API Project"
                {...register('name')}
              />
              {errors.name && (
                <p className="text-sm text-red-600">{errors.name.message}</p>
              )}
            </div>
            
            <div className="grid gap-2">
              <Label htmlFor="platform">Platform</Label>
              <Input
                id="platform"
                placeholder="e.g. go, javascript, python"
                {...register('platform')}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                placeholder="Describe your API project..."
                rows={3}
                {...register('description')}
              />
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
            <Button type="submit" disabled={updateMutation.isPending}>
              {updateMutation.isPending ? 'Saving...' : 'Save Changes'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
