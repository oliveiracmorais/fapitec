"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import {
  ResponsiveContainer,
  LineChart,
  Line,
  CartesianGrid,
  Tooltip as RechartsTooltip,
} from "recharts";
import GraficosSection from "../../../components/graficos-section";
import { useAuth } from "../../../context/auth-context";

type DadoLinha = {
  mes: string;
  valor: number;
};

type Indicador = {
  id: string;
  nome: string;
  valor: number;
  tipo: string;
  cor: string;
  dados?: DadoLinha[];
};

type Modulo = {
  id: string;
  nome: string;
  descricao: string;
  status: string;
  icone: string;
  rota: string;
};

const CORES: Record<string, string> = {
  violeta: "#7C3AED",
  verde: "#10B981",
  azul: "#3B82F6",
  laranja: "#F97316",
};

const modulos: Modulo[] = [
  {
    id: "gestao-editais",
    nome: "Gestão de Editais",
    descricao: "Gerenciamento de editais de pesquisa e inovação",
    status: "Em desenvolvimento",
    icone: "📋",
    rota: "/editais",
  },
  {
    id: "inscricao-selecao",
    nome: "Inscrição e Seleção de Propostas",
    descricao: "Inscrição em chamadas e seleção de propostas",
    status: "Em planejamento",
    icone: "📝",
    rota: "/inscricao",
  },
  {
    id: "contratacao-concessao",
    nome: "Contratação e Concessão de Apoios",
    descricao: "Contratação e concessão de apoios institucionais",
    status: "Em planejamento",
    icone: "📑",
    rota: "/contratacao",
  },
  {
    id: "financeiro-pagamentos",
    nome: "Financeiro e Pagamentos",
    descricao: "Gestão financeira e processamento de pagamentos",
    status: "Em planejamento",
    icone: "💰",
    rota: "/financeiro",
  },
  {
    id: "prestacao-contas",
    nome: "Prestação de Contas",
    descricao: "Prestação de contas financeiras dos projetos",
    status: "Em planejamento",
    icone: "🧾",
    rota: "/prestacao-contas",
  },
  {
    id: "tomada-contas-especial",
    nome: "Tomada de Contas Especial",
    descricao: "Processos de tomada de contas especial",
    status: "Em planejamento",
    icone: "🔍",
    rota: "/tomada-contas",
  },
  {
    id: "comunicacao-institucional",
    nome: "Comunicação Institucional",
    descricao: "Comunicação e divulgação institucional",
    status: "Em planejamento",
    icone: "📢",
    rota: "/comunicacao",
  },
  {
    id: "relatorios-indicadores",
    nome: "Relatórios e Painéis de Indicadores",
    descricao: "Relatórios gerenciais e painéis de indicadores",
    status: "Em planejamento",
    icone: "📊",
    rota: "/relatorios",
  },
  {
    id: "gestao-documental",
    nome: "Gestão Documental",
    descricao: "Gestão de documentos e arquivos digitais",
    status: "Em planejamento",
    icone: "📁",
    rota: "/gestao-documental",
  },
  {
    id: "infraestrutura-suporte",
    nome: "Infraestrutura, Suporte e Segurança",
    descricao: "Infraestrutura, suporte técnico e segurança da informação",
    status: "Em planejamento",
    icone: "🛡️",
    rota: "/infraestrutura",
  },
];

function MiniPieChart({ valor, cor }: { valor: number; cor: string }) {
  const radius = 32;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (valor / 100) * circumference;
  const strokeColor = CORES[cor] || "#7C3AED";

  return (
    <div className="flex justify-center" aria-label={`${valor}%`}>
      <svg width={80} height={80} viewBox="0 0 80 80" role="img">
        <circle
          cx="40" cy="40" r={radius}
          fill="none" stroke="#E5E7EB" strokeWidth={8}
        />
        <circle
          cx="40" cy="40" r={radius}
          fill="none" stroke={strokeColor} strokeWidth={8}
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          transform="rotate(-90 40 40)"
          strokeLinecap="round"
        />
        <text
          x="40" y="40"
          textAnchor="middle" dominantBaseline="central"
          fontSize="14" fontWeight="bold"
          fill={strokeColor}
        >
          {valor}%
        </text>
      </svg>
    </div>
  );
}

