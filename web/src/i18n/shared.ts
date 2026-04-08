import type { Messages } from './modules';

export type TranslationValues = Record<string, string | number | Date>;

export type Translators = {
  [K in keyof Messages]: (
    key: DotNotationKeys<Messages[K]>,
    values?: TranslationValues
  ) => string;
};

type DotNotationKeys<T, Prefix extends string = ''> = T extends object
  ? {
      [K in keyof T]: K extends string
        ? T[K] extends object
          ? `${Prefix}${K}` | DotNotationKeys<T[K], `${Prefix}${K}.`>
          : `${Prefix}${K}`
        : never;
    }[keyof T]
  : never;

export type AllTranslationKeys = DotNotationKeys<Messages>;
export type AllScopePaths = DotNotationKeys<Messages>;

type ShiftingKeys<T, P extends string> = P extends `${infer Head}.${infer Tail}`
  ? Head extends keyof T ? ShiftingKeys<T[Head], Tail> : never
  : P extends keyof T ? DotNotationKeys<T[P]> : never;

export type UnifiedTranslations = {
  (key: AllTranslationKeys, values?: TranslationValues): string;
} & Translators;

export type ScopedTranslations<P extends string> = (
  key: ShiftingKeys<Messages, P>,
  values?: TranslationValues
) => string;
