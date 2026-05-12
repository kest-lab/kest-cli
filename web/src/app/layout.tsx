import type { Metadata } from "next";
import { JetBrains_Mono, Noto_Sans } from "next/font/google";
import { NextIntlClientProvider } from "next-intl";
import { getLocale, getMessages } from "next-intl/server";
import { cn } from "@/utils";
import { ThemeProvider } from "@/providers/theme-provider";
import { QueryProvider } from "@/providers/query-provider";
import { AuthProvider } from "@/providers/auth-provider";
import { LocalizedErrorBoundary } from "@/components/localized-error-boundary";
import { Toaster } from "@/components/ui/sonner";
import { Analytics } from "@/components/analytics";
import { getT } from "@/i18n/server";
import "./globals.css";

const notoSans = Noto_Sans({
  subsets: ["latin"],
  weight: ["400", "500", "600"],
  variable: "--font-noto-sans",
});
const mono = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-mono-fallback",
});

export const metadata: Metadata = {
  title: {
    default: "kest",
    template: "%s | kest",
  },
  description: "Open-source API testing and collaboration with context-aware flows and AI-powered diagnosis.",
  metadataBase: new URL("http://localhost:3000"), // Purified: Use hardcoded base
  icons: {
    icon: [
      { url: "/icon.svg", type: "image/svg+xml" },
      { url: "/favicon.ico", sizes: "any" },
    ],
  },
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const locale = await getLocale();
  const messages = await getMessages();
  const t = await getT();

  return (
    <html lang={locale} suppressHydrationWarning>
      <body className={cn(notoSans.variable, mono.variable, "min-h-screen antialiased")}>
        <NextIntlClientProvider locale={locale} messages={messages}>
          <LocalizedErrorBoundary
            fallbackTitle={t.common("errorBoundaryTitle")}
            fallbackDescription={t.common("errorBoundaryDescription")}
            retryLabel={t.common("errorBoundaryRetry")}
          >
            <ThemeProvider
              attribute="class"
              defaultTheme="light"
              enableSystem={false}
              disableTransitionOnChange
            >
              <QueryProvider>
                <AuthProvider>
                  {children}
                </AuthProvider>
              </QueryProvider>
              <Toaster richColors position="top-right" />
            </ThemeProvider>
          </LocalizedErrorBoundary>
        </NextIntlClientProvider>
        <Analytics />
      </body>
    </html>
  );
}
