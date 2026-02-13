import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/utils"

const textareaVariants = cva(
    "placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground flex min-h-[80px] w-full rounded-lg px-3 py-2 text-base shadow-xs transition-all disabled:cursor-not-allowed disabled:opacity-50 md:text-sm input-depth focus-border resize-none",
    {
        variants: {
            variant: {
                outline: "border border-border bg-background hover:border-border-strong focus-visible:border-primary dark:bg-input/10",
                filled: "border border-transparent bg-muted/30 hover:bg-muted/50 focus-visible:bg-background focus-visible:border-primary focus-visible:shadow-sm",
            },
            error: {
                true: "border-destructive focus-visible:border-destructive text-destructive placeholder:text-destructive/50",
            }
        },
        defaultVariants: {
            variant: "outline",
        },
    }
)

export interface TextareaProps
    extends React.TextareaHTMLAttributes<HTMLTextAreaElement>,
    VariantProps<typeof textareaVariants> {
    error?: boolean
    errorText?: React.ReactNode
    root?: boolean
}

const Textarea = React.forwardRef<HTMLTextAreaElement, TextareaProps>(
    ({ className, variant, error, errorText, root, ...props }, ref) => {
        const isError = error || !!errorText

        const textareaNode = (
            <textarea
                className={cn(textareaVariants({ variant, error: isError, className }))}
                ref={ref}
                {...props}
            />
        )

        if (root || errorText) {
            return (
                <div className="flex flex-col gap-0.5 w-full">
                    {textareaNode}
                    {errorText && (
                        <p className="text-xs font-medium text-destructive animate-in fade-in slide-in-from-top-1 duration-200">
                            {errorText}
                        </p>
                    )}
                </div>
            )
        }

        return textareaNode
    }
)
Textarea.displayName = "Textarea"

export { Textarea }
