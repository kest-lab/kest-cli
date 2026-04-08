"use client"

import * as React from "react"
import { Button } from "@/components/ui/button"

export function DepthShowcase() {
  return (
    <section className="space-y-6">
      <h2 className="text-2xl font-semibold">Depth & Glassmorphism</h2>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 p-8 bg-linear-to-br from-muted/20 via-transparent to-primary/5 rounded-2xl relative overflow-hidden border">
        {/* Animated Background Decor */}
        <div className="absolute top-[-20%] left-[-10%] size-72 bg-primary/15 rounded-full blur-[80px]" />
        <div className="absolute bottom-[-20%] right-[-10%] size-72 bg-purple-500/15 rounded-full blur-[80px]" />
        
        <div className="glass-panel p-6 rounded-xl flex flex-col gap-4 z-10">
          <div className="space-y-2">
            <h3 className="text-lg font-semibold">Glass Panel (Light)</h3>
            <p className="text-sm text-muted-foreground leading-relaxed">
              Utilizes <code className="text-xs bg-muted px-1 py-0.5 rounded">backdrop-filter: blur</code> and <code className="text-xs bg-muted px-1 py-0.5 rounded">color-mix</code> for a translucent, modern feel.
            </p>
          </div>
          <div className="flex gap-3 mt-auto">
            <Button variant="outline" size="sm">Secondary</Button>
            <Button size="sm">Primary</Button>
          </div>
        </div>

        <div className="glass-panel-dark p-6 rounded-xl flex flex-col gap-4 z-10">
          <div className="space-y-2">
            <h3 className="text-lg font-semibold text-white">Glass Panel (Dark)</h3>
            <p className="text-sm text-white/70 leading-relaxed">
              Optimized for high-contrast environments or dark mode surfaces, maintaining readability through adjusted transparency.
            </p>
          </div>
          <div className="flex gap-3 mt-auto">
            <Button variant="ghost" size="sm" className="border border-white/20 text-white! hover:bg-white/10 hover:text-white!">
              Ghost
            </Button>
            <Button size="sm" className="bg-white text-black hover:bg-white/90">
              Action
            </Button>
          </div>
        </div>
      </div>
    </section>
  )
}
