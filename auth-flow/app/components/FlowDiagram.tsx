export default function FlowDiagram() {
    return (
        <div className="p-6 bg-gray-50 border border-gray-200 font-mono text-xs">
            <div className="flex justify-between items-center mb-8">
                <div className="diagram-box">Client</div>
                <div className="text-center w-full border-b border-dashed border-gray-400 relative">
                    <span className="absolute -top-3 bg-gray-50 px-2 left-1/2 -translate-x-1/2">
                        OAuth 2.0
                    </span>
                </div>
                <div className="diagram-box" style={{ background: "#000", color: "#fff" }}>Authflow-Go</div>
            </div>
            <div className="flex justify-center mb-2">
                <div className="h-8 border-l border-dashed border-gray-400" />
            </div>
            <div className="flex justify-between items-center mb-8">
                <div className="diagram-box w-24">Redis</div>
                <div className="w-full border-b border-dashed border-gray-400 relative">
                    <span className="absolute -top-3 bg-gray-50 px-2 left-1/2 -translate-x-1/2">
                        Revocation
                    </span>
                </div>
                <div className="diagram-box w-24">Services</div>
            </div>

            <div className="mt-6 pt-4 border-t border-gray-200">
                <span className="text-gray-500">{"// Token Mint Request"}</span>
                <pre className="text-gray-800 mt-2">
                    {`POST /mint
{
  "sub": "username",
  "provider": "google",
  "email": "user@example.com"
}`}
                </pre>
            </div>
        </div>
    );
}
