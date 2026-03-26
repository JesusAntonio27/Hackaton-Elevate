# CONTEXT_IA.md — NPS Sentinel · Guía de integración y consolidación

> **Propósito:** Este documento le da a una IA toda la información arquitectónica necesaria para integrar el backend del proyecto `nps-hck` con el frontend refactorizado de `Maquetado-FrontEnd`, consolidando todo en el repositorio raíz `Hackaton-Elevate`. No es necesario leer el código fuente de cada archivo.

---

## 1. Visión general del repositorio

```
Hackaton-Elevate/                  ← REPO PRINCIPAL (destino final)
├── .env                           ← API key (ANTHROPIC_API_KEY=...)
├── Logo-NPS-Sentinel.svg          ← Asset de branding principal
├── Logo-NPS-Sentinel.png          ← Versión PNG del logo
├── datos_analizados.json          ← Dataset NPS procesado por IA (~619 KB)
├── datos_nps.json                 ← Dataset NPS crudo (~404 KB)
├── procesador_asincrono.py        ← Script Python que llama Anthropic API para etiquetar reseñas como "problemáticas"
├── test_nps_label.py              ← Script de prueba del etiquetador
├── JSONS-Encuestas/               ← JSONs de encuestas individuales (fuente de datos)
│
├── Maquetado-FrontEnd/            ← FRONT MÁS ACTUALIZADO (fuente de verdad del UI)
│   ├── dashboard.html             ← Estructura HTML del dashboard (486 líneas, limpio)
│   ├── dashboard.css              ← Todos los estilos del dashboard (refactorizado)
│   ├── dashboard.js               ← Toda la lógica JS del dashboard (476 líneas, refactorizado)
│   ├── chatbot.html               ← Interfaz del chatbot RAG (~27 KB)
│   ├── encuesta.html              ← Formulario de encuesta NPS
│   ├── confirmacion.html          ← Pantalla de confirmación post-encuesta
│   ├── Logo-NPS-Sentinel.svg      ← Copia del logo (misma que raíz)
│   ├── Logo-NPS-Sentinel.png      ← Copia PNG del logo
│   ├── README_IA_API.md           ← Documentación de la API de IA para chatbot
│   ├── Maquetado-FrontEnd/        ← Subcarpeta adicional (ignorar, es residual)
│   └── ui-elements/               ← Assets de UI adicionales
│
└── nps-hck/                       ← VERSION ANTERIOR integrada con backend Go
    ├── main.go                    ← Backend Go completo (1432 líneas)
    ├── go.mod / go.sum            ← Dependencias Go
    ├── dashboard_general.html     ← VERSIÓN ANTIGUA del dashboard (3121 líneas, TODO inline)
    ├── respaldo_resenas.json      ← Backup de reseñas (~205 KB)
    ├── Material.md                ← Documentación/notas del hackathon
    └── dashboard/
        ├── index.html             ← Dashboard prototipo temprano
        ├── datos_nps.json         ← Copia local de datos NPS
        ├── insertar_datos.py      ← Script para insertar datos en MongoDB
        └── main copy.go           ← Versión anterior del backend
```

---

## 2. El problema: dos versiones del dashboard

| Aspecto | `nps-hck/dashboard_general.html` | `Maquetado-FrontEnd/dashboard.html` |
|---|---|---|
| **Estado** | Versión "integrada" (antigua) | **Fuente de verdad del UI** (más reciente) |
| **Tamaño** | 3121 líneas — CSS+JS inline monolítico | 486 líneas HTML + archivos separados |
| **CSS** | Inline en `<style>` (>2000 líneas) | Separado en `dashboard.css` |
| **JS** | Inline en `<script>` (integrado con backend) | Separado en `dashboard.js` |
| **Logo** | Logo antiguo "Coppel TI" con texto en HTML | Logo SVG `Logo-NPS-Sentinel.svg` correcto |
| **Funcionalidades exclusivas** | Ver §3 | Ver §4 |

---

## 3. Funcionalidades presentes en `dashboard_general.html` que faltan en el front nuevo

Estas son las integraciones que se deben portar al front de `Maquetado-FrontEnd/`:

### 3.1 Conexión real con el backend Go (API en `localhost:8080`)
- Las llamadas fetch en el dashboard viejo apuntan a `http://localhost:8080/{endpoint}`
- En `dashboard.js` del front nuevo, solo se llama a `/pareto` y a `../datos_analizados.json` (fetch local)
- **Integración pendiente:** conectar los endpoints reales del backend Go

### 3.2 Tokens de diseño adicionales en el CSS viejo
El CSS del `dashboard_general.html` define variables CSS que el `dashboard.css` nuevo puede no tener o puede tener diferente:
- `--shadow-hover`, `--shadow-primary` usados en cards
- `.issue-subname` (clase de nombre de subcaussa en Pareto, **ausente en el HTML nuevo**)
- `.nav-badge` (badge de notificación en sidebar)
- `.active-focus` (resaltado de issue-item activo en drawer)

