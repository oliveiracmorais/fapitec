import "./globals.css";

export const metadata = {
  title: "FAPITEC-SE",
  description: "Plataforma modular da FAPITEC-SE"
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="pt-BR">
      <body>{children}</body>
    </html>
  );
}
