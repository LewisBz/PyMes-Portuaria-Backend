# Feature Specification: API de gestión y trazabilidad de cargas (PyMEs portuarias)

**Feature Branch**: `001-cargo-traceability`

**Created**: 2026-09-13

**Status**: Draft

**Input**: User description: "Sistema web SaaS para gestión y trazabilidad de cargas en PyMEs logísticas portuarias (Universidad Libre Seccional Barranquilla, 2026). Backend Go monolito modular, PostgreSQL, JWT, sin integraciones con DIAN/terminales/navieras."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Autenticación y administración de usuarios (Priority: P1)

El administrador inicia sesión, crea usuarios con uno de los cuatro roles internos (Administrador, Operador logístico, Gerente, Cliente) y puede desactivar cuentas. Cada persona entra con correo y contraseña y solo ve lo que su rol permite.

**Why this priority**: Sin identidad y roles no hay operación segura ni aislamiento de datos entre clientes.

**Independent Test**: Crear un administrador semilla, iniciar sesión, crear un operador y un cliente, y verificar que un token de cliente no puede listar usuarios (`GET /api/v1/usuarios` → 403).

**Acceptance Scenarios**:

1. **Given** un usuario activo con credenciales válidas en una organización activa, **When** inicia sesión, **Then** recibe un JWT cuyos claims incluyen `user_id`, `organizacion_id` y `rol`.
2. **Given** un administrador autenticado, **When** crea un usuario con rol Operador, **Then** el operador puede iniciar sesión.
3. **Given** un cliente autenticado, **When** intenta crear o listar usuarios, **Then** el sistema rechaza la petición.
4. **Given** credenciales inválidas o cuenta desactivada, **When** intenta iniciar sesión, **Then** el sistema rechaza el acceso sin revelar si el correo existe.

---

### User Story 2 - Clientes de la PyME y registro de cargas (Priority: P1)

El operador registra clientes de la PyME y altas de carga (referencia, cliente dueño, origen/destino u otros datos operativos mínimos) y consulta el listado de cargas activas.

**Why this priority**: Es el núcleo operativo: sin cargas no hay trazabilidad ni documentos.

**Independent Test**: Autenticado como operador, crear un cliente, crear una carga asociada y listarla; un cliente autenticado no ve cargas de otros clientes; `POST /api/v1/cargas` con token de cliente → 403; `GET /api/v1/cargas/{id}` de una carga de otro cliente → **404** (no 403).

**Acceptance Scenarios**:

1. **Given** un operador autenticado, **When** registra un cliente y una carga de ese cliente, **Then** la carga queda en un estado inicial y aparece en el listado.
2. **Given** un cliente autenticado, **When** lista cargas, **Then** solo ve las suyas.
3. **Given** un operador, **When** intenta crear una carga sin cliente válido, **Then** el sistema rechaza la operación.

---

### User Story 3 - Estados de carga y línea de tiempo (Priority: P1)

El operador actualiza el estado de una carga. Cada cambio queda en una línea de tiempo con fecha/hora, usuario responsable, estado anterior y estado nuevo. El cliente y el gerente consultan esa trazabilidad.

**Why this priority**: Resuelve el problema central (historial disperso, sin responsables ni incidencias ligadas al tiempo).

**Independent Test**: Cambiar el estado de una carga dos veces y obtener la línea de tiempo ordenada; un **gerente** obtiene **200** en `GET /api/v1/cargas/{id}/trazabilidad`; un cliente no puede cambiar estados.

**Acceptance Scenarios**:

1. **Given** una carga existente, **When** el operador cambia su estado, **Then** el estado actual se actualiza y se añade un evento de trazabilidad inmutable.
2. **Given** un cliente dueño de la carga, **When** consulta la línea de tiempo, **Then** ve los eventos en orden cronológico.
3. **Given** un gerente autenticado de la misma organización, **When** consulta `GET /api/v1/cargas/{id}/trazabilidad`, **Then** recibe **200** y la línea de tiempo.
4. **Given** un cliente, **When** intenta cambiar el estado, **Then** el sistema rechaza la petición.

---

### User Story 4 - Incidencias operativas (Priority: P2)

