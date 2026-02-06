# 🧩 Splitmeet API

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version"/>
  <img src="https://img.shields.io/badge/Gin-Framework-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Gin"/>
  <img src="https://img.shields.io/badge/PostgreSQL-15+-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL"/>
  <img src="https://img.shields.io/badge/JWT-Auth-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white" alt="JWT"/>
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
| **Grupos de Amigos** | Organiza tus contactos en grupos para crear salidas recurrentes |
| **Invitaciones** | Sistema de invitación con estados (pendiente, aceptado, rechazado) |
| **Productos Predefinidos** | Catálogo de productos por categoría (restaurante, cine, etc.) |
| **División Inteligente** | Múltiples modos de división de gastos |
| **Tracking de Pagos** | Semáforo visual de quién ha pagado y quién debe |
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
| **Infrastructure** | Implementaciones concretas | PostgreSQL repos, JWT adapter |
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
           ▼
    ┌──────────────┐         ┌──────────────┐
    │   CATEGORY   │────────>│   PRODUCT    │
    │   ✅ HECHO   │         │   ✅ HECHO   │
    └──────────────┘         └──────────────┘
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
                ┌──────────────┐
                │   PAYMENT    │
                │   ✅ HECHO   │
                └──────────────┘
```

### Descripción de Módulos

#### 1. 👤 User (Implementado)
Gestión completa de usuarios con autenticación JWT.

```
Funcionalidades:
├── Registro de usuario
├── Login con JWT
├── Obtener perfil
├── Actualizar perfil
├── Buscar por username
└── Eliminar cuenta
```

#### 2. 🏷️ Category (Implementado)
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

#### 3. 📦 Product (Implementado)
Catálogo de productos predefinidos y personalizados.

```
Funcionalidades:
├── Listar productos por categoría
├── Buscar productos
├── Crear producto personalizado
└── Productos predefinidos (sin precio fijo)

Características del producto:
├── Nombre
├── Presentación (Jarra, Vaso, Bolsa, etc.)
├── Tamaño (Chico, Mediano, Grande)
└── Precio (opcional en predefinidos)
```

#### 4. 👥 Group (Implementado)
Grupos de amigos para organizar salidas.

```
Funcionalidades:
├── Crear grupo
├── Listar mis grupos
├── Ver detalle de grupo
├── Invitar miembro
├── Responder invitación (aceptar/rechazar)
├── Eliminar miembro
├── Actualizar grupo
└── Eliminar grupo

Estados de membresía:
├── 🟡 Pendiente (pending)
├── 🟢 Aceptado (accepted)
└── 🔴 Rechazado (rejected)
```

#### 5. 🎉 Outing (Implementado)
Módulo central para gestión de salidas/eventos.

```
Funcionalidades:
├── Crear salida (con o sin grupo)
├── Listar mis salidas
├── Ver detalle de salida
├── Actualizar salida
├── Eliminar salida
├── Agregar participante
├── Confirmar participación
├── Agregar producto/item
├── Actualizar item
├── Eliminar item
├── Dividir item entre personas
└── Calcular montos automáticamente

