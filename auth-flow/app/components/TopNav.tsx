import Link from "next/link";

const navLinks = [
    { label: "Docs [D]", href: "#" },
    { label: "Architecture [A]", href: "#architecture" },
    { label: "Quick Start [Q]", href: "#quickstart" },
    { label: "API Reference [R]", href: "#api" },
    { label: "Security [S]", href: "#security" },
];

export default function TopNav() {
    return (
        <nav className="top-nav">
            <div
                style={{
                    fontWeight: 900,
                    fontSize: "18px",
                    marginRight: "auto",
                    letterSpacing: "-1px",
                }}
            >
                AUTHFLOW<span style={{ color: "#666" }}>_go</span>
            </div>
            <div className="hidden md:flex">
                {navLinks.map((link) => (
                    <Link key={link.label} href={link.href} className="nav-link">
                        {link.label}
                    </Link>
                ))}
                <a
                    href="https://github.com/abhiram-karanth-core/auth-flow"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="nav-link"
                >
                    GitHub [G]
                </a>
            </div>
            <Link href="#quickstart" className="btn-primary" style={{ marginLeft: "30px" }}>
                Get Started
            </Link>
        </nav>
    );
}
