import type { Metadata } from 'next';
import { notFound } from 'next/navigation';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { apiExternalBaseUrl } from '@/config/api';
import { getT } from '@/i18n/server';
import type { PublicApiSpecShare } from '@/types/api-spec';
import { formatDate } from '@/utils';

interface ApiSpecSharePageProps {
  params: Promise<{
    slug: string;
  }>;
}

const getPublicShareUrl = (slug: string) =>
  `${apiExternalBaseUrl}/public/api-spec-shares/${encodeURIComponent(slug)}`;

async function fetchPublicShare(slug: string): Promise<PublicApiSpecShare> {
  const response = await fetch(getPublicShareUrl(slug), {
    cache: 'no-store',
    headers: {
      Accept: 'application/json',
    },
  });

  if (response.status === 404) {
    notFound();
  }

  if (!response.ok) {
    throw new Error(`Failed to load public API spec share: ${response.status}`);
  }

  const payload = (await response.json()) as { data?: PublicApiSpecShare };
  if (!payload?.data) {
    throw new Error('Public API spec share response is missing data');
  }

  return payload.data;
}

export async function generateMetadata({
  params,
}: ApiSpecSharePageProps): Promise<Metadata> {
  const { slug } = await params;

  try {
    const t = await getT('project');
    const share = await fetchPublicShare(slug);
    return {
      title: `${share.method} ${share.path}`,
      description:
        share.summary || share.description || t('share.publicDescription', {
          method: share.method,
          path: share.path,
        }),
    };
  } catch {
    const t = await getT('project');
    return {
      title: t('share.sharedApiSpec'),
      description: t('share.sharedApiSpecDescription'),
    };
  }
}

function CodePanel({
  title,
  value,
  emptyLabel,
}: {
  title: string;
  value?: string;
  emptyLabel: string;
}) {
  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between gap-3">
        <h3 className="text-sm font-semibold tracking-tight text-slate-900">{title}</h3>
      </div>
      {value?.trim() ? (
        <pre className="overflow-x-auto rounded-3xl border border-slate-200 bg-slate-950 p-5 text-xs leading-6 text-slate-100">
          <code>{value}</code>
        </pre>
      ) : (
        <div className="rounded-3xl border border-dashed border-slate-300 bg-white/70 p-5 text-sm text-slate-500">
          {emptyLabel}
        </div>
      )}
    </div>
  );
}

function JsonPanel({
  title,
  value,
  emptyLabel,
}: {
  title: string;
  value: unknown;
  emptyLabel: string;
}) {
  return (
    <CodePanel
      title={title}
      value={value === undefined || value === null ? '' : JSON.stringify(value, null, 2)}
      emptyLabel={emptyLabel}
    />
  );
}

