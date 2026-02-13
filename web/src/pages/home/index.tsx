import { Link } from 'react-router-dom'
import {
  ArrowRight,
  Terminal,
  Cpu,
  ShieldCheck,
  Zap,
  Code2,
  Sparkles,
  Command,
  Box,
  Share2,
  Workflow,
  Globe,
  Activity,
  Copy
} from 'lucide-react'
import { Button } from '@/components/ui/button'

export function HomePage() {
  return (
    <div className="min-h-screen bg-[#020617] text-slate-200 selection:bg-blue-500/30">
      {/* Cinematic Background */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none -z-10">
        <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-blue-600/10 blur-[120px] rounded-full" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-purple-600/10 blur-[120px] rounded-full" />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full h-full bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-20 mix-blend-overlay" />
      </div>

      {/* Navigation (Overlay) */}
      <nav className="container mx-auto px-6 py-6 flex justify-between items-center border-b border-slate-800/50 backdrop-blur-md sticky top-0 z-50">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 bg-gradient-to-tr from-blue-600 to-indigo-500 rounded-lg flex items-center justify-center shadow-lg shadow-blue-500/20">
            <Zap className="w-5 h-5 text-white fill-current" />
          </div>
          <span className="text-xl font-bold tracking-tighter text-white">KEST <span className="text-blue-500">OS</span></span>
        </div>
        <div className="hidden md:flex gap-8 text-sm font-medium text-slate-400">
          <a href="#features" className="hover:text-white transition-colors">Features</a>
          <a href="#ai-workflow" className="hover:text-white transition-colors">AI Context</a>
          <a href="#cli" className="hover:text-white transition-colors">CLI</a>
        </div>
        <Link to="/projects">
          <Button variant="outline" className="border-slate-700 hover:bg-slate-800 text-slate-200">
            Console
          </Button>
        </Link>
      </nav>

      {/* Hero Section */}
      <section className="container mx-auto px-6 pt-24 pb-32 text-center relative">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-500/5 border border-blue-500/20 text-blue-400 text-xs font-mono mb-8 animate-pulse">
          <Sparkles className="w-3 h-3" />
          <span>v1.2: AI Context Engine Deployment Complete</span>
        </div>

        <h1 className="text-7xl md:text-8xl font-black tracking-tighter text-white mb-8 leading-none">
          API Logic <br />
          <span className="bg-gradient-to-b from-white to-slate-500 bg-clip-text text-transparent">On Autopilot.</span>
        </h1>

        <p className="max-w-2xl mx-auto text-lg text-slate-400 mb-12 leading-relaxed">
          Kest is the bridge between your source code and AI-driven development.
          Synchronize API contracts instantly, generate deep-context for LLMs, and enforce zero-drift assertions.
        </p>

        <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
          <Link to="/projects">
            <Button size="lg" className="h-14 px-10 rounded-full bg-blue-600 hover:bg-blue-500 text-white font-bold text-lg shadow-[0_0_40px_-10px_rgba(37,99,235,0.5)] transition-all">
              Initialize Project
              <ArrowRight className="ml-2 w-5 h-5" />
            </Button>
          </Link>
          <Button size="lg" variant="ghost" className="h-14 px-10 rounded-full text-slate-300 hover:text-white hover:bg-slate-800/50">
            Read Specs
          </Button>
        </div>

        {/* Dashboard Preview Overlay */}
        <div className="mt-20 relative group">
          <div className="absolute inset-0 bg-blue-500/20 blur-[100px] opacity-20 group-hover:opacity-40 transition-opacity" />
          <div className="relative rounded-2xl border border-slate-800 bg-slate-900/50 backdrop-blur-xl p-2 shadow-2xl overflow-hidden shadow-blue-500/5">
            <div className="absolute top-0 left-0 right-0 h-8 border-b border-slate-800 bg-slate-900/80 flex items-center px-4 gap-2">
              <div className="w-2.5 h-2.5 rounded-full bg-red-500/50" />
              <div className="w-2.5 h-2.5 rounded-full bg-yellow-500/50" />
              <div className="w-2.5 h-2.5 rounded-full bg-green-500/50" />
              <div className="ml-4 text-[10px] font-mono text-slate-500 uppercase tracking-widest">Kest Web Dashboard - dev.env</div>
            </div>
            <div className="pt-10 pb-6 px-8 grid grid-cols-12 gap-6">
              <div className="col-span-3 border-r border-slate-800 pr-6 hidden md:block">
                <div className="h-4 w-32 bg-slate-800 rounded mb-4" />
                <div className="space-y-2">
                  {[1, 2, 3, 4, 5].map(i => <div key={i} className="h-2 w-full bg-slate-800/50 rounded" />)}
                </div>
              </div>
              <div className="col-span-full md:col-span-9">
                <div className="flex justify-between items-center mb-6">
                  <div className="h-6 w-48 bg-slate-800 rounded" />
                  <div className="h-6 w-12 bg-blue-600/30 rounded-full" />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div className="h-32 bg-slate-800/40 rounded-xl" />
                  <div className="h-32 bg-slate-800/40 rounded-xl" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* AI Editor Integration Section */}
      <section id="ai-workflow" className="py-32 bg-slate-950/50 relative border-y border-slate-900">
        <div className="container mx-auto px-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-20 items-center">
            <div>
              <div className="flex items-center gap-2 text-blue-500 mb-6 uppercase tracking-[0.3em] text-xs font-bold">
                <Cpu className="w-4 h-4" />
                AI-Driven Engineering
              </div>
              <h2 className="text-4xl md:text-5xl font-bold text-white mb-8 tracking-tight">
                Source to Docs <br />
                <span className="text-blue-500">In v0.6.0.</span>
              </h2>
              <p className="text-slate-400 text-lg mb-10 leading-relaxed">
                Kest CLI now scans your Go source code and generates hallucination-free API documentation.
                Perfectly aligned with your DTOs, including Mermaid sequence diagrams and security requirements.
              </p>

              <ul className="space-y-6">
                <li className="flex gap-4">
                  <div className="flex-shrink-0 w-6 h-6 rounded-full bg-blue-500/20 flex items-center justify-center text-blue-500">
                    <CheckCircle2 className="w-4 h-4" />
                  </div>
                  <div>
                    <span className="text-white font-bold block">Zero-Hallucination Scanning</span>
                    <p className="text-sm text-slate-500">Use `kest doc --ai` to generate precise request/response samples from your Go struct tags.</p>
                  </div>
                </li>
                <li className="flex gap-4">
                  <div className="flex-shrink-0 w-6 h-6 rounded-full bg-purple-500/20 flex items-center justify-center text-purple-500">
                    <CheckCircle2 className="w-4 h-4" />
                  </div>
                  <div>
                    <span className="text-white font-bold block">Agentic Flow Analysis</span>
                    <p className="text-sm text-slate-500">Automatically generates Mermaid diagrams showing your Handler {"->"} Service {"->"} Repo architecture.</p>
                  </div>
                </li>
              </ul>
            </div>

            <div className="relative group">
              <div className="absolute inset-0 bg-gradient-to-br from-blue-500/20 to-purple-500/20 blur-3xl opacity-30" />
              <div className="bg-[#0f172a] border border-slate-800 rounded-xl p-6 font-mono text-sm leading-relaxed shadow-2xl">
                <div className="flex items-center gap-2 mb-4 text-slate-500">
                  <Command className="w-4 h-4" />
                  <span>Windsurf / Cursor AI Prompt</span>
                </div>
                <div className="text-slate-200">
                  "I just modified the <span className="text-blue-400">auth-service</span>. <br />
                  Run <code className="bg-slate-800 px-1.5 py-0.5 rounded text-indigo-400">kest sync</code> to update the spec, <br />
                  then check <code className="bg-slate-800 px-1.5 py-0.5 rounded text-indigo-400">kest-spec.json</code> <br />
                  to ensure the <span className="text-purple-400">login</span> endpoint implementation <br />
                  matches our security assertions."
                </div>
                <div className="mt-6 border-t border-slate-800 pt-4 text-[10px] text-slate-500 uppercase tracking-widest">
                  Agentic Feedback Loop Active
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Feature Grids */}
      <section id="features" className="py-32 container mx-auto px-6">
        <div className="text-center mb-20">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">Core Capabilities</h2>
          <p className="text-slate-500">Engineered for speed, built for reliability.</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {[
            {
              icon: Code2,
              title: "Doc-as-Test",
              desc: "Write specs that act as live tests. Zero drift allowed.",
              color: "text-blue-500",
              bg: "bg-blue-500/5"
            },
            {
              icon: ShieldCheck,
              title: "Hard Assertions",
              desc: "Status, JSON Schema, Latency & Header integrity.",
              color: "text-green-500",
              bg: "bg-green-500/5"
            },
            {
              icon: Workflow,
              title: "Team Sync",
              desc: "Unified workspace for backend, frontend & QA alignment.",
              color: "text-purple-500",
              bg: "bg-purple-500/5"
            },
            {
              icon: Terminal,
              title: "CLI First",
              desc: "Seamless integration with CI/CD and local toolchains.",
              color: "text-indigo-500",
              bg: "bg-indigo-500/5"
            }
          ].map((f, i) => (
            <div key={i} className={`p-8 rounded-2xl border border-slate-800 ${f.bg} hover:border-slate-700 transition-all hover:-translate-y-1`}>
              <f.icon className={`w-8 h-8 ${f.color} mb-6`} />
              <h3 className="text-lg font-bold text-white mb-3">{f.title}</h3>
              <p className="text-sm text-slate-500 leading-relaxed">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* CLI Visualizer Section */}
      <section id="cli" className="py-32 relative overflow-hidden">
        <div className="container mx-auto px-6">
          <div className="bg-gradient-to-b from-slate-900 to-[#020617] rounded-[2rem] border border-slate-800 p-8 md:p-16 flex flex-col md:flex-row items-center gap-16">
            <div className="flex-1">
              <h2 className="text-4xl font-extrabold text-white mb-6">Zero-Configuration Sync</h2>
              <p className="text-slate-400 mb-8 max-w-lg leading-relaxed">
                Connect your codebase to the cloud in one command. Kest automatically detects routes,
                detects changes, and prompts for versioning.
              </p>
              <div className="flex flex-col gap-6">
                <div className="flex flex-wrap gap-4">
                  <div className="flex items-center gap-2 px-4 py-2 rounded-lg bg-slate-800/50 border border-slate-700">
                    <Box className="w-4 h-4 text-blue-500" />
                    <span className="text-sm font-mono text-slate-300">Go / Gin</span>
                  </div>
                  <div className="flex items-center gap-2 px-4 py-2 rounded-lg bg-slate-800/50 border border-slate-700">
                    <Box className="w-4 h-4 text-teal-500" />
                    <span className="text-sm font-mono text-slate-300">FastAPI</span>
                  </div>
                  <div className="flex items-center gap-2 px-4 py-2 rounded-lg bg-slate-800/50 border border-slate-700">
                    <Box className="w-4 h-4 text-emerald-500" />
                    <span className="text-sm font-mono text-slate-300">Spring Boot</span>
                  </div>
                </div>

                <div className="flex flex-col gap-3">
                  <div className="flex items-center gap-3">
                    <Button className="bg-white text-black hover:bg-slate-200 font-bold px-6">
                      <Zap className="w-4 h-4 mr-2" />
                      Download CLI v0.6.0
                    </Button>
                    <span className="text-xs text-slate-500 font-mono">macOS / Linux / Windows</span>
                  </div>
                  <div className="bg-black/50 border border-slate-800 rounded-lg px-4 py-2 font-mono text-[10px] text-blue-400 flex justify-between items-center group">
                    <code>curl -sSL https://kest.dev/install.sh | sh</code>
                    <Copy className="w-3 h-3 text-slate-600 group-hover:text-blue-400 cursor-pointer" />
                  </div>
                </div>
              </div>
            </div>

            <div className="flex-1 w-full max-w-md">
              <div className="bg-black rounded-lg p-6 shadow-2xl font-mono text-xs border border-slate-800">
                <div className="flex gap-2 mb-6">
                  <div className="w-3 h-3 rounded-full bg-slate-800" />
                  <div className="w-3 h-3 rounded-full bg-slate-800" />
                </div>
                <div className="space-y-2">
                  <div className="flex gap-2">
                    <span className="text-green-500">$</span>
                    <span className="text-white">kest sync --dir .</span>
                  </div>
                  <div className="text-slate-500">Searching for API routes...</div>
                  <div className="text-blue-400">Found: [GET] /users, [POST] /auth/login</div>
                  <div className="text-slate-500">Comparing with remote spec (ver: 1.2.0)...</div>
                  <div className="text-yellow-400">! Deviation detected in /auth/login (Schema mismatch)</div>
                  <div className="text-white">Pushing 2 updates to Kest OS...</div>
                  <div className="text-green-500">✓ Sync complete. <span className="underline cursor-pointer">View Delta</span></div>
                  <div className="text-slate-500 mt-4 pt-2 border-t border-slate-800 text-[10px]">
                    📄 Logs stored at: .kest/logs/sync_20240204.log
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Global Scaling / Distribution Section */}
      <section className="py-32 container mx-auto px-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-12">
        <div className="flex flex-col gap-4">
          <div className="w-10 h-10 rounded-xl bg-orange-500/10 flex items-center justify-center text-orange-500">
            <Globe className="w-5 h-5" />
          </div>
          <h3 className="text-xl font-bold text-white">Multi-Region Verification</h3>
          <p className="text-sm text-slate-500">Run tests from edge nodes in US-East, EU-West, and AP-Singapore simultaneously.</p>
        </div>
        <div className="flex flex-col gap-4">
          <div className="w-10 h-10 rounded-xl bg-blue-500/10 flex items-center justify-center text-blue-500">
            <Share2 className="w-5 h-5" />
          </div>
          <h3 className="text-xl font-bold text-white">Recursive Discovery</h3>
          <p className="text-sm text-slate-500">Deep-scan code dependencies to find hidden endpoints often missed by static parsers.</p>
        </div>
        <div className="flex flex-col gap-4">
          <div className="w-10 h-10 rounded-xl bg-emerald-500/10 flex items-center justify-center text-emerald-500">
            <Activity className="w-5 h-5" />
          </div>
          <h3 className="text-xl font-bold text-white">Latency Profiling</h3>
          <p className="text-sm text-slate-500">Advanced waterfall charts for every request to identify bottlenecking server-side logic.</p>
        </div>
      </section>

      {/* Final CTA */}
      <footer className="py-32 relative text-center">
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-1/2 h-px bg-gradient-to-r from-transparent via-slate-800 to-transparent" />
        <h2 className="text-4xl font-bold text-white mb-8">Elevate Your Engineering Standard.</h2>
        <Link to="/register">
          <Button size="lg" className="h-14 px-12 rounded-full bg-white text-black hover:bg-slate-200 font-black tracking-tight">
            GET STARTED FOR FREE
          </Button>
        </Link>
        <p className="mt-8 text-slate-500 text-sm">No credit card required. API-first forever.</p>

        <div className="mt-24 pt-12 border-t border-slate-900/50 flex flex-col md:flex-row justify-between items-center gap-6 px-12 text-slate-500 text-[10px] tracking-widest uppercase">
          <div>© 2026 Kest Platform / Engineered for AI-Native Development</div>
          <div className="flex gap-8">
            <a href="#" className="hover:text-white transition-colors">Twitter</a>
            <a href="#" className="hover:text-white transition-colors">GitHub</a>
            <a href="#" className="hover:text-white transition-colors">Discord</a>
          </div>
        </div>
      </footer>
    </div>
  )
}

function CheckCircle2({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="3"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
    >
      <path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z" />
      <path d="m9 12 2 2 4-4" />
    </svg>
  )
}
