"use client"

import * as React from "react"
import { LucideIcon, TrendingUp, TrendingDown } from "lucide-react"
import { cn } from "@/utils"

interface StatCardProps {
    title: string
    value: string | number
    description?: string
    trend?: {
        value: string
        isPositive: boolean
    }
    icon: LucideIcon
    variant?: "default" | "primary" | "success" | "warning"
    className?: string
}

export function StatCard({
    title,
    value,
    description,
    trend,
    icon: Icon,
    variant = "default",
    className,
}: StatCardProps) {
    const variantStyles = {
        default: "bg-bg-surface border-border-subtle",
        primary: "bg-linear-to-br from-primary/10 via-primary/5 to-transparent border-primary/20",
        success: "bg-linear-to-br from-success/10 via-success/5 to-transparent border-success/20",
        warning: "bg-linear-to-br from-warning/10 via-warning/5 to-transparent border-warning/20",
    }

    const iconStyles = {
        default: "bg-muted text-text-muted",
        primary: "bg-primary/10 text-primary",
        success: "bg-success/10 text-success",
        warning: "bg-warning/10 text-warning",
    }

    return (
        <div
            className={cn(
                "relative overflow-hidden rounded-xl border p-6 transition-all duration-300",
                "hover:shadow-premium hover:-translate-y-0.5",
                variantStyles[variant],
                className
            )}
        >
            <div className="absolute inset-0 bg-linear-to-br from-white/5 to-transparent pointer-events-none" />

            <div className="relative flex items-start justify-between">
                <div className="space-y-2">
                    <p className="text-sm font-medium text-text-muted">{title}</p>
                    <div className="flex items-baseline gap-2">
                        <span className="text-3xl font-bold tracking-tight text-text-main">
                            {value}
                        </span>
                        {trend && (
                            <span
                                className={cn(
                                    "flex items-center gap-0.5 text-xs font-medium",
                                    trend.isPositive ? "text-success" : "text-destructive"
                                )}
                            >
                                {trend.isPositive ? (
                                    <TrendingUp className="h-3 w-3" />
                                ) : (
                                    <TrendingDown className="h-3 w-3" />
                                )}
                                {trend.value}
                            </span>
                        )}
                    </div>
                    {description && (
                        <p className="text-xs text-text-muted">{description}</p>
                    )}
                </div>

                <div
                    className={cn(
                        "flex h-10 w-10 items-center justify-center rounded-lg",
                        iconStyles[variant]
                    )}
                >
                    <Icon className="h-5 w-5" />
                </div>
            </div>
        </div>
    )
}

export function StatCardSkeleton() {
    return (
        <div className="relative overflow-hidden rounded-xl border border-border-subtle bg-bg-surface p-6">
            <div className="flex items-start justify-between">
                <div className="space-y-2">
                    <div className="h-4 w-24 animate-pulse rounded bg-muted" />
                    <div className="h-8 w-32 animate-pulse rounded bg-muted" />
                    <div className="h-3 w-20 animate-pulse rounded bg-muted" />
                </div>
                <div className="h-10 w-10 animate-pulse rounded-lg bg-muted" />
            </div>
        </div>
    )
}
