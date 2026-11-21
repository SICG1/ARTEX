#!/bin/sh
# Llama a la API de ARTEX autenticándose con la clave guardada en ARTEX_PASS_FILE.
# Uso: ./artex-api.sh GET /api/tasks
#      ./artex-api.sh POST /api/tasks '{"description":"...","goal":"..."}'
# La clave y el token nunca se imprimen.
set -u
PASS_FILE="${ARTEX_PASS_FILE:-/tmp/opencode/artex-pass}"
BASE="${ARTEX_URL:-http://127.0.0.1:8787}"
METHOD="${1:?uso: artex-api.sh METODO RUTA [JSON]}"
PATHQ="${2:?falta la ruta}"
BODY="${3:-}"

[ -r "$PASS_FILE" ] || { echo "no puedo leer $PASS_FILE" >&2; exit 1; }

TOKEN=$(python3 - "$BASE" "$PASS_FILE" <<'PY'
import json,sys,urllib.request
base,passfile=sys.argv[1],sys.argv[2]
pw=open(passfile).read().strip()
req=urllib.request.Request(base+"/api/auth/login",
    data=json.dumps({"username":"ARTEX","password":pw}).encode(),
    headers={"Content-Type":"application/json"})
try:
    r=json.load(urllib.request.urlopen(req,timeout=10))
except Exception as e:
    print("LOGIN_FAIL",file=sys.stderr); sys.exit(1)
tok=r.get("token") or r.get("access_token") or ""
if not tok:
    print("SIN_TOKEN",file=sys.stderr); sys.exit(1)
print(tok)
PY
) || { echo "login falló (¿clave correcta?)" >&2; exit 1; }

if [ -n "$BODY" ]; then
  curl -s -m 120 -X "$METHOD" "$BASE$PATHQ" \
    -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d "$BODY"
else
  curl -s -m 120 -X "$METHOD" "$BASE$PATHQ" -H "Authorization: Bearer $TOKEN"
fi
