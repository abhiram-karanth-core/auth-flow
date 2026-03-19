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
# Authflow — Integration Guide

This guide covers how to integrate Authflow into your service as either a **client** (calling `/mint` to issue tokens) or a **downstream service** (verifying Authflow-issued JWTs on protected routes).

---

## Two Integration Roles

| Role | What you do |
|------|-------------|
| **Client** | Authenticate your user, then call `/mint` to get an Authflow JWT |
| **Downstream service** | Accept Authflow JWTs on protected routes and verify them via JWKS |

A single service can be both. The Flask example below is both — it calls `/mint` after login, and protects its own routes using the returned JWT.

---

## Step 0 — Register Your Application

Before calling `/mint`, register your app with Authflow to get a `client_id` and `client_secret`.

```bash
curl -X POST https://authflow-go.onrender.com/clients \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-app",
    "redirect_uri": "https://my-app.com/callback"
  }'
```

**Response — save these, the secret is shown only once:**
```json
{
  "client_id": "550e8400-e29b-41d4-a716-446655440000",
  "client_secret": "7f3d9a1b-...",
  "redirect_uri": "https://my-app.com/callback"
}
```

Store `AUTHFLOW_CLIENT_ID` and `AUTHFLOW_CLIENT_SECRET` as environment variables. Never hardcode them.

---

## How `/mint` Works

`/mint` is a **token issuance endpoint, not an authentication endpoint**. It assumes your service has already verified the user's identity. You are responsible for authentication — Authflow is responsible for the token.

```
Your service         Authflow
    |                    |
    |-- verify user ---→ |  (your logic: check password, OAuth, SSO)
    |                    |
    |-- POST /mint ----→ |  (with client_id + client_secret + sub)
    |                    |
    |←-- JWT ----------- |
    |                    |
    |-- return JWT --→ client
```

**Request:**
```json
POST /mint
Content-Type: application/json

{
  "client_id": "your-client-id",
  "client_secret": "your-client-secret",
  "sub": "username-or-user-id",
  "provider": "local"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIs..."
}
```

---

## How JWT Verification Works

Downstream services verify Authflow JWTs using the public key served at:

```
GET https://authflow-go.onrender.com/.well-known/jwks.json
```

No shared secret is needed. Each service fetches the public key and verifies signatures locally — no call to Authflow per request.

**Important:** Cache the public key with a TTL (e.g. 1 hour) and retry the JWKS fetch on verification failure. This handles key rotation without breaking your service.

---

## Python — Flask

**Dependencies:**
```bash
pip install flask flask-jwt-extended PyJWT requests cryptography
```

**Setup and protected route:**
```python
import os
import time
import requests
from flask import Flask, request, jsonify
from flask_jwt_extended import JWTManager, jwt_required, get_jwt_identity
from jwt.algorithms import RSAAlgorithm

app = Flask(__name__)

# Flask session secret — not the JWT key
app.config["SECRET_KEY"] = os.environ["FLASK_SECRET_KEY"]

# Tell flask-jwt-extended to use RS256 and Authflow's public key
app.config["JWT_ALGORITHM"] = "RS256"
app.config["JWT_PUBLIC_KEY"] = _get_public_key()

jwt = JWTManager(app)

# --- JWKS with TTL cache ---
_jwks_cache = {"key": None, "fetched_at": 0}
JWKS_URL = "https://authflow-go.onrender.com/.well-known/jwks.json"
JWKS_TTL = 3600  # 1 hour

def _get_public_key():
    now = time.time()
    if _jwks_cache["key"] and (now - _jwks_cache["fetched_at"]) < JWKS_TTL:
        return _jwks_cache["key"]
    jwks = requests.get(JWKS_URL).json()
    key = RSAAlgorithm.from_jwk(jwks["keys"][0])
    _jwks_cache.update({"key": key, "fetched_at": now})
    return key

# --- Login: authenticate, then mint ---
@app.route("/login", methods=["POST"])
def login():
    data = request.json
    username = data.get("username")
    password = data.get("password")

    if not username or not password:
        return jsonify({"error": "Missing credentials"}), 400

    user = User.query.filter_by(username=username).first()
    if not user or not user.check_password(password):
        return jsonify({"error": "Invalid credentials"}), 401

    resp = requests.post(
        "https://authflow-go.onrender.com/mint",
        json={
            "sub": username,
            "provider": "local",
            "client_id": os.environ["AUTHFLOW_CLIENT_ID"],
            "client_secret": os.environ["AUTHFLOW_CLIENT_SECRET"],
        }
    )
    if resp.status_code != 200:
        return jsonify({"error": "Auth service unavailable"}), 500

    return resp.json(), 200

# --- Protected route ---
@app.route("/upload", methods=["POST"])
@jwt_required()
def upload():
    current_user = get_jwt_identity()  # returns sub claim
    # your logic here
    return jsonify({"user": current_user})
```

