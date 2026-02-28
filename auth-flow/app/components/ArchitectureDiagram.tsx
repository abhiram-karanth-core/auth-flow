export default function ArchitectureDiagram() {
    return (
        <div className="grid-2">
            <div className="p-4 border border-dashed border-gray-400">
                <h3 className="border-b border-gray-300 pb-2 mb-4">Key Design Principle</h3>
                <p className="font-mono text-sm">
                    Authentication and authorization are intentionally decoupled.
                    Platforms decide how users authenticate. Authflow-Go decides how
                    access is granted and revoked. This enables a single JWT issuer,
                    consistent logout semantics, and stateless downstream services.
                </p>
                <ul className="methodology-list mt-4">
                    <li>Authentication: Handled by OAuth Providers</li>
                    <li>Authorization: Centralized in Authflow-Go</li>
                    <li>Token Lifecycle: Controlled in one place</li>
                    <li>Revocation: Redis-backed with TTL cleanup</li>
                </ul>
            </div>
            <div className="bg-gray-50 border border-gray-200 p-8 flex items-center justify-center">
                <div className="relative w-full h-full min-h-[200px] font-mono text-xs">
                    <div className="absolute top-0 left-0 border-2 border-black p-2 bg-white z-10">
                        Client
                    </div>
                    <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 border-2 border-black p-4 z-20 shadow-xl" style={{ background: "#000", color: "#fff" }}>
                        AUTHFLOW
                    </div>
                    <div className="absolute bottom-0 right-0 border-2 border-black p-2 bg-white z-10">
                        Services
                    </div>
                    <div className="absolute top-0 right-0 border-2 border-dashed border-gray-500 p-2 bg-gray-100">
                        Google
                        <br />
                        OAuth
                    </div>
                    <div className="absolute bottom-0 left-0 border-2 border-dashed border-gray-500 p-2 bg-gray-100">
                        Redis
                    </div>

                    <svg
                        className="absolute inset-0 w-full h-full pointer-events-none"
                        style={{ zIndex: 1 }}
                    >
                        <line
                            x1="50"
                            y1="30"
                            x2="50%"
                            y2="50%"
                            stroke="black"
                            strokeWidth="1"
                        />
                        <line
                            x1="50%"
                            y1="50%"
                            x2="90%"
                            y2="30"
                            stroke="black"
                            strokeWidth="1"
                            strokeDasharray="4"
                        />
                        <line
                            x1="50%"
                            y1="50%"
                            x2="80%"
                            y2="90%"
                            stroke="black"
                            strokeWidth="1"
                        />
                        <line
                            x1="50%"
                            y1="50%"
                            x2="20%"
                            y2="90%"
                            stroke="black"
                            strokeWidth="1"
                            strokeDasharray="4"
                        />
                    </svg>
                </div>
            </div>
        </div>
    );
}
