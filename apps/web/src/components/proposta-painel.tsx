"use client";

import Link from "next/link";
import type { PropostaResumo } from "../lib/api-propostas";
import { formatarBRL, formatarStatus, corStatus } from "../lib/api-propostas";

type Props = {
  propostas: PropostaResumo[];
  carregando: boolean;
};

export default function PropostaPainel({ propostas, carregando }: Props) {
  if (carregando) {
    return <p className="text-center text-gray-600">Carregando...</p>;
  }

  if (propostas.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-gray-300 bg-white p-12 text-center">
        <p className="text-gray-600">Nenhuma proposta encontrada.</p>
        <Link
          href="/editais"
          className="mt-2 inline-block text-sm font-medium text-brand-600 hover:text-brand-700"
        >
          Ver editais disponíveis
        </Link>
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b bg-gray-50 text-left text-gray-600">
            <th className="px-4 py-3 font-medium">Protocolo</th>
            <th className="px-4 py-3 font-medium">Edital</th>
            <th className="px-4 py-3 font-medium">Valor Solicitado</th>
            <th className="px-4 py-3 font-medium">Status</th>
            <th className="px-4 py-3 font-medium">Atualizado em</th>
            <th className="px-4 py-3" />
          </tr>
        </thead>
        <tbody>
          {propostas.map((p) => (
            <tr key={p.id} className="border-b last:border-0 hover:bg-gray-50">
              <td className="px-4 py-3 font-mono text-xs text-gray-700">
                {p.protocolo || "—"}
              </td>
              <td className="px-4 py-3 text-gray-700">
                {p.edital_nome || `Edital #${p.edital_id}`}
              </td>
              <td className="px-4 py-3 text-gray-700">
                {formatarBRL(p.valor_total_solicitado)}
              </td>
              <td className="px-4 py-3">
                <span className={`inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${corStatus(p.status)}`}>
                  {formatarStatus(p.status)}
                </span>
              </td>
              <td className="px-4 py-3 text-gray-500 text-xs">
                {new Date(p.data_atualizacao).toLocaleDateString("pt-BR")}
              </td>
              <td className="px-4 py-3 text-right">
                <Link
                  href={{ pathname: `/minhas-propostas/${p.id}` }}
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
  );
}