---

## Python — FastAPI

**Dependencies:**
```bash
pip install fastapi python-jose[cryptography] httpx
```

```python
import time
import httpx
from fastapi import FastAPI, Depends, HTTPException
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from jose import jwt, JWTError
from jose.backends import RSAKey
import json

app = FastAPI()
security = HTTPBearer()

JWKS_URL = "https://authflow-go.onrender.com/.well-known/jwks.json"
AUTHFLOW_ISSUER = "authflow-go"
JWKS_TTL = 3600

_jwks_cache = {"key": None, "fetched_at": 0}

def get_public_key():
    now = time.time()
    if _jwks_cache["key"] and (now - _jwks_cache["fetched_at"]) < JWKS_TTL:
        return _jwks_cache["key"]
    resp = httpx.get(JWKS_URL)
    jwks = resp.json()
    _jwks_cache.update({"key": jwks["keys"][0], "fetched_at": now})
    return _jwks_cache["key"]

def verify_token(credentials: HTTPAuthorizationCredentials = Depends(security)):
    token = credentials.credentials
    try:
        public_key = get_public_key()
        claims = jwt.decode(
            token,
            public_key,
            algorithms=["RS256"],
            issuer=AUTHFLOW_ISSUER,
            options={"require_exp": True}
        )
        return claims["sub"]
    except JWTError:
        raise HTTPException(status_code=401, detail="Invalid token")

# --- Login: authenticate, then mint ---
@app.post("/login")
async def login(body: dict):
    username = body.get("username")
    password = body.get("password")

    # verify credentials against your DB here
    if not authenticate_user(username, password):
        raise HTTPException(status_code=401, detail="Invalid credentials")

    async with httpx.AsyncClient() as client:
        resp = await client.post(
            "https://authflow-go.onrender.com/mint",
            json={
                "sub": username,
                "provider": "local",
                "client_id": os.environ["AUTHFLOW_CLIENT_ID"],
                "client_secret": os.environ["AUTHFLOW_CLIENT_SECRET"],
            }
        )
    if resp.status_code != 200:
        raise HTTPException(status_code=500, detail="Auth service unavailable")

    return resp.json()

# --- Protected route ---
@app.post("/upload")
async def upload(current_user: str = Depends(verify_token)):
    # current_user is the sub claim
    return {"user": current_user}
```

---

## Node.js — Express

**Dependencies:**
```bash
npm install express jwks-rsa express-oauth2-jwt-bearer axios
```

