import type { Metadata } from "next";
import GridBackground from "../components/GridBackground";
import TopNav from "../components/TopNav";
import HeroBanner from "../components/HeroBanner";
import ContentCard, { TitleBlock } from "../components/ContentCard";
import FlowDiagram from "../components/FlowDiagram";
import CodeTabs from "../components/CodeTabs";
import ArchitectureDiagram from "../components/ArchitectureDiagram";
import ApiTable from "../components/ApiTable";
import Link from "next/link";

export const metadata: Metadata = {
    title: "AUTHFLOW | DOCUMENTATION",
    description:
        "Centralized authorization service written in Go. Issues application-level JWTs, integrates with OAuth providers, and enforces secure logout using Redis-backed token revocation.",
};

const features = [
    {
        num: "01",
        title: "Centralized JWT Issuance",
        desc: "Single /mint endpoint issues signed JWT access tokens. One authority for your entire system.",
    },
    {
        num: "02",
        title: "OAuth 2.0 Authentication",
        desc: "Google OAuth via Goth for identity verification. Establishes secure browser sessions during OAuth flows.",
    },
    {
        num: "03",
        title: "Redis Token Revocation",
        desc: "Immediate logout via Redis. Stores revoked:<jti> with TTL matching token expiry. Automatic cleanup.",
    },
    {
        num: "04",
        title: "Stateless Validation",
        desc: "Downstream services verify JWT signatures locally. Middleware checks Redis for revocation. No token minting downstream.",
    },
    {
        num: "05",
        title: "Microservice Ready",
        desc: "Clean Go project structure. Provider-agnostic token design. Pluggable into any web application.",
    },
    {
        num: "06",
        title: "Decoupled Auth",
        desc: "Platforms decide how users authenticate. Authflow decides how access is granted and revoked. Single JWT issuer.",
    },
];

const steps = [
    { num: "01", label: "OAuth\nInitiation", highlight: false },
    { num: "02", label: "Provider\nAuth", highlight: false },
    { num: "03", label: "Callback\nValidation", highlight: false },
    { num: "04", label: "JWT\nIssuance", highlight: true },
    { num: "05", label: "Token\nRevocation", highlight: false },
];

