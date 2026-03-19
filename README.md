# Authflow

Authflow is an API-first identity and token service written in Go. It issues application-level JWTs, supports OAuth-based identity verification, and enforces secure logout using Redis-backed token revocation.

Authflow is designed for platforms that want a single JWT issuer for all downstream services, while keeping authentication methods flexible. Applications can integrate with Authflow using standard HTTP redirects and API calls — no client SDK is required.

---

## What Authflow Is

Authflow acts as a centralized trust layer for modern applications and microservices.

**Authorization Service — Primary Role**
- Issues signed JWT access tokens (RS256)
- Enforces token revocation using Redis
- Acts as the single JWT authority for the system
- Exposes a JWKS endpoint for downstream verification

**Authentication Integrator — Secondary Role**
- Supports OAuth 2.0 login for identity verification
- Accepts upstream-authenticated identity via `/mint`
- Keeps authentication method decoupled from token authority

### Design Goals

- Single JWT authority across the system
- Provider-agnostic authentication inputs — OAuth, local login, or SSO
- Stateless downstream verification using JWKS
- Immediate logout semantics using Redis-backed revocation
- SDK-less integration through standard HTTP endpoints

---

## Why Authflow

Most authentication products focus primarily on login UI, sessions, and frontend SDKs.

Authflow is centered on a different responsibility:

> No matter how a user authenticates, Authflow becomes the single service responsible for issuing, validating, and revoking application JWTs.

This gives downstream services a stable and simple trust model:

- Trust one issuer
- Verify signatures locally
- Check revocation centrally
- Stay independent of the original login method

---

## Key Architecture Principle

Authentication and token authority are intentionally separated.

- Platforms decide how users authenticate — OAuth, username/password, SSO, or another trusted method
- Authflow decides how access is granted and revoked

This lets applications evolve their login methods without changing downstream authorization logic.

---

## Features

- Centralized JWT issuance (`/mint`)
- RS256 signing with public key distribution via JWKS
- OAuth 2.0 login support
- Redis-backed JWT revocation
- Stateless JWT verification for downstream services
- Multi-client registration with redirect URI validation
- API-first, SDK-less integration model
- Microservice-friendly design

---

## SDK-less Integration Model

Authflow does not require platform-specific SDKs.

Applications integrate directly with Authflow using:

- Browser redirects for OAuth flows
- HTTP API calls for token issuance and logout
- Bearer tokens for authenticated requests
- JWKS for downstream JWT verification

This makes Authflow usable from React, MERN, Django, Spring Boot, FastAPI, mobile apps, and other HTTP-capable clients without framework lock-in.

---

## Core Endpoints

```
GET  /auth/{provider}
GET  /auth/{provider}/callback

POST /mint
POST /logout

GET  /.well-known/jwks.json
```

---

## Supported Flows

Authflow supports two broad integration patterns.

### 1. OAuth Login Flow

Use Authflow as the OAuth entry point for browser-based login.

**Flow:**

1. Client redirects the user to `GET /auth/{provider}`
2. The OAuth provider authenticates the user
3. Authflow validates the callback
4. Authflow issues an application JWT
5. Authflow redirects the user to the registered `redirect_uri`

This is useful for browser apps that want Authflow to handle OAuth and return an application token.

### 2. Token Issuance After Upstream Authentication

Use Authflow to mint an application JWT after identity has already been verified by the calling platform.

**Flow:**

1. The client platform authenticates the user using its own method
2. The platform calls `POST /mint`
3. Authflow issues an application JWT
4. Downstream services trust only the Authflow-issued token

This is useful when the platform wants to control the login experience but delegate token authority and revocation to Authflow.

---

## Token Issuance (`/mint`)

`POST /mint` issues an Authflow JWT for an already-authenticated user.

**Request**
```json
POST /mint
Content-Type: application/json

{
  "sub": "username",
  "provider": "google | local | sso",
  "email": "user@example.com"
}
```

