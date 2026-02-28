export default function ApiTable() {
    const endpoints = [
        {
            method: "GET",
            color: "text-blue-600",
            endpoint: "/auth/{provider}",
            description:
                "Initiate OAuth flow. Redirects to provider authorization endpoint (e.g., Google).",
        },
        {
            method: "GET",
            color: "text-blue-600",
            endpoint: "/auth/{provider}/callback",
            description:
                "OAuth callback. Validates authorization code & state, exchanges for credentials, issues JWT.",
        },
        {
            method: "POST",
            color: "text-green-600",
            endpoint: "/mint",
            description:
                "Issue a JWT. Accepts sub, provider, and optional email. Does not authenticate — assumes identity is already verified.",
        },
        {
            method: "POST",
            color: "text-red-600",
            endpoint: "/logout",
            description:
                "Revoke a JWT. Validates token, extracts jti, writes revoked:<jti> to Redis with TTL.",
        },
    ];

    return (
        <table className="spreadsheet-table">
            <thead>
                <tr>
                    <th style={{ width: "15%" }}>METHOD</th>
                    <th style={{ width: "35%" }}>ENDPOINT</th>
                    <th style={{ width: "50%" }}>DESCRIPTION</th>
                </tr>
            </thead>
            <tbody>
                {endpoints.map((ep, i) => (
                    <tr key={i}>
                        <td className={`font-bold ${ep.color}`}>{ep.method}</td>
                        <td>{ep.endpoint}</td>
                        <td>{ep.description}</td>
                    </tr>
                ))}
            </tbody>
        </table>
    );
}