El operador registra incidencias ligadas a una carga (descripción, `tipo`, estado abierta/cerrada). Quedan visibles en el detalle de la carga y en la trazabilidad.

**Why this priority**: Completa el historial operativo; es independiente del cambio de estado pero complementario.

**Independent Test**: Abrir una incidencia sobre una carga y listarla en el detalle; el cliente la ve, no la cierra.

**Acceptance Scenarios**:

1. **Given** una carga activa, **When** el operador abre una incidencia, **Then** queda asociada a la carga y aparece un evento de trazabilidad.
2. **Given** una incidencia abierta, **When** el operador la cierra, **Then** el estado de la incidencia cambia y se registra el evento.

---

### User Story 5 - Documentos adjuntos (Priority: P2)

El operador adjunta documentos (p. ej. PDF) a una carga. El cliente solo descarga los documentos que la PyME haya marcado como autorizados para él.

**Why this priority**: Reduce incertidumbre del cliente sin sustituir sistemas oficiales de aduana.

**Independent Test**: Subir un PDF a una carga, marcarlo como visible al cliente, y descargarlo como ese cliente; un documento no autorizado no es descargable por el cliente.

**Acceptance Scenarios**:

1. **Given** una carga, **When** el operador sube un documento autorizado para el cliente, **Then** el cliente puede listarlo y descargarlo.
2. **Given** un documento interno no autorizado al cliente, **When** el cliente intenta descargarlo, **Then** el sistema rechaza la petición.
3. **Given** un archivo que no es un tipo permitido o supera el tamaño máximo, **When** se intenta subir, **Then** el sistema rechaza la carga.

---

### User Story 6 - Tablero de métricas para gerencia (Priority: P3)

El gerente (y el administrador) consulta indicadores: cargas activas, retrasos (`retrasada` derivado), tiempo promedio total de entrega en horas e incidencias abiertas. El cliente no accede al tablero global.

**Why this priority**: Valor estratégico; depende de datos de P1/P2.

**Independent Test**: Con varias cargas e incidencias de prueba, el endpoint de métricas devuelve agregados coherentes; un cliente recibe 403.

**Acceptance Scenarios**:

1. **Given** un gerente autenticado, **When** consulta el tablero, **Then** recibe métricas agregadas de la organización.
2. **Given** un cliente autenticado, **When** consulta el tablero global, **Then** el sistema rechaza la petición.

---

### Edge Cases

- Token JWT expirado o manipulado: rechazar con 401, sin filtrar datos.
- Un cliente intenta acceder por ID a una carga de otro cliente o fuera de su tenant: **404 Not Found** (nunca 403; no revelar existencia).
- Cambio de estado al mismo estado: rechazar o no generar evento duplicado (se rechaza).
- Eventos de trazabilidad no se editan ni borran (solo append).
- Organización (PyME) inactiva o usuario desactivado: bloquear login y peticiones autenticadas.
- Subida concurrente o petición cancelada: respetar cancelación de la petición hacia la base de datos.
- Fuera de alcance: no hay sincronización con DIAN, terminales portuarios ni navieras.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST autenticar usuarios internos con correo y contraseña y emitir JWT cuyos claims mandatorios son `user_id`, `organizacion_id` y `rol`. El middleware JWT MUST rechazar peticiones si `organizaciones.activa` es false (además de usuario inactivo).
- **FR-002**: El sistema MUST soportar roles Administrador, Operador logístico, Gerente y Cliente, y aplicar autorización por ruta y por recurso (el cliente solo sus cargas).
- **FR-003**: Las contraseñas MUST almacenarse con hash irreversible (nunca en texto plano).
- **FR-004**: El administrador MUST poder crear, listar, actualizar y desactivar usuarios de su organización.
- **FR-005**: El operador MUST poder registrar y consultar clientes de la PyME.
- **FR-006**: El operador MUST poder crear, listar y consultar cargas asociadas a un cliente.
- **FR-007**: El operador MUST poder cambiar el estado de una carga; cada cambio MUST generar un evento de trazabilidad con fecha/hora, usuario, estado anterior y nuevo.
- **FR-008**: Cliente, operador, gerente y administrador MUST poder consultar la línea de tiempo de una carga que les corresponda.
- **FR-009**: El operador MUST poder abrir y cerrar incidencias asociadas a una carga; MUST quedar reflejado en trazabilidad.
- **FR-010**: El operador MUST poder adjuntar documentos a una carga y marcar visibilidad para el cliente.
- **FR-011**: El cliente MUST poder listar y descargar solo documentos autorizados de sus cargas.
- **FR-012**: Gerente y administrador MUST poder consultar métricas operativas agregadas: cargas activas, cargas retrasadas (flag derivado), incidencias abiertas, y **tiempo promedio total de entrega** en horas (`fecha_entrega - created_at` solo sobre cargas con estado `entregada`). v1 MUST NOT calcular intervalos intermedios entre estados.
- **FR-013**: El sistema MUST NO integrar ni sustituir sistemas oficiales (DIAN, terminales, navieras); el alcance es operativa interna de la PyME.
- **FR-014**: Errores HTTP MUST devolver un JSON uniforme (`error`, `code`).
- **FR-015**: Cada registro operativo MUST pertenecer a una organización (tenant SaaS) para aislamiento entre PyMEs.

