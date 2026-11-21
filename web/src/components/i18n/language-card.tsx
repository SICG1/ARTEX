"use client";

import * as React from "react";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { type Locale, LOCALES, locale, setLocale } from "@/i18n";
import { idiomaDelNavegador, hayPreferencia } from "@/lib/i18n";

export function LanguageCard() {
  const [valor, setValor] = React.useState<string>("auto");

  React.useEffect(() => {
    setValor(hayPreferencia() ? locale : "auto");
  }, []);

  const cambiar = (v: string) => {
    if (v === "auto") {
      try {
        window.localStorage.removeItem("artex.locale");
        window.localStorage.removeItem("artex.lang");
      } catch {
        /* almacenamiento bloqueado */
      }
      window.location.reload();
      return;
    }
    setLocale(v as Locale);
    setValor(v);
  };

  return (
    <Card className="mb-4 break-inside-avoid md:mb-6">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-base">Idioma de la interfaz</CardTitle>
        <CardDescription>
          Por defecto se usa el idioma del navegador. Puedes forzar español, inglés u otro idioma disponible.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex items-center justify-between gap-4">
        <div className="text-sm text-muted-foreground">
          Idioma detectado en el navegador: <b>{idiomaDelNavegador() === "es" ? "español" : "inglés"}</b>
        </div>
        <Select value={valor} onValueChange={cambiar}>
          <SelectTrigger className="w-56">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="auto">Automático (navegador)</SelectItem>
            {LOCALES.map((l) => (
              <SelectItem key={l.code} value={l.code}>
                {l.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </CardContent>
    </Card>
  );
}
