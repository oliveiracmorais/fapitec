import type { Metadata } from "next";
import "./globals.css";
import { AuthWrapper } from "../context/auth-wrapper";

export const metadata: Metadata = {
  title: "FAPITEC-SE",
  description:
    "Plataforma integrada de gestão institucional da Fundação de Apoio à Pesquisa e à Inovação Tecnológica do Estado de Sergipe",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="pt-BR">
      <head>
        <link rel="icon" type="image/jpeg" href="/Logo-1.jpeg" />
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link
          href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap"
          rel="stylesheet"
        />
      </head>
      <body>
        <AuthWrapper>{children}</AuthWrapper>
      </body>
    </html>
  );
}
