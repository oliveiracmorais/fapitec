"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { listarPropostas, formatarStatus, corStatus, formatarBRL } from "../../../../lib/api-propostas";
import type { PropostaResumo } from "../../../../lib/api-propostas";

export default function AdminAvaliacoesPage() {
  const [propostas, setPropostas] = useState<PropostaResumo[]>([]);
  const [carregando, setCarregando] = useState(true);
  const [filtroEdital, setFiltroEdital] = useState("");
  const [filtroStatus, setFiltroStatus] = useState("");

  useEffect(() => {
    const params: { edital_id?: number } = {};
    if (filtroEdital) params.edital_id = Number(filtroEdital);

    setCarregando(true);
    listarPropostas(params)
      .then(setPropostas)
      .catch(() => setPropostas([]))
      .finally(() => setCarregando(false));
  }, [filtroEdital]);

  const filtradas = filtroStatus
    ? propostas.filter((p) => p.status === filtroStatus)
    : propostas;

  const editais = [...new Set(propostas.map((p) => p.edital_id))].sort();

  return (
    <main className="mx-auto max-w-6xl px-4 py-6">
      <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Acompanhamento de Avaliações</h1>
          <p className="mt-1 text-sm text-gray-600">
            Visualize o andamento das avaliações por edital e proposta.
          </p>
        </div>
        <Link
          href={{ pathname: "/admin/avaliacoes/classificacao" }}
          className="inline-flex items-center gap-1.5 rounded-lg bg-brand-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-700"
        >
          Classificação Final
        </Link>
      </div>

      <div className="mb-6 flex gap-2">
        <select
          value={filtroEdital}
          onChange={(e) => setFiltroEdital(e.target.value)}
          className="rounded-lg border border-gray-300 px-4 py-2 text-sm outline-none focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
        >
          <option value="">Todos os editais</option>
          {editais.map((id) => (
            <option key={id} value={id}>
              Edital #{id}
            </option>
          ))}
        </select>
        <select
          value={filtroStatus}
          onChange={(e) => setFiltroStatus(e.target.value)}
          className="rounded-lg border border-gray-300 px-4 py-2 text-sm outline-none focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
        >
          <option value="">Todos os status</option>
          <option value="submetida">Submetida</option>
          <option value="em_avaliacao">Em Avaliação</option>
          <option value="aprovada">Aprovada</option>
          <option value="reprovada">Reprovada</option>
        </select>
      </div>

      {carregando ? (
        <p className="text-center text-gray-600">Carregando...</p>
      ) : filtradas.length === 0 ? (
        <div className="rounded-xl border border-dashed border-gray-300 bg-white p-12 text-center">
          <p className="text-gray-600">Nenhuma proposta encontrada.</p>
        </div>
      ) : (
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b bg-gray-50 text-left text-gray-600">
                <th className="px-4 py-3 font-medium">Protocolo</th>
                <th className="px-4 py-3 font-medium">Edital</th>
                <th className="px-4 py-3 font-medium">Status</th>
                <th className="px-4 py-3 font-medium">Valor</th>
                <th className="px-4 py-3 font-medium">Versão</th>
                <th className="px-4 py-3" />
              </tr>
            </thead>
            <tbody>
              {filtradas.map((p) => (
                <tr key={p.id} className="border-b last:border-0 hover:bg-gray-50">
                  <td className="px-4 py-3 font-medium text-gray-900">{p.protocolo}</td>
                  <td className="px-4 py-3 text-gray-600">#{p.edital_id}</td>
                  <td className="px-4 py-3">
                    <span
                      className={`inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${corStatus(p.status)}`}
                    >
                      {formatarStatus(p.status)}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-gray-600">{formatarBRL(p.valor_total_solicitado)}</td>
                  <td className="px-4 py-3 text-gray-600">v{p.versao}</td>
                  <td className="px-4 py-3 text-right">
                    <Link
                      href={{ pathname: `/avaliacoes/${p.id}` }}
                      className="text-sm text-brand-600 hover:text-brand-700"
                    >
                      Ver →
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </main>
  );
}
