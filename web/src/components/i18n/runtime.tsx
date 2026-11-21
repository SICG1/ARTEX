"use client";

import * as React from "react";

import { EVENTO_IDIOMA, idiomaActual, type Lang } from "@/lib/i18n";
import { DICCIONARIO_ES } from "@/lib/i18n/diccionario";

const ATRIBUTOS = ["placeholder", "title", "aria-label", "alt"] as const;
// INPUT/TEXTAREA/SELECT no tienen nodos de texto; se dejan fuera para poder traducir sus atributos.
const ETIQUETAS_IGNORADAS = new Set(["SCRIPT", "STYLE", "CODE", "PRE", "KBD", "SAMP", "SVG", "PATH"]);
const MAX_LARGO = 120;

function traducirTexto(texto: string): string | null {
  const limpio = texto.trim();
  if (!limpio || limpio.length > MAX_LARGO) return null;
  const destino = DICCIONARIO_ES[limpio];
  if (!destino || destino === limpio) return null;
  return texto.replace(limpio, destino);
}

function debeIgnorar(el: Element | null): boolean {
  let actual = el;
  while (actual) {
    if (ETIQUETAS_IGNORADAS.has(actual.tagName)) return true;
    if (actual.hasAttribute("data-no-i18n")) return true;
    if (actual.getAttribute("contenteditable") === "true") return true;
    actual = actual.parentElement;
  }
  return false;
}

function traducirNodo(nodo: Node) {
  if (nodo.nodeType === Node.TEXT_NODE) {
    const padre = nodo.parentElement;
    if (!padre || debeIgnorar(padre)) return;
    const nuevo = traducirTexto(nodo.nodeValue || "");
    if (nuevo && nuevo !== nodo.nodeValue) nodo.nodeValue = nuevo;
    return;
  }
  if (nodo.nodeType !== Node.ELEMENT_NODE) return;
  const el = nodo as Element;
  if (debeIgnorar(el)) return;
  for (const attr of ATRIBUTOS) {
    const valor = el.getAttribute(attr);
    if (!valor) continue;
    const nuevo = traducirTexto(valor);
    if (nuevo && nuevo !== valor) el.setAttribute(attr, nuevo);
  }
  const walkerEl = document.createTreeWalker(el, NodeFilter.SHOW_ELEMENT);
  let elHijo = walkerEl.nextNode();
  while (elHijo) {
    const hijo = elHijo as Element;
    if (!debeIgnorar(hijo)) {
      for (const attr of ATRIBUTOS) {
        const valor = hijo.getAttribute(attr);
        if (!valor) continue;
        const nuevo = traducirTexto(valor);
        if (nuevo && nuevo !== valor) hijo.setAttribute(attr, nuevo);
      }
    }
    elHijo = walkerEl.nextNode();
  }
  const walker = document.createTreeWalker(el, NodeFilter.SHOW_TEXT);
  let actual = walker.nextNode();
  while (actual) {
    const padre = actual.parentElement;
    if (padre && !debeIgnorar(padre)) {
      const nuevo = traducirTexto(actual.nodeValue || "");
      if (nuevo && nuevo !== actual.nodeValue) actual.nodeValue = nuevo;
    }
    actual = walker.nextNode();
  }
}

/**
 * Traduce la interfaz al español cuando el idioma activo es "es".
 * Trabaja sobre el DOM porque la app original tiene los textos fijos en ingles.
 */
export function I18nRuntime() {
  const [lang, setLang] = React.useState<Lang>("en");

  React.useEffect(() => {
    const aplicar = (nuevo: Lang) => {
      setLang(nuevo);
      document.documentElement.lang = nuevo;
    };
    aplicar(idiomaActual());
    const alCambiar = (e: Event) => aplicar(((e as CustomEvent).detail as Lang) || idiomaActual());
    window.addEventListener(EVENTO_IDIOMA, alCambiar);
    return () => window.removeEventListener(EVENTO_IDIOMA, alCambiar);
  }, []);

  React.useEffect(() => {
    if (lang !== "es") return;
    document.documentElement.lang = "es";
    const pendientes = new Set<Node>();
    let timer: ReturnType<typeof setTimeout> | null = null;
    const vaciar = () => {
      timer = null;
      for (const n of pendientes) traducirNodo(n);
      pendientes.clear();
    };
    const programar = (n: Node) => {
      pendientes.add(n);
      if (!timer) timer = setTimeout(vaciar, 60);
    };
    programar(document.body);
    const obs = new MutationObserver((mutaciones) => {
      for (const m of mutaciones) {
        if (m.type === "characterData" && m.target) programar(m.target.parentElement || m.target);
        if (m.type === "attributes" && m.target) programar(m.target);
        for (const n of Array.from(m.addedNodes)) programar(n);
      }
    });
    obs.observe(document.body, {
      childList: true,
      subtree: true,
      characterData: true,
      attributes: true,
      attributeFilter: ["placeholder", "title", "aria-label", "alt"]
    });
    return () => {
      obs.disconnect();
      if (timer) clearTimeout(timer);
    };
  }, [lang]);

  return null;
}