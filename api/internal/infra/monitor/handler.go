package monitor

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the monitor routes
func RegisterRoutes(r *gin.Engine) {
	g := r.Group("/monitor")
	// Add Basic Auth or generic protection here if needed
	// g.Use(gin.BasicAuth(...))

	g.GET("", Dashboard)
	g.GET("/stats", StatsAPI)
}

func Dashboard(c *gin.Context) {
	html := fmt.Sprintf(dashboardHTML, "Monitor")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func StatsAPI(c *gin.Context) {
	c.JSON(http.StatusOK, GetStats())
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>%s - System Monitor</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js" defer></script>
    <style>
        body { background-color: #0f172a; color: #e2e8f0; font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif; }
        .card { background-color: #1e293b; border: 1px solid #334155; border-radius: 0.75rem; padding: 1.5rem; }
        .stat-value { font-size: 1.875rem; font-weight: 700; color: #f8fafc; }
        .stat-label { font-size: 0.875rem; font-weight: 500; color: #94a3b8; margin-top: 0.25rem; }
        .badge { display: inline-flex; align-items: center; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.75rem; font-weight: 600; }
        .badge-success { background-color: rgba(16, 185, 129, 0.1); color: #34d399; }
        .badge-error { background-color: rgba(239, 68, 68, 0.1); color: #f87171; }
        .badge-warning { background-color: rgba(245, 158, 11, 0.1); color: #fbbf24; }
        .animate-pulse { animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite; }
        @keyframes pulse { 0%%, 100%% { opacity: 1; } 50%% { opacity: .5; } }
    </style>
</head>
<body>
    <div x-data="monitor()" x-init="fetchStats()" class="min-h-screen p-6 max-w-7xl mx-auto">
        
        <!-- Header -->
        <div class="flex items-center justify-between mb-8">
            <div class="flex items-center space-x-3">
                <div class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></div>
                <h1 class="text-2xl font-bold text-white tracking-tight">System Monitor</h1>
            </div>
            <div class="flex items-center space-x-4 text-sm text-slate-400">
                <span x-text="stats.app.version"></span>
                <span class="text-slate-600">|</span>
                <span x-text="stats.app.environment" class="uppercase"></span>
            </div>
        </div>

        <!-- Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
            
            <!-- Uptime -->
            <div class="card">
                <div class="flex items-center justify-between">
                    <span class="text-slate-400 text-sm font-medium">Uptime</span>
                    <svg class="w-5 h-5 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                </div>
                <div class="stat-value mt-2" x-text="stats.app.uptime || '-'"></div>
            </div>

            <!-- Goroutines -->
            <div class="card">
                <div class="flex items-center justify-between">
                    <span class="text-slate-400 text-sm font-medium">Goroutines</span>
                    <svg class="w-5 h-5 text-indigo-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
                </div>
                <div class="stat-value mt-2" x-text="stats.system.goroutines || 0"></div>
            </div>

            <!-- Memory -->
            <div class="card">
                <div class="flex items-center justify-between">
                    <span class="text-slate-400 text-sm font-medium">Memory Usage</span>
                    <svg class="w-5 h-5 text-purple-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" /></svg>
                </div>
                <div class="stat-value mt-2" x-text="stats.system.memory || '0 B'"></div>
            </div>

            <!-- CPU Codes -->
            <div class="card">
                <div class="flex items-center justify-between">
                    <span class="text-slate-400 text-sm font-medium">CPU Cores</span>
                    <svg class="w-5 h-5 text-pink-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                </div>
                <div class="stat-value mt-2" x-text="stats.system.cpu || 1"></div>
            </div>

        </div>

        <!-- Infrastructure Status -->
        <h2 class="text-lg font-semibold text-white mb-4">Infrastructure</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            
            <!-- Database -->
            <div class="card flex items-center justify-between">
                <div class="flex items-center space-x-4">
                    <div class="p-3 bg-slate-800 rounded-lg">
                        <svg class="w-6 h-6 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" /></svg>
                    </div>
                    <div>
                        <div class="text-base font-medium text-white">Database</div>
                        <div class="text-sm text-slate-400" x-text="stats.database.latency || '0ms'"></div>
                    </div>
                </div>
                <span class="badge" :class="getStatusClass(stats.database.status)" x-text="stats.database.status"></span>
            </div>

            <!-- Redis (Placeholder) -->
             <div class="card flex items-center justify-between opacity-50">
                <div class="flex items-center space-x-4">
                    <div class="p-3 bg-slate-800 rounded-lg">
                        <svg class="w-6 h-6 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" /></svg>
                    </div>
                    <div>
                        <div class="text-base font-medium text-white">Redis</div>
                        <div class="text-sm text-slate-400">Not Configured</div>
                    </div>
                </div>
                <span class="badge badge-warning">disabled</span>
            </div>

        </div>

    </div>

    <script>
        function monitor() {
            return {
                stats: {
                    app: {},
                    system: {},
                    database: { status: 'loading' }
                },
                fetchStats() {
                    fetch('/monitor/stats')
                        .then(res => res.json())
                        .then(data => {
                            this.stats = data;
                        })
                        .catch(err => console.error(err))
                        .finally(() => {
                            setTimeout(() => this.fetchStats(), 2000); // Poll every 2s
                        });
                },
                getStatusClass(status) {
                    if (status === 'operational') return 'badge-success';
                    if (status === 'error') return 'badge-error';
                    return 'badge-warning';
                }
            }
        }
    </script>
</body>
</html>`
