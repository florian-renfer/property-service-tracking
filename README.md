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
