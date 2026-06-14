"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

type Edital = {
  id: number;
  nome: string;
  descricao: string;
  data_inicio: string;
  data_fim: string;
  status: string;
  tipo_chamada: string;
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

export default function EditaisPage() {
  const router = useRouter();
  const [editais, setEditais] = useState<Edital[]>([]);
  const [carregando, setCarregando] = useState(true);

  const [filtroTitulo, setFiltroTitulo] = useState("");
  const [filtroStatus, setFiltroStatus] = useState("");
  const [filtroTipo, setFiltroTipo] = useState("");

  useEffect(() => {
    carregarEditais();
  }, []);

  async function carregarEditais() {
    setCarregando(true);
    try {
      const res = await fetch("/api/v1/editais");
      const data = await res.json();
      setEditais(data.editais ?? []);
    } catch {
      setEditais([]);
    } finally {
      setCarregando(false);
    }
  }

  const tipos = useMemo(() => {
    const unique = new Set(editais.map((e) => e.tipo_chamada).filter(Boolean));
    return Array.from(unique).sort();
  }, [editais]);

  const filtrados = useMemo(() => {
    return editais.filter((e) => {
      if (
        filtroTitulo &&
        !e.nome.toLowerCase().includes(filtroTitulo.toLowerCase())
      )
        return false;
      if (filtroStatus && e.status !== filtroStatus) return false;
      if (filtroTipo && e.tipo_chamada !== filtroTipo) return false;
      return true;
    });
  }, [editais, filtroTitulo, filtroStatus, filtroTipo]);

  function limparFiltros() {
    setFiltroTitulo("");
    setFiltroStatus("");
    setFiltroTipo("");
  }

  function labelStatus(status: string) {
    return statusLabel[status] || status;
  }

  return (
    <main className="mx-auto max-w-6xl px-4 py-6">
      <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Editais</h1>
        <Link
          href="/editais/novo"
          className="inline-flex items-center gap-1.5 rounded-lg bg-brand-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-700"
        >
          + Novo Edital
        </Link>
      </div>

      <div className="mb-6 rounded-xl border border-gray-200 bg-white p-4 shadow-sm">
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
          <div>
            <label
              htmlFor="filtro-titulo"
              className="mb-1 block text-xs font-medium text-gray-500"
            >
              Título
            </label>
            <input
              id="filtro-titulo"
              type="text"
              value={filtroTitulo}
              onChange={(e) => setFiltroTitulo(e.target.value)}
              placeholder="Buscar por título..."
              className="block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
            />
          </div>
          <div>
            <label
              htmlFor="filtro-status"
              className="mb-1 block text-xs font-medium text-gray-500"
            >
              Status
            </label>
            <select
              id="filtro-status"
              value={filtroStatus}
              onChange={(e) => setFiltroStatus(e.target.value)}
              className="block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
            >
              <option value="">Todos</option>
              <option value="ativo">Ativo</option>
              <option value="encerrado">Encerrado</option>
              <option value="em_avaliacao">Em Avaliação</option>
            </select>
          </div>
          <div>
            <label
              htmlFor="filtro-tipo"
              className="mb-1 block text-xs font-medium text-gray-500"
            >
              Tipo de Chamada
            </label>
            <select
              id="filtro-tipo"
              value={filtroTipo}
              onChange={(e) => setFiltroTipo(e.target.value)}
              className="block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
            >
              <option value="">Todos</option>
              {tipos.map((t) => (
                <option key={t} value={t}>
                  {t}
                </option>
              ))}
            </select>
          </div>
          <div className="flex items-end">
            <button
              onClick={limparFiltros}
              className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm text-gray-600 transition-colors hover:bg-gray-100"
            >
              Limpar Filtros
            </button>
          </div>
        </div>
      </div>

      {carregando ? (
        <p className="text-center text-gray-600">Carregando...</p>
      ) : filtrados.length === 0 ? (
        <div className="rounded-xl border border-dashed border-gray-300 bg-white p-12 text-center">
          <p className="text-gray-600">Nenhum edital encontrado.</p>
          <Link
            href="/editais/novo"
            className="mt-2 inline-block text-sm font-medium text-brand-600 hover:text-brand-700"
          >
            Criar primeiro edital
          </Link>
        </div>
      ) : (
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b bg-gray-50 text-left text-gray-600">
                  <th className="whitespace-nowrap px-4 py-3 font-medium">
                    Título
                  </th>
                  <th className="whitespace-nowrap px-4 py-3 font-medium">
                    Tipo
                  </th>
                  <th className="whitespace-nowrap px-4 py-3 font-medium">
                    Vigência
                  </th>
                  <th className="whitespace-nowrap px-4 py-3 font-medium">
                    Status
                  </th>
                  <th className="whitespace-nowrap px-4 py-3 font-medium">
                    Ações
                  </th>
                </tr>
              </thead>
              <tbody>
                {filtrados.map((edital) => (
                  <tr
                    key={edital.id}
                    className="border-b last:border-0 hover:bg-gray-50"
                  >
                    <td className="px-4 py-3">
                      <Link
                        href={`/editais/${edital.id}`}
                        className="font-medium text-brand-600 hover:text-brand-700"
                      >
                        {edital.nome}
                      </Link>
                      <p className="mt-0.5 text-xs text-gray-600 line-clamp-1">
                        {edital.descricao}
                      </p>
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 text-gray-600">
                      {edital.tipo_chamada || "—"}
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 text-gray-600">
                      {new Date(edital.data_inicio).toLocaleDateString(
                        "pt-BR",
                      )}{" "}
                      a{" "}
                      {new Date(edital.data_fim).toLocaleDateString("pt-BR")}
                    </td>
                    <td className="px-4 py-3">
                      <span
                        className={`inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${statusColor[edital.status] || "bg-gray-100 text-gray-600"}`}
                      >
                        {labelStatus(edital.status)}
                      </span>
                    </td>
                    <td className="whitespace-nowrap px-4 py-3">
                      <div className="flex items-center gap-1">
                        <Link
                          href={`/editais/${edital.id}`}
                          aria-label={`Visualizar propostas do edital ${edital.nome}`}
                          className="rounded-md bg-brand-50 px-2.5 py-1.5 text-xs font-medium text-brand-700 transition-colors hover:bg-brand-100"
                        >
                          Propostas
                        </Link>
                        <Link
                          href={`/editais/${edital.id}#avaliacao`}
                          aria-label={`Visualizar formulário de avaliação do edital ${edital.nome}`}
                          className="rounded-md bg-amber-50 px-2.5 py-1.5 text-xs font-medium text-amber-700 transition-colors hover:bg-amber-100"
                        >
                          Avaliação
                        </Link>
                        <Link
                          href={`/avaliadores`}
                          aria-label={`Visualizar avaliadores do edital ${edital.nome}`}
                          className="rounded-md bg-blue-50 px-2.5 py-1.5 text-xs font-medium text-blue-700 transition-colors hover:bg-blue-100"
                        >
                          Avaliadores
                        </Link>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </main>
  );
}
