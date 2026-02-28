import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "AUTHFLOW | DOCUMENTATION",
  description: "Centralized authorization for modern platforms.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
