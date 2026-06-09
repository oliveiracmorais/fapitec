"use client";

import type { DadosProponente, DadosAcademicos, ItemOrcamentario } from "../lib/api-propostas";
import { formatarBRL } from "../lib/api-propostas";
import PropostaStatusBadge from "./proposta-status-badge";

type Props = {
  dadosProponente: DadosProponente;
  dadosAcademicos: DadosAcademicos;
  itensOrcamentarios: ItemOrcamentario[];
  status?: string;
  protocolo?: string;
  valorGlobalEdital: number;
};

function Campo({ label, valor }: { label: string; valor: string }) {
  return (
    <div>
      <span className="text-xs font-medium text-gray-500">{label}</span>
      <p className="text-sm text-gray-900">{valor || "—"}</p>
    </div>
  );
}

export default function PropostaResumo({ dadosProponente, dadosAcademicos, itensOrcamentarios, status, protocolo, valorGlobalEdital }: Props) {
  const valorTotal = itensOrcamentarios.reduce((acc, item) => acc + item.valor_total, 0);
  const excede = valorGlobalEdital > 0 && valorTotal > valorGlobalEdital;

  return (
    <div className="space-y-6">
      <h2 className="text-lg font-semibold text-gray-900">Resumo da Proposta</h2>

      {protocolo && (
        <div className="rounded-lg bg-brand-50 p-4 text-center">
          <p className="text-xs text-brand-700">Protocolo</p>
          <p className="text-xl font-bold text-brand-800">{protocolo}</p>
        </div>
      )}

      {status && (
        <div className="flex items-center justify-center gap-2">
          <span className="text-sm text-gray-600">Status:</span>
          <PropostaStatusBadge status={status} />
        </div>
      )}

      <div className="rounded-lg border border-gray-200 p-4">
        <h3 className="mb-3 text-sm font-semibold text-gray-700">Dados do Proponente</h3>
        <div className="grid grid-cols-2 gap-3">
          <Campo label="Nome" valor={dadosProponente.nome} />
          <Campo label="CPF" valor={dadosProponente.cpf} />
          <Campo label="E-mail" valor={dadosProponente.email} />
          <Campo label="Telefone" valor={dadosProponente.telefone} />
        </div>
      </div>

      <div className="rounded-lg border border-gray-200 p-4">
        <h3 className="mb-3 text-sm font-semibold text-gray-700">Formação Acadêmica</h3>
        <div className="grid grid-cols-2 gap-3">
          <Campo label="Maior Titulação" valor={dadosAcademicos.maior_titulacao} />
          <Campo label="Curso" valor={dadosAcademicos.curso} />
          <Campo label="Instituição" valor={dadosAcademicos.instituicao} />
          <Campo label="Ano Conclusão" valor={dadosAcademicos.ano_conclusao ? String(dadosAcademicos.ano_conclusao) : ""} />
        </div>
      </div>

      <div className="rounded-lg border border-gray-200 p-4">
        <h3 className="mb-3 text-sm font-semibold text-gray-700">Orçamento</h3>
        {itensOrcamentarios.length === 0 ? (
          <p className="text-sm text-gray-500">Nenhum item orçamentário.</p>
        ) : (
          <div className="space-y-1">
            {itensOrcamentarios.map((item, i) => (
              <div key={i} className="flex justify-between text-sm">
                <span className="text-gray-600">
                  {item.descricao} ({item.tipo}) — {item.quantidade}x {formatarBRL(item.valor_unitario)}
                </span>
                <span className="font-medium text-gray-900">{formatarBRL(item.valor_total)}</span>
              </div>
            ))}
            <div className={`border-t pt-2 text-sm font-bold ${excede ? "text-red-700" : "text-gray-900"}`}>
              <div className="flex justify-between">
                <span>Total</span>
                <span>{formatarBRL(valorTotal)}</span>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
