"use client";

import Image from "next/image";
import Link from "next/link";
import { useAuth } from "../context/auth-context";
import { useState, useMemo } from "react";

const navegacao: { href: string; label: string }[] = [
  { href: "/dashboard", label: "Dashboard" },
  { href: "/editais", label: "Editais" },
  { href: "/avaliadores", label: "Avaliadores" },
  { href: "/avaliacoes", label: "Avaliações" },
  { href: "/admin/avaliacoes", label: "Acompanhamento" },
  { href: "/minhas-propostas", label: "Minhas Propostas" },
];

function iniciais(nome: string): string {
  return nome
    .split(" ")
    .filter((p) => p.length > 2)
    .slice(0, 2)
    .map((p) => p[0].toUpperCase())
    .join("");
}

export default function Header() {
  const { usuario, logout } = useAuth();
  const [menuAberto, setMenuAberto] = useState(false);

  const avatar = useMemo(
    () => (usuario?.nome ? iniciais(usuario.nome) : "U"),
    [usuario?.nome],
  );

  return (
    <header className="border-b border-gray-200 bg-white shadow-sm">
      {/* Linha superior: marca + avatar/sair (sempre visível) */}
      <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-1.5">
        <div className="flex shrink-0 items-center gap-3">
          <Link href="/dashboard" aria-label="Voltar ao início">
            <Image
              src="/logo-2.png"
              alt="FAPITEC-SE"
              width={140}
              height={40}
              className="h-9 w-auto"
            />
          </Link>
          <span className="hidden text-sm font-semibold text-brand-700 lg:inline xl:text-base">
            Plataforma Integrada de Gestão Institucional
          </span>
        </div>

        <div className="hidden items-center gap-3 sm:flex">
          {usuario && (
            <>
              <div
                className="flex h-8 w-8 items-center justify-center rounded-full bg-brand-600 text-xs font-bold text-white"
                title={usuario.nome}
              >
                {avatar}
              </div>
              <button
                onClick={logout}
                className="whitespace-nowrap rounded-lg bg-gray-100 px-3 py-1.5 text-sm text-gray-700 transition-colors hover:bg-gray-200"
                aria-label="Sair da plataforma"
              >
                Sair
              </button>
            </>
          )}
        </div>

        <button
          className="flex items-center gap-2 rounded-lg p-2 text-gray-600 hover:bg-gray-100 lg:hidden"
          onClick={() => setMenuAberto(!menuAberto)}
          aria-label={menuAberto ? "Fechar menu" : "Abrir menu"}
          aria-expanded={menuAberto}
        >
          <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d={
                menuAberto
                  ? "M6 18L18 6M6 6l12 12"
                  : "M4 6h16M4 12h16M4 18h16"
              }
            />
          </svg>
          {usuario && (
            <span className="max-w-[100px] truncate text-sm text-gray-700 sm:hidden">
              {usuario.nome}
            </span>
          )}
        </button>
      </div>

      {/* Linha inferior: navegação (lg+) / menu mobile */}
      <div className="hidden border-t border-gray-100 bg-gray-50/50 lg:block">
        <nav className="mx-auto flex max-w-6xl items-center gap-0.5 px-4 py-1">
          {navegacao.map((item) => (
            <Link
              key={item.href}
              href={{ pathname: item.href }}
              className="whitespace-nowrap rounded-md px-2.5 py-1.5 text-sm text-gray-600 transition-colors hover:bg-brand-50 hover:text-brand-700"
            >
              {item.label}
            </Link>
          ))}
        </nav>
      </div>

      {menuAberto && (
        <div className="border-t border-gray-200 bg-white lg:hidden">
          <nav className="space-y-1 px-4 py-3">
            {navegacao.map((item) => (
              <Link
                key={item.href}
                href={{ pathname: item.href }}
                onClick={() => setMenuAberto(false)}
                className="block rounded-md px-3 py-2 text-sm text-gray-600 transition-colors hover:bg-brand-50 hover:text-brand-700"
              >
                {item.label}
              </Link>
            ))}
            {usuario && (
              <div className="border-t border-gray-100 pt-2">
                <div className="px-3 py-2 text-sm text-gray-500">
                  {usuario.nome}
                </div>
                <button
                  onClick={() => {
                    logout();
                    setMenuAberto(false);
                  }}
                  className="flex w-full items-center rounded-md px-3 py-2 text-sm text-red-600 transition-colors hover:bg-red-50"
                >
                  Sair
                </button>
              </div>
            )}
          </nav>
        </div>
      )}
    </header>
  );
}