### 3.3 Lógica del drawer de "Foco Rojo" (ya portada en dashboard.js nuevo)
- `abrirDetalleFoco()` / `cerrarDetalleFoco()` — ✅ ya existe en `dashboard.js`
- El drawer HTML (`.detail-overlay`, `.detail-drawer`) — ✅ ya existe en `dashboard.html`

### 3.4 Area-grid horizontal slider
- ✅ Ya portado en `dashboard.js` (función `initAreaGridSlider`)

---

## 4. Funcionalidades exclusivas del front nuevo (`Maquetado-FrontEnd/`)

- Logo NPS Sentinel SVG correcto
- CSS separado limpio y mantenible
- JS separado limpio (`dashboard.js`)
- Carga de `datos_analizados.json` con `calcularStats()`, `tendenciaMensual()`, `actualizarDashboard()`
- Dropdown de selección de área (Soporte Técnico / Vista General) funcional
- Historial de reseñas renderizado desde JSON real (`renderHistorialResenas`)
- Chatbot (`chatbot.html`) con navegación sidebar funcional

---

## 5. Backend Go — `nps-hck/main.go`

### Stack
- **Lenguaje:** Go
- **Router:** gorilla/mux
- **Bases de datos:** MongoDB (datos estructurados) + Qdrant (búsqueda vectorial)
- **LLM:** Gemini 3 Flash vía API liteLLM proxy (`api.genius.coppel.services`)
- **Puerto:** `localhost:8080`
- **CORS:** habilitado para todos los orígenes (`*`)

### Modelos de datos clave

**`Resena`** — Entidad principal (colección MongoDB dinámica):
```
_id, collection, centro{nombre}, categoria{nombre}, historia, problema_pricipal,
problema_secundario, fecha_captura, fecha_cierre, nsp_score, estatus{nombre},
nivel_1{categoria, sub_categoria}, nivel_2{categoria, sub_categoria},
sentimientoIA{nombre}, problemaIA, npsIA, razonNpsIA, razonIA, clasificacionNPS{nombre},
iniciativa, iniciativaIA, centroSoporte{nombre}, respuesta, respuestaIA
```

**`IAEnrichment`** — Output del LLM al enriquecer una reseña:
```
id_respuesta, problema_principal, problema_secundario, npsIA, razonNpsIA, razonIA,
iniciativaIA, respuestaIA, nivel1_categoria, nivel1_subcat, nivel2_categoria,
nivel2_subcat, sentimiento_nombre, clasificacion_nombre, centro_inferido,
centro_soporte_dependencia
```

### Endpoints disponibles (todos en `localhost:8080`)

| Método | Ruta | Propósito |
|---|---|---|
| `GET` | `/nps/stats` | KPIs globales de NPS |
| `GET` | `/nps/stats/historico` | Tendencia mensual histórica |
| `GET` | `/nps/analisis/pareto` | Análisis de Pareto (top causas) |
| `POST` | `/nps/prediccion-estrategica` | Predicción IA por centro+categoría |
| `GET` | `/nps/resenas` | Listar reseñas (con filtros query params) |
| `DELETE` | `/nps/resenas/{id}` | Eliminar reseña |
| `POST` | `/nps/insertar` | Insertar reseña individual (enriquece con IA) |
| `POST` | `/nps/buscar` | Búsqueda semántica vectorial (RAG) |
| `POST` | `/nps/recomendar` | Recomendación estratégica IA por tema |
| `POST` | `/nps/procesar-lote` | Insertar lote de reseñas con enriquecimiento IA |
| `GET/POST/PUT/DELETE` | `/centros/{id?}` | CRUD catálogo centros |
| `GET/POST/PUT/DELETE` | `/categorias/{id?}` | CRUD catálogo categorías |
| `GET/POST/PUT/DELETE` | `/estatuses/{id?}` | CRUD catálogo estatus |
| `GET/POST/PUT/DELETE` | `/sentimientos/{id?}` | CRUD catálogo sentimientos |
| `GET/POST/PUT/DELETE` | `/clasificacionesNPS/{id?}` | CRUD clasificaciones NPS |
| `GET/POST` | `/colecciones` | Listar/registrar colecciones Qdrant |
| `DELETE` | `/colecciones/{nombre}` | Eliminar colección Qdrant |
| `DELETE` | `/admin/limpiar` | Borrar toda la base de datos |

### Filtros disponibles en `GET /nps/resenas`
Query params opcionales: `collection`, `categoria` (ObjectID), `estatus` (ObjectID), `sentimientoIA` (ObjectID), `clasificacionNPS` (ObjectID)

### Response de `/nps/analisis/pareto`
```json
{
  "items": [
    { "categoria": "str", "subcategoria": "str", "frecuencia": 10,
      "porcentaje": 88.0, "porcentaje_acumulado": 88.0, "areas_afectadas": ["Centro A"] }
  ]
}
```
> ⚠️ El `dashboard.js` nuevo ya consume `/pareto` en `cargarPareto()`. Solo hay que cambiar la URL a `http://localhost:8080/nps/analisis/pareto` y adaptar el shape de respuesta.

