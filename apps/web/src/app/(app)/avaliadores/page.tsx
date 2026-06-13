"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { listarAvaliadores, formatarEstado, corEstado, formatarDataISO } from "../../../lib/api-avaliadores";
import type { AvaliadorSaida } from "../../../lib/api-avaliadores";

export default function AvaliadoresPage() {
  const [avaliadores, setAvaliadores] = useState<AvaliadorSaida[]>([]);
  const [filtro, setFiltro] = useState({ nome: "", estado: "" });
  const [carregando, setCarregando] = useState(true);

  useEffect(() => {
    carregar();
  }, []);

  async function carregar(params?: { nome?: string; estado?: string }) {
    setCarregando(true);
    try {
      const data = await listarAvaliadores(params);
      setAvaliadores(data);
    } catch {
      setAvaliadores([]);
    } finally {
      setCarregando(false);
    }
  }

  function handleFiltrar(e: React.FormEvent) {
    e.preventDefault();
    const params: { nome?: string; estado?: string } = {};
    if (filtro.nome) params.nome = filtro.nome;
    if (filtro.estado) params.estado = filtro.estado;
    carregar(params);
  }

  return (
    <main className="mx-auto max-w-6xl px-4 py-6">
      <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Avaliadores</h1>
        <Link
          href={{ pathname: "/avaliadores/novo" }}
          className="inline-flex items-center gap-1.5 rounded-lg bg-brand-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-700"
        >
          + Novo Avaliador
        </Link>
      </div>

      <form onSubmit={handleFiltrar} className="mb-6 flex gap-2">
        <input
          type="text"
          value={filtro.nome}
          onChange={(e) => setFiltro((prev) => ({ ...prev, nome: e.target.value }))}
          placeholder="Buscar por nome..."
          className="flex-1 rounded-lg border border-gray-300 px-4 py-2 text-sm outline-none focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
        />
        <select
          value={filtro.estado}
          onChange={(e) => setFiltro((prev) => ({ ...prev, estado: e.target.value }))}
          className="rounded-lg border border-gray-300 px-4 py-2 text-sm outline-none focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
        >
          <option value="">Todos</option>
          <option value="ativo">Ativo</option>
          <option value="inativo">Inativo</option>
        </select>
        <button
          type="submit"
          className="rounded-lg bg-brand-600 px-4 py-2 text-sm font-medium text-white hover:bg-brand-700"
        >
          Buscar
        </button>
      </form>

      {carregando ? (
        <p className="text-center text-gray-600">Carregando...</p>
      ) : avaliadores.length === 0 ? (
        <div className="rounded-xl border border-dashed border-gray-300 bg-white p-12 text-center">
          <p className="text-gray-600">Nenhum avaliador encontrado.</p>
          <Link
            href={{ pathname: "/avaliadores/novo" }}
            className="mt-2 inline-block text-sm font-medium text-brand-600 hover:text-brand-700"
          >
            Cadastrar primeiro avaliador
          </Link>
        </div>
      ) : (
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b bg-gray-50 text-left text-gray-600">
                <th className="px-4 py-3 font-medium">Nome</th>
                <th className="px-4 py-3 font-medium">CPF</th>
                <th className="px-4 py-3 font-medium">E-mail</th>
                <th className="px-4 py-3 font-medium">Titulação</th>
                <th className="px-4 py-3 font-medium">Estado</th>
                <th className="px-4 py-3 font-medium">Cadastro</th>
                <th className="px-4 py-3" />
              </tr>
            </thead>
            <tbody>
              {avaliadores.map((av) => (
                <tr key={av.id} className="border-b last:border-0 hover:bg-gray-50">
                  <td className="px-4 py-3">
                    <Link
                      href={{ pathname: `/avaliadores/${av.id}` }}
                      className="font-medium text-brand-600 hover:text-brand-700"
                    >
                      {av.nome}
                    </Link>
                  </td>
                  <td className="px-4 py-3 text-gray-600">{av.cpf}</td>
                  <td className="px-4 py-3 text-gray-600">{av.email}</td>
                  <td className="px-4 py-3 text-gray-600">{av.titulacao_maxima || "—"}</td>
                  <td className="px-4 py-3">
                    <span
                      className={`inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${corEstado(av.estado)}`}
                    >
                      {formatarEstado(av.estado)}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-gray-600">{formatarDataISO(av.data_cadastro)}</td>
                  <td className="px-4 py-3 text-right">
                    <Link
                      href={{ pathname: `/avaliadores/${av.id}` }}
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