### Key Entities *(include if feature involves data)*

- **Organización**: PyME suscriptora (tenant); agrupa usuarios, clientes y cargas.
- **Usuario**: Persona con rol; pertenece a una organización; puede estar activo/inactivo.
- **Cliente**: Empresa o persona atendida por la PyME; dueña de cargas; puede tener un usuario de portal.
- **Carga**: Unidad operativa trazada; tiene estado actual, cliente y metadatos logísticos internos.
- **Evento de trazabilidad**: Registro inmutable de un cambio (estado, incidencia, documento) con actor y timestamps.
- **Incidencia**: Problema operativo ligado a una carga (abierta/cerrada).
- **Documento**: Archivo adjunto a una carga, con flag de visibilidad para el cliente.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un operador puede registrar un cliente, una carga, un cambio de estado y ver la línea de tiempo en un único flujo autenticado (API), sin hojas de cálculo.
- **SC-002**: Un cliente autenticado obtiene estado, trazabilidad y documentos autorizados de sus cargas y no puede ver ni mutar datos de otros clientes (verificado por pruebas de autorización).
- **SC-003**: El 100% de los cambios de estado persistidos tienen evento de trazabilidad con actor y estados anterior/nuevo.
- **SC-004**: Las pruebas de contrato de la API cubren autenticación, flujo de alta, listado y detalle de cargas (create/list/get), trazabilidad, incidencias, documentos y métricas para los cuatro roles.
- **SC-005**: Un gerente obtiene métricas agregadas sin acceder a integraciones externas.

## Assumptions

- Este repositorio implementa **solo el backend** (API REST + PostgreSQL + almacenamiento de archivos); el frontend web es otro repositorio o fase posterior.
- El modelo SaaS es **multi-tenant por `organizacion_id`** en una misma base de datos (no un esquema por PyME en v1).
- **v1 multi-tenant opera bajo una sola organización aprovisionada mediante seed/despliegue inicial** (no hay API de alta de una segunda PyME).
- Facturación de suscripción mensual/anual queda **fuera de v1** (se asume organización activa).
- HTTP JSON REST; no hay notificaciones push ni email transaccional en v1 salvo error de API.
- Estados de carga en v1 son un **enum SQL inmutable**: `registrada`, `en_transito`, `en_puerto`, `entregada`, `cancelada`. No hay catálogo editable en runtime. `retrasada` **no** es un estado.
- Flag derivado `retrasada` (no persistido): `fecha_comprometida < NOW() AND estado != 'entregada'`. Se expone en API y en métricas (`cargas_retrasadas`).
- Tiempo promedio total de entrega (v1): media de `fecha_entrega - created_at` en horas, solo cargas `entregada`. Sin intervalos intermedios entre estados.
- Tamaño máximo de documento y tipos permitidos: PDF, JPEG, PNG; máximo 10 MB por archivo.
- Datos de prueba son sintéticos (amenaza de no acceder a datos reales de puerto).
- Rige la constitution **v1.0.0** (monolito modular, contrato HTTP `/api/v1`, tests de contrato antes de implementar).