**Response**
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIs..."
}
```

> **Important:** `/mint` is a token issuance endpoint, not an authentication endpoint. It assumes the caller has already verified the user's identity using a trusted flow. This allows platforms to use any authentication strategy while keeping one token authority.

---

## JWKS Endpoint

Downstream services verify Authflow-issued JWTs using the public key served at:

```
GET /.well-known/jwks.json
```

No shared secret is required. Each downstream service fetches the public key from Authflow and verifies JWT signatures locally. This keeps services stateless and avoids calling the auth server on every request.

---

## Logout and Token Revocation

Authflow implements logout as JWT revocation, not just browser session destruction.

**Flow:**

1. Client calls `POST /logout` with a Bearer token
2. Authflow validates the JWT
3. Authflow extracts the `jti` claim
4. Authflow stores `revoked:<jti>` in Redis
5. The Redis entry expires automatically when the token would normally expire

**Why Redis:**

- Immediate logout
- Stateless JWT validation
- Automatic cleanup via TTL
- No persistent blacklist storage

**Revocation model:**
```
revoked:<jti> -> "1"   (TTL = token_expiry - now)
```

A revoked token is rejected even if its JWT signature and expiry are otherwise valid.

---

## JWT Claims

| Claim | Description |
|-------|-------------|
| `sub` | Subject (user identifier) |
| `iss` | Issuer (`authflow-go`) |
| `jti` | Unique token ID |
| `iat` | Issued at |
| `exp` | Expiration time |
| `email` | User email |
| `provider` | Authentication source (`google`, `local`, `sso`) |

---

## Authentication and Authorization Model

Authflow focuses on token trust, not application-specific permission logic.

Authflow answers:

- Who issued this token?
- Is the token valid?
- Has the token been revoked?

The client platform still owns business authorization decisions such as:

- Is this user an admin?
- Can this user upload files?
- Can this user access tenant A but not tenant B?

This separation keeps application roles and business rules independent from the identity provider.

---

## Plugging Into Your Service

Downstream services only need two things:

- The Authflow issuer (`authflow-go`)
- The Authflow JWKS URL

No shared secret is required.

### Python — Flask
```python
from jwt.algorithms import RSAAlgorithm
import requests

jwks = requests.get('https://authflow-go.onrender.com/.well-known/jwks.json').json()
app.config['JWT_ALGORITHM'] = 'RS256'
app.config['JWT_PUBLIC_KEY'] = RSAAlgorithm.from_jwk(jwks['keys'][0])

jwt = JWTManager(app)  # must come after config

@app.route('/upload', methods=['POST'])
@jwt_required()
def upload():
    current_user = get_jwt_identity()  # returns sub
```

### Python — FastAPI
```python
def verify_token(token = Depends(security)):
    claims = jwt.decode(
        token.credentials,
        get_public_key(),
        algorithms=['RS256'],
        issuer='authflow-go'
    )
    return claims['sub']

@app.post('/upload')
def upload(current_user: str = Depends(verify_token)):
    ...
```

### Node.js — Express
```javascript
const checkJwt = expressjwt({
  secret: jwksRsa.expressJwtSecret({
    jwksUri: 'https://authflow-go.onrender.com/.well-known/jwks.json',
    cache: true,
  }),
  algorithms: ['RS256'],
  issuer: 'authflow-go',
});

app.post('/upload', checkJwt, (req, res) => {
  const currentUser = req.auth.sub;
});
```

### Java — Spring Boot
```yaml
spring:
  security:
    oauth2:
      resourceserver:
        jwt:
          jwk-set-uri: https://authflow-go.onrender.com/.well-known/jwks.json
          issuer-uri: https://authflow-go.onrender.com
```
```java
@PostMapping("/upload")
public ResponseEntity upload(@AuthenticationPrincipal Jwt jwt) {
    String currentUser = jwt.getSubject();
}
```

---

## Security Model

| Property | Status |
|----------|--------|
| Single JWT issuer |  Authflow only |
| Asymmetric signing (RS256) |  Private key never leaves Authflow |
| Public key distribution |  Via JWKS |
| Immediate logout |  Redis-backed revocation |
| Stateless downstream validation |  No auth server call per request |
| Downstream services mint tokens |  Never |
| Shared secret across services |  Not required |

---

## Good Fit For

Authflow is a strong fit for:

- Microservice architectures
- MERN and SPA backends that need one JWT authority
- Polyglot systems using Node, Python, Java, or Go
- Platforms that want provider-agnostic authentication inputs
- Teams that want local JWT verification with central revocation
- Products that prefer API-first auth over SDK-heavy integration

---

## Live Integration

Authflow is currently used in production as the authentication and token layer for:

**[rag-works.vercel.app](https://rag-works.vercel.app)**

You can test the full flow by creating an account with Google OAuth or username/password, logging in, accessing protected routes, and logging out.