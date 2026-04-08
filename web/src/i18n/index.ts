// i18n barrel export
export { 
  locales, 
  defaultLocale, 
  localeNames, 
  localeMapping,
  isLocaleSwitcherEnabled,
  type Locale 
} from './config';

export { messages, type Messages } from './modules';
export { loadAllModules, loadModule } from './loader';
export type { ScopedTranslations, UnifiedTranslations } from './shared';
