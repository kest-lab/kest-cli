# 🚀 LlamaFront AI Scaffold

A modern, AI-powered frontend application scaffold designed for the AI era. Built specifically for vibe coding and AI-assisted development, LlamaFront provides everything you need to build intelligent, scalable, and performant frontend applications with maximum developer productivity.

## ✨ Frontend-First AI Features

- 🎨 **Component-Driven**: Extensive UI component library with Radix UI and custom designs
- 🚀 **Performance Optimized**: Next.js 16.1.1 App Router with automatic code splitting
- 🌙 **Theme System**: Beautiful dark/light themes with CSS variables
- 📱 **Mobile-First**: Responsive design for all screen sizes
- 🔍 **TypeScript**: Full type safety and excellent DX (TypeScript 5.9+)
- ⚡ **Hot Reload**: Instant development feedback with Next.js Turbopack
- 🤖 **AI-Ready**: Clean patterns for AI code generation and "vibe coding"
- 🔐 **Auth Integration**: Built-in authentication patterns with persistent storage
- 📊 **State Management**: Zustand 5.0 for predictable, granular state handling
- 🌍 **I18n**: Full internationalization support with `next-intl`
- 🎨 **Styleguide Explorer**: Pre-built component gallery and design system playground
- 🛠️ **Developer Tools**: Pre-configured ESLint 9, Prettier, and Vitest

## 🤖 AI Developer Experience

- **AI-Friendly Code Structure**: Clean, predictable patterns that AI tools (like Windsurf, Cursor, Bolt) understand
- **Smart Component Design**: Components designed for AI generation and modification, utilizing Atomic Design principles
- **Type Safety**: Comprehensive TypeScript types for better AI code completion and error prevention
- **Documentation**: Rich JSDoc comments for AI context understanding
- **Error Handling**: Standardized error handling patterns for AI debugging assistance

## 🆕 Latest Updates (v2.1.0)

- ✅ **Next.js 16.1.1** - Latest stable version
- ✅ **React 19.2.3** - Full support for React 19 features
- ✅ **Tailwind CSS 4.1.18** - Modern utility-first CSS
- ✅ **Next-Intl 4.6** - Comprehensive i18n solution
- ✅ **Zustand 5.0** - Optimized state management
- ✅ **Architecture Guide** - Comprehensive guide for building scalable AI-ready apps

👉 Check out the [Optimization Summary Report](docs/OPTIMIZATION_SUMMARY.md) for details.

## 🛠️ Frontend-Optimized Tech Stack

- **Framework**: Next.js 16.1.1 (App Router)
- **Library**: React 19.2.3
- **Language**: TypeScript 5.9.3 (Strict Mode)
- **Styling**: Tailwind CSS 4.1.18 + PostCSS
- **UI Components**: Radix UI + Lucide Icons
- **State Management**: Zustand 5.0.9
- **Data Fetching**: TanStack Query v5
- **Forms**: React Hook Form 7.69 + Zod 4.2
- **Theming**: Next-Themes 0.4
- **Testing**: Vitest 4.0 + Testing Library

## 🚀 Quick Start

### Prerequisites

- Node.js 18+
- pnpm 10+ (Recommended)

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/llamacto/llamafront-ai-scaffold.git
   cd llamafront-ai-scaffold
   ```

2. **Install dependencies**

   ```bash
   pnpm install
   ```

3. **Set up environment variables**

   ```bash
   cp .env.example .env.local
   # Edit .env.local with your configuration
   ```

4. **Run the development server**

   ```bash
   pnpm dev
   ```

5. **Open your browser**

   Navigate to [http://localhost:3000](http://localhost:3000)

## 📁 Project Structure

```text
src/
├── app/                    # Next.js App Router
│   ├── (auth)/            # Authentication routes
│   ├── (site)/            # Marketing/Public pages
│   ├── (normal)/          # Authenticated layout group
│   │   ├── console/       # Admin dashboard
│   │   └── styleguide/    # Component explorer
│   └── api/               # API Route handlers
├── components/            # Reusable components
│   ├── ui/               # Base UI library (Shadcn-like)
│   ├── features/         # Domain-specific components
│   │   ├── auth/         # Login, Register, Guards
│   │   ├── console/      # Dashboard widgets
│   │   └── styleguide/   # Styleguide specific components
│   └── common/           # Shared layout components
├── hooks/                # Custom React hooks
├── services/             # API/Business logic layer
├── store/                # Zustand state slices
├── i18n/                 # Translation files
├── providers/            # React context providers
├── utils/                # Utility functions
└── types/                # TypeScript definitions
```

## 📊 Features in Depth

### 🎨 **Styleguide**
A built-in styleguide available at `/styleguide` allows you to explore all UI components, colors, and typography in isolation. This is perfect for maintaining visual consistency.

### 🔐 **Authentication**
Complete auth flow out-of-the-box:
- Login/Register pages with validation
- JWT & Session management in Zustand
- `AuthGuard` component for route protection
- Middleware-level protection patterns

### 🌍 **Internationalization**
Powered by `next-intl`, supporting:
- Multi-language routing
- Type-safe translation keys
- Dynamic language switching

### 📈 **Dashboard & Analytics**
A ready-to-use console layout includes:
- Performance monitoring charts (Recharts)
- Activity timelines
- Stats summaries
- Responsive sidebar navigation

## 🚀 Deployment

### Vercel (Recommended)
Deployment is seamless on Vercel with zero configuration.

### Docker
```bash
docker build -t llamafront-web .
docker run -p 3000:3000 llamafront-web
```

## 🧪 Scripts

```bash
pnpm dev          # Start development with Turbopack
pnpm build        # Build for production
pnpm start        # Start production server
pnpm lint         # Run ESLint & Type-check
pnpm test         # Run unit tests
pnpm format       # Format with Prettier
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

---

Made with ❤️ by Llamacto Team for the AI Development Community
