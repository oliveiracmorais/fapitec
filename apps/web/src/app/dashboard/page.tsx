"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Image from "next/image";

type UsuarioSessao = {
  id: number;
  nome: string;
  documento: string;
  email: string;
  estrangeiro: boolean;
};

type Sessao = {
  usuario: UsuarioSessao;
  timestamp: number;
};

type Modulo = {
  id: string;
  nome: string;
  descricao: string;
  status: string;
  icone: string;
};

const modulosDisponiveis: Modulo[] = [
  {
    id: "editais",
    nome: "Editais",
    descricao: "Gerenciamento de editais de pesquisa e inovação",
    status: "Em desenvolvimento",
    icone: "📋",
  },
  {
    id: "inscricao",
    nome: "Inscrição",
    descricao: "Inscrição em projetos e chamadas",
    status: "Em planejamento",
    icone: "📝",
  },
  {
    id: "bolsistas",
    nome: "Bolsistas",
    descricao: "Gestão de bolsistas e pesquisadores",
    status: "Em planejamento",
    icone: "👥",
  },
  {
    id: "prestacao-contas",
    nome: "Prestação de Contas",
    descricao: "Prestação de contas financeiras",
    status: "Em planejamento",
    icone: "💰",
  },
];

function obterSessao(): Sessao | null {
  if (typeof window === "undefined") return null;
  const raw = localStorage.getItem("sessao");
  if (!raw) return null;
  try {
    return JSON.parse(raw) as Sessao;
  } catch {
    return null;
  }
}

function limparSessao() {
  localStorage.removeItem("sessao");
}

export default function DashboardPage() {
  const router = useRouter();
  const [sessao, setSessao] = useState<Sessao | null>(null);

  useEffect(() => {
    const s = obterSessao();
    if (!s) {
      router.replace("/");
      return;
    }
    setSessao(s);
  }, [router]);

  function handleLogout() {
    limparSessao();
    router.push("/");
  }

  if (!sessao) return null;

  const { usuario } = sessao;

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="border-b border-gray-200 bg-white shadow-sm">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
          <div className="flex items-center gap-3">
            <Image
              src="/logo-2.png"
              alt="FAPITEC-SE"
              width={140}
              height={40}
              className="h-9 w-auto"
            />
            <span className="hidden text-sm text-gray-400 sm:inline">
              Plataforma integrada de gestão institucional
            </span>
          </div>
          <div className="flex items-center gap-3 text-sm">
            <span className="text-gray-700">{usuario.nome}</span>
            <button
              onClick={handleLogout}
              className="rounded-lg bg-gray-100 px-3 py-1.5 text-gray-700 transition-colors hover:bg-gray-200"
            >
              Sair
            </button>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-4 py-8">
        <div className="mb-8 rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <div className="flex items-center gap-4">
            <div className="flex h-14 w-14 items-center justify-center rounded-full bg-brand-100 text-xl font-bold text-brand-700">
              {usuario.nome.charAt(0).toUpperCase()}
            </div>
            <div>
              <h1 className="text-xl font-bold text-gray-900">
                Bem-vindo, {usuario.nome}
              </h1>
              <p className="text-sm text-gray-500">
                Gerencie seus projetos e editais na plataforma
              </p>
            </div>
          </div>
          <div className="mt-4 grid grid-cols-1 gap-3 border-t border-gray-100 pt-4 text-sm text-gray-600 sm:grid-cols-3">
            <div>
              <span className="font-medium text-gray-900">Documento:</span>{" "}
              {usuario.documento}
            </div>
            <div>
              <span className="font-medium text-gray-900">E-mail:</span>{" "}
              {usuario.email}
            </div>
            <div>
              <span className="font-medium text-gray-900">Tipo:</span>{" "}
              {usuario.estrangeiro ? "Estrangeiro" : "Brasileiro"}
            </div>
          </div>
        </div>

        <h2 className="mb-4 text-lg font-bold text-gray-900">Módulos</h2>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {modulosDisponiveis.map((modulo) => (
            <div
              key={modulo.id}
              className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition-all hover:-translate-y-0.5 hover:shadow-md"
            >
              <div className="mb-3 text-2xl">{modulo.icone}</div>
              <h3 className="font-bold text-gray-900">{modulo.nome}</h3>
              <p className="mt-1 text-sm text-gray-500">{modulo.descricao}</p>
              <div className="mt-4">
                <span className="rounded-full bg-yellow-100 px-2.5 py-0.5 text-xs font-medium text-yellow-800">
                  {modulo.status}
                </span>
              </div>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}
