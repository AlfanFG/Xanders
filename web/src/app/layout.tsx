import type { Metadata } from "next";
import "./globals.css";
import { AppProvider } from "./context/AppContext";
import Providers from "../components/Providers";
import Navbar from "./Navbar";

export const metadata: Metadata = {
  title: "Xanders | Cinematic Video Generator",
  description: "B2B AI Video Generation SaaS for High-End Industries",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body suppressHydrationWarning>
        <Providers>
          <AppProvider>
            <Navbar />
            <main>{children}</main>
          </AppProvider>
        </Providers>
      </body>
    </html>
  );
}
