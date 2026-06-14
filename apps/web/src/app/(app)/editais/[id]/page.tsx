"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import Link from "next/link";

type Edital = {
  id: number;
  nome: string;
  descricao: string;
  data_inicio: string;
  data_fim: string;
  status: string;
  tipo_chamada: string;
  nota_de_corte: number;
  valor_global: number;
  modelo_formulario: number;
  relatorios_exigidos: string[];
  titulo_minimo_elegibilidade: string;
  exige_empresa: boolean;
  porte_empresa: string[];
  enquadramento_empresa: string[];
  documentos_obrigatorios: string[];
  criado_em: string;
};

const statusLabel: Record<string, string> = {
  ativo: "Ativo",
  encerrado: "Encerrado",
  em_avaliacao: "Em Avaliação",
};

const statusColor: Record<string, string> = {
  ativo: "bg-green-100 text-green-800",
  encerrado: "bg-gray-100 text-gray-600",
  em_avaliacao: "bg-yellow-100 text-yellow-800",
};

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="border-b border-gray-100 pb-6 last:border-b-0">
      <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-gray-500">
        {title}
      </h2>
      <div className="grid grid-cols-2 gap-4 text-sm">
        {children}
      </div>
    </div>
  );
}

function Field({ label, value }: { label: string; value: string | number | boolean | string[] | null | undefined }) {
  let display: string;
  if (Array.isArray(value)) {
    display = value.length > 0 ? value.join(", ") : "—";
  } else if (typeof value === "boolean") {
    display = value ? "Sim" : "Não";
  } else if (value === null || value === undefined || value === "") {
    display = "—";
  } else {
    display = String(value);
  }

  return (
    <div>
      <span className="font-medium text-gray-900">{label}</span>
      <p className="mt-0.5 text-gray-600">{display}</p>
    </div>
  );
}

export default function EditalDetalhePage() {
  const router = useRouter();
  const params = useParams();
  const [edital, setEdital] = useState<Edital | null>(null);
  const [carregando, setCarregando] = useState(true);

  useEffect(() => {
    if (!params?.id) return;

    fetch(`/api/v1/editais/${params.id}`)
      .then((res) => {
        if (!res.ok) throw new Error("Não encontrado");
        return res.json();
      })
      .then((data) => setEdital(data))
      .catch(() => router.push("/editais"))
      .finally(() => setCarregando(false));
  }, [params, router]);

  if (carregando) {
    return (
      <div className="flex items-center justify-center py-20 text-gray-600">
        Carregando...
      </div>
    );
  }

  if (!edital) return null;

  return (
    <main className="mx-auto max-w-3xl px-4 py-6">
      <div className="rounded-xl border border-gray-200 bg-white p-8 shadow-sm">
        <div className="mb-6 flex items-start justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              {edital.nome}
            </h1>
            <p className="mt-1 text-sm text-gray-600">{edital.descricao}</p>
          </div>
          <span
            className={`rounded-full px-3 py-1 text-xs font-medium ${statusColor[edital.status] || "bg-gray-100 text-gray-600"}`}
          >
            {statusLabel[edital.status] || edital.status}
          </span>
        </div>

        <div className="space-y-6">
          <Section title="Informações Gerais">
            <Field label="Tipo de Chamada" value={edital.tipo_chamada} />
            <Field label="Vigência" value={`${new Date(edital.data_inicio).toLocaleDateString("pt-BR")} a ${new Date(edital.data_fim).toLocaleDateString("pt-BR")}`} />
            <Field label="Modelo de Formulário" value={edital.modelo_formulario || "Nenhum"} />
            <Field label="Nota de Corte" value={edital.nota_de_corte} />
            <Field label="Valor Global" value={edital.valor_global ? `R$ ${(edital.valor_global / 100).toFixed(2)}` : null} />
            <Field label="Criado em" value={new Date(edital.criado_em).toLocaleDateString("pt-BR")} />
          </Section>

          <Section title="Proponente">
            <Field label="Titulação Mínima" value={edital.titulo_minimo_elegibilidade} />
          </Section>

          <Section title="Empresa">
            <Field label="Exige Empresa Vinculada" value={edital.exige_empresa} />
            <Field label="Porte Permitido" value={edital.porte_empresa} />
            <Field label="Enquadramento" value={edital.enquadramento_empresa} />
          </Section>

          <Section title="Documentos Obrigatórios">
            <Field label="Documentos" value={edital.documentos_obrigatorios} />
          </Section>
        </div>

        <div className="mt-6 flex flex-wrap items-center gap-4 border-t border-gray-100 pt-6">
          <Link
            href="/editais"
            className="text-sm font-medium text-brand-600 hover:text-brand-700"
          >
            ← Voltar para lista
          </Link>
          {edital.status === "ativo" && (
            <Link
              href={{ pathname: `/editais/${edital.id}/inscrever` }}
              className="rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-green-700"
            >
              Inscrever-se neste edital
            </Link>
          )}
        </div>
      </div>
    </main>
  );
}
