"use client"

import * as React from "react"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip"
import { Button } from "@/components/ui/button"

export function TooltipShowcase() {
  return (
    <section className="space-y-8">
      <h2 className="text-2xl font-semibold">Tooltips</h2>
      
      <div className="space-y-4">
        <h3 className="text-lg font-medium">Variants</h3>
        <div className="flex flex-wrap gap-8 p-10 border rounded-xl bg-card items-center justify-center">
          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="outline">Default Glass</Button>
            </TooltipTrigger>
            <TooltipContent variant="default">
              Glassmorphism look (New)
            </TooltipContent>
          </Tooltip>

          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="outline">Inverted Dark</Button>
            </TooltipTrigger>
            <TooltipContent variant="inverted">
              High contrast (Classic)
            </TooltipContent>
          </Tooltip>

          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="outline">Primary Color</Button>
            </TooltipTrigger>
            <TooltipContent variant="primary">
              Branded theme style
            </TooltipContent>
          </Tooltip>

          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="outline">Destructive</Button>
            </TooltipTrigger>
            <TooltipContent variant="destructive">
              Critical warning message
            </TooltipContent>
          </Tooltip>
        </div>
      </div>

      <div className="space-y-4">
        <h3 className="text-lg font-medium">Positions</h3>
        <div className="flex flex-wrap gap-8 p-10 border rounded-xl bg-card items-center justify-center">
          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="outline">Top</Button>
            </TooltipTrigger>
            <TooltipContent side="top">Tooltip on Top</TooltipContent>
          </Tooltip>

          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="outline">Bottom</Button>
            </TooltipTrigger>
            <TooltipContent side="bottom">Tooltip on Bottom</TooltipContent>
          </Tooltip>

          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="outline">Left</Button>
            </TooltipTrigger>
            <TooltipContent side="left">Tooltip on Left</TooltipContent>
          </Tooltip>

          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="outline">Right</Button>
            </TooltipTrigger>
            <TooltipContent side="right">Tooltip on Right</TooltipContent>
          </Tooltip>
        </div>
      </div>
    </section>
  )
}
