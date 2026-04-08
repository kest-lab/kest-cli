/**
 * @component Skeleton
 * @category UI
 * @status Stable
 * @description Provides a placeholder while content is loading, typically used for cards, lists, or headers.
 * @usage Use to reduce perceived latency by showing the structure of the content before it loads.
 * @example
 * <div className="space-y-4">
 *   <Skeleton className="h-4 w-[250px]" />
 *   <Skeleton className="h-4 w-[200px]" />
 * </div>
 */
import { cn } from "@/utils"

function Skeleton({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="skeleton"
      className={cn("bg-accent animate-pulse rounded-md", className)}
      {...props}
    />
  )
}

export { Skeleton }
