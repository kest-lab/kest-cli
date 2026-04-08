"use client"

import * as React from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { 
  Breadcrumb, 
  BreadcrumbItem, 
  BreadcrumbLink, 
  BreadcrumbList, 
  BreadcrumbPage, 
  BreadcrumbSeparator 
} from "@/components/ui/breadcrumb"
import { Card, CardContent } from "@/components/ui/card"

export function NavigationShowcase() {
  return (
    <section className="space-y-8">
      <h2 className="text-2xl font-semibold">Navigation & Tabs</h2>

      <div className="space-y-4">
        <h3 className="text-lg font-medium">Tabs (Spring Transitions)</h3>
        <Tabs defaultValue="account" className="w-full">
          <TabsList className="grid w-full grid-cols-3 lg:w-[450px]">
            <TabsTrigger value="account" className="interactive-subtle">Account</TabsTrigger>
            <TabsTrigger value="password">Password</TabsTrigger>
            <TabsTrigger value="settings" disabled>Disabled</TabsTrigger>
          </TabsList>
          <TabsContent value="account" className="mt-4 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 animate-in fade-in-50 duration-500">
            <Card className="border-border/50">
              <CardContent className="pt-6">
                <p className="text-sm text-text-muted">
                  Account settings and profile information. The trigger uses the <code>interactive-subtle</code> utility for a light bounce.
                </p>
              </CardContent>
            </Card>
          </TabsContent>
          <TabsContent value="password" className="mt-4 animate-in fade-in-50 duration-500">
            <Card className="border-border/50">
              <CardContent className="pt-6">
                <p className="text-sm text-text-muted">
                  Change your password and security settings here.
                </p>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>

      <div className="space-y-4">
        <h3 className="text-lg font-medium">Breadcrumbs</h3>
        <div className="p-6 border rounded-xl bg-card">
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem>
                <BreadcrumbLink href="/" className="interactive-subtle px-2 py-1 rounded-md">Home</BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbLink href="/styleguide" className="interactive-subtle px-2 py-1 rounded-md">Components</BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbPage className="px-2 py-1">Navigation</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </div>
      </div>
    </section>
  )
}
