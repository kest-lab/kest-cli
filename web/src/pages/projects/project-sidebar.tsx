import { useMemo } from 'react'
import { Link } from 'react-router-dom'
import { ChevronRight, FolderTree, FileText } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import type { APISpec, CategoryTree } from '@/types/kest-api'

interface ProjectSidebarProps {
  projectId: number
  categories: CategoryTree[]
  apiSpecs: APISpec[]
  isLoading?: boolean
  categoriesLoading?: boolean
}

type FlatCategory = {
  category: CategoryTree
  depth: number
}

const METHOD_COLORS: Record<string, string> = {
  GET: 'bg-blue-100 text-blue-800',
  POST: 'bg-green-100 text-green-800',
  PUT: 'bg-yellow-100 text-yellow-800',
  PATCH: 'bg-orange-100 text-orange-800',
  DELETE: 'bg-red-100 text-red-800',
}

const flattenCategories = (items: CategoryTree[], depth = 0, acc: FlatCategory[] = []) => {
  items.forEach((category) => {
    acc.push({ category, depth })
    if (category.children && category.children.length > 0) {
      flattenCategories(category.children, depth + 1, acc)
    }
  })
  return acc
}

export function ProjectSidebar({
  projectId,
  categories,
  apiSpecs,
  isLoading,
  categoriesLoading,
}: ProjectSidebarProps) {
  const flatCategories = useMemo(() => flattenCategories(categories), [categories])

  const specsByCategory = useMemo(() => {
    const map = new Map<number, APISpec[]>()
    apiSpecs.forEach((spec) => {
      if (!spec.category_id) return
      const bucket = map.get(spec.category_id) || []
      bucket.push(spec)
      map.set(spec.category_id, bucket)
    })
    return map
  }, [apiSpecs])

  const uncategorizedSpecs = useMemo(
    () => apiSpecs.filter((spec) => !spec.category_id),
    [apiSpecs]
  )

  const renderSpecList = (specs: APISpec[]) => {
    if (specs.length === 0) {
      return <p className="text-xs text-muted-foreground">No API specifications yet.</p>
    }

    return (
      <div className="space-y-1">
        {specs
          .slice()
          .sort((a, b) => a.path.localeCompare(b.path))
          .map((spec) => (
            <Link
              key={spec.id}
              to={`/projects/${projectId}/api-specs/${spec.id}`}
              className="flex items-center gap-2 text-xs text-muted-foreground hover:text-foreground"
            >
              <Badge className={`${METHOD_COLORS[spec.method] || 'bg-gray-100 text-gray-800'} text-[10px]`}>
                {spec.method}
              </Badge>
              <span className="font-mono truncate">{spec.path}</span>
            </Link>
          ))}
      </div>
    )
  }

  return (
    <Card className="h-fit lg:sticky lg:top-6">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg">Project Sidebar</CardTitle>
        <CardDescription>Categories and API specs</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {(categoriesLoading || isLoading) && (
          <div className="space-y-3">
            {[1, 2, 3, 4].map((i) => (
              <Skeleton key={i} className="h-6 w-full" />
            ))}
          </div>
        )}

        {!categoriesLoading && !isLoading && flatCategories.length === 0 && uncategorizedSpecs.length === 0 && (
          <div className="text-sm text-muted-foreground">
            No categories or API specifications yet.
          </div>
        )}

        {!categoriesLoading && !isLoading && (flatCategories.length > 0 || uncategorizedSpecs.length > 0) && (
          <div className="space-y-3 max-h-[65vh] overflow-auto pr-1">
            {flatCategories.map(({ category, depth }) => {
              const specs = specsByCategory.get(category.id) || []
              return (
                <details key={category.id} className="group rounded-md border border-border px-2 py-1">
                  <summary
                    className="flex items-center justify-between gap-2 cursor-pointer select-none list-none [&::-webkit-details-marker]:hidden"
                  >
                    <div className="flex items-center gap-2 py-1" style={{ paddingLeft: `${depth * 12}px` }}>
                      <ChevronRight className="h-4 w-4 text-muted-foreground transition-transform group-open:rotate-90" />
                      <FolderTree className="h-4 w-4 text-muted-foreground" />
                      <span className="text-sm font-medium">{category.name}</span>
                    </div>
                    <Badge variant="outline" className="text-[10px]">
                      {specs.length}
                    </Badge>
                  </summary>
                  <div className="mt-2 space-y-2 pl-6">
                    {renderSpecList(specs)}
                  </div>
                </details>
              )
            })}

            {uncategorizedSpecs.length > 0 && (
              <details className="group rounded-md border border-dashed border-border px-2 py-1">
                <summary
                  className="flex items-center justify-between gap-2 cursor-pointer select-none list-none [&::-webkit-details-marker]:hidden"
                >
                  <div className="flex items-center gap-2 py-1">
                    <ChevronRight className="h-4 w-4 text-muted-foreground transition-transform group-open:rotate-90" />
                    <FileText className="h-4 w-4 text-muted-foreground" />
                    <span className="text-sm font-medium">Uncategorized</span>
                  </div>
                  <Badge variant="outline" className="text-[10px]">
                    {uncategorizedSpecs.length}
                  </Badge>
                </summary>
                <div className="mt-2 space-y-2 pl-6">
                  {renderSpecList(uncategorizedSpecs)}
                </div>
              </details>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