Reglas de negocio:
├── Solo editable si status = 'active'
├── Se bloquea cuando todos pagan
└── Cálculos automáticos al agregar items
```

#### 6. 💳 Payment (Implementado)
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

Validaciones:
├── No permite pagar más del saldo restante
├── No permite pagar si no hay items
├── Auto-cancela pagos pendientes cuando se alcanza el total
└── Al pagar todos → outing.is_editable = false
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
      │ owner_id    │          │ status       │          │ category_id │
      │ is_active   │          │ invited_by   │          │ group_id    │
      └─────────────┘          │ invited_at   │          │ creator_id  │
                               │ responded_at │          │ outing_date │
                               └──────────────┘          │ split_type  │
                                                         │ total_amount│
      ┌─────────────┐                                    │ status      │
      │ CATEGORIES  │                                    │ is_editable │
      ├─────────────┤                                    └──────┬──────┘
      │ id          │◄───────────────────────────────────────────┤
      │ name        │                                            │
      │ icon        │                                            │
      │ is_active   │          ┌────────────────────┐            │
      └──────┬──────┘          │OUTING_PARTICIPANTS │            │
             │                 ├────────────────────┤            │
             │                 │ id                 │◄───────────┤
             ▼                 │ outing_id          │            │
      ┌─────────────┐          │ user_id            │            │
      │  PRODUCTS   │          │ status             │            │
      ├─────────────┤          │ amount_owed        │            │
      │ id          │          │ custom_amount      │            │
      │ category_id │          │ joined_at          │            │
      │ name        │          └─────────┬──────────┘            │
      │ presentation│                    │                       │
      │ size        │                    │                       │
      │ default_price          ┌─────────────────┐               │
      │ is_predefined│         │  ITEM_SPLITS    │               │
      │ created_by  │          ├─────────────────┤               │
      └──────┬──────┘          │ id              │               │
             │                 │ outing_item_id  │               │
             │                 │ participant_id  │               │
             │                 │ split_amount    │               │
             │                 │ percentage      │               │
             │                 └─────────────────┘               │
             │                          ▲                        │
             │                          │                        │
             │                 ┌─────────────────┐               │
             └────────────────>│  OUTING_ITEMS   │◄──────────────┘
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

                               ┌─────────────────┐
                               │    PAYMENTS     │
                               ├─────────────────┤
                               │ id              │
                               │ outing_id       │
                               │ participant_id  │
                               │ amount          │
                               │ status          │
                               │ paid_at         │  (fecha confirmación)
                               │ confirmed_by    │
                               │ notes           │
                               └─────────────────┘
```

### Enums del Sistema

```sql
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
        │  2. POST /groups/:id/invite     │                                │
        │  {username: "userB"}            │                                │
        │────────────────────────────────>│                                │
        │                                 │                                │
        │                                 │  Crear member con              │
        │                                 │  status: 'pending'             │
        │                                 │                                │
        │                                 │  [Notificación]                │
        │                                 │────────────────────────────────>
        │                                 │                                │
        │                                 │  3. PATCH /groups/:id/invitation
        │                                 │  {action: "accept"}            │
        │                                 │<────────────────────────────────
        │                                 │                                │
        │                                 │  Actualizar status:            │
        │                                 │  'accepted'                    │
        │                                 │                                │
        │  ◄── B es ahora miembro ────────│                                │
        │                                 │                                │
        ▼                                 ▼                                ▼
```

