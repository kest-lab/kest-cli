/**
 * @component Toaster
 * @category UI
 * @status Stable
 * @description A toast notification component powered by Sonner, with customized themes for the scaffold.
 * @usage Include at the root of the application to enable global toast notifications.
 * @example
 * <Toaster />
 * // Usage in code:
 * import { toast } from "sonner";
 * toast.success("Operation successful");
 */
"use client"

import { useTheme } from "@/providers/theme-provider"
import { Toaster as Sonner, ToasterProps } from "sonner"

const Toaster = ({ ...props }: ToasterProps) => {
  const { resolvedTheme } = useTheme()

  return (
    <Sonner
      theme={(resolvedTheme ?? "light") as ToasterProps["theme"]}
      className="toaster group"
      style={
        {
          "--normal-bg": "var(--popover)",
          "--normal-text": "var(--popover-foreground)",
          "--normal-border": "var(--border)",
        } as React.CSSProperties
      }
      {...props}
    />
  )
}

export { Toaster }
