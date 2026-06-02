"use client";

import { useEffect, useState } from "react";
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
  const [filtro, setFiltro] = useState("");
  const [carregando, setCarregando] = useState(true);

  useEffect(() => {
    carregarEditais();
  }, []);

  async function carregarEditais(titulo?: string) {
    setCarregando(true);
    try {
      const params = titulo ? `?titulo=${encodeURIComponent(titulo)}` : "";
      const res = await fetch(`/api/v1/editais${params}`);
      const data = await res.json();
      setEditais(data.editais ?? []);
    } catch {
      setEditais([]);
    } finally {
      setCarregando(false);
    }
  }

  function handleFiltrar(e: React.FormEvent) {
    e.preventDefault();
    carregarEditais(filtro);
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

      <form onSubmit={handleFiltrar} className="mb-6 flex gap-2">
        <input
          type="text"
          value={filtro}
          onChange={(e) => setFiltro(e.target.value)}
          placeholder="Buscar por título..."
          className="flex-1 rounded-lg border border-gray-300 px-4 py-2 text-sm outline-none focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
        />
        <button
          type="submit"
          className="rounded-lg bg-brand-600 px-4 py-2 text-sm font-medium text-white hover:bg-brand-700"
        >
          Buscar
        </button>
      </form>

      {carregando ? (
        <p className="text-center text-gray-600">Carregando...</p>
      ) : editais.length === 0 ? (
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
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b bg-gray-50 text-left text-gray-600">
                <th className="px-4 py-3 font-medium">Título</th>
                <th className="px-4 py-3 font-medium">Tipo</th>
                <th className="px-4 py-3 font-medium">Vigência</th>
                <th className="px-4 py-3 font-medium">Status</th>
                <th className="px-4 py-3" />
              </tr>
            </thead>
            <tbody>
              {editais.map((edital) => (
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
                  <td className="px-4 py-3 text-gray-600">
                    {edital.tipo_chamada || "—"}
                  </td>
                  <td className="px-4 py-3 text-gray-600">
                    {edital.data_inicio} a {edital.data_fim}
                  </td>
                  <td className="px-4 py-3">
                    <span
                      className={`inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${statusColor[edital.status] || "bg-gray-100 text-gray-600"}`}
                    >
                      {statusLabel[edital.status] || edital.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-right">
                    <Link
                      href={`/editais/${edital.id}`}
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