```javascript
const express = require('express');
const { auth } = require('express-oauth2-jwt-bearer');
const axios = require('axios');

const app = express();
app.use(express.json());

// --- JWT middleware ---
const checkJwt = auth({
  jwksUri: 'https://authflow-go.onrender.com/.well-known/jwks.json',
  issuer: 'authflow-go',
  algorithms: ['RS256'],
});

// --- Login: authenticate, then mint ---
app.post('/login', async (req, res) => {
  const { username, password } = req.body;
  if (!username || !password)
    return res.status(400).json({ error: 'Missing credentials' });

  const user = await db.findUser(username);
  if (!user || !user.verifyPassword(password))
    return res.status(401).json({ error: 'Invalid credentials' });

  try {
    const { data } = await axios.post(
      'https://authflow-go.onrender.com/mint',
      {
        sub: username,
        provider: 'local',
        client_id: process.env.AUTHFLOW_CLIENT_ID,
        client_secret: process.env.AUTHFLOW_CLIENT_SECRET,
      }
    );
    return res.json(data);
  } catch {
    return res.status(500).json({ error: 'Auth service unavailable' });
  }
});

// --- Protected route ---
app.post('/upload', checkJwt, (req, res) => {
  const currentUser = req.auth.payload.sub;
  // your logic here
  res.json({ user: currentUser });
});
```

---

## Java — Spring Boot

**`pom.xml` dependency:**
```xml
<dependency>
  <groupId>org.springframework.boot</groupId>
  <artifactId>spring-boot-starter-oauth2-resource-server</artifactId>
</dependency>
```

**`application.yml`:**
```yaml
spring:
  security:
    oauth2:
      resourceserver:
        jwt:
          jwk-set-uri: https://authflow-go.onrender.com/.well-known/jwks.json
          issuer-uri: https://authflow-go.onrender.com
```

**Security config:**
```java
@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        http
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/login", "/register").permitAll()
                .anyRequest().authenticated()
            )
            .oauth2ResourceServer(oauth2 -> oauth2.jwt(Customizer.withDefaults()));
        return http.build();
    }
}
```

**Login endpoint (calls `/mint`):**
```java
@RestController
public class AuthController {

    @Value("${authflow.client-id}")
    private String clientId;

    @Value("${authflow.client-secret}")
    private String clientSecret;

    private final RestTemplate restTemplate = new RestTemplate();

    @PostMapping("/login")
    public ResponseEntity<?> login(@RequestBody LoginRequest body) {
        User user = userRepository.findByUsername(body.getUsername())
            .orElseThrow(() -> new ResponseStatusException(HttpStatus.UNAUTHORIZED));

        if (!passwordEncoder.matches(body.getPassword(), user.getPasswordHash()))
            throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Invalid credentials");

        Map<String, String> mintBody = Map.of(
            "sub", body.getUsername(),
            "provider", "local",
            "client_id", clientId,
            "client_secret", clientSecret
        );

        ResponseEntity<Map> resp = restTemplate.postForEntity(
            "https://authflow-go.onrender.com/mint", mintBody, Map.class
        );

        return ResponseEntity.ok(resp.getBody());
    }

    // Protected route
    @PostMapping("/upload")
    public ResponseEntity<?> upload(@AuthenticationPrincipal Jwt jwt) {
        String currentUser = jwt.getSubject(); // sub claim
        // your logic here
        return ResponseEntity.ok(Map.of("user", currentUser));
    }
}
```

**`application.yml` for client credentials:**
```yaml
authflow:
  client-id: ${AUTHFLOW_CLIENT_ID}
  client-secret: ${AUTHFLOW_CLIENT_SECRET}
```

---

## Go — chi

**Dependencies:**
```bash
go get github.com/go-chi/chi/v5
go get github.com/golang-jwt/jwt/v5
go get github.com/MicahParks/keyfunc/v2
```

