#!/bin/sh
# Arranque local de ARTEX en esta máquina (Kali, sin GPU).
# Uso: ./run-local.sh [modelo]
set -u
cd "$(dirname "$0")" || exit 1

MODEL="${1:-${ARTEX_LLM_MODEL:-qwen3:8b}}"

# Aísla del grupo de procesos de la terminal: sobrevive al cierre del shell.
setsid env \
  ARTEX_LLM_PROVIDER=ollama \
  ARTEX_LLM_MODEL="$MODEL" \
  ./start.sh -proxy 127.0.0.1:8798 \
  > artex.log 2>&1 < /dev/null &

sleep 6
if curl -s -m 5 http://127.0.0.1:8787/api/health | grep -q '"ok":true'; then
  echo "[run-local] ARTEX arriba en http://localhost:8787 (modelo $MODEL)"
else
  echo "[run-local] no respondió; revisa artex.log" >&2
  tail -5 artex.log >&2
  exit 1
fi
