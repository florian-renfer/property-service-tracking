# Status - 2026-06-16 21:42 UTC+2

Simplified Docker Compose startup ordering by removing keycloak-config service and using OIDC endpoint health checks instead of management port, with services gated: postgres (healthy) → keycloak (OIDC ready) → api, eliminating kcadm complexity; next: verify all services start cleanly, test /me endpoint with valid token, then address aud vs. azp claim validation.
