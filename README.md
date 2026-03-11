# 🧩 Splitmeet API

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version"/>
  <img src="https://img.shields.io/badge/Gin-Framework-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Gin"/>
  <img src="https://img.shields.io/badge/PostgreSQL-15+-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL"/>
  <img src="https://img.shields.io/badge/JWT-Auth-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white" alt="JWT"/>
  <img src="https://img.shields.io/badge/SSE-Realtime-FF6B6B?style=for-the-badge" alt="SSE"/>
  <img src="https://img.shields.io/badge/Architecture-Hexagonal-FF6B6B?style=for-the-badge" alt="Hexagonal"/>
</p>

> **Splitmeet** es una aplicación móvil diseñada para resolver el eterno problema de *"¿cuánto te debo?"* en las reuniones sociales. No es solo una calculadora, es un **gestor colaborativo de salidas** que permite desglosar gastos, asignar deudas y monitorear quién ya pagó y quién sigue debiendo.

---

## 📋 Tabla de Contenidos

- [Características](#-características)
- [Arquitectura](#-arquitectura)
- [Módulos del Sistema](#-módulos-del-sistema)
- [Modelo de Datos](#-modelo-de-datos)
- [Flujos de Negocio](#-flujos-de-negocio)
- [API Endpoints](#-api-endpoints)
- [Estructura del Proyecto](#-estructura-del-proyecto)
- [Instalación](#-instalación)
- [Configuración](#-configuración)
- [Desarrollo](#-desarrollo)

---

## ✨ Características

### Funcionalidades Principales

| Característica | Descripción |
|----------------|-------------|
| **Gestión de Salidas** | Crea eventos con nombre, fecha, categoría y descripción |
| **Grupos de Amigos** | Organiza tus contactos en grupos con roles (owner, admin, member) |
| **Invitaciones** | Sistema de invitación con estados (pendiente, aceptado, rechazado) |
| **Notificaciones en Tiempo Real** | SSE (Server-Sent Events) para notificaciones push al móvil |
| **Productos Predefinidos** | Catálogo de productos por categoría (restaurante, cine, etc.) |
| **División Inteligente** | Múltiples modos de división de gastos |
| **Tracking de Pagos** | Semáforo visual de quién ha pagado y quién debe |
| **Paginación** | Paginación y búsqueda en listados (grupos, notificaciones) |
| **Roles en Grupos** | Owner, admin y member con permisos diferenciados |
| **Auto-invitación** | Al crear salida con grupo, se invita automáticamente a todos los miembros |
| **Historial** | Registro completo de todas las salidas para referencia futura |

### Tipos de División de Gastos

```
┌─────────────────────────────────────────────────────────────────┐
│                    MODOS DE DIVISIÓN                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  💰 EQUITATIVO (equal)                                         │
│     Total ÷ Número de personas                                 │
│     Ejemplo: $1000 ÷ 4 = $250 c/u                              │
│                                                                 │
│  🎯 CANTIDAD FIJA (custom_fixed)                               │
│     Una persona paga monto fijo, resto equitativo              │
│     Ejemplo: Pedro paga $400, resto divide $600 ÷ 3            │
│                                                                 │
│  🍽️ POR CONSUMO (per_consumption)                              │
│     Cada quien paga exactamente lo que pidió                   │
│     Ideal para cuentas detalladas                              │
│                                                                 │
│  👤 UN PAGADOR (single_payer)                                  │
│     Una persona paga todo (recordatorio de deuda)              │
│     Útil cuando alguien "invita" temporalmente                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### División de Productos Compartidos

Los productos pueden dividirse entre un subconjunto de participantes:

```
Ejemplo: Salida al Cine - 5 personas

🍿 Palomitas Grandes ($150)
   └── Dividido entre: Pedro, Ana (2 personas)
   └── Cada quien: $75

🥤 Combo Pareja ($200)  
   └── Dividido entre: Carlos, María (2 personas)
   └── Cada quien: $100

🍿 Nachos ($80)
   └── Dividido entre: Pedro, Ana, Luis (3 personas)
   └── Cada quien: $26.67
```

---

## 🏗 Arquitectura

### Arquitectura Hexagonal (Puertos y Adaptadores)

Splitmeet implementa una **arquitectura hexagonal** que separa claramente las responsabilidades y permite una alta mantenibilidad y testabilidad.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           ARQUITECTURA HEXAGONAL                            │
└─────────────────────────────────────────────────────────────────────────────┘

                              ┌─────────────────┐
                              │   Controllers   │  ◄── Adaptadores de Entrada
                              │   (HTTP/REST)   │      (Gin Framework)
                              └────────┬────────┘
                                       │
                                       ▼
                    ┌──────────────────────────────────────┐
                    │              PUERTOS                 │
                    │  ┌─────────────────────────────────┐ │
                    │  │         Interfaces de           │ │
                    │  │         Entrada (API)           │ │
                    │  └─────────────────────────────────┘ │
                    └──────────────────┬───────────────────┘
                                       │
                                       ▼
          ┌────────────────────────────────────────────────────────┐
          │                                                        │
          │                    🎯 DOMINIO                          │
          │                                                        │
          │   ┌─────────────┐    ┌─────────────┐    ┌──────────┐  │
          │   │  Entities   │    │  Use Cases  │    │  Ports   │  │
          │   │  (Models)   │    │    (App)    │    │(Interfaces│  │
          │   └─────────────┘    └─────────────┘    └──────────┘  │
          │                                                        │
          └────────────────────────────┬───────────────────────────┘
                                       │
                                       ▼
                    ┌──────────────────────────────────────┐
                    │              PUERTOS                 │
                    │  ┌─────────────────────────────────┐ │
                    │  │        Interfaces de            │ │
                    │  │        Salida (Repos)           │ │
                    │  └─────────────────────────────────┘ │
                    └──────────────────┬───────────────────┘
                                       │
                                       ▼
                              ┌─────────────────┐
                              │   Repositories  │  ◄── Adaptadores de Salida
                              │  (PostgreSQL)   │      (Database)
                              └─────────────────┘
```

### Capas del Sistema

| Capa | Responsabilidad | Ejemplos |
|------|-----------------|----------|
| **Domain** | Reglas de negocio puras | Entities, Value Objects |
| **Application** | Casos de uso | CreateOuting, AddItem, RegisterPayment |
| **Infrastructure** | Implementaciones concretas | PostgreSQL repos, JWT adapter, SSE Hub |
| **Interfaces** | Puntos de entrada/salida | HTTP Controllers, Routes |

---

## 📦 Módulos del Sistema

### Diagrama de Módulos

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           MÓDULOS SPLITMEET                                 │
└─────────────────────────────────────────────────────────────────────────────┘

    ┌──────────────┐
    │    USER      │  ◄── Gestión de usuarios y autenticación
    │   ✅ HECHO   │
    └──────┬───────┘
           │
           ├──────────────────────────────┐
           ▼                              ▼
    ┌──────────────┐         ┌──────────────────┐
    │   CATEGORY   │────────>│   PRODUCT        │
    │   ✅ HECHO   │         │   ✅ HECHO       │
    └──────────────┘         └──────────────────┘
           │                        │
           │    ┌──────────────┐    │
           └───>│    GROUP     │<───┘
                │   ✅ HECHO   │
                └──────┬───────┘
                       │
                       ▼
                ┌──────────────┐
                │   OUTING     │  ◄── Módulo central
                │   ✅ HECHO   │
                └──────┬───────┘
                       │
                       ▼
                ┌──────────────┐       ┌───────────────────┐
                │   PAYMENT    │       │  NOTIFICATION     │
                │   ✅ HECHO   │       │  ✅ HECHO (SSE)   │
                └──────────────┘       └───────────────────┘
```

### Descripción de Módulos

#### 1. 👤 User
Gestión completa de usuarios con autenticación JWT.

```
Funcionalidades:
├── Registro de usuario
├── Login con JWT
├── Obtener perfil propio
├── Actualizar perfil
├── Buscar por username
├── Buscar usuarios (parcial)
├── Ver invitaciones pendientes
└── Eliminar cuenta
```

#### 2. 🏷️ Category
Categorías predefinidas para clasificar salidas.

```
Funcionalidades:
├── Listar todas las categorías
└── Obtener categoría por ID

Categorías iniciales:
├── 🍽️ Restaurante
├── 🎬 Cine
├── 🍺 Bar
├── ✈️ Viaje
├── 🛒 Supermercado
└── 📦 Otro
```

#### 3. 📦 Product
Catálogo de productos predefinidos y personalizados.

```
Funcionalidades:
├── Listar productos por categoría
├── Buscar productos
├── Crear producto personalizado
└── Productos predefinidos (sin precio fijo)
```

#### 4. 👥 Group
Grupos de amigos con roles y permisos.

```
Funcionalidades:
├── Crear grupo (creador = owner)
├── Listar mis grupos (paginado, con búsqueda)
├── Ver detalle de grupo
├── Invitar miembro (envía notificación)
├── Responder invitación (notifica al owner)
├── Eliminar miembro (owner/admin)
├── Actualizar grupo (owner/admin)
├── Eliminar grupo (owner)
├── Transferir ownership a otro miembro
└── Cambiar rol de miembro (admin/member)

Roles de membresía:
├── 👑 Owner   – Control total del grupo
├── 🛡️ Admin   – Puede editar grupo y remover miembros
└── 👤 Member  – Participante estándar

Estados de invitación:
├── 🟡 Pendiente (pending)
├── 🟢 Aceptado (accepted)
└── 🔴 Rechazado (rejected)
```

#### 5. 🎉 Outing
Módulo central para gestión de salidas/eventos.

```
Funcionalidades:
├── Crear salida (con o sin grupo)
├── Auto-invitación de miembros del grupo
├── Listar mis salidas
├── Ver detalle de salida
├── Actualizar salida
├── Eliminar salida
├── Agregar participante (envía notificación)
├── Confirmar participación (notifica al creador)
├── Agregar producto/item
├── Actualizar item
├── Eliminar item
├── Dividir item entre personas
└── Calcular montos automáticamente

Reglas de negocio:
├── Al crear con group_id: se invita automáticamente a todos los miembros aceptados
├── Solo editable si status = 'active'
├── Se bloquea cuando todos pagan
└── Cálculos automáticos al agregar items
```

#### 6. 💳 Payment
Sistema de tracking de pagos con validaciones.

```
Funcionalidades:
├── Registrar pago (con validación de monto)
├── Confirmar pago (por el organizador)
├── Ver pagos de una salida
├── Ver resumen de pagos
└── Auto-cancelación de pagos pendientes

Estados de pago:
├── 🟡 Pendiente (pending)
├── 🟢 Pagado (paid)
└── 🔴 Cancelado (cancelled)
```

#### 7. 🔔 Notification
Notificaciones en tiempo real via SSE (Server-Sent Events).

```
Funcionalidades:
├── Stream SSE en tiempo real (GET /notifications/stream)
├── Listar notificaciones (paginado)
├── Marcar notificación como leída
├── Marcar todas como leídas
└── Contador de no leídas

Tipos de notificación:
├── 📨 group_invitation     – Te invitaron a un grupo
├── 🎉 outing_invitation    – Te invitaron a una salida
├── ✅ invitation_accepted   – Alguien aceptó tu invitación
└── ❌ invitation_rejected   – Alguien rechazó tu invitación

Flujo SSE:
├── Cliente se conecta a GET /notifications/stream
├── Servidor mantiene conexión abierta
├── Cuando se crea una notificación, se envía al cliente en tiempo real
└── Keep-alive cada 30 segundos para mantener la conexión
```

---

## 🗄 Modelo de Datos

### Diagrama Entidad-Relación

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         DIAGRAMA ENTIDAD-RELACIÓN                           │
└─────────────────────────────────────────────────────────────────────────────┘

                                 ┌─────────────┐
                                 │   USERS     │
                                 ├─────────────┤
                                 │ id          │
                                 │ username    │
                                 │ name        │
                                 │ email       │
                                 │ phone       │
                                 │ password    │
                                 └──────┬──────┘
                                        │
              ┌─────────────────────────┼─────────────────────────┐
              │                         │                         │
              ▼                         ▼                         ▼
      ┌─────────────┐          ┌──────────────┐          ┌─────────────┐
      │   GROUPS    │          │GROUP_MEMBERS │          │  OUTINGS    │
      ├─────────────┤          ├──────────────┤          ├─────────────┤
      │ id          │◄────────>│ id           │          │ id          │
      │ name        │          │ group_id     │          │ name        │
      │ description │          │ user_id      │          │ description │
      │ owner_id    │          │ role         │          │ category_id │
      │ is_active   │          │ status       │          │ group_id    │
      └─────────────┘          │ invited_by   │          │ creator_id  │
                               │ invited_at   │          │ outing_date │
                               │ responded_at │          │ split_type  │
                               └──────────────┘          │ total_amount│
                                                         │ status      │
      ┌─────────────┐                                    │ is_editable │
      │ CATEGORIES  │                                    └──────┬──────┘
      ├─────────────┤                                           │
      │ id          │◄──────────────────────────────────────────┤
      │ name        │                                           │
      │ icon        │                                           │
      │ is_active   │          ┌────────────────────┐           │
      └──────┬──────┘          │OUTING_PARTICIPANTS │           │
             │                 ├────────────────────┤           │
             │                 │ id                 │◄──────────┤
             ▼                 │ outing_id          │           │
      ┌─────────────┐          │ user_id            │           │
      │  PRODUCTS   │          │ invited_by         │           │
      ├─────────────┤          │ status             │           │
      │ id          │          │ amount_owed        │           │
      │ category_id │          │ custom_amount      │           │
      │ name        │          │ joined_at          │           │
      │ presentation│          └─────────┬──────────┘           │
      │ size        │                    │                      │
      │ default_price          ┌─────────────────┐              │
      │ is_predefined│         │  ITEM_SPLITS    │              │
      │ created_by  │          ├─────────────────┤              │
      └──────┬──────┘          │ id              │              │
             │                 │ outing_item_id  │              │
             │                 │ participant_id  │              │
             │                 │ split_amount    │              │
             │                 │ percentage      │              │
             │                 └─────────────────┘              │
             │                          ▲                       │
             │                          │                       │
             │                 ┌─────────────────┐              │
             └────────────────>│  OUTING_ITEMS   │◄─────────────┘
                               ├─────────────────┤
                               │ id              │
                               │ outing_id       │
                               │ product_id      │
                               │ custom_name     │
                               │ custom_presentation
                               │ quantity        │
                               │ unit_price      │
                               │ subtotal        │  (GENERATED)
                               │ is_shared       │
                               └─────────────────┘

      ┌─────────────────┐             ┌─────────────────┐
      │    PAYMENTS     │             │  NOTIFICATIONS  │
      ├─────────────────┤             ├─────────────────┤
      │ id              │             │ id              │
      │ outing_id       │             │ user_id         │
      │ participant_id  │             │ type            │
      │ amount          │             │ title           │
      │ status          │             │ message         │
      │ paid_at         │             │ reference_id    │
      │ confirmed_by    │             │ inviter_name    │
      │ notes           │             │ group_name      │
      └─────────────────┘             │ outing_name     │
                                      │ is_read         │
                                      │ created_at      │
                                      └─────────────────┘
```

### Enums del Sistema

```sql
-- Roles de membresía en grupo
member_role: 'owner' | 'admin' | 'member'

-- Estados de membresía en grupo
member_status: 'pending' | 'accepted' | 'rejected'

-- Tipos de división de gastos
split_type: 'equal' | 'custom_fixed' | 'per_consumption' | 'single_payer'

-- Estados de una salida
outing_status: 'active' | 'completed' | 'cancelled'

-- Estados de participación en salida
participant_status: 'pending' | 'confirmed' | 'declined'

-- Estados de pago
payment_status: 'pending' | 'paid' | 'cancelled'

-- Tipos de notificación
notification_type: 'group_invitation' | 'outing_invitation' | 'invitation_accepted' | 'invitation_rejected'
```

---

## 🔄 Flujos de Negocio

### Flujo 1: Crear Grupo e Invitar Amigos

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    FLUJO: CREAR GRUPO E INVITAR                             │
└─────────────────────────────────────────────────────────────────────────────┘

    USUARIO A                          SISTEMA                         USUARIO B
        │                                 │                                │
        │  1. POST /groups                │                                │
        │  {name, description}            │                                │
        │────────────────────────────────>│                                │
        │                                 │                                │
        │  ◄── Grupo creado (owner: A)    │                                │
        │                                 │                                │
        │  2. POST /groups/:id/members    │                                │
        │  {username: "userB"}            │                                │
        │────────────────────────────────>│                                │
        │                                 │                                │
        │                                 │  Crear member con              │
        │                                 │  status: 'pending'             │
        │                                 │  role: 'member'                │
        │                                 │                                │
        │                                 │  [SSE Notification]            │
        │                                 │  type: group_invitation        │
        │                                 │────────────────────────────────>
        │                                 │                                │
        │                                 │  3. PATCH /groups/:id/members/respond
        │                                 │  {action: "accept"}            │
        │                                 │<────────────────────────────────
        │                                 │                                │
        │                                 │  Actualizar status: 'accepted' │
        │                                 │                                │
        │  [SSE Notification]             │                                │
        │  type: invitation_accepted      │                                │
        │  ◄──────────────────────────────│                                │
        │                                 │                                │
        ▼                                 ▼                                ▼
```

### Flujo 2: Crear Salida desde un Grupo (Auto-invitación)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    FLUJO: CREAR SALIDA DESDE GRUPO                          │
└─────────────────────────────────────────────────────────────────────────────┘

    CREADOR                            SISTEMA                      MIEMBROS
        │                                 │                             │
        │  1. POST /outings               │                             │
        │  {                              │                             │
        │    name: "Cena viernes",        │                             │
        │    group_id: 1,                 │                             │
        │    ...                          │                             │
        │  }                              │                             │
        │────────────────────────────────>│                             │
        │                                 │                             │
        │                                 │  Crear outing               │
        │                                 │  Auto-invitar a todos       │
        │                                 │  los miembros aceptados     │
        │                                 │  del grupo                  │
        │                                 │                             │
        │                                 │  [SSE Notifications]        │
        │                                 │  type: outing_invitation    │
        │                                 │  a cada miembro             │
        │                                 │─────────────────────────────>
        │                                 │                             │
        │                                 │  2. PATCH /outings/:id/     │
        │                                 │     participants/confirm    │
        │                                 │  {accept: true}             │
        │                                 │<─────────────────────────────
        │                                 │                             │
        │  [SSE Notification]             │  Actualizar status:         │
        │  type: invitation_accepted      │  'confirmed'                │
        │  ◄──────────────────────────────│                             │
        │                                 │                             │
        ▼                                 ▼                             ▼
```

### Flujo 3: Agregar Productos y Calcular División

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    FLUJO: AGREGAR PRODUCTOS                                 │
└─────────────────────────────────────────────────────────────────────────────┘

    USUARIO                            SISTEMA
        │                                 │
        │  1. GET /products?category_id=1 │
        │────────────────────────────────>│
        │                                 │
        │  ◄── Productos de Restaurante   │
        │      (predefinidos)             │
        │                                 │
        │  2. POST /outings/:id/items     │
        │  {                              │
        │    product_id: 5,               │  // Jarra de Agua Grande
        │    quantity: 2,                 │
        │    unit_price: 50.00            │
        │  }                              │
        │────────────────────────────────>│
        │                                 │
        │                                 │  ┌─────────────────────┐
        │                                 │  │ Crear outing_item   │
        │                                 │  │ subtotal = 2 × 50   │
        │                                 │  │ subtotal = $100     │
        │                                 │  └─────────────────────┘
        │                                 │
        │  3. POST /outings/:id/items/:itemId/splits
        │  {                              │
        │    participant_ids: [1, 3]      │  // Solo Pedro y Ana
        │  }                              │
        │────────────────────────────────>│
        │                                 │
        │                                 │  ┌─────────────────────┐
        │                                 │  │ Dividir $200 ÷ 2    │
        │                                 │  │ Pedro: $100         │
        │                                 │  │ Ana: $100           │
        │                                 │  └─────────────────────┘
        │                                 │
        ▼                                 ▼
```

### Flujo 4: Proceso de Pago

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    FLUJO: PROCESO DE PAGO                                   │
└─────────────────────────────────────────────────────────────────────────────┘

    PARTICIPANTE                       SISTEMA                         CREADOR
        │                                 │                                │
        │  1. GET /payments/outing/:id    │                                │
        │     /summary                    │                                │
        │────────────────────────────────>│                                │
        │                                 │                                │
        │  ◄── Estado de todos los pagos  │                                │
        │      Mi deuda: $250             │                                │
        │                                 │                                │
        │  2. POST /payments              │                                │
        │  {outing_id, amount}            │                                │
        │────────────────────────────────>│                                │
        │                                 │                                │
        │                                 │  Crear payment pending         │
        │                                 │                                │
        │                                 │  3. PATCH /payments/:id/confirm│
        │                                 │<────────────────────────────────
        │                                 │                                │
        │                                 │  confirmed_by = creador_id     │
        │                                 │                                │
        │                                 │  ┌─────────────────────────┐   │
        │                                 │  │ ¿Todos pagaron?         │   │
        │                                 │  │ SI → outing.status =    │   │
        │                                 │  │      'completed'        │   │
        │                                 │  │      is_editable = false│   │
        │                                 │  │ NO → Mantener 'active'  │   │
        │                                 │  └─────────────────────────┘   │
        │                                 │                                │
        ▼                                 ▼                                ▼
```

---

## 🌐 API Endpoints

### Autenticación
Todos los endpoints (excepto registro y login) requieren header:
```
Authorization: Bearer <JWT_TOKEN>
```

### Paginación
Los endpoints paginados aceptan query parameters:
```
?page=1&limit=10&search=texto&sort=created_at&order=desc
```

### Endpoints por Módulo

#### 👤 Users

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/users` | Registrar usuario |
| `POST` | `/users/login` | Iniciar sesión (con email) |
| `GET` | `/users/profile` | Obtener mi perfil |
| `PATCH` | `/users/update` | Actualizar mi perfil |
| `GET` | `/users/get/:id` | Obtener usuario por ID |
| `GET` | `/users/username/:username` | Buscar por username exacto |
| `GET` | `/users/search?username=xxx` | Buscar usuarios (parcial) |
| `GET` | `/users/invitations` | Ver invitaciones pendientes |
| `DELETE` | `/users/delete` | Eliminar mi cuenta |

#### 🏷️ Categories

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/categories` | Listar todas las categorías |
| `GET` | `/categories/:id` | Obtener categoría por ID |

#### 📦 Products

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/products/category/:id` | Listar productos por categoría |
| `GET` | `/products/:id` | Obtener producto por ID |
| `POST` | `/products` | Crear producto personalizado |

#### 👥 Groups

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/groups` | Crear grupo (creador = owner) |
| `GET` | `/groups?page=1&limit=10&search=` | Listar mis grupos (paginado) |
| `GET` | `/groups/:id` | Obtener detalle de grupo |
| `PATCH` | `/groups/:id` | Actualizar grupo (owner/admin) |
| `DELETE` | `/groups/:id` | Eliminar grupo (owner) |
| `GET` | `/groups/:id/members` | Listar miembros (con rol) |
| `POST` | `/groups/:id/members` | Invitar usuario (envía notificación SSE) |
| `PATCH` | `/groups/:id/members/respond` | Responder invitación (notifica al owner) |
| `DELETE` | `/groups/:id/members/:userId` | Remover miembro (owner/admin) |
| `GET` | `/groups/invitations` | Ver invitaciones pendientes |
| `PATCH` | `/groups/:id/transfer` | Transferir ownership |
| `PATCH` | `/groups/:id/members/:userId/role` | Cambiar rol (owner only) |

#### 🎉 Outings

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/outings` | Crear salida (auto-invita grupo si `group_id`) |
| `GET` | `/outings/me` | Listar mis salidas |
| `GET` | `/outings/:id` | Obtener detalle de salida |
| `GET` | `/outings/group/:groupId` | Salidas de un grupo |
| `PATCH` | `/outings/:id` | Actualizar salida |
| `DELETE` | `/outings/:id` | Eliminar salida |
| `GET` | `/outings/:id/participants` | Listar participantes |
| `POST` | `/outings/:id/participants` | Agregar participante (con notificación) |
| `PATCH` | `/outings/:id/participants/confirm` | Confirmar asistencia (notifica creador) |
| `DELETE` | `/outings/:id/participants/:userId` | Remover participante |
| `GET` | `/outings/:id/items` | Listar items |
| `POST` | `/outings/:id/items` | Agregar item |
| `PATCH` | `/outings/:id/items/:itemId` | Actualizar item |
| `DELETE` | `/outings/:id/items/:itemId` | Eliminar item |
| `POST` | `/outings/:id/items/:itemId/splits` | Dividir item |
| `GET` | `/outings/:id/items/:itemId/splits` | Ver división de item |
| `GET` | `/outings/:id/calculate` | Calcular montos |

#### 💳 Payments

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/payments` | Registrar pago (validado) |
| `GET` | `/payments/:id` | Obtener detalle de pago |
| `GET` | `/payments/outing/:id` | Ver pagos de una salida |
| `GET` | `/payments/outing/:id/summary` | Resumen de pagos |
| `PATCH` | `/payments/:id/confirm` | Confirmar pago recibido |
| `DELETE` | `/payments/:id` | Eliminar pago pendiente |

#### 🔔 Notifications

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/notifications?page=1&limit=20` | Listar notificaciones (paginado) |
| `GET` | `/notifications/stream` | Stream SSE en tiempo real |
| `PATCH` | `/notifications/:id/read` | Marcar notificación como leída |
| `PATCH` | `/notifications/read-all` | Marcar todas como leídas |

**SSE Stream:**
```
GET /notifications/stream
Authorization: Bearer <JWT_TOKEN>

# Respuesta (text/event-stream):
event: notification
data: {"id":1,"type":"group_invitation","title":"Invitación a grupo","message":"Juan te invitó al grupo Amigos","reference_id":5,"inviter_name":"Juan","group_name":"Amigos","is_read":false,"created_at":"2025-01-15T10:30:00Z"}

event: keep-alive
data: ping
```

---

## 📁 Estructura del Proyecto

```
splitmeet-api/
│
├── main.go                          # Punto de entrada
├── go.mod                           # Dependencias
├── go.sum
├── schema.sql                       # Script de base de datos
├── README.md                        # Este archivo
│
├── internal/
│   │
│   ├── core/                        # Configuración central
│   │   ├── cors.go                  # Configuración CORS
│   │   ├── postgresql.go            # Conexión a BD
│   │   ├── pagination.go            # Utilidad de paginación
│   │   ├── sse_hub.go               # Hub SSE para notificaciones
│   │   └── notification_service.go  # Servicio compartido de notificaciones
│   │
│   ├── middleware/                   # Middlewares globales
│   │   └── auth.go                  # Middleware JWT
│   │
│   ├── user/                        # MÓDULO USER
│   │   ├── app/                     # Casos de uso
│   │   │   ├── create_user.go
│   │   │   ├── delete_user.go
│   │   │   ├── get_by_username.go
│   │   │   ├── get_pending_invitations.go
│   │   │   ├── get_profile.go
│   │   │   ├── get_user.go
│   │   │   ├── login.go
│   │   │   ├── search_users.go
│   │   │   └── update_my_profile.go
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   ├── ports/
│   │   │   └── repository/
│   │   └── infra/
│   │       ├── dependencies.go
│   │       ├── adapters/
│   │       ├── controllers/
│   │       ├── repository/
│   │       ├── routes/
│   │       └── services/
│   │
│   ├── category/                    # MÓDULO CATEGORY
│   │   ├── app/
│   │   ├── domain/
│   │   └── infra/
│   │
│   ├── product/                     # MÓDULO PRODUCT
│   │   ├── app/
│   │   ├── domain/
│   │   └── infra/
│   │
│   ├── group/                       # MÓDULO GROUP
│   │   ├── app/
│   │   │   ├── create_group.go
│   │   │   ├── delete_group.go
│   │   │   ├── get_group.go
│   │   │   ├── get_members.go
│   │   │   ├── get_my_groups.go
│   │   │   ├── get_pending_invitations.go
│   │   │   ├── invite_member.go
│   │   │   ├── remove_member.go
│   │   │   ├── respond_invitation.go
│   │   │   ├── set_member_role.go
│   │   │   ├── transfer_ownership.go
│   │   │   └── update_group.go
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   │   ├── group.go
│   │   │   │   └── group_member.go
│   │   │   └── repository/
│   │   │       └── group_repository.go
│   │   └── infra/
│   │
│   ├── outing/                      # MÓDULO OUTING
│   │   ├── app/
│   │   │   ├── create_outing.go     # Auto-invita miembros de grupo
│   │   │   ├── add_participant.go   # Con notificación SSE
│   │   │   ├── confirm_participation.go  # Notifica al creador
│   │   │   └── ...
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   │   ├── outing.go
│   │   │   │   ├── outing_participant.go
│   │   │   │   ├── outing_item.go
│   │   │   │   └── item_split.go
│   │   │   └── repository/
│   │   └── infra/
│   │
│   ├── payment/                     # MÓDULO PAYMENT
│   │   ├── app/
│   │   ├── domain/
│   │   └── infra/
│   │
│   └── notification/                # MÓDULO NOTIFICATION (SSE)
│       ├── app/
│       │   ├── get_notifications.go
│       │   ├── mark_as_read.go
│       │   └── create_notification.go
│       ├── domain/
│       │   ├── entities/
│       │   │   └── notification.go
│       │   └── repository/
│       │       └── notification_repository.go
│       └── infra/
│           ├── dependencies.go
│           ├── controllers/
│           │   ├── get_notifications.go
│           │   ├── mark_as_read.go
│           │   └── sse_stream.go
│           ├── repository/
│           │   └── notification_postgresql.go
│           └── routes/
│               └── notification_router.go
```

---

## 🚀 Instalación

### Prerrequisitos

- Go 1.25 o superior
- PostgreSQL 15 o superior
- Git

### Pasos

```bash
# Clonar el repositorio
git clone https://github.com/JosephAntonyDev/Splitmeet-API.git
cd Splitmeet-API

# Instalar dependencias
go mod download

# Configurar variables de entorno (ver sección Configuración)

# Ejecutar migraciones de base de datos
psql -U tu_usuario -d splitmeet -f schema.sql

# Ejecutar el servidor
go run main.go
```

---

## ⚙️ Configuración

### Variables de Entorno

Crear archivo `.env` en la raíz del proyecto:

```env
# Servidor
PORT=8080
GIN_MODE=debug

# Base de Datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_password
DB_NAME=splitmeet
DB_SSLMODE=disable

# JWT
JWT_SECRET=tu_secreto_super_seguro
JWT_EXPIRATION_HOURS=24
```

---

## 🛠 Desarrollo

### Estado de Implementación

Todos los módulos están completamente implementados:

```
1. ✅ User          (autenticación, perfil, búsqueda, invitaciones)
2. ✅ Category      (listado de categorías predefinidas)
3. ✅ Product       (catálogo por categoría)
4. ✅ Group         (grupos con roles, invitaciones, transferencia)
5. ✅ Outing        (salidas, auto-invitación, participantes, items, splits)
6. ✅ Payment       (pagos con validaciones y auto-cancelación)
7. ✅ Notification   (SSE en tiempo real, marcar leídas, paginación)
```

### Convenciones de Código

- **Naming**: camelCase para variables, PascalCase para tipos exportados
- **Errores**: Siempre manejar y propagar errores apropiadamente
- **Arquitectura**: Hexagonal con capas domain → app → infra
- **Notificaciones**: Usar `core.NotificationService` para enviar desde cualquier módulo

### Estructura de un Módulo

Cada módulo sigue esta estructura:

```go
// domain/entities/example.go
type Example struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

// domain/repository/example_repository.go
type ExampleRepository interface {
    Create(e *entities.Example) error
    FindByID(id int64) (*entities.Example, error)
}

// app/create_example.go
type CreateExampleUseCase struct {
    repo repository.ExampleRepository
}

func (uc *CreateExampleUseCase) Execute(input CreateExampleInput) (*entities.Example, error) {
    // Lógica de negocio
}
```

---

## 📄 Licencia

Este proyecto es privado y está desarrollado para fines académicos.

---

## 👥 Equipo

- **Backend Developer**: Joseph Antony
- **Curso**: Desarrollo de Aplicaciones Móviles
- **Universidad**: 8vo Cuatrimestre

---

<p align="center">
  <strong>Splitmeet</strong> - Dividir cuentas nunca fue tan fácil 💰
</p>
