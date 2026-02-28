"use client";

import { useState } from "react";

interface Tab {
    id: string;
    label: string;
    content: string;
}

const tabs: Tab[] = [
    {
        id: "mint",
        label: "Mint Token",
        content: `# Issue a JWT after authentication
POST /mint
Content-Type: application/json

{
  "sub": "username",
  "provider": "google | local | sso",
  "email": "user@example.com"
}

# Response
{
  "access_token": "eyJhbGciOiJIUzI1NiIs..."
}

# ⚠️ /mint does not authenticate users
# It assumes the caller has already verified identity`,
    },
    {
        id: "logout",
        label: "Logout / Revoke",
        content: `# Revoke a token (Redis-backed)
POST /logout
Authorization: Bearer <JWT>

# Flow:
# 1. Validates the JWT
# 2. Extracts the jti claim
# 3. Writes revoked:<jti> to Redis
# 4. Redis entry expires when JWT would expire

# Redis revocation model:
# revoked:<jti> → "1" (TTL = token_expiry − now)

# Any request using a revoked token is
# rejected by middleware, even if the
# JWT is otherwise valid.`,
    },
    {
        id: "oauth",
        label: "OAuth Flow",
        content: `# 1. Initiate OAuth login
GET /auth/google

# 2. Google authenticates the user
#    User sees consent screen

# 3. Callback handling
GET /auth/google/callback
# Authorization code + state validated
# Exchanged for provider credentials

# 4. Authflow issues application JWT
# Secure redirect with JWT + user identity

# 5. Client uses JWT for API access
Authorization: Bearer <access_token>`,
    },
    {
        id: "middleware",
        label: "Middleware (Go)",
        content: `// Protected resource access middleware
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := extractToken(r)
        
        // 1. Verify signature & expiry
        token, err := jwt.Parse(tokenString, getKey)
        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized", 401)
            return
        }
        
        // 2. Check Redis for revoked:<jti>
        jti := token.Claims["jti"].(string)
        if isRevoked(jti) {
            http.Error(w, "Token revoked", 401)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}`,
    },
];

export default function CodeTabs() {
    const [activeTab, setActiveTab] = useState("mint");

    return (
        <div>
            <div className="code-tab-header">
                {tabs.map((tab) => (
                    <div
                        key={tab.id}
                        className={`code-tab ${activeTab === tab.id ? "active" : ""}`}
                        onClick={() => setActiveTab(tab.id)}
                    >
                        {tab.label}
                    </div>
                ))}
            </div>
            {tabs.map((tab) => (
                <div
                    key={tab.id}
                    className="code-block"
                    style={{ display: activeTab === tab.id ? "block" : "none" }}
                >
                    {tab.content}
                </div>
            ))}
        </div>
    );
}
