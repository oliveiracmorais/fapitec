"use client";

import Image from "next/image";
import Link from "next/link";
import { useAuth } from "../context/auth-context";

export default function Header() {
  const { usuario, logout } = useAuth();

  return (
    <header className="border-b border-gray-200 bg-white shadow-sm">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
        <div className="flex items-center gap-3">
          <Link href="/dashboard" aria-label="Voltar ao início">
            <Image
              src="/logo-2.png"
              alt="FAPITEC-SE"
              width={140}
              height={40}
              className="h-9 w-auto"
            />
          </Link>
          <span className="hidden text-sm text-gray-600 sm:inline">
            Plataforma Integrada de Gestão Institucional
          </span>
        </div>
        <div className="flex items-center gap-3 text-sm">
          {usuario && (
            <>
              <span className="hidden text-gray-700 sm:inline">
                {usuario.nome}
              </span>
              <button
                onClick={logout}
                className="rounded-lg bg-gray-100 px-3 py-1.5 text-gray-700 transition-colors hover:bg-gray-200"
                aria-label="Sair da plataforma"
              >
                Sair
              </button>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
