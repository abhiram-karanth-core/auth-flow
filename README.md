# Authflow

Authflow is a centralized authorization service written in Go. It issues application-level JWTs, integrates with OAuth providers for identity verification, and enforces secure logout using Redis-backed token revocation.

The service is designed to be used by platforms that manage how users authenticate (OAuth, username/password, SSO, etc.) while delegating authorization, token lifecycle, and revocation to a single trusted service.

---

## What This Project Is

Authflow functions as:

**Authorization Service (Primary Role)**
- Issues signed JWT access tokens (RS256)
- Enforces token revocation using Redis
- Acts as the single JWT authority for the system
- Exposes a JWKS endpoint for downstream verification

**Authentication Integrator (Secondary Role)**
- Supports OAuth 2.0 login (Google via Goth) for identity verification
- Establishes secure browser sessions during OAuth flows

### Key Design Principle

Authentication and authorization are intentionally decoupled.

- Platforms decide how users authenticate
- Authflow-Go decides how access is granted and revoked

This enables a single JWT issuer, consistent logout semantics, stateless downstream services, and centralized security guarantees.

---

## Features

- Centralized JWT issuance (`/mint`)
- RS256 signing with public key distribution via JWKS (`/.well-known/jwks.json`)
- OAuth 2.0 authentication (Google via Goth)
- Redis-backed JWT revocation (secure logout)
- Stateless JWT validation middleware
- Provider-agnostic token design
- Microservice-friendly architecture

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

## Token Issuance (`/mint`)

Authflow-Go exposes a token minting endpoint that issues JWTs after authentication has already occurred.

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

> **Important:** `/mint` does not authenticate users. It assumes the caller has already verified identity. This allows platforms to use any authentication strategy.

---

## JWKS Endpoint

Downstream services verify tokens using the public key served at:

```
GET /.well-known/jwks.json
```

No secret is shared. Downstream services fetch the public key from this endpoint and verify JWT signatures locally. This is the only configuration a downstream service needs to integrate with Authflow-Go.

See the [Integration Guide](#plugging-into-your-service) below.

---

## Logout & Token Revocation

Logout is implemented as JWT revocation, not session destruction.

**Flow**
1. Client calls `POST /logout` with a Bearer token
2. Authflow-Go validates the JWT, extracts the `jti` claim, and writes `revoked:<jti>` to Redis
3. Redis entry expires automatically when the JWT would expire

**Why Redis?**
- Immediate logout
- Stateless JWT validation
- Automatic cleanup via TTL
- No persistent blacklist storage

**Revocation model**
```
revoked:<jti> → "1"  (TTL = token_expiry − now)
```

Any request using a revoked token is rejected by middleware, even if the JWT is otherwise valid.

---

## JWT Claims

| Claim | Description |
|-------|-------------|
| `sub` | Subject (user identifier) |
| `iss` | Issuer (`authflow-go`) |
| `jti` | Unique token ID |
| `iat` | Issued at |
| `exp` | Expiration time |
| `email` | User email (OAuth logins) |
| `provider` | Authentication source (`google`, `local`) |

---

## Authentication & Authorization Flow

### OAuth Flow

1. User initiates OAuth login via `GET /auth/{provider}`
2. OAuth provider (Google) authenticates the user
3. Authflow-Go validates the callback, exchanges the authorization code, and establishes a session
4. Authflow-Go issues an RS256-signed application JWT
5. Client is redirected to the registered `redirect_uri` with the token

### Protected Resource Access

1. Client sends `Authorization: Bearer <JWT>`
2. Middleware verifies signature using the JWKS public key, checks expiry, and checks Redis for `revoked:<jti>`
3. Request proceeds only if the token is valid and not revoked

---

## Plugging Into Your Service

Downstream services only need two things: the JWKS URL and the issuer name. No shared secrets.

### Python — Flask

```python
from jwt.algorithms import RSAAlgorithm
import requests

jwks = requests.get('https://authflow-go.onrender.com/.well-known/jwks.json').json()
app.config['JWT_ALGORITHM'] = 'RS256'
app.config['JWT_PUBLIC_KEY'] = RSAAlgorithm.from_jwk(jwks['keys'][0])

jwt = JWTManager(app)  # must come after config

# protect routes
@app.route('/upload', methods=['POST'])
@jwt_required()
def upload():
    current_user = get_jwt_identity()  # returns sub
```

### Python — FastAPI

```python
def verify_token(token = Depends(security)):
    claims = jwt.decode(token.credentials, get_public_key(), algorithms=['RS256'], issuer='authflow-go')
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
# application.yml
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
| Single JWT issuer | ✅ Authflow-Go only |
| Asymmetric signing (RS256) | ✅ Private key never leaves auth server |
| Public key distribution | ✅ Via JWKS endpoint |
| Immediate logout | ✅ Redis revocation |
| Stateless downstream validation | ✅ No auth server call per request |
| Downstream services mint tokens | ❌ Never |

---

## Live Integration

Authflow-Go is running in production as the authentication layer for:

**[rag-works.vercel.app](https://rag-works.vercel.app)**

You can test the full flow by creating an account (Google OAuth or username + password), logging in, accessing protected routes, and logging out. All auth operations are powered by Authflow-Go.