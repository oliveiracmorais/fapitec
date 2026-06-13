"use client";

import { useState } from "react";
import { finalizarAvaliacao } from "../../../../../lib/api-avaliacoes";
import type { ClassificacaoSaida } from "../../../../../lib/api-avaliacoes";

export default function ClassificacaoPage() {
  const [editalId, setEditalId] = useState("");
  const [notaDeCorte, setNotaDeCorte] = useState(70);
  const [classificacao, setClassificacao] = useState<ClassificacaoSaida[] | null>(null);
  const [carregando, setCarregando] = useState(false);
  const [erro, setErro] = useState("");

  async function handleFinalizar(e: React.FormEvent) {
    e.preventDefault();
    setErro("");

    if (!editalId) {
      setErro("Informe o ID do edital.");
      return;
    }

    setCarregando(true);
    try {
      const resultado = await finalizarAvaliacao(Number(editalId), {
        nota_de_corte: notaDeCorte,
      });
      setClassificacao(resultado);
    } catch (err) {
      setErro(err instanceof Error ? err.message : "Erro ao finalizar avaliação.");
    } finally {
      setCarregando(false);
    }
  }

  return (
    <main className="mx-auto max-w-4xl px-4 py-6">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Classificação Final</h1>
        <p className="mt-1 text-sm text-gray-600">
          Finalize a avaliação de um edital e gere a classificação das propostas.
        </p>
      </div>

      <div className="mb-8 rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <form onSubmit={handleFinalizar} className="flex flex-wrap items-end gap-4">
          <div>
            <label className="mb-1 block text-sm font-medium text-gray-900">
              ID do Edital
            </label>
            <input
              type="number"
              min={1}
              value={editalId}
              onChange={(e) => setEditalId(e.target.value)}
              placeholder="Ex: 1"
              className="w-32 rounded-lg border border-gray-300 px-4 py-2 text-sm outline-none focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium text-gray-900">
              Nota de Corte
            </label>
            <input
              type="number"
              min={0}
              max={100}
              value={notaDeCorte}
              onChange={(e) => setNotaDeCorte(Number(e.target.value))}
              className="w-32 rounded-lg border border-gray-300 px-4 py-2 text-sm outline-none focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
            />
          </div>
          <button
            type="submit"
            disabled={carregando}
            className="rounded-lg bg-brand-600 px-6 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-700 disabled:opacity-50"
          >
            {carregando ? "Finalizando..." : "Finalizar Avaliação"}
          </button>
        </form>

        {erro && (
          <div className="mt-4 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">
            {erro}
          </div>
        )}
      </div>

      {classificacao && (
        <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="border-b border-gray-200 px-6 py-4">
            <h2 className="text-lg font-semibold text-gray-900">
              Resultado — Edital #{editalId}
            </h2>
            <p className="text-sm text-gray-600">
              Nota de corte: {notaDeCorte}
            </p>
          </div>
          {classificacao.length === 0 ? (
            <div className="p-12 text-center text-gray-600">
              Nenhuma proposta encontrada para este edital.
            </div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b bg-gray-50 text-left text-gray-600">
                  <th className="px-6 py-3 font-medium">#</th>
                  <th className="px-6 py-3 font-medium">Protocolo</th>
                  <th className="px-6 py-3 font-medium">Nota Final</th>
                  <th className="px-6 py-3 font-medium">Status</th>
                </tr>
              </thead>
              <tbody>
                {classificacao
                  .sort((a, b) => b.nota_final - a.nota_final)
                  .map((item, i) => (
                    <tr
                      key={item.proposta_id}
                      className="border-b last:border-0 hover:bg-gray-50"
                    >
                      <td className="px-6 py-3 font-medium text-gray-900">
                        {i + 1}º
                      </td>
                      <td className="px-6 py-3 text-gray-600">{item.protocolo}</td>
                      <td className="px-6 py-3">
                        <span className="font-semibold text-gray-900">
                          {item.nota_final}
                        </span>
                      </td>
                      <td className="px-6 py-3">
                        <span
                          className={`inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${
                            item.status === "aprovada"
                              ? "bg-green-100 text-green-800"
                              : "bg-red-100 text-red-800"
                          }`}
                        >
                          {item.status === "aprovada" ? "Aprovada" : "Reprovada"}
                        </span>
                      </td>
                    </tr>
                  ))}
              </tbody>
            </table>
          )}
        </div>
      )}
    </main>
  );
}