### Flujo 2: Crear Salida desde un Grupo

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    FLUJO: CREAR SALIDA DESDE GRUPO                          │
└─────────────────────────────────────────────────────────────────────────────┘

    CREADOR                            SISTEMA                      MIEMBROS
        │                                 │                             │
        │  1. GET /groups/:id/members     │                             │
        │────────────────────────────────>│                             │
        │                                 │                             │
        │  ◄── Lista de miembros          │                             │
        │      aceptados                  │                             │
        │                                 │                             │
        │  2. POST /outings               │                             │
        │  {                              │                             │
        │    name: "Cena viernes",        │                             │
        │    group_id: 1,                 │                             │
        │    category_id: 1,              │                             │
        │    outing_date: "2026-02-10",   │                             │
        │    participant_ids: [2, 3, 4]   │                             │
        │  }                              │                             │
        │────────────────────────────────>│                             │
        │                                 │                             │
        │                                 │  Crear outing               │
        │                                 │  Crear participants         │
        │                                 │  con status: 'pending'      │
        │                                 │                             │
        │                                 │  [Notificaciones]           │
        │                                 │─────────────────────────────>
        │                                 │                             │
        │                                 │  3. PATCH /outings/:id/     │
        │                                 │     participants/:userId/   │
        │                                 │     confirm                 │
        │                                 │<─────────────────────────────
        │                                 │                             │
        │                                 │  Actualizar status:         │
        │                                 │  'confirmed'                │
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
        │  3. POST /outings/:id/items     │
        │  {                              │
        │    custom_name: "Gambas",       │  // Producto personalizado
        │    custom_presentation: "Plato",│
        │    quantity: 1,                 │
        │    unit_price: 200.00           │
        │  }                              │
        │────────────────────────────────>│
        │                                 │
        │  ◄── Item creado                │
        │                                 │
        │  4. POST /outings/:id/items/:itemId/splits
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
        │                                 │  Recalcular amount_owed
        │                                 │  de cada participante
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
        │  1. GET /outings/:id/payments   │                                │
        │────────────────────────────────>│                                │
        │                                 │                                │
        │  ◄── Estado de todos los pagos  │                                │
        │      Mi deuda: $250             │                                │
        │      Status: pending            │                                │
        │                                 │                                │
        │  [Realiza pago físico/transfer] │                                │
        │                                 │                                │
        │  2. PATCH /payments/:id/pay     │                                │
        │────────────────────────────────>│                                │
        │                                 │                                │
        │                                 │  Actualizar payment:           │
        │                                 │  status = 'paid'               │
        │                                 │  paid_at = NOW()               │
        │                                 │                                │
        │                                 │  [Notificación al creador]     │
        │                                 │────────────────────────────────>
        │                                 │                                │
        │                                 │  3. PATCH /payments/:id/confirm
        │                                 │<────────────────────────────────
        │                                 │                                │
        │                                 │  confirmed_by = creador_id     │
        │                                 │                                │
        │                                 │  ┌─────────────────────────┐   │
        │                                 │  │ ¿Todos pagaron?         │   │
        │                                 │  │                         │   │
        │                                 │  │ SI → outing.status =    │   │
        │                                 │  │      'completed'        │   │
        │                                 │  │      is_editable = false│   │
        │                                 │  │                         │   │
        │                                 │  │ NO → Mantener 'active'  │   │
        │                                 │  └─────────────────────────┘   │
        │                                 │                                │
        ▼                                 ▼                                ▼
```

### Flujo 5: Cálculo Automático de Montos

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ALGORITMO: CÁLCULO DE MONTOS                             │
└─────────────────────────────────────────────────────────────────────────────┘

    Ejemplo: Salida con 4 participantes, split_type = 'per_consumption'

    ┌─────────────────────────────────────────────────────────────────┐
    │                        OUTING_ITEMS                             │
    ├─────────────────────────────────────────────────────────────────┤
    │ Item 1: Jarra Agua    │ $100  │ is_shared: TRUE  │ → Split     │
    │ Item 2: Hamburguesa   │ $150  │ is_shared: FALSE │ → 1 persona │
    │ Item 3: Pizza         │ $200  │ is_shared: TRUE  │ → Split     │
    │ Item 4: Refresco      │ $30   │ is_shared: FALSE │ → 1 persona │
    └─────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
    ┌─────────────────────────────────────────────────────────────────┐
    │                        ITEM_SPLITS                              │
    ├─────────────────────────────────────────────────────────────────┤
    │ Item 1 ($100) → Pedro: $25, Ana: $25, Luis: $25, María: $25    │
    │ Item 2 ($150) → Pedro: $150                                     │
    │ Item 3 ($200) → Ana: $100, Luis: $100                          │
    │ Item 4 ($30)  → María: $30                                      │
    └─────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
    ┌─────────────────────────────────────────────────────────────────┐
    │                    RESUMEN POR PERSONA                          │
    ├─────────────────────────────────────────────────────────────────┤
    │ Pedro:  $25 + $150         = $175                               │
    │ Ana:    $25 + $100         = $125                               │
    │ Luis:   $25 + $100         = $125                               │
    │ María:  $25 + $30          = $55                                │
    ├─────────────────────────────────────────────────────────────────┤
    │ TOTAL:                       $480                               │
    └─────────────────────────────────────────────────────────────────┘
```

---

## 🌐 API Endpoints

### Autenticación
Todos los endpoints (excepto registro y login) requieren header:
```
Authorization: Bearer <JWT_TOKEN>
```

### Endpoints por Módulo

#### 👤 Users (Implementado)

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

