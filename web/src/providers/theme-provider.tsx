'use client';

import * as React from 'react';

type Theme = 'light' | 'dark' | 'system';
type ResolvedTheme = Exclude<Theme, 'system'>;
type Attribute = `data-${string}` | 'class';

interface ValueObject {
  [themeName: string]: string;
}

export interface ThemeProviderProps extends React.PropsWithChildren {
  themes?: Theme[];
  forcedTheme?: Theme;
  enableSystem?: boolean;
  disableTransitionOnChange?: boolean;
  enableColorScheme?: boolean;
  storageKey?: string;
  defaultTheme?: Theme;
  attribute?: Attribute | Attribute[];
  value?: ValueObject;
}

export interface UseThemeProps {
  themes: Theme[];
  forcedTheme?: Theme;
  setTheme: React.Dispatch<React.SetStateAction<Theme>>;
  theme?: Theme;
  resolvedTheme?: ResolvedTheme;
  systemTheme?: ResolvedTheme;
}

const ThemeContext = React.createContext<UseThemeProps | undefined>(undefined);

const DEFAULT_THEMES: Theme[] = ['light', 'dark'];
const DEFAULT_STORAGE_KEY = 'theme';

function getSystemTheme(): ResolvedTheme {
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function setAttribute(
  attribute: Attribute,
  theme: ResolvedTheme,
  themes: Theme[],
  value?: ValueObject
) {
  const root = document.documentElement;
  const mappedTheme = value?.[theme] ?? theme;
  const mappedThemes = themes.map((entry) => value?.[entry] ?? entry);

  if (attribute === 'class') {
    root.classList.remove(...mappedThemes);
    root.classList.add(mappedTheme);
    return;
  }

  root.setAttribute(attribute, mappedTheme);
}

function setColorScheme(theme: Theme, enableColorScheme: boolean, defaultTheme: Theme) {
  if (!enableColorScheme) {
    return;
  }

  const root = document.documentElement;
  const effectiveTheme =
    theme === 'system'
      ? getSystemTheme()
      : theme === undefined
        ? defaultTheme === 'system'
          ? getSystemTheme()
          : defaultTheme
        : theme;

  root.style.colorScheme = effectiveTheme;
}

function disableTransitionsTemporarily() {
  const style = document.createElement('style');
  style.appendChild(
    document.createTextNode(
      '*,*::before,*::after{-webkit-transition:none!important;transition:none!important}'
    )
  );
  document.head.appendChild(style);

  return () => {
    window.getComputedStyle(document.body);
    window.setTimeout(() => {
      document.head.removeChild(style);
    }, 1);
  };
}

export function ThemeProvider({
  children,
  themes = DEFAULT_THEMES,
  forcedTheme,
  enableSystem = true,
  disableTransitionOnChange = false,
  enableColorScheme = true,
  storageKey = DEFAULT_STORAGE_KEY,
  defaultTheme = enableSystem ? 'system' : 'light',
  attribute = 'data-theme',
  value,
}: ThemeProviderProps) {
  const [theme, setThemeState] = React.useState<Theme | undefined>(undefined);
  const [systemTheme, setSystemTheme] = React.useState<ResolvedTheme | undefined>(undefined);

  const attributes = React.useMemo(
    () => (Array.isArray(attribute) ? attribute : [attribute]),
    [attribute]
  );

  const setTheme = React.useCallback<React.Dispatch<React.SetStateAction<Theme>>>(
    (nextTheme) => {
      setThemeState((currentTheme) => {
        const resolvedTheme = typeof nextTheme === 'function' ? nextTheme(currentTheme ?? defaultTheme) : nextTheme;

        try {
          window.localStorage.setItem(storageKey, resolvedTheme);
        } catch {}

        return resolvedTheme;
      });
    },
    [defaultTheme, storageKey]
  );

  React.useEffect(() => {
    setSystemTheme(getSystemTheme());

    try {
      const storedTheme = window.localStorage.getItem(storageKey) as Theme | null;
      setTheme(storedTheme ?? defaultTheme);
    } catch {
      setTheme(defaultTheme);
    }

    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handleChange = () => {
      const nextSystemTheme = getSystemTheme();
      setSystemTheme(nextSystemTheme);
    };

    handleChange();
    mediaQuery.addEventListener('change', handleChange);

    return () => mediaQuery.removeEventListener('change', handleChange);
  }, [defaultTheme, setTheme, storageKey]);

  React.useEffect(() => {
    const configuredTheme = forcedTheme ?? theme ?? defaultTheme;
    const activeTheme =
      configuredTheme === 'system' ? systemTheme ?? getSystemTheme() : configuredTheme;
    const cleanup = disableTransitionOnChange ? disableTransitionsTemporarily() : null;

    attributes.forEach((entry) => setAttribute(entry, activeTheme, themes, value));
    setColorScheme(configuredTheme, enableColorScheme, defaultTheme);

    cleanup?.();
  }, [
    attributes,
    defaultTheme,
    disableTransitionOnChange,
    enableColorScheme,
    forcedTheme,
    systemTheme,
    theme,
    themes,
    value,
  ]);

  const contextValue = React.useMemo<UseThemeProps>(() => {
    const resolvedTheme =
      forcedTheme && forcedTheme !== 'system'
        ? forcedTheme
        : theme === 'system'
          ? systemTheme
          : theme === 'light' || theme === 'dark'
            ? theme
            : systemTheme;

    return {
      theme,
      setTheme,
      forcedTheme,
      resolvedTheme,
      systemTheme,
      themes: enableSystem ? [...themes, 'system'] : themes,
    };
  }, [enableSystem, forcedTheme, setTheme, systemTheme, theme, themes]);

  return <ThemeContext.Provider value={contextValue}>{children}</ThemeContext.Provider>;
}

export function useTheme(): UseThemeProps {
  const context = React.useContext(ThemeContext);

  if (context) {
    return context;
  }

  return {
    setTheme: () => undefined,
    themes: [...DEFAULT_THEMES, 'system'],
    theme: 'system',
    resolvedTheme: undefined,
    systemTheme: undefined,
    forcedTheme: undefined,
  };
}