```go
package main

import (
    "context"
    "encoding/json"
    "net/http"
    "os"
    "time"

    "github.com/MicahParks/keyfunc/v2"
    "github.com/go-chi/chi/v5"
    "github.com/golang-jwt/jwt/v5"
)

const jwksURL = "https://authflow-go.onrender.com/.well-known/jwks.json"
const authflowIssuer = "authflow-go"

// JWKS with automatic refresh
var jwks *keyfunc.JWKS

func init() {
    var err error
    jwks, err = keyfunc.NewDefault([]string{jwksURL})
    if err != nil {
        panic("failed to fetch JWKS: " + err.Error())
    }
}

// JWTAuth middleware
func JWTAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenStr := r.Header.Get("Authorization")
        if len(tokenStr) < 8 || tokenStr[:7] != "Bearer " {
            http.Error(w, "missing token", http.StatusUnauthorized)
            return
        }
        tokenStr = tokenStr[7:]

        token, err := jwt.Parse(tokenStr, jwks.Keyfunc,
            jwt.WithIssuer(authflowIssuer),
            jwt.WithValidMethods([]string{"RS256"}),
        )
        if err != nil || !token.Valid {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        ctx := context.WithValue(r.Context(), "sub", claims["sub"])
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Login: authenticate, then mint
func loginHandler(w http.ResponseWriter, r *http.Request) {
    var body struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    json.NewDecoder(r.Body).Decode(&body)

    // verify credentials against your DB here
    if !authenticateUser(body.Username, body.Password) {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }

    payload, _ := json.Marshal(map[string]string{
        "sub":           body.Username,
        "provider":      "local",
        "client_id":     os.Getenv("AUTHFLOW_CLIENT_ID"),
        "client_secret": os.Getenv("AUTHFLOW_CLIENT_SECRET"),
    })

    resp, err := http.Post(
        "https://authflow-go.onrender.com/mint",
        "application/json",
        bytes.NewReader(payload),
    )
    if err != nil || resp.StatusCode != 200 {
        http.Error(w, "auth service unavailable", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    w.Header().Set("Content-Type", "application/json")
    io.Copy(w, resp.Body)
}

// Protected route
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    currentUser := r.Context().Value("sub").(string)
    // your logic here
    json.NewEncoder(w).Encode(map[string]string{"user": currentUser})
}

func main() {
    r := chi.NewRouter()
    r.Post("/login", loginHandler)
    r.Group(func(pr chi.Router) {
        pr.Use(JWTAuth)
        pr.Post("/upload", uploadHandler)
    })
    http.ListenAndServe(":8080", r)
}
```

---

## Logout

Send the user's Authflow JWT to invalidate it server-side. The token is revoked in Redis immediately — it will be rejected even if it hasn't expired yet.

```bash
curl -X POST https://authflow-go.onrender.com/logout \
  -H "Authorization: Bearer <access_token>"
```

All language examples — call this endpoint with the user's token in the `Authorization` header. No other cleanup is needed on the Authflow side.

---

## OAuth Flow

For browser-based login via Google, GitHub, etc., redirect the user to Authflow instead of handling OAuth yourself.

**1. Redirect user to Authflow:**
```
GET https://authflow-go.onrender.com/auth/{provider}
    ?client_id=your-client-id
    &redirect_uri=https://your-app.com/callback
```

**2. Authflow handles OAuth, then redirects back:**
```
GET https://your-app.com/callback?token=eyJhbGci...
```

**3. Your callback endpoint receives the token and stores it:**
```python
# Flask example
@app.route("/callback")
def callback():
    token = request.args.get("token")
    if not token:
        return "Missing token", 400
    # store token in session or return to frontend
    return jsonify({"access_token": token})
```

The token is a standard Authflow JWT — verify it the same way as any other.

---

## Environment Variables Reference

| Variable | Used by | Description |
|----------|---------|-------------|
| `AUTHFLOW_CLIENT_ID` | Client services | Issued by `POST /clients` |
| `AUTHFLOW_CLIENT_SECRET` | Client services | Issued by `POST /clients`, shown once |
| `AUTHFLOW_ISSUER` | Downstream services | Always `authflow-go` |
| `AUTHFLOW_JWKS_URL` | Downstream services | `https://authflow-go.onrender.com/.well-known/jwks.json` |

---

## JWT Claims Reference

| Claim | Description |
|-------|-------------|
| `sub` | Subject — user identifier passed to `/mint` |
| `iss` | Issuer — always `authflow-go` |
| `jti` | Unique token ID — used for revocation |
| `iat` | Issued at |
| `exp` | Expiry |
| `email` | Email, if provided |
| `provider` | Auth source — `local`, `google`, `sso` |
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