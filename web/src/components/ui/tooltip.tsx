/**
 * @component Tooltip
 * @category UI
 * @status Stable
 * @description A popup that displays information related to an element when the element receives keyboard focus or the mouse hovers over it.
 * @usage Use to provide brief, descriptive labels or information for icons or buttons.
 * @example
 * <Tooltip>
 *   <TooltipTrigger asChild>
 *     <Button variant="outline">Hover</Button>
 *   </TooltipTrigger>
 *   <TooltipContent>
 *     <p>Add to library</p>
 *   </TooltipContent>
 * </Tooltip>
 */
"use client"

import * as React from "react"
import * as TooltipPrimitive from "@radix-ui/react-tooltip"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/utils"

function TooltipProvider({
  delayDuration = 0,
  ...props
}: React.ComponentProps<typeof TooltipPrimitive.Provider>) {
  return (
    <TooltipPrimitive.Provider
      data-slot="tooltip-provider"
      delayDuration={delayDuration}
      {...props}
    />
  )
}

function Tooltip({
  ...props
}: React.ComponentProps<typeof TooltipPrimitive.Root>) {
  return (
    <TooltipProvider>
      <TooltipPrimitive.Root data-slot="tooltip" {...props} />
    </TooltipProvider>
  )
}

function TooltipTrigger({
  ...props
}: React.ComponentProps<typeof TooltipPrimitive.Trigger>) {
  return <TooltipPrimitive.Trigger data-slot="tooltip-trigger" {...props} />
}

const tooltipVariants = cva(
  "z-50 w-fit origin-(--radix-tooltip-content-transform-origin) rounded-lg px-4 py-2 text-xs font-medium shadow-tooltip transition-all duration-500 " +
  "data-[state=delayed-open]:data-[side=top]:animate-[float-in-top_0.4s_ease-out] " +
  "data-[state=delayed-open]:data-[side=bottom]:animate-[float-in-bottom_0.4s_ease-out] " +
  "data-[state=delayed-open]:data-[side=left]:animate-[float-in-left_0.4s_ease-out] " +
  "data-[state=delayed-open]:data-[side=right]:animate-[float-in-right_0.4s_ease-out] " +
  "data-[state=closed]:data-[side=top]:animate-[float-out-top_0.2s_ease-in] " +
  "data-[state=closed]:data-[side=bottom]:animate-[float-out-bottom_0.2s_ease-in] " +
  "data-[state=closed]:data-[side=left]:animate-[float-out-left_0.2s_ease-in] " +
  "data-[state=closed]:data-[side=right]:animate-[float-out-right_0.2s_ease-in] " +
  "data-[state=closed]:opacity-0",
  {
    variants: {
      variant: {
        default: "bg-popover/90 text-popover-foreground backdrop-blur-md border",
        inverted: "bg-zinc-950/90 text-zinc-50 backdrop-blur-md border border-white/10 shadow-2xl",
        primary: "bg-primary/90 text-primary-foreground backdrop-blur-md border border-primary/20 shadow-lg shadow-primary/10",
        destructive: "bg-destructive/10 text-destructive backdrop-blur-md border border-destructive/20 shadow-lg shadow-destructive/5",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
)

function TooltipContent({
  className,
  variant,
  sideOffset = 10,
  children,
  ...props
}: React.ComponentProps<typeof TooltipPrimitive.Content> & VariantProps<typeof tooltipVariants>) {
  return (
    <TooltipPrimitive.Portal>
      <TooltipPrimitive.Content
        data-slot="tooltip-content"
        sideOffset={sideOffset}
        className={cn(tooltipVariants({ variant }), className)}
        {...props}
      >
        {children}
      </TooltipPrimitive.Content>
    </TooltipPrimitive.Portal>
  )
}

export { Tooltip, TooltipTrigger, TooltipContent, TooltipProvider }
