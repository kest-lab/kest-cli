"use client"

import * as React from "react"
import { Separator } from "@/components/ui/separator"
import { ButtonShowcase } from "@/components/features/styleguide/button-showcase"
import { FormShowcase } from "@/components/features/styleguide/form-showcase"
import { DataDisplayShowcase } from "@/components/features/styleguide/data-display-showcase"
import { FeedbackShowcase } from "@/components/features/styleguide/feedback-showcase"
import { DepthShowcase } from "@/components/features/styleguide/depth-showcase"
import { TooltipShowcase } from "@/components/features/styleguide/tooltip-showcase"
import { OverlayShowcase } from "@/components/features/styleguide/overlay-showcase"
import { NavigationShowcase } from "@/components/features/styleguide/navigation-showcase"

export default function StyleguidePage() {
  return (
    <div className="container max-w-7xl mx-auto py-10 space-y-20 pb-40">
      {/* Header */}
      <section className="space-y-4">
        <div className="space-y-2">
          <h1 className="text-5xl font-extrabold tracking-tight bg-clip-text text-transparent bg-linear-to-r from-primary to-purple-600">
            UI Styleguide
          </h1>
          <p className="text-xl text-text-subtle max-w-2xl font-medium">
            A premium component library powered by Tailwind v4, Spring animations, and an advanced elevation system.
          </p>
        </div>
        <Separator className="bg-border/50" />
      </section>

      {/* Components */}
      <ButtonShowcase />
      <FormShowcase />
      <DataDisplayShowcase />
      <FeedbackShowcase />
      <TooltipShowcase />
      <NavigationShowcase />    
      <DepthShowcase />
      <OverlayShowcase />

      {/* Footer Info */}
      <footer className="pt-20 border-t flex flex-col md:flex-row justify-between gap-8 text-text-muted">
        <div className="space-y-2">
          <p className="font-semibold text-text-main">Design Tokens</p>
          <ul className="text-sm space-y-1">
            <li>Easing: <code className="bg-muted px-1 rounded text-xs">--ease-spring</code></li>
            <li>Shadows: <code className="bg-muted px-1 rounded text-xs">--shadow-premium</code></li>
            <li>Glass: <code className="bg-muted px-1 rounded text-xs">--glass-blur</code></li>
          </ul>
        </div>
        <div className="text-sm max-w-xs leading-relaxed italic">
          &quot;The details are not the details. They make the design.&quot; - Charles Eames
        </div>
      </footer>
    </div>
  )
}