export default async function ApiSpecSharePage({
  params,
}: ApiSpecSharePageProps) {
  const t = await getT('project');
  const { slug } = await params;
  const share = await fetchPublicShare(slug);
  const hasDocSections = Boolean(
    share.doc_markdown || share.doc_markdown_zh || share.doc_markdown_en
  );

  return (
    <main className="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(14,165,233,0.18),_transparent_34%),linear-gradient(180deg,_#f8fbff_0%,_#eef6ff_48%,_#f7fafc_100%)] px-4 py-10 md:px-8 md:py-14">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-6">
        <section className="overflow-hidden rounded-[2rem] border border-sky-200/70 bg-white/80 shadow-[0_18px_60px_rgba(15,23,42,0.08)] backdrop-blur">
          <div className="border-b border-sky-100 bg-[linear-gradient(135deg,rgba(14,165,233,0.12),rgba(255,255,255,0.82))] px-6 py-6 md:px-8 md:py-8">
            <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
              <div className="space-y-3">
                <div className="flex flex-wrap items-center gap-2">
                  <Badge className="bg-sky-600 text-white hover:bg-sky-600">{share.method}</Badge>
                  <Badge variant="outline" className="border-slate-300 bg-white/70 text-slate-700">
                    v{share.version}
                  </Badge>
                  <Badge variant="outline" className="border-slate-300 bg-white/70 text-slate-700">
                    {t('share.shared')}
                  </Badge>
                </div>
                <div>
                  <h1 className="font-mono text-2xl font-semibold tracking-tight text-slate-950 md:text-3xl">
                    {share.path}
                  </h1>
                  <p className="mt-3 max-w-3xl text-sm leading-7 text-slate-600">
                    {share.summary || share.description || t('share.noSummary')}
                  </p>
                </div>
              </div>

              <div className="grid gap-3 sm:grid-cols-2">
                <div className="rounded-3xl border border-white/60 bg-white/75 px-4 py-3">
                  <p className="text-xs uppercase tracking-[0.16em] text-slate-500">{t('share.published')}</p>
                  <p className="mt-2 text-sm font-medium text-slate-900">
                    {formatDate(share.shared_at, 'YYYY-MM-DD HH:mm')}
                  </p>
                </div>
                <div className="rounded-3xl border border-white/60 bg-white/75 px-4 py-3">
                  <p className="text-xs uppercase tracking-[0.16em] text-slate-500">{t('share.updated')}</p>
                  <p className="mt-2 text-sm font-medium text-slate-900">
                    {formatDate(share.updated_at, 'YYYY-MM-DD HH:mm')}
                  </p>
                </div>
              </div>
            </div>

            {share.tags.length ? (
              <div className="mt-5 flex flex-wrap gap-2">
                {share.tags.map((tag) => (
                  <Badge
                    key={tag}
                    variant="outline"
                    className="rounded-full border-slate-300 bg-white/75 text-slate-700"
                  >
                    {tag}
                  </Badge>
                ))}
              </div>
            ) : null}
          </div>

          <div className="grid gap-4 px-6 py-6 md:grid-cols-3 md:px-8">
            <div className="rounded-3xl border border-slate-200 bg-slate-50/80 p-4">
              <p className="text-xs uppercase tracking-[0.16em] text-slate-500">{t('share.docSource')}</p>
              <p className="mt-2 text-sm font-medium text-slate-900">{share.doc_source || 'manual'}</p>
            </div>
            <div className="rounded-3xl border border-slate-200 bg-slate-50/80 p-4">
              <p className="text-xs uppercase tracking-[0.16em] text-slate-500">{t('common.path')}</p>
              <p className="mt-2 font-mono text-sm text-slate-900">{share.path}</p>
            </div>
            <div className="rounded-3xl border border-slate-200 bg-slate-50/80 p-4">
              <p className="text-xs uppercase tracking-[0.16em] text-slate-500">{t('common.method')}</p>
              <p className="mt-2 text-sm font-medium text-slate-900">{share.method}</p>
            </div>
          </div>
        </section>

        <section className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
          <div className="space-y-6">
            <Card className="rounded-[2rem] border-slate-200/80 bg-white/85 shadow-[0_16px_40px_rgba(15,23,42,0.06)]">
              <CardHeader>
                <CardTitle>{t('share.descriptionTitle')}</CardTitle>
                <CardDescription>{t('share.descriptionDescription')}</CardDescription>
              </CardHeader>
              <CardContent className="text-sm leading-7 text-slate-700">
                {share.description || t('share.noDetailedDescription')}
              </CardContent>
            </Card>

            {hasDocSections ? (
              <Card className="rounded-[2rem] border-slate-200/80 bg-white/85 shadow-[0_16px_40px_rgba(15,23,42,0.06)]">
                <CardHeader>
                  <CardTitle>{t('share.publishedDocs')}</CardTitle>
                  <CardDescription>{t('share.publishedDocsDescription')}</CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                  <CodePanel
                    title={t('share.defaultMarkdown')}
                    value={share.doc_markdown}
                    emptyLabel={t('share.noDefaultMarkdown')}
                  />
                  <CodePanel
                    title={t('share.chineseMarkdown')}
                    value={share.doc_markdown_zh}
                    emptyLabel={t('share.noChineseMarkdown')}
                  />
                  <CodePanel
                    title={t('share.englishMarkdown')}
                    value={share.doc_markdown_en}
                    emptyLabel={t('share.noEnglishMarkdown')}
                  />
                </CardContent>
              </Card>
            ) : null}
          </div>

          <div className="space-y-6">
            <Card className="rounded-[2rem] border-slate-200/80 bg-white/85 shadow-[0_16px_40px_rgba(15,23,42,0.06)]">
              <CardHeader>
                <CardTitle>{t('share.requestSchema')}</CardTitle>
                <CardDescription>{t('share.requestSchemaDescription')}</CardDescription>
              </CardHeader>
              <CardContent>
                <JsonPanel
                  title={t('common.requestBody')}
                  value={share.request_body}
                  emptyLabel={t('share.noRequestBody')}
                />
              </CardContent>
            </Card>

            <Card className="rounded-[2rem] border-slate-200/80 bg-white/85 shadow-[0_16px_40px_rgba(15,23,42,0.06)]">
              <CardHeader>
                <CardTitle>{t('common.parameters')}</CardTitle>
                <CardDescription>{t('share.parametersDescription')}</CardDescription>
              </CardHeader>
              <CardContent>
                <JsonPanel
                  title={t('common.parameters')}
                  value={share.parameters}
                  emptyLabel={t('share.noParameterSchema')}
                />
              </CardContent>
            </Card>

            <Card className="rounded-[2rem] border-slate-200/80 bg-white/85 shadow-[0_16px_40px_rgba(15,23,42,0.06)]">
              <CardHeader>
                <CardTitle>{t('common.responses')}</CardTitle>
                <CardDescription>{t('share.responsesDescription')}</CardDescription>
              </CardHeader>
              <CardContent>
                <JsonPanel
                  title={t('common.responses')}
                  value={share.responses}
                  emptyLabel={t('share.noResponseSchema')}
                />
              </CardContent>
            </Card>
          </div>
        </section>
      </div>
    </main>
  );
}
