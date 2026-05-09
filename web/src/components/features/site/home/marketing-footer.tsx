import Link from 'next/link';
import { FaDiscord, FaGithub, FaXTwitter } from 'react-icons/fa6';
import { Logo } from '@/components/ui/icons';
import type { MarketingFooterContent } from './types';

/**
 * @component MarketingFooter
 * @category Feature
 * @status Stable
 * @description Displays the multi-column footer used by the public marketing site.
 * @usage Use in the site layout for marketing pages.
 * @example
 * <MarketingFooter brandName="KestFlow" content={footer} />
 */
export interface MarketingFooterProps {
  brandName: string;
  content: MarketingFooterContent;
}

const socialIcons = {
  GitHub: FaGithub,
  Discord: FaDiscord,
  X: FaXTwitter,
} as const;

// 判断 footer 链接是否为外部地址。
// 作用：API Docs 这类外链需要显式设置新标签页打开。
function isExternalHref(href: string) {
  return href.startsWith('http://') || href.startsWith('https://');
}

function FooterLink({ label, href, placeholder = false }: { label: string; href?: string; placeholder?: boolean }) {
  if (href) {
    // footer 同时承载站内锚点和外部文档站链接，这里统一处理跳转行为。
    const external = isExternalHref(href);

    return (
      <Link
        href={href}
        target={external ? '_blank' : undefined}
        rel={external ? 'noreferrer' : undefined}
        className="text-sm text-text-main transition-colors duration-200 hover:text-text-subtle"
      >
        {label}
      </Link>
    );
  }

  return (
    <span className="inline-flex items-center gap-2 text-sm text-text-muted">
      {label}
      {placeholder ? <span className="size-1.5 rounded-full bg-text-muted/60" /> : null}
    </span>
  );
}

export function MarketingFooter({ brandName, content }: MarketingFooterProps) {
  return (
    <footer className="border-t border-border-main bg-bg-canvas">
      <div className="container py-16 sm:py-20">
        <div className="grid gap-10 lg:grid-cols-[1.05fr_1.95fr]">
          <div className="max-w-sm">
            <div>
              <Logo className="h-9 w-[111px] text-black" role="img" aria-label={brandName} />
              <p className="mt-4 text-base leading-7 text-text-subtle">{content.tagline}</p>
            </div>
            <div className="mt-6">
              <p className="figma-caption text-text-muted">{content.socialsTitle}</p>
              <div className="mt-3 flex flex-wrap gap-3">
                {content.socials.map((social) => {
                  const Icon = socialIcons[social.label as keyof typeof socialIcons];

                  return (
                    <span
                      key={social.label}
                      className="inline-flex items-center gap-2 rounded-pill border border-border-main bg-bg-canvas px-3 py-1.5 text-sm text-text-main"
                    >
                      {Icon ? <Icon className="size-4 shrink-0" /> : null}
                      <span>{social.label}</span>
                    </span>
                  );
                })}
              </div>
            </div>
          </div>

          <div className="grid gap-8 sm:grid-cols-2 xl:grid-cols-5">
            {content.columns.map((column) => (
              <div key={column.title}>
                <p className="figma-caption text-text-muted">{column.title}</p>
                <div className="mt-4 space-y-3">
                  {column.links.map((link) => (
                    <FooterLink
                      key={link.label}
                      label={link.label}
                      href={link.href}
                      placeholder={link.placeholder}
                    />
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="mt-12 flex flex-col gap-3 border-t border-border-main pt-6 text-sm text-text-muted sm:flex-row sm:items-center sm:justify-between">
          <p>© {new Date().getFullYear()} {brandName}.</p>
          <p>{content.tagline}</p>
        </div>
      </div>
    </footer>
  );
}
