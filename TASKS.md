# MVP Commit Plan

## 1. Project foundation

- [ ] Initialize Go module
- [ ] Add basic backend folder structure
- [ ] Add `.gitignore`
- [ ] Add `.env.example`
- [ ] Add `Makefile`
- [ ] Add Docker Compose with PostgreSQL
- [ ] Add health endpoint
- [ ] Add basic HTTP server startup
- [ ] Add graceful shutdown

## 2. Database foundation

- [ ] Add migration tool setup
- [ ] Add initial migration for `users`
- [ ] Add initial migration for `roles`
- [ ] Add initial migration for `user_roles`
- [ ] Add initial migration for `properties`
- [ ] Add initial migration for `service_providers`
- [ ] Add initial migration for `property_service_providers`
- [ ] Add initial migration for `tasks`
- [ ] Add initial migration for `service_entries`
- [ ] Add initial migration for `service_entry_tasks`
- [ ] Add database indexes and constraints
- [ ] Add seed data for initial roles

## 3. SQL and persistence

- [ ] Add sqlc configuration
- [ ] Add user queries
- [ ] Add role queries
- [ ] Add property queries
- [ ] Add service provider queries
- [ ] Add task queries
- [ ] Add service entry queries
- [ ] Generate sqlc code
- [ ] Add database connection package

## 4. API foundation

- [ ] Add API error response format
- [ ] Add request validation helpers
- [ ] Add JSON response helpers
- [ ] Add route registration structure
- [ ] Add API version prefix `/api/v1`

## 5. Security

- [ ] Add password hashing utility
- [ ] Add user login endpoint
- [ ] Add session creation
- [ ] Add session persistence
- [ ] Add auth middleware
- [ ] Add logout endpoint
- [ ] Add `/me` endpoint
- [ ] Add role middleware
- [ ] Add property ownership authorization helper

## 6. Property management

- [ ] Add create property endpoint
- [ ] Add list properties endpoint
- [ ] Add get property endpoint
- [ ] Add update property endpoint
- [ ] Add deactivate property endpoint
- [ ] Add QR token generation for properties
- [ ] Add QR code response endpoint

## 7. Service providers

- [ ] Add create service provider endpoint
- [ ] Add list service providers endpoint
- [ ] Add update service provider endpoint
- [ ] Add deactivate service provider endpoint
- [ ] Add assign service provider to property endpoint
- [ ] Add remove service provider from property endpoint
- [ ] Add list service providers for property endpoint

## 8. Tasks

- [ ] Add create task endpoint
- [ ] Add list tasks by property endpoint
- [ ] Add update task endpoint
- [ ] Add deactivate task endpoint
- [ ] Add task sorting support

## 9. Public check-in flow

- [ ] Add resolve QR token endpoint
- [ ] Add public check-in response model
- [ ] Add start service entry endpoint
- [ ] Add complete service entry endpoint
- [ ] Add task completion persistence
- [ ] Add validation for assigned service provider
- [ ] Add validation for active property and active tasks

## 10. History

- [ ] Add list service entries by property endpoint
- [ ] Add get service entry detail endpoint
- [ ] Add duration calculation
- [ ] Add filters for date range
- [ ] Add filters for service provider

## 11. Testing

- [ ] Add test database setup
- [ ] Add auth handler tests
- [ ] Add property authorization tests
- [ ] Add property CRUD tests
- [ ] Add service provider assignment tests
- [ ] Add task CRUD tests
- [ ] Add QR token resolution tests
- [ ] Add service entry lifecycle tests
- [ ] Add service history tests

## 12. Frontend foundation

- [ ] Initialize React app
- [ ] Add routing
- [ ] Add API client
- [ ] Add auth state handling
- [ ] Add mobile-first layout shell
- [ ] Add login page
- [ ] Add property list page
- [ ] Add property detail page
- [ ] Add task management UI
- [ ] Add service provider management UI
- [ ] Add QR code display UI
- [ ] Add public check-in page
- [ ] Add service history UI