---

## 6. Dataset local — `datos_analizados.json`

Archivo JSON array en la **raíz del repo**. Cada elemento tiene:
```json
{
  "id_respuesta": "int",
  "fecha": "YYYY-MM-DD",
  "comentario": "string",
  "nps_score": 0-10,
  "area_colaborador": "Soporte Técnico | ...",
  "categoria_problema": "string",
  "clasificacion_nps": "Promotor | Pasivo | Detractor",
  "problematico": true | false   ← agregado por procesador_asincrono.py
}
```

`dashboard.js` lo carga con `fetch('../datos_analizados.json')` (ruta relativa desde `Maquetado-FrontEnd/`).

---

## 7. Python — `procesador_asincrono.py`

- Lee `datos_nps.json`
- Llama a Anthropic API (`claude-haiku-4-5-20251001`) vía `ANTHROPIC_API_KEY` del `.env`
- Agrega el campo `"problematico": bool` a cada reseña
- Escribe output en `datos_analizados.json`
- **No modifica el backend Go; opera como pipeline offline**

---

## 8. Tarea de consolidación — qué debe hacer la IA

### Objetivo
Consolidar en `Hackaton-Elevate/Maquetado-FrontEnd/` un dashboard que tenga:
1. El **diseño y estructura** del front nuevo (`dashboard.html` + `dashboard.css` + `dashboard.js`)
2. Las **integraciones reales** que tenía `dashboard_general.html` (backend Go en `localhost:8080`)
3. El repositorio raíz `Hackaton-Elevate` como único origen de verdad (sin depender de `nps-hck/`)

### Archivos que NO deben modificarse (solo lectura/referencia)
- `nps-hck/main.go` — el backend funciona tal cual
- `nps-hck/dashboard_general.html` — solo referencia para extraer integraciones

### Archivos que SÍ deben modificarse / crearse
| Archivo | Acción | Descripción |
|---|---|---|
| `Maquetado-FrontEnd/dashboard.js` | **Modificar** | Actualizar URLs de fetch para apuntar a `http://localhost:8080/nps/analisis/pareto` y adaptar shape de respuesta |
| `Maquetado-FrontEnd/dashboard.css` | **Revisar/complementar** | Verificar que tenga todas las clases usadas por `dashboard.html` (especialmente `.active-focus`, `.focus-chip`, `.drawer-*`, `.view-selector`, `.custom-dropdown`) |
| `Maquetado-FrontEnd/dashboard.html` | **Revisar** | Confirmar que el drawer HTML (`.detail-overlay`, `.detail-drawer`) y el area-grid existan |

### Verificaciones de integración
- `cargarPareto()` en `dashboard.js` → `GET http://localhost:8080/nps/analisis/pareto` → `data.items[]`
- `cargarDatos()` en `dashboard.js` → `fetch('../datos_analizados.json')` (ya funciona si corres desde `Hackaton-Elevate/`)
- El chatbot RAG en `chatbot.html` → usa su propia lógica (ver `README_IA_API.md`)

---

## 9. Estado actual del sistema (qué funciona, qué no)

| Componente | Estado |
|---|---|
| Backend Go (`main.go`) | ✅ Implementado, listo (requiere MongoDB + Qdrant corriendo) |
| Frontend dashboard (HTML/CSS/JS) | ✅ Refactorizado y limpio en `Maquetado-FrontEnd/` |
| Carga de `datos_analizados.json` | ✅ Funciona (fetch local) |
| Endpoint `/pareto` en dashboard | ⚠️ Apunta a `/pareto` sin host — debe corregirse a `localhost:8080/nps/analisis/pareto` |
| Chatbot RAG | ✅ Funcional con su propia lógica |
| Encuesta + Confirmación | ✅ HTML estático listo |
| `procesador_asincrono.py` | ✅ Funcional (requiere `ANTHROPIC_API_KEY` en `.env`) |
| Logo NPS Sentinel | ✅ SVG en raíz y en `Maquetado-FrontEnd/` |

---

## 10. Notas adicionales para la IA

- **No crear archivos nuevos innecesarios.** Todo el front ya existe en `Maquetado-FrontEnd/`.
- **`nps-hck/dashboard_general.html` debe eliminarse** después de verificar que todas sus integraciones estén portadas.
- El servidor de desarrollo local corre con `python3 -m http.server 8765` desde `Hackaton-Elevate/`. Las rutas relativas de los assets asumen esto.
- La variable `currentMode` en `dashboard.js` controla si se muestra `soporte` o `general`. La función `actualizarDashboard(modo)` filtra `allNPSData` por `area_colaborador === 'Soporte Técnico'`.
- Las clases del Pareto (`.issue-item`, `.issue-rank.top`, `.btn-details`) existen tanto en el CSS viejo como en el nuevo — son compatibles.
- El campo `problematico` de `datos_analizados.json` está disponible en cada reseña pero aún **no se usa en el dashboard** — es una extensión futura para los "focos rojos".