#### 🏷️ Categories (Implementado)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/categories` | Listar todas las categorías |
| `GET` | `/categories/:id` | Obtener categoría por ID |

#### 📦 Products (Implementado)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/products/category/:id` | Listar productos por categoría |
| `GET` | `/products/:id` | Obtener producto por ID |
| `POST` | `/products` | Crear producto personalizado |

#### 👥 Groups (Implementado)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/groups` | Crear grupo |
| `GET` | `/groups` | Listar mis grupos |
| `GET` | `/groups/:id` | Obtener detalle de grupo |
| `PATCH` | `/groups/:id` | Actualizar grupo |
| `DELETE` | `/groups/:id` | Eliminar grupo |
| `GET` | `/groups/:id/members` | Listar miembros |
| `POST` | `/groups/:id/members` | Invitar usuario |
| `PATCH` | `/groups/:id/members/respond` | Responder invitación (usa `{"status": "accepted"}`) |
| `DELETE` | `/groups/:id/members/:userId` | Remover miembro |

#### 🎉 Outings (Implementado)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/outings` | Crear salida |
| `GET` | `/outings/me` | Listar mis salidas |
| `GET` | `/outings/:id` | Obtener detalle de salida |
| `GET` | `/outings/group/:id` | Salidas de un grupo |
| `PATCH` | `/outings/:id` | Actualizar salida |
| `DELETE` | `/outings/:id` | Eliminar salida |
| `GET` | `/outings/:id/participants` | Listar participantes |
| `POST` | `/outings/:id/participants` | Agregar participante |
| `PATCH` | `/outings/:id/participants/confirm` | Confirmar asistencia (usa `{"accept": true}`) |
| `DELETE` | `/outings/:id/participants/:userId` | Remover participante |
| `GET` | `/outings/:id/items` | Listar items |
| `POST` | `/outings/:id/items` | Agregar item |
| `PATCH` | `/outings/:id/items/:itemId` | Actualizar item |
| `DELETE` | `/outings/:id/items/:itemId` | Eliminar item |
| `POST` | `/outings/:id/items/:itemId/splits` | Dividir item |
| `GET` | `/outings/:id/items/:itemId/splits` | Ver división de item |
| `GET` | `/outings/:id/calculate` | Calcular montos |

#### 💳 Payments (Implementado)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/payments` | Registrar pago (validado) |
| `GET` | `/payments/:id` | Obtener detalle de pago |
| `GET` | `/payments/outing/:id` | Ver pagos de una salida |
| `GET` | `/payments/outing/:id/summary` | Resumen de pagos |
| `PATCH` | `/payments/:id/confirm` | Confirmar pago recibido |
| `DELETE` | `/payments/:id` | Eliminar pago pendiente |

---

## 📁 Estructura del Proyecto

