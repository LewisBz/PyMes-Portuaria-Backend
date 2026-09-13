# Spec Kit — Guia paso a paso (Cursor + opencode)

## Que es Spec-Driven Development?

Spec-Driven Development (SDD) es un enfoque donde **las especificaciones se vuelven ejecutables**. En vez de escribir specs que se descartan despues, Spec Kit las usa para generar codigo funcional directamente a traves de un agente de IA.

En resumen: defines **que** quieres construir, y el agente se encarga del **como**.

Este repo tiene Spec Kit inicializado. La integracion **default es Cursor** (`cursor-agent`). **opencode sigue instalado** para usarlo desde ese agente.

| Agente | Clave CLI | Archivos | Como invocar |
|--------|-----------|----------|--------------|
| Cursor (default) | `cursor-agent` | `.cursor/skills/speckit-*/SKILL.md` | `/speckit-specify`, `/speckit-plan`, ... |
| opencode | `opencode` | `.opencode/commands/speckit.*.md` | `/speckit.specify`, `/speckit.plan`, ... |

---

## Prerrequisitos

Antes de empezar necesitas tener instalado:

| Herramienta | Version minima | Verificar |
|-------------|---------------|-----------|
| [uv](https://docs.astral.sh/uv/) | Cualquier version reciente | `uv --version` |
| [Python](https://www.python.org/downloads/) | 3.11+ | `python --version` |
| [Git](https://git-scm.com/downloads) | 2.20+ | `git --version` |
| [Cursor](https://cursor.com) | IDE con Agent | Abrir este proyecto en Cursor |
| [opencode](https://opencode.ai) (opcional) | Instalado y funcionando | Solo si trabajas desde opencode |

> **Nota:** Si usas `uv`, no necesitas instalar Python manualmente. `uv` lo descarga automaticamente cuando lo necesita.

---

## Paso 1: Instalar uv

`uv` es un package manager rapido de Python. Si aun no lo tienes:

### Windows (PowerShell)

```powershell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
irm https://astral.sh/uv/install.ps1 | iex
```

Despues de instalar, **reinicia la shell** o ejecuta:

```powershell
$env:Path = "C:\Users\<tu-usuario>\.local\bin;$env:Path"
```

### macOS / Linux

```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
```

Verifica la instalacion:

```bash
uv --version
```

---

## Paso 2: Instalar Spec Kit CLI

Instala `specify` como herramienta global con `uv tool install`:

```bash
uv tool install specify-cli --from git+https://github.com/github/spec-kit.git
```

> **Alternativa (sin instalar):** Si solo quieres probarlo una vez:
> ```bash
> uvx --from git+https://github.com/github/spec-kit.git specify init mi-proyecto
> ```

Verifica que funciona:

```bash
specify version
specify check
```

### Actualizar Spec Kit

Para actualizar a la ultima version:

```bash
uv tool install specify-cli --force --from git+https://github.com/github/spec-kit.git
```

### Desinstalar

```bash
uv tool uninstall specify-cli
```

---

## Paso 3: Inicializar o cambiar de agente

### Proyecto nuevo (Cursor)

```bash
specify init mi-proyecto --integration cursor-agent --script ps
cd mi-proyecto
```

### Proyecto existente (directorio actual)

```bash
specify init . --here --integration cursor-agent --script ps
```

> **Cuidado:** `--force` en `specify init` puede pisar `.specify/` (templates, scripts, constitution). En un repo que **ya** tiene Spec Kit, no reinicialices: usa los comandos de integracion.

### Este proyecto ya esta inicializado

No hace falta `specify init`. Para anadir Cursor **sin quitar** opencode:

```powershell
specify integration install cursor-agent --script ps --force
specify integration use cursor-agent
specify integration list
```

`--force` en `install` solo autoriza tener dos integraciones a la vez (opencode no es multi-install-safe). `--force` en `use` **no** es necesario si quieres conservar scripts y templates ya customizados.

Esto deja:

```
.specify/                 # Infra compartida (no se reescribe al anadir Cursor)
.cursor/skills/           # Skills de Spec Kit para Cursor
.opencode/commands/       # Slash commands de Spec Kit para opencode
```

Estructura de `.specify/`:

```
.specify/
├── integrations/       # Manifests de cada agente
├── memory/             # Constitution y memoria del proyecto
├── scripts/            # Scripts de automatizacion (PowerShell en este repo)
├── templates/          # Templates para specs, plans y tasks
├── workflows/          # Workflow de Spec Kit
└── init-options.json   # Opciones de inicializacion
```

Para volver a opencode como default (sin desinstalar Cursor):

```powershell
specify integration use opencode
```

---

## Paso 4: El workflow completo

En **Cursor**, abre Agent y usa los skills en este orden (separador `-`). En **opencode**, los mismos pasos con punto (`/speckit.specify`).

### 1. `/speckit-constitution` — Establecer principios

Crea el documento fundacional del proyecto. Define estandares de codigo, testing, arquitectura y gobernanza.

```
/speckit-constitution Define principios para nuestro proyecto: codigo limpio,
testeo obligatorio, documentacion clara, y arquitectura modular.
```

### 2. `/speckit-specify` — Definir que construir

Describe la feature o funcionalidad que quieres implementar.

```
/speckit-specify Quiero un endpoint REST para gestionar usuarios con
autenticacion JWT, CRUD completo, y validacion de datos.
```

### 3. `/speckit-clarify` — Resolver ambiguedades (opcional)

El agente te hace preguntas estructuradas para aclarar puntos antes de planificar.

```
/speckit-clarify
```

### 4. `/speckit-plan` — Plan tecnico

Genera la arquitectura, stack tecnico, y estructura del codigo a implementar.

```
/speckit-plan
```

### 5. `/speckit-tasks` — Dividir en tareas

Descompone el plan en tareas accionables y ordenadas.

```
/speckit-tasks
```

### 6. `/speckit-implement` — Ejecutar

El agente valida todos los artefactos y ejecuta la implementacion en orden.

```
/speckit-implement
```

### 7. `/speckit-converge` — Evaluar gaps (opcional)

Evalua el codigo generado contra el spec original y agrega tareas pendientes si hay gaps.

```
/speckit-converge
```

### Flujo visual

```
constitution → specify → clarify → plan → tasks → implement → converge
     ↑                          ↑                 ↑
  (obligatorio)            (opcional)         (opcional)
```

---

## Comandos utiles

| Comando | Descripcion |
|---------|-------------|
| `specify version` | Mostrar version instalada |
| `specify check` | Verificar dependencias y configuracion |
| `specify self upgrade` | Auto-actualizar a la ultima version |
| `specify integration list` | Listar integraciones (marca default e instaladas) |
| `specify integration install <key>` | Instalar un agente (usa `--force` si ya hay otro no multi-install-safe) |
| `specify integration use <key>` | Cambiar el default **sin** desinstalar el otro |
| `specify integration switch <key>` | Cambia de agente **desinstalando** el anterior — no usar si quieres ambos |
| `specify extension search` | Buscar extensions comunitarias |

---

## Tips

- **Escribe specs claros:** Cuanto mas especifico seas en `/speckit-specify`, mejor sera el resultado.
- **Revisa el plan antes de implementar:** Si el plan no te convence, edita el spec y vuelve a `/speckit-plan`.
- **Usa `/speckit-clarify`** si tu spec tiene areas ambiguas antes de generar el plan.
- **La constitution es opcional pero recomendada:** Define estandares que el agente respetara durante toda la implementacion.
- **Agrega `.specify/` a `.gitignore`** si contiene datos sensibles del agente.
- **No uses `specify init --here --force`** en este repo para “pasar a Cursor”: pisa infra compartida. Usa `integration install` + `integration use`.

---

## Links utiles

- [Documentacion oficial](https://github.github.io/spec-kit/)
- [Repositorio GitHub](https://github.com/github/spec-kit)
- [Releases](https://github.com/github/spec-kit/releases)
- [Integraciones soportadas](https://github.github.io/spec-kit/reference/integrations.html)
- [uv docs](https://docs.astral.sh/uv/)
