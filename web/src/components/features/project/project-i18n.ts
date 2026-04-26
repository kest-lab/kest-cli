import type { ProjectWorkspaceModuleI18nKey } from '@/components/features/project/project-navigation';
import type { ScopedTranslations } from '@/i18n/shared';

type ProjectT = ScopedTranslations<'project'>;
type ProjectKey = Parameters<ProjectT>[0];

export function getProjectModuleCopy(
  t: ProjectT,
  moduleKey: ProjectWorkspaceModuleI18nKey,
  field: 'label' | 'shortLabel' | 'description'
) {
  return t(`modules.${moduleKey}.${field}` as ProjectKey);
}