```
splitmeet-api/
│
├── main.go                          # Punto de entrada
├── go.mod                           # Dependencias
├── go.sum
├── README.md                        # Este archivo
│
├── internal/
│   │
│   ├── core/                        # Configuración central
│   │   ├── cors.go                  # Configuración CORS
│   │   └── postgresql.go            # Conexión a BD
│   │
│   ├── middleware/                  # Middlewares globales
│   │   └── auth.go                  # Middleware JWT
│   │
│   ├── user/                        # ✅ MÓDULO USER
│   │   ├── app/                     # Casos de uso
│   │   │   ├── create_user.go
│   │   │   ├── delete_user.go
│   │   │   ├── get_by_username.go
│   │   │   ├── get_profile.go
│   │   │   ├── get_user.go
│   │   │   ├── login.go
│   │   │   └── update_my_profile.go
│   │   │
│   │   ├── domain/                  # Dominio
│   │   │   ├── entities/
│   │   │   │   └── user.go
│   │   │   ├── ports/
│   │   │   │   ├── bcrypt_port.go
│   │   │   │   └── token_port.go
│   │   │   └── repository/
│   │   │       └── user_repository.go
│   │   │
│   │   └── infra/                   # Infraestructura
│   │       ├── dependencies.go
│   │       ├── adapters/
│   │       │   ├── bcrypt_adapter.go
│   │       │   └── jwt_manager.go
│   │       ├── controllers/
│   │       ├── repository/
│   │       ├── routes/
│   │       └── services/
│   │
│   ├── category/                    # 📦 MÓDULO CATEGORY
│   │   ├── app/
│   │   │   ├── get_all_categories.go
│   │   │   └── get_category.go
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   │   └── category.go
│   │   │   └── repository/
│   │   │       └── category_repository.go
│   │   └── infra/
│   │       ├── dependencies.go
│   │       ├── controllers/
│   │       ├── repository/
│   │       └── routes/
│   │
│   ├── product/                     # 📦 MÓDULO PRODUCT
│   │   ├── app/
│   │   │   ├── get_products_by_category.go
│   │   │   ├── get_product.go
│   │   │   ├── search_products.go
│   │   │   └── create_custom_product.go
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   │   └── product.go
│   │   │   └── repository/
│   │   │       └── product_repository.go
│   │   └── infra/
│   │
│   ├── group/                       # 📦 MÓDULO GROUP
│   │   ├── app/
│   │   │   ├── create_group.go
│   │   │   ├── get_group.go
│   │   │   ├── get_my_groups.go
│   │   │   ├── update_group.go
│   │   │   ├── delete_group.go
│   │   │   ├── invite_member.go
│   │   │   ├── respond_invitation.go
│   │   │   ├── get_members.go
│   │   │   └── remove_member.go
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   │   ├── group.go
│   │   │   │   └── group_member.go
│   │   │   └── repository/
│   │   │       └── group_repository.go
│   │   └── infra/
│   │
│   ├── outing/                      # 📦 MÓDULO OUTING
│   │   ├── app/
│   │   │   ├── create_outing.go
│   │   │   ├── get_outing.go
│   │   │   ├── get_my_outings.go
│   │   │   ├── update_outing.go
│   │   │   ├── delete_outing.go
│   │   │   ├── add_participant.go
│   │   │   ├── confirm_participation.go
│   │   │   ├── remove_participant.go
│   │   │   ├── add_item.go
│   │   │   ├── update_item.go
│   │   │   ├── remove_item.go
│   │   │   ├── split_item.go
│   │   │   └── calculate_splits.go
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   │   ├── outing.go
│   │   │   │   ├── outing_participant.go
│   │   │   │   ├── outing_item.go
│   │   │   │   └── item_split.go
│   │   │   └── repository/
│   │   │       └── outing_repository.go
│   │   └── infra/
│   │
│   └── payment/                     # 📦 MÓDULO PAYMENT
│       ├── app/
│       │   ├── get_outing_payments.go
│       │   ├── register_payment.go
│       │   ├── confirm_payment.go
│       │   └── get_my_pending_payments.go
│       ├── domain/
│       │   ├── entities/
│       │   │   └── payment.go
│       │   └── repository/
│       │       └── payment_repository.go
│       └── infra/
│
└── docs/                            # Documentación adicional
    └── database.sql                 # Script de base de datos
```

---

## 🚀 Instalación

### Prerrequisitos

- Go 1.21 o superior
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
psql -U tu_usuario -d splitmeet -f docs/database.sql

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
4. ✅ Group         (grupos con invitaciones)
5. ✅ Outing        (salidas, participantes, items, splits)
6. ✅ Payment       (pagos con validaciones y auto-cancelación)
```

### Convenciones de Código

- **Naming**: camelCase para variables, PascalCase para tipos exportados
- **Errores**: Siempre manejar y propagar errores apropiadamente
- **Comentarios**: Documentar funciones públicas con GoDoc
- **Tests**: Escribir tests unitarios para casos de uso

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
    Create(ctx context.Context, e *entities.Example) error
    FindByID(ctx context.Context, id int64) (*entities.Example, error)
    // ... otros métodos
}

// app/create_example.go
type CreateExampleUseCase struct {
    repo repository.ExampleRepository
}

func (uc *CreateExampleUseCase) Execute(ctx context.Context, input CreateExampleInput) (*entities.Example, error) {
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
