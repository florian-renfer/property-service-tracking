🏢 Property Service Tracking

A mobile-first platform for documenting and verifying property maintenance services.

🎯 Goal

Property managers often rely on trust when working with external service providers such as:

* 🧹 Cleaning companies
* 🌳 Garden maintenance providers
* 🏠 Caretakers
* ❄️ Winter services
* 🔧 Facility management companies

This project provides a simple way to document:

* When a service was performed
* How long it took
* Which tasks were completed
* Which service provider performed the work

The goal is to create transparency and accountability without introducing unnecessary complexity for service workers.

⸻

🚀 MVP Scope

Property Manager

* 🔐 Login
* 🏢 Create and manage properties
* 🤝 Create and assign service providers
* ✅ Define tasks per property
* 🔳 Generate QR codes
* 📋 Review service history

Service Worker

* 📱 Scan QR code
* 🤝 Select service provider
* ✍️ Enter name
* ▶️ Start service
* ✅ Complete tasks
* ⏹️ Finish service

System

* ⏱️ Store start time
* ⏱️ Store end time
* 📊 Calculate duration
* 📝 Store completed tasks
* 📚 Maintain service history

⸻

❌ Out of Scope

The following features are intentionally excluded from the MVP:

* 📍 GPS tracking
* 📷 Photo uploads
* 🔔 Push notifications
* 🐞 Defect reporting
* 📄 Invoice management
* 👥 Employee management
* 📱 Native mobile apps
* ⚙️ Microservices

⸻

👤 Roles

Global Admin

* Manage users
* Manage roles
* Access all properties

Property Manager

* Manage owned properties
* Manage service providers
* Manage tasks
* Review service history

Service Worker

No account required in the MVP.

⸻

🏗️ Architecture

Backend

* Go
* Chi
* PostgreSQL
* pgx
* sqlc
* golang-migrate

Frontend

* React
* TypeScript
* Mobile-first
* PWA-ready

Backend structure currently provides foundation only:

* `internal/domain` - domain entities and invariants
* `internal/app` - application services and use cases
* `internal/app/ports` - application ports (interfaces) for adapters
* `internal/adapters/inbound/web` - HTTP router and REST handlers
* `internal/adapters/outbound` - infrastructure adapters

Current API surface:

* `GET /api/v1/health`
* `POST /api/v1/properties`
* `GET /api/v1/properties`
* `GET /api/v1/properties/{id}`
* `PATCH /api/v1/properties/{id}`

⸻

🔐 Security

* Email/password authentication
* Session cookies
* Role-based access control
* Object-level authorization
* Password hashing with bcrypt

⸻

♿ Accessibility

The application targets:

* WCAG 2.1 AA
* EN 301 549

⸻

📋 MVP Checklist

Foundation

* Project setup
* PostgreSQL setup
* Database migrations
* sqlc configuration
* Health endpoint

Security

* Authentication
* Authorization
* Session management

Property Management

* Property CRUD
* Service Provider CRUD
* Property-Service Provider assignment

Task Management

* Task CRUD

Check-In Flow

* QR code generation
* QR code resolution
* Service start
* Task completion
* Service completion

History

* Service history
* Service details

⸻

🎉 Definition of Done

The MVP is complete when a property manager can:

1. Create a property
2. Assign service providers
3. Define tasks
4. Generate a QR code

And a service worker can:

1. Scan the QR code
2. Select a service provider
3. Enter their name
4. Complete a service entry

And the property manager can:

1. Review the completed service history

⸻

🛠️ Developer Setup

1. Create local env:
   * `cp .env.example .env`
2. Install tooling:
   * `make tools`
3. Start dependencies:
   * `docker compose up -d`
4. Start API:
   * `make run`

⸻

🧭 Feature Implementation Blueprint

New features should follow the hexagonal architecture already used in the codebase:

* Start in `internal/domain` when the feature introduces business state, invariants, or behavior.
* Add use cases in `internal/app` and define required ports in `internal/app/ports`.
* Implement inbound delivery mechanisms in `internal/adapters/inbound`, such as HTTP handlers.
* Implement outbound infrastructure in `internal/adapters/outbound`, such as PostgreSQL or in-memory repositories.
* Wire concrete adapters only at composition roots such as `cmd/api/main.go`.

Implementation standards:

* Keep domain code independent from HTTP, database, environment variables, framework types, and infrastructure errors.
* Keep application services focused on orchestration, transaction boundaries, authorization decisions, and port calls.
* Keep adapters responsible for transport, persistence, serialization, request parsing, and external system details.
* Prefer explicit input/output DTOs at application and adapter boundaries.
* Return named sentinel errors where callers need stable error handling with `errors.Is`.
* Do not let clients provide audit metadata directly; derive it from the authenticated principal or temporary request context.
* Avoid real infrastructure in application tests; use in-memory or fake port implementations.

Testing blueprint:

* Domain tests cover every invariant, constructor, behavior method, and mutation failure path.
* Application tests cover success paths, validation propagation, repository errors, not-found behavior, and no-side-effect guarantees.
* Adapter tests cover request/response contracts, JSON validation, status codes, error mapping, and persistence behavior.
* Integration-style HTTP tests should use `httptest` and real in-memory adapters, not external services.
* Every new feature should run with `go test ./...` and maintain high package-level coverage.
* Aim for 100% coverage in domain, application, and simple infrastructure adapters.
* For HTTP adapters, prioritize meaningful branch coverage over artificial tests, but keep coverage high and explain any intentionally uncovered defensive branch.

Production-grade checklist:

* Validate all external input at the adapter boundary and enforce business invariants in the domain.
* Use context-aware ports for request cancellation and future timeouts.
* Protect shared in-memory state with synchronization.
* Keep errors safe for clients; log internal detail at the boundary when logging exists.
* Preserve deterministic behavior in tests by controlling fakes, fixtures, and expected error paths.
* Keep feature changes small, cohesive, and covered before wiring them into runtime composition.

⸻

✅ Quality Gates

* Format: `make fmt` / check only: `make fmt-check`
* Lint: `make lint` (revive)
* Tests: `make test`
* Static checks: `make vet`

CI runs these checks on pushes to `main` and pull requests.
