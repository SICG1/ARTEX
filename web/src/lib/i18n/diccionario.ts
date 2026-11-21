import es from "@/i18n/locales/es.json";

// Spanish UI dictionary.
//
// The source of truth is the "es" locale in i18n/locales/es.json. This module keeps
// the dictionary exported for the DOM translation runtime (components/i18n/runtime.tsx),
// which covers the few strings that still render without going through t().
export const DICCIONARIO_ES: Record<string, string> = es as Record<string, string>;
