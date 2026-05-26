"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import GraficosSection from "../../../components/graficos-section";

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

type Indicador = {
  id: string;
  nome: string;
  valor: number;
  tipo: string;
  cor: string;
};

type Modulo = {
  id: string;
  nome: string;
  descricao: string;
  status: string;
  icone: string;
};

const modulos: Modulo[] = [
  {
    id: "gestao-editais",
    nome: "Gestão de Editais",
    descricao: "Gerenciamento de editais de pesquisa e inovação",
    status: "Em desenvolvimento",
    icone: "📋",
  },
  {
    id: "inscricao-selecao",
    nome: "Inscrição e Seleção de Propostas",
    descricao: "Inscrição em chamadas e seleção de propostas",
    status: "Em planejamento",
    icone: "📝",
  },
  {
    id: "contratacao-concessao",
    nome: "Contratação e Concessão de Apoios",
    descricao: "Contratação e concessão de apoios institucionais",
    status: "Em planejamento",
    icone: "📑",
  },
  {
    id: "financeiro-pagamentos",
    nome: "Financeiro e Pagamentos",
    descricao: "Gestão financeira e processamento de pagamentos",
    status: "Em planejamento",
    icone: "💰",
  },
  {
    id: "prestacao-contas",
    nome: "Prestação de Contas",
    descricao: "Prestação de contas financeiras dos projetos",
    status: "Em planejamento",
    icone: "🧾",
  },
  {
    id: "tomada-contas-especial",
    nome: "Tomada de Contas Especial",
    descricao: "Processos de tomada de contas especial",
    status: "Em planejamento",
    icone: "🔍",
  },
  {
    id: "comunicacao-institucional",
    nome: "Comunicação Institucional",
    descricao: "Comunicação e divulgação institucional",
    status: "Em planejamento",
    icone: "📢",
  },
  {
    id: "relatorios-indicadores",
    nome: "Relatórios e Painéis de Indicadores",
    descricao: "Relatórios gerenciais e painéis de indicadores",
    status: "Em planejamento",
    icone: "📊",
  },
  {
    id: "gestao-documental",
    nome: "Gestão Documental",
    descricao: "Gestão de documentos e arquivos digitais",
    status: "Em planejamento",
    icone: "📁",
  },
  {
    id: "infraestrutura-suporte",
    nome: "Infraestrutura, Suporte e Segurança",
    descricao: "Infraestrutura, suporte técnico e segurança da informação",
    status: "Em planejamento",
    icone: "🛡️",
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

function corProgresso(cor: string): string {
  const mapa: Record<string, string> = {
    violeta: "bg-purple-500",
    verde: "bg-green-500",
    azul: "bg-blue-500",
    laranja: "bg-orange-500",
  };
  return mapa[cor] || "bg-brand-500";
}

function corTextoProgresso(cor: string): string {
  const mapa: Record<string, string> = {
    violeta: "text-purple-700",
    verde: "text-green-700",
    azul: "text-blue-700",
    laranja: "text-orange-700",
  };
  return mapa[cor] || "text-brand-700";
}

function corFundoProgresso(cor: string): string {
  const mapa: Record<string, string> = {
    violeta: "bg-purple-100",
    verde: "bg-green-100",
    azul: "bg-blue-100",
    laranja: "bg-orange-100",
  };
  return mapa[cor] || "bg-brand-100";
}

function formatarValor(valor: number, tipo: string): string {
  if (tipo === "numero") {
    return valor.toLocaleString("pt-BR");
  }
  return `${valor}%`;
}

export default function DashboardPage() {
  const router = useRouter();
  const [sessao, setSessao] = useState<Sessao | null>(null);
  const [indicadores, setIndicadores] = useState<Indicador[]>([]);
  const [carregandoIndicadores, setCarregandoIndicadores] = useState(true);

  useEffect(() => {
    const s = obterSessao();
    if (!s) {
      router.replace("/");
      return;
    }
    setSessao(s);
  }, [router]);

  useEffect(() => {
    fetch("/api/v1/dashboard/indicadores")
      .then((res) => res.json())
      .then((data) => setIndicadores(data.indicadores || []))
      .catch(() => setIndicadores([]))
      .finally(() => setCarregandoIndicadores(false));
  }, []);

  if (!sessao) return null;

  const { usuario } = sessao;

  return (
    <main className="mx-auto max-w-6xl px-4 py-8">
      <div className="mb-8 rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <div className="flex items-center gap-4">
          <div className="flex h-14 w-14 items-center justify-center rounded-full bg-brand-100 text-xl font-bold text-brand-700">
            {usuario.nome.charAt(0).toUpperCase()}
          </div>
          <div className="flex-1">
            <h1 className="text-xl font-bold text-gray-900">
              Bem-vindo, {usuario.nome}
            </h1>
            <p className="text-sm text-gray-600">
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

      {!carregandoIndicadores && indicadores.length > 0 && (
        <section aria-label="Indicadores do sistema" className="mb-8">
          <h2 className="mb-4 text-lg font-bold text-gray-900">
            Indicadores
          </h2>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {indicadores.map((ind) => (
              <div
                key={ind.id}
                className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition-all hover:shadow-md"
                title={ind.nome}
              >
                <p className="text-sm font-medium text-gray-500">
                  {ind.nome}
                </p>
                {ind.tipo === "numero" ? (
                  <p className="mt-2 text-3xl font-bold text-gray-900">
                    {formatarValor(ind.valor, ind.tipo)}
                  </p>
                ) : (
                  <>
                    <p
                      className={`mt-2 text-3xl font-bold ${corTextoProgresso(ind.cor)}`}
                    >
                      {formatarValor(ind.valor, ind.tipo)}
                    </p>
                    <div
                      className={`mt-3 h-3 w-full rounded-full ${corFundoProgresso(ind.cor)}`}
                      role="progressbar"
                      aria-valuenow={ind.valor}
                      aria-valuemin={0}
                      aria-valuemax={100}
                      aria-label={`${ind.nome}: ${ind.valor}%`}
                    >
                      <div
                        className={`h-3 rounded-full ${corProgresso(ind.cor)} transition-all duration-500`}
                        style={{ width: `${Math.min(ind.valor, 100)}%` }}
                      />
                    </div>
                  </>
                )}
                <div className="mt-3 flex items-center gap-2">
                  <span
                    className={`inline-block h-2 w-2 rounded-full ${corProgresso(ind.cor)}`}
                  />
                  <span className="text-xs text-gray-400">
                    Atualizado recentemente
                  </span>
                </div>
              </div>
            ))}
          </div>
        </section>
      )}

      <GraficosSection />

      <h2 className="mb-4 text-lg font-bold text-gray-900">Módulos</h2>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {modulos.map((modulo) => (
          <div
            key={modulo.id}
            className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition-all hover:-translate-y-0.5 hover:shadow-md"
            title={modulo.descricao}
          >
            <div className="mb-3 text-2xl" aria-hidden="true">
              {modulo.icone}
            </div>
            <h3 className="font-bold text-gray-900">{modulo.nome}</h3>
            <p className="mt-1 text-sm text-gray-600">
              {modulo.descricao}
            </p>
            <div className="mt-4">
              <span className="rounded-full bg-yellow-100 px-2.5 py-0.5 text-xs font-medium text-yellow-800">
                {modulo.status}
              </span>
            </div>
          </div>
        ))}
      </div>
    </main>
  );
}
