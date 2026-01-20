# Authflow-Go

Authflow-Go is a lightweight **authentication and authorization service** written in Go.  
It handles **user authentication via OAuth providers** (Google) and **issues signed JWTs** that can be used by downstream services for authorization.

This project is designed to act as a **central auth service** in a microservice or API-based architecture.

---

## What This Project Is

Authflow-Go functions as:

 **Authentication Server**  
- Handles OAuth login (Google)
- Manages secure user sessions

 **Authorization Token Issuer**  
- Issues signed JWT access tokens
- Tokens can be consumed by other services for authorization

➡️ In practice, this makes Authflow-Go an **Authorization Server for your system**, though it is **not a full OAuth 2.0 Authorization Server implementation** (like Keycloak or Auth0).

---

## 🚫 What This Project Is NOT

- ❌ Not a complete OAuth 2.0 / OpenID Connect provider
- ❌ Does not issue authorization codes to third-party clients
- ❌ Does not manage users in a database (yet)

Instead, it focuses on **simple, secure auth flows** for modern backend systems.

---

## Features

- Google OAuth 2.0 authentication (via Goth)
- Secure cookie-based session management
- JWT access token generation (HS256)
- Environment-aware security settings
- Clean Go project structure
- Ready for microservice integration

---


## Authentication Flow

1. User initiates OAuth login
2. Google authenticates the user
3. Authflow-Go creates a secure session
4. A JWT access token is generated
5. Client uses JWT to access protected APIs

---

## JWT Claims

Tokens issued by Authflow-Go contain:

- `email` – authenticated user identifier
- `sub` – subject (user)
- `iss` – issuer (`authflow-go`)
- `iat` – issued at
- `exp` – expiration time (24 hours)

---