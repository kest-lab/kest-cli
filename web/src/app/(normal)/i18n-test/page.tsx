'use client';

import React from 'react';
import { useT } from '@/i18n/client';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Layers, ChevronRight, CheckCircle2 } from 'lucide-react';

export default function I18nTestPage() {
  // 1. Root level t (Global)
  const tGlobal = useT();
  
  // 2. Scoped level 1
  const tLevel1 = useT('test.level1');
  
  // 3. Scoped level 2
  const tLevel2 = useT('test.level1.level2');
  
  // 4. Scoped level 4 (Deeply nested)
  const tLevel4 = useT('test.level1.level2.level3.level4');

  return (
    <div className="container max-w-4xl py-10 space-y-10 animate-in fade-in duration-500">
      <div className="space-y-2">
        <h1 className="text-4xl font-extrabold tracking-tight lg:text-5xl bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">
          {tGlobal('test.title')}
        </h1>
        <p className="text-muted-foreground text-lg">
          Testing multi-layer scoped translation type safety and runtime resolution.
        </p>
      </div>

      <div className="grid gap-6">
        {/* Level 1 Scope */}
        <Card className="border-primary/10 bg-muted/50 overflow-hidden group hover:border-primary/30 transition-all duration-300">
          <CardHeader className="bg-primary/5 border-b border-primary/10">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="p-2 bg-primary rounded-lg text-primary-foreground shadow-lg shadow-primary/20">
                  <Layers className="w-5 h-5" />
                </div>
                <CardTitle className="text-xl">Scope: <code className="text-primary">test.level1</code></CardTitle>
              </div>
              <Badge variant="outline" className="bg-background/50 border-primary/20">Level 1</Badge>
            </div>
          </CardHeader>
          <CardContent className="pt-6 space-y-4">
            <div className="flex items-center gap-3 p-3 rounded-lg bg-background/50 border border-border/50 group-hover:border-primary/20 transition-colors">
              <CheckCircle2 className="w-5 h-5 text-success" />
              <div className="flex-1">
                <p className="text-sm font-semibold text-muted-foreground mb-1 uppercase tracking-wider">{"t('title')"}</p>
                <p className="text-xl font-medium">{tLevel1('title')}</p>
              </div>
            </div>
            <div className="bg-muted p-3 rounded-md border border-border/30">
               <code className="text-xs text-muted-foreground">{"useT('test.level1') -> t('title')"}</code>
            </div>
          </CardContent>
        </Card>

        {/* Level 2 Scope */}
        <Card className="border-info/10 bg-info/5 overflow-hidden group hover:border-info/30 transition-all duration-300">
          <CardHeader className="bg-info/10 border-b border-info/10">
             <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="p-2 bg-info rounded-lg text-info-foreground shadow-lg shadow-info/20">
                  <Layers className="w-5 h-5" />
                </div>
                <CardTitle className="text-xl">Scope: <code className="text-info font-bold">test...level2</code></CardTitle>
              </div>
              <Badge variant="outline" className="bg-background/50 border-info/20">Level 2</Badge>
            </div>
          </CardHeader>
          <CardContent className="pt-6 space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="p-4 rounded-lg bg-background/50 border border-border/50 group-hover:border-info/20 transition-colors">
                <p className="text-xs font-bold text-muted-foreground mb-1 uppercase tracking-tighter">{"t('title')"}</p>
                <p className="text-lg font-semibold">{tLevel2('title')}</p>
              </div>
              <div className="p-4 rounded-lg bg-background/50 border border-border/50 group-hover:border-info/20 transition-colors">
                <p className="text-xs font-bold text-muted-foreground mb-1 uppercase tracking-tighter">{"t('message')"}</p>
                <p className="text-lg font-medium">{tLevel2('message')}</p>
              </div>
            </div>
             <div className="bg-muted p-3 rounded-md border border-border/30 text-center">
               <code className="text-xs text-muted-foreground">{"t.test('level1.level2.message') -> "}{tGlobal.test('level1.level2.message')}</code>
            </div>
          </CardContent>
        </Card>

        {/* Level 4 Scope - Deepest */}
        <Card className="border-warning/20 bg-warning/5 overflow-hidden group hover:border-warning/40 transition-all duration-300 relative">
          <div className="absolute top-0 right-0 p-4 opacity-5 group-hover:opacity-10 transition-opacity">
            <Layers className="w-40 h-40 -mr-10 -mt-10 rotate-12" />
          </div>
          <CardHeader className="bg-warning/10 border-b border-warning/10">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="p-2 bg-warning rounded-lg text-warning-foreground shadow-lg shadow-warning/20">
                  <Layers className="w-5 h-5" />
                </div>
                <CardTitle className="text-xl">Scope: <code className="text-warning-foreground font-mono">level4</code></CardTitle>
              </div>
              <Badge variant="outline" className="bg-background/50 border-warning/30">Level 4 Deep</Badge>
            </div>
          </CardHeader>
          <CardContent className="pt-6 space-y-6">
            <div className="bg-background/40 p-6 rounded-xl border-2 border-dashed border-warning/20 group-hover:border-warning/40 transition-all">
              <div className="space-y-4">
                <div className="flex items-center gap-2 text-warning-foreground font-semibold">
                  <span>{tLevel4('title')}</span>
                  <ChevronRight className="w-4 h-4" />
                </div>
                <h3 className="text-2xl font-bold tracking-tight">{tLevel4('message')}</h3>
                <p className="text-muted-foreground bg-muted/80 p-3 rounded-md border border-border/20 font-mono text-sm">
                  {tLevel4('deepValue', { value: 'Successfully Injected!' })}
                </p>
              </div>
            </div>
            
            <Separator className="bg-warning/10" />
            
            <div className="flex items-center justify-center gap-4 text-xs font-medium text-muted-foreground">
              <span className="px-2 py-1 bg-muted rounded">TypeScript Verified</span>
              <span className="px-2 py-1 bg-muted rounded">Runtime Resolved</span>
              <span className="px-2 py-1 bg-muted rounded">Dynamic Loaded</span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
