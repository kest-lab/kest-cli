// Type augmentation for next-intl v4
// This enables type-safe translations with IDE auto-completion

import type { Messages } from './modules';

// Declare module augmentation for next-intl
// In v4, we augment the AppConfig interface instead of global IntlMessages
declare module 'next-intl' {
  interface AppConfig {
    Messages: Messages;
  }
}
