"use client";

import { useEffect, useState } from "react";

const ZONA = "America/Santiago";
const formatoHora = new Intl.DateTimeFormat("es-CL", {
  timeZone: ZONA,
  hour: "2-digit",
  minute: "2-digit",
  second: "2-digit",
  timeZoneName: "short",
});
const formatoFecha = new Intl.DateTimeFormat("es-CL", {
  timeZone: ZONA,
  weekday: "long",
  day: "numeric",
  month: "long",
  year: "numeric",
});

// Hora de Chile en la cabecera. Se actualiza cada segundo, solo con la pestana visible.
export function RelojChile() {
  const [ahora, setAhora] = useState<Date | null>(null);

  useEffect(() => {
    const tic = () => {
      if (document.visibilityState === "visible") setAhora(new Date());
    };
    tic();
    const t = setInterval(tic, 1000);
    document.addEventListener("visibilitychange", tic);
    return () => {
      clearInterval(t);
      document.removeEventListener("visibilitychange", tic);
    };
  }, []);

  // Render nothing until mounted, so the server and the client markup match.
  if (!ahora) return null;
  return (
    <span className="font-medium text-muted-foreground text-xs tabular-nums" title={formatoFecha.format(ahora)}>
      {formatoHora.format(ahora)}
    </span>
  );
}