function MiniLineChart({ dados, cor }: { dados: DadoLinha[]; cor: string }) {
  return (
    <div className="mt-2" aria-label="Gráfico de linha">
      <ResponsiveContainer width="100%" height={60}>
        <LineChart data={dados}>
          <CartesianGrid strokeDasharray="3 3" stroke="#E5E7EB" />
          <RechartsTooltip
            contentStyle={{ fontSize: 12, padding: "4px 8px" }}
          />
          <Line
            type="monotone"
            dataKey="valor"
            stroke={CORES[cor] || "#3B82F6"}
            strokeWidth={2}
            dot={false}
            activeDot={{ r: 4 }}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
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

function corFundoProgresso(cor: string): string {
  const mapa: Record<string, string> = {
    violeta: "bg-purple-100",
    verde: "bg-green-100",
    azul: "bg-blue-100",
    laranja: "bg-orange-100",
  };
  return mapa[cor] || "bg-brand-100";
}

function formatarValor(valor: number): string {
  return valor.toLocaleString("pt-BR");
}

export default function DashboardPage() {
  const router = useRouter();
  const { usuario, carregando } = useAuth();
  const [indicadores, setIndicadores] = useState<Indicador[]>([]);
  const [carregandoIndicadores, setCarregandoIndicadores] = useState(true);

  useEffect(() => {
    fetch("/api/v1/dashboard/indicadores")
      .then((res) => res.json())
      .then((data) => setIndicadores(data.indicadores || []))
      .catch(() => setIndicadores([]))
      .finally(() => setCarregandoIndicadores(false));
  }, []);

  const pushUrl = (url: string) => (router.push as (u: string) => void)(url);

  if (carregando) return null;
  if (!usuario) return null;

  return (
    <main className="mx-auto max-w-6xl px-4 py-8">
      <div
        className="mb-8 rounded-xl border border-gray-200 bg-white p-6 shadow-sm"
        title="Informações do usuário"
      >
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
          <button
            type="button"
            onClick={() => pushUrl("/perfil")}
            className="rounded-lg border border-brand-300 bg-brand-50 px-4 py-2 text-sm font-medium text-brand-700 transition-colors hover:bg-brand-100"
            aria-label="Atualizar Perfil"
          >
            Atualizar Perfil
          </button>
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
                role="region"
                aria-label={`Indicador: ${ind.nome}`}
              >
                <p className="text-sm font-medium text-gray-500">
                  {ind.nome}
                </p>
                {ind.tipo === "pizza" ? (
                  <>
                    <MiniPieChart valor={ind.valor} cor={ind.cor} />
                    <div className="mt-3 flex items-center gap-2">
                      <span
                        className="inline-block h-2 w-2 rounded-full"
                        style={{ backgroundColor: CORES[ind.cor] || "#7C3AED" }}
                      />
                      <span className="text-xs text-gray-400">
                        Atualizado recentemente
                      </span>
                    </div>
                  </>
                ) : ind.tipo === "linha" ? (
                  <>
                    <p className="mt-2 text-3xl font-bold text-gray-900">
                      {formatarValor(ind.valor)}
                    </p>
                    {ind.dados && ind.dados.length > 0 && (
                      <MiniLineChart dados={ind.dados} cor={ind.cor} />
                    )}
                    <div className="mt-1 flex items-center gap-2">
                      <span
                        className="inline-block h-2 w-2 rounded-full"
                        style={{ backgroundColor: CORES[ind.cor] || "#3B82F6" }}
                      />
                      <span className="text-xs text-gray-400">
                        Atualizado recentemente
                      </span>
                    </div>
                  </>
                ) : ind.tipo === "numero" ? (
                  <>
                    <p className="mt-2 text-3xl font-bold text-gray-900">
                      {formatarValor(ind.valor)}
                    </p>
                    <div className="mt-3 flex items-center gap-2">
                      <span
                        className={`inline-block h-2 w-2 rounded-full ${corProgresso(ind.cor)}`}
                      />
                      <span className="text-xs text-gray-400">
                        Atualizado recentemente
                      </span>
                    </div>
                  </>
                ) : (
                  <>
                    <p className="mt-2 text-3xl font-bold text-gray-900">
                      {ind.valor}%
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
                    <div className="mt-3 flex items-center gap-2">
                      <span
                        className={`inline-block h-2 w-2 rounded-full ${corProgresso(ind.cor)}`}
                      />
                      <span className="text-xs text-gray-400">
                        Atualizado recentemente
                      </span>
                    </div>
                  </>
                )}
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
            role="button"
            tabIndex={0}
            onClick={() => pushUrl(modulo.rota)}
            onKeyDown={(e) => {
              if (e.key === "Enter" || e.key === " ") {
                e.preventDefault();
                pushUrl(modulo.rota);
              }
            }}
            className="cursor-pointer rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition-all hover:-translate-y-0.5 hover:shadow-md"
            title={modulo.descricao}
            aria-label={`Módulo: ${modulo.nome} — ${modulo.descricao}`}
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
