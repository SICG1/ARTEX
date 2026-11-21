"use client";

// Language helpers used by the custom Spanish/English UI (language card + DOM runtime).
//
// The single source of truth is the upstream translation layer in "@/i18n": this module
// only exposes the small compatibility surface the custom components were written
// against, mapped onto the upstream locale state ("es" is a registered locale).

import { locale, setLocale } from "@/i18n";

export type Lang = "es" | "en";

const CLAVE = "artex.lang";
const CLAVE_UPSTREAM = "artex.locale";
export const EVENTO_IDIOMA = "artex:lang";

export function idiomaDelNavegador(): Lang {
  if (typeof navigator === "undefined") return "en";
  const idiomas = navigator.languages ?? [navigator.language];
  return idiomas.some((l) => (l || "").toLowerCase().startsWith("es")) ? "es" : "en";
}

/** Idioma activo: "es" solo cuando la locale efectiva de la app es español. */
export function idiomaActual(): Lang {
  return locale === "es" ? "es" : "en";
}

/** Fija el idioma de la interfaz (recarga la página, igual que el control de preferencias). */
export function guardarIdioma(lang: Lang) {
  try {
    window.localStorage.setItem(CLAVE, lang);
  } catch {
    /* almacenamiento bloqueado: setLocale lo volverá a intentar */
  }
  window.dispatchEvent(new CustomEvent(EVENTO_IDIOMA, { detail: lang }));
  setLocale(lang);
}

/** Vuelve a la detección automática por idioma del navegador. */
export function usarIdiomaDelNavegador() {
  try {
    window.localStorage.removeItem(CLAVE);
    window.localStorage.removeItem(CLAVE_UPSTREAM);
  } catch {
    /* almacenamiento bloqueado */
  }
  window.location.reload();
}

export function hayPreferencia(): boolean {
  if (typeof window === "undefined") return false;
  try {
    return window.localStorage.getItem(CLAVE_UPSTREAM) !== null;
  } catch {
    return false;
  }
}