export default function DocsPage() {
    return (
        <>
            <GridBackground />
            <TopNav />

            <div className="main-stage">
                {/* Hero Banner */}
                <HeroBanner />

                {/* Hero Content Card */}
                <ContentCard>
                    <div className="grid-2">
                        <div>
                            <span className="meta-data">
                                AUTHORIZATION SERVICE // WRITTEN IN GO
                            </span>
                            <h1 className="mb-6">
                                Centralized Authorization
                                <br />
                                for Modern Platforms.
                            </h1>
                            <p className="text-xl text-gray-600 mb-8 font-mono">
                                Issue application-level JWTs, integrate with OAuth providers for
                                identity verification, and enforce secure logout using
                                Redis-backed token revocation.
                            </p>
                            <div className="flex gap-4">
                                <Link href="#quickstart" className="btn-primary">
                                    Quick Start_
                                </Link>
                                <Link href="#architecture" className="btn-secondary">
                                    View Architecture_
                                </Link>
                            </div>
                        </div>
                        <div className="border-l border-gray-300 pl-8 hidden md:block">
                            <span className="meta-data mb-4">FLOW DIAGRAM</span>
                            <FlowDiagram />
                        </div>
                    </div>
                </ContentCard>

                {/* Feature Cards */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
                    {features.map((f) => (
                        <ContentCard key={f.num}>
                            <div className="feature-icon">{f.num}</div>
                            <h3>{f.title}</h3>
                            <p className="text-sm font-mono text-gray-600">{f.desc}</p>
                        </ContentCard>
                    ))}
                </div>

                {/* How It Works */}
                <ContentCard>
                    <TitleBlock meta="AUTHENTICATION & AUTHORIZATION FLOW" title="How It Works" />
                    <div className="flex flex-col md:flex-row justify-between items-center gap-4 text-center font-mono text-xs overflow-x-auto py-4">
                        {steps.map((step, i) => (
                            <div key={step.num} className="contents">
                                <div
                                    className={`p-4 border border-black min-w-[120px] ${step.highlight ? "bg-black text-white" : ""
                                        }`}
                                >
                                    <div className="font-bold mb-1">STEP {step.num}</div>
                                    {step.label.split("\n").map((line, j) => (
                                        <span key={j}>
                                            {line}
                                            {j < step.label.split("\n").length - 1 && <br />}
                                        </span>
                                    ))}
                                </div>
                                {i < steps.length - 1 && (
                                    <div className="text-2xl hidden md:block">→</div>
                                )}
                            </div>
                        ))}
                    </div>
                </ContentCard>

                {/* JWT Claims */}
                <ContentCard id="claims">
                    <TitleBlock meta="TOKEN STRUCTURE" title="JWT Claims" />
                    <table className="spreadsheet-table">
                        <thead>
                            <tr>
                                <th style={{ width: "20%" }}>CLAIM</th>
                                <th style={{ width: "80%" }}>DESCRIPTION</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr>
                                <td className="font-bold">sub</td>
                                <td>Subject (user identifier)</td>
                            </tr>
                            <tr>
                                <td className="font-bold">iss</td>
                                <td>Issuer (authflow-go)</td>
                            </tr>
                            <tr>
                                <td className="font-bold">jti</td>
                                <td>Unique token ID</td>
                            </tr>
                            <tr>
                                <td className="font-bold">iat</td>
                                <td>Issued at</td>
                            </tr>
                            <tr>
                                <td className="font-bold">exp</td>
                                <td>Expiration time</td>
                            </tr>
                            <tr>
                                <td className="font-bold">provider</td>
                                <td>Authentication source (e.g., google, local)</td>
                            </tr>
                        </tbody>
                    </table>
                </ContentCard>

                {/* Quick Start */}
                <ContentCard id="quickstart">
                    <TitleBlock meta="USAGE" title="Quick Start" />
                    <CodeTabs />
                </ContentCard>

                {/* Architecture */}
                <ContentCard id="architecture">
                    <TitleBlock meta="SYSTEM DESIGN" title="Architecture" />
                    <ArchitectureDiagram />
                </ContentCard>

                {/* API Reference */}
                <ContentCard id="api">
                    <TitleBlock meta="CORE ENDPOINTS" title="API Reference" />
                    <ApiTable />
                </ContentCard>

                {/* Security Model */}
                <ContentCard id="security">
                    <TitleBlock meta="HARDENING" title="Security Model" />
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                        <div>
                            <h3 className="text-sm font-bold border-b border-gray-300 mb-2 pb-1">
                                SINGLE JWT ISSUER
                            </h3>
                            <p className="font-mono text-xs text-gray-600 mb-4">
                                Authflow-Go is the sole JWT authority. Downstream services never
                                mint tokens. All token lifecycle operations are centralized.
                            </p>
                            <h3 className="text-sm font-bold border-b border-gray-300 mb-2 pb-1">
                                REDIS REVOCATION MODEL
                            </h3>
                            <p className="font-mono text-xs text-gray-600">
                                {"revoked:<jti> → \"1\" (TTL = token_expiry − now)."}{" "}
                                Any request using a revoked token is rejected by middleware,
                                even if the JWT is otherwise valid. Automatic cleanup via TTL.
                            </p>
                        </div>
                        <div>
                            <h3 className="text-sm font-bold border-b border-gray-300 mb-2 pb-1">
                                STATELESS AUTHORIZATION
                            </h3>
                            <p className="font-mono text-xs text-gray-600 mb-4">
                                JWTs enable stateless authorization. Middleware verifies
                                signature & expiry, then checks Redis for revoked jti. No
                                database hits for standard validation.
                            </p>
                            <h3 className="text-sm font-bold border-b border-gray-300 mb-2 pb-1">
                                DECOUPLED AUTH
                            </h3>
                            <p className="font-mono text-xs text-gray-600">
                                Authentication strategies vary per platform. Authorization is
                                centralized. Token lifecycle is controlled in one place. Logout
                                is enforceable across services.
                            </p>
                        </div>
                    </div>
                </ContentCard>

                {/* Live Integration */}
                <ContentCard id="live">
                    <TitleBlock meta="PRODUCTION" title="Live Integration" />
                    <p className="font-mono text-sm text-gray-600">
                        AuthFlow is currently running in production as the authentication layer for:
                    </p>
                    <a
                        href="https://rag-works.vercel.app"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="btn-primary"
                        style={{ marginBottom: "16px" }}
                    >
                        rag-works.vercel.app →
                    </a>
                    <ul className="methodology-list mt-4">
                        <li>Create an account</li>
                        <li>Log in via OAuth</li>
                        <li>Access protected routes</li>
                        <li>Log out (token revoked via Redis)</li>
                    </ul>
                    <p className="font-mono text-sm text-gray-600">
                        All auth operations in rag-works are powered by AuthFlow.
                    </p>
                    
                </ContentCard>

                {/* Footer */}
                <ContentCard style={{ marginBottom: 0 }}>
                    <div className="flex flex-col md:flex-row justify-between items-center font-mono text-xs text-gray-500">
                        <div>
                            © 2026 AUTHFLOW OPEN SOURCE
                           
                            
                        </div>
                        <div className="flex gap-4 mt-4 md:mt-0">
                            <a
                                href="https://github.com/abhiram-karanth-core/auth-flow"
                                target="_blank"
                                rel="noopener noreferrer"
                                className="hover:text-black hover:underline"
                            >
                                GITHUB
                            </a>
                            
                        </div>
                        <div className="mt-4 md:mt-0">
                            MAINTAINED BY @ABHIRAM-KARANTH
                        </div>
                    </div>
                </ContentCard>
            </div>
        </>
    );
}
