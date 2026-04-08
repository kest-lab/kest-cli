"use client"

import * as React from "react"
import { Skeleton } from "@/components/ui/skeleton"
import { Card, CardContent, CardHeader } from "@/components/ui/card"

/**
 * @component LoadingShowcase
 * @category Feature
 * @status Stable
 * @description Demonstrates skeleton loading patterns for various UI components.
 * @usage Reference this showcase when implementing loading states in the application.
 */
export function LoadingShowcase() {
  return (
    <section className="space-y-8">
      <h2 className="text-2xl font-semibold">Loading States & Skeletons</h2>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {/* Basic Shapes */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Basic Shapes</h3>
          <div className="p-6 border rounded-xl bg-card space-y-4">
            <div className="flex items-center gap-4">
              <Skeleton className="size-12 rounded-full" />
              <div className="space-y-2 flex-1">
                <Skeleton className="h-4 w-3/4" />
                <Skeleton className="h-3 w-1/2" />
              </div>
            </div>
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-5/6" />
            <Skeleton className="h-4 w-4/6" />
          </div>
        </div>

        {/* Card Loading Pattern */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Card Pattern</h3>
          <Card>
            <CardHeader className="space-y-2">
              <Skeleton className="h-5 w-1/3" />
              <Skeleton className="h-4 w-2/3" />
            </CardHeader>
            <CardContent className="space-y-3">
              <Skeleton className="h-32 w-full rounded-lg" />
              <div className="flex gap-2">
                <Skeleton className="h-8 w-20 rounded-md" />
                <Skeleton className="h-8 w-20 rounded-md" />
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Avatar + Text Pattern */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">List Pattern</h3>
          <div className="p-6 border rounded-xl bg-card space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="flex items-center gap-3">
                <Skeleton className="size-10 rounded-full shrink-0" />
                <div className="space-y-1.5 flex-1">
                  <Skeleton className="h-4 w-1/2" />
                  <Skeleton className="h-3 w-3/4" />
                </div>
                <Skeleton className="h-6 w-16 rounded-full" />
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Table Loading Pattern */}
      <div className="space-y-4">
        <h3 className="text-lg font-medium">Table Pattern</h3>
        <div className="border rounded-xl bg-card overflow-hidden">
          <div className="bg-muted/30 p-4 border-b">
            <div className="flex gap-4">
              <Skeleton className="h-4 w-16" />
              <Skeleton className="h-4 w-32" />
              <Skeleton className="h-4 w-24" />
              <Skeleton className="h-4 w-20 ml-auto" />
            </div>
          </div>
          {[1, 2, 3].map((i) => (
            <div key={i} className="p-4 border-b last:border-0">
              <div className="flex gap-4 items-center">
                <Skeleton className="h-4 w-12" />
                <Skeleton className="h-4 w-40" />
                <Skeleton className="h-6 w-20 rounded-full" />
                <Skeleton className="h-4 w-24 ml-auto" />
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Inline Loading States */}
      <div className="space-y-4">
        <h3 className="text-sm font-medium text-text-muted uppercase tracking-wider">Inline Variations</h3>
        <div className="flex flex-wrap gap-6 p-6 border rounded-xl bg-card">
          <div className="space-y-2">
            <span className="text-xs text-text-muted">Button</span>
            <Skeleton className="h-9 w-24 rounded-lg" />
          </div>
          <div className="space-y-2">
            <span className="text-xs text-text-muted">Input</span>
            <Skeleton className="h-9 w-48 rounded-lg" />
          </div>
          <div className="space-y-2">
            <span className="text-xs text-text-muted">Badge</span>
            <Skeleton className="h-5 w-16 rounded-full" />
          </div>
          <div className="space-y-2">
            <span className="text-xs text-text-muted">Avatar</span>
            <Skeleton className="size-8 rounded-full" />
          </div>
        </div>
      </div>
    </section>
  )
}
