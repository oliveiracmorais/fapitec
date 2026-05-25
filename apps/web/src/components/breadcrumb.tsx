"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";

type Rota = "/" | "/dashboard" | "/editais" | `/editais/${string}` | `/editais/novo`;

const rotas: Record<string, string> = {
  dashboard: "Início",
  editais: "Editais",
  novo: "Novo Edital",
};

export default function Breadcrumb() {
  const pathname = usePathname();
  const segmentos = pathname.split("/").filter(Boolean);

  if (segmentos.length === 0) return null;

  const itens = segmentos.map((seg, i) => {
    const caminho = "/" + segmentos.slice(0, i + 1).join("/");
    const nome = rotas[seg] || seg;
    const isUltimo = i === segmentos.length - 1;

    if (isUltimo) {
      return (
        <span key={caminho} className="text-sm font-medium text-gray-900">
          {nome}
        </span>
      );
    }

    return (
      <span key={caminho} className="flex items-center gap-1.5">
        <Link
          href={caminho as Rota}
          className="text-sm text-gray-500 transition-colors hover:text-brand-600"
        >
          {nome}
        </Link>
        <svg
          className="h-3.5 w-3.5 text-gray-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 5l7 7-7 7"
          />
        </svg>
      </span>
    );
  });

  return (
    <nav aria-label="Breadcrumb" className="mx-auto max-w-6xl px-4 pt-4">
      <div className="flex items-center gap-1.5 text-sm">
        <Link
          href="/dashboard"
          className="text-gray-500 transition-colors hover:text-brand-600"
          aria-label="Início"
        >
          <svg
            className="h-4 w-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
            />
          </svg>
        </Link>
        {segmentos.length > 0 && (
          <svg
            className="h-3.5 w-3.5 text-gray-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9 5l7 7-7 7"
            />
          </svg>
        )}
        {itens}
      </div>
    </nav>
  );
}
