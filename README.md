## OAuth Flow Verification

> OAuth callback validation relies on provider-issued authorization
> codes and session cookies established during browser-based authentication.
> As such, callback verification is demonstrated via browser network traces
> rather than Postman.

---

### 1️⃣ OAuth Initiation (Postman)

<p align="center">
  <img src="result/token.png" alt="OAuth Redirect to Google" width="800"/>
</p>

<p align="center">
  <em>
    AuthFlow initiates OAuth 2.0 authentication by redirecting the client to
    Google’s authorization endpoint, establishing a secure session for
    callback validation.
  </em>
</p>

---

### 2️⃣ Google Authentication & Consent

<p align="center">
  <img src="result/choose-account.png" alt="Google Account Selection" width="800"/>
</p>

<p align="center">
  <em>
    User authentication and consent are handled by Google, which acts as the
    external identity provider.
  </em>
</p>

---

### 3️⃣ OAuth Callback Handling (Authorization Code Exchange)

<p align="center">
  <img src="result/callback-state.png" alt="OAuth Callback Handling" width="800"/>
</p>

<p align="center">
  <em>
    Google redirects back to AuthFlow’s callback endpoint, where the
    authorization code and state are validated and exchanged for provider
    credentials.
  </em>
</p>

---

### 4️⃣ Application JWT Issuance

<p align="center">
  <img src="result/callback-token.png" alt="JWT Issuance via Callback" width="800"/>
</p>

<p align="center">
  <em>
    After successful OAuth validation, AuthFlow issues an application-scoped
    JWT and securely redirects it to the client application along with the
    mapped user identity.
  </em>
</p>

---

## Authentication Architecture

<p align="center">
  <img src="result/oauth-flowchart.svg" alt="AuthFlow OAuth Architecture" width="900"/>
</p>

<p align="center">
  <em>
    High-level authentication flow illustrating OAuth-based identity
    verification and application-level JWT issuance.
  </em>
</p>
