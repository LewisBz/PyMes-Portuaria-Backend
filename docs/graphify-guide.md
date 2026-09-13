# Graphify — Guia de instalacion y uso

Graphify convierte el repo (codigo, docs, papers, imagenes) en un **grafo de conocimiento** consultable: nodos, aristas, comunidades y un informe auditable.

En este proyecto el CLI **ya esta instalado**. El grafo **aun no se ha construido**: no existe `graphify-out/` hasta que corras el primer mapeo.

---

## Que produce

Tras un mapeo completo, los artefactos quedan en `graphify-out/`:

| Archivo | Que es |
|---------|--------|
| `graph.html` | Grafo interactivo (abrir en el navegador) |
| `GRAPH_REPORT.md` | Informe: god nodes, conexiones sorprendentes, preguntas sugeridas |
| `graph.json` | Datos crudos (GraphRAG, queries, path, explain) |

Opcionales: vault Obsidian (`--obsidian`), wiki (`--wiki`), SVG, GraphML, Neo4j/FalkorDB.

---

## Prerrequisitos

| Herramienta | Para que | Verificar |
|-------------|----------|-----------|
| [uv](https://docs.astral.sh/uv/) | Instalar el CLI | `uv --version` |
| Python 3.11+ | Runtime (uv lo puede bajar) | `python --version` |
| [Cursor](https://cursor.com) | Skill `/graphify` y regla del proyecto | Abrir este repo en Cursor |

**API key:** no hace falta para mapear **codigo**. Graphify extrae AST de forma estructural, sin LLM.

`GEMINI_API_KEY` o `GOOGLE_API_KEY` solo sirven para extraccion semantica de **docs, papers e imagenes**. Si no hay key, el agente de Cursor hace esa parte. Graphify **no** lee `ANTHROPIC_API_KEY` ni `OPENAI_API_KEY`.

---

## Paso 1: Instalar el CLI

El paquete en PyPI se llama **`graphifyy`**; el comando es **`graphify`**.

```powershell
uv tool install --upgrade graphifyy
```

Sin uv:

```powershell
pip install graphifyy
```

Verifica:

```powershell
graphify --help
```

> El `python` del sistema puede no importar `graphify`. El interprete correcto es el de `uv tool` (junto al ejecutable `graphify`).

### Actualizar

```powershell
uv tool install --upgrade graphifyy
```

### Desinstalar

```powershell
uv tool uninstall graphifyy
```

---

## Paso 2: Cablear Cursor (este repo)

En la raiz del proyecto:

```powershell
graphify install --platform cursor
```

Equivale a `graphify cursor install` y escribe [`.cursor/rules/graphify.mdc`](../.cursor/rules/graphify.mdc): el agente usa el grafo **cuando exista** `graphify-out/graph.json`. Si el grafo no esta, sigue con Read/Grep/Glob.

Para quitar la regla:

```powershell
graphify cursor uninstall
```

Si Cursor avisa que el skill global esta desfasado respecto al paquete:

```powershell
graphify install --platform claude
```

(actualiza el skill en `~/.claude/skills/graphify`; un `graphify install` solo refresca la plataforma detectada.)

---

## Paso 3: Primer mapeo (cuando quieras)

Aun no se ha corrido. Cuando quieras mapear el backend:

1. Abre Agent en Cursor en la **raiz del repo**.
2. Invoca el skill:

```
/graphify .
```

El agente detecta archivos, extrae AST (codigo), agrupa comunidades y genera `graphify-out/`. Un corpus solo de codigo no necesita API key.

Variantes utiles:

```
/graphify . --mode deep
/graphify . --update
/graphify . --no-viz
```

Si hay **mas de ~500 archivos** o **>2M palabras**, Graphify pedira acotar a una subcarpeta. `--code-only` (CLI headless) indexa solo codigo.

### CLI headless (CI / scripts)

```powershell
graphify extract . --code-only
graphify export html
```

---

## Consultar el grafo

Solo tiene sentido **despues** de que exista `graphify-out/graph.json`.

En chat (el agente debe usar el grafo, no reexplorar el repo a ciegas):

```
/graphify query "Como fluye una request HTTP hasta la base de datos?"
/graphify path "AuthModule" "Database"
/graphify explain "NombreDeUnSimbolo"
```

Desde PowerShell, en la raiz del repo:

```powershell
graphify query "Como funciona la autenticacion?"
graphify query "trace X" --dfs --budget 1500
graphify path "SimboloA" "SimboloB"
graphify explain "Simbolo"
graphify god-nodes --top 10
```

Tras cambiar codigo:

```powershell
graphify update .
```

Eso re-extrae AST de archivos nuevos o modificados, sin LLM.

---

## Git

Hoy **no hay** `.gitignore` en la raiz. Cuando exista `graphify-out/`, anade:

```
graphify-out/
```

Es salida generada (grande, local). No hace falta versionarla salvo que el equipo quiera compartir el grafo.

---

## Comandos frecuentes

| Comando | Descripcion |
|---------|-------------|
| `/graphify .` | Pipeline completo sobre el directorio actual |
| `/graphify . --update` | Incremental (codigo cambiado) |
| `/graphify . --cluster-only` | Reagrupar comunidades sobre el grafo existente |
| `graphify query "..."` | BFS sobre `graph.json` |
| `graphify path "A" "B"` | Camino mas corto entre dos nodos |
| `graphify explain "X"` | Vecinos y contexto de un nodo |
| `graphify export html` | Regenerar `graph.html` |
| `graphify hook install` | Post-commit: actualiza el grafo al hacer commit |
| `graphify cursor install` | Regla `.cursor/rules/graphify.mdc` |

---

## Tips

- **Codigo primero:** el primer mapeo de este backend puede ser `--code-only` o el skill `/graphify .` sin Gemini.
- **No pidas API key** para un corpus solo de codigo.
- **No abras `graph.html` en grafos de >5000 nodos** sin `--no-viz` o aviso: el visor se pone pesado.
- **Honesty:** Graphify etiqueta aristas EXTRACTED / INFERRED / AMBIGUOUS; no inventes aristas al interpretar el informe.

---

## Links utiles

- [Repositorio Graphify](https://github.com/safishamsi/graphify)
- [Sponsors](https://github.com/sponsors/safishamsi)
- Skill local (Claude): `~/.claude/skills/graphify/SKILL.md`
- Spec Kit en este repo: [spec-kit-guide.md](spec-kit-guide.md)
