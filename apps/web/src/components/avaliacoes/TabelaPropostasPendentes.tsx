"use client";

import Link from "next/link";
import type { PropostaParaAvaliarSaida } from "../../lib/api-avaliacoes";
import {
  formatarStatusProposta,
  corStatusProposta,
  formatarDataISO,
  formatarMoeda,
} from "../../lib/api-avaliacoes";

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="border-b border-gray-100 pb-6 last:border-b-0">
      <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-gray-500">{title}</h2>
      <div className="grid grid-cols-2 gap-4 text-sm">{children}</div>
    </div>
  );
}

function Field({ label, value }: { label: string; value: string | number | null | undefined }) {
  const display = value === null || value === undefined || value === "" ? "—" : String(value);
  return (
    <div>
      <span className="font-medium text-gray-900">{label}</span>
      <p className="mt-0.5 text-gray-600">{display}</p>
    </div>
  );
}

export default function TabelaPropostasPendentes({
  propostas,
  carregando,
}: {
  propostas: PropostaParaAvaliarSaida[];
  carregando: boolean;
}) {
  if (carregando) {
    return <p className="text-center text-gray-600">Carregando...</p>;
  }

  if (propostas.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-gray-300 bg-white p-12 text-center">
        <p className="text-gray-600">Nenhuma proposta pendente de avaliação.</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {propostas.map((p) => (
        <div
          key={p.id}
          className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm"
        >
          <div className="mb-4 flex items-start justify-between">
            <div>
              <div className="flex items-center gap-2">
                <h3 className="text-lg font-semibold text-gray-900">
                  {p.dados_proponente.nome}
                </h3>
                <span
                  className={`inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${corStatusProposta(p.status)}`}
                >
                  {formatarStatusProposta(p.status)}
                </span>
              </div>
              <p className="mt-1 text-sm text-gray-600">
                Protocolo: {p.protocolo} — Edital #{p.edital_id}
              </p>
            </div>
            <Link
              href={{ pathname: `/avaliacoes/${p.id}` }}
              className="rounded-lg bg-brand-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-700"
            >
              Avaliar
            </Link>
          </div>

          <Section title="Proponente">
            <Field label="Nome" value={p.dados_proponente.nome} />
            <Field label="CPF" value={p.dados_proponente.cpf} />
            <Field label="E-mail" value={p.dados_proponente.email} />
          </Section>

          <Section title="Perfil Acadêmico">
            <Field label="Maior Titulação" value={p.dados_academicos.maior_titulacao} />
            <Field label="Área de Conhecimento" value={p.dados_academicos.area_conhecimento} />
          </Section>

          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-600">
              Submetida em {formatarDataISO(p.data_submissao)}
            </span>
            <span className="font-medium text-gray-900">
              Valor solicitado: {formatarMoeda(p.valor_total)}
            </span>
          </div>

          {p.pareceres.length > 0 && (
            <div className="mt-4 border-t border-gray-100 pt-4">
              <p className="text-xs font-medium text-gray-500">
                {p.pareceres.length} parecer(es) já emitido(s)
              </p>
            </div>
          )}
        </div>
      ))}
    </div>
  );
}
