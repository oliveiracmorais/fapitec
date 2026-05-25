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

export default function EditalDetalhePage() {
  const router = useRouter();
  const params = useParams();
  const [edital, setEdital] = useState<Edital | null>(null);
  const [carregando, setCarregando] = useState(true);

  useEffect(() => {
    const sessao =
      typeof window !== "undefined" && localStorage.getItem("sessao");
    if (!sessao) {
      router.replace("/");
      return;
    }

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

        <div className="grid grid-cols-2 gap-6 border-t border-gray-100 pt-6 text-sm">
          <div>
            <span className="font-medium text-gray-900">Tipo de Chamada</span>
            <p className="mt-0.5 text-gray-600">
              {edital.tipo_chamada || "—"}
            </p>
          </div>
          <div>
            <span className="font-medium text-gray-900">Vigência</span>
            <p className="mt-0.5 text-gray-600">
              {edital.data_inicio} a {edital.data_fim}
            </p>
          </div>
          <div>
            <span className="font-medium text-gray-900">Criado em</span>
            <p className="mt-0.5 text-gray-600">
              {new Date(edital.criado_em).toLocaleDateString("pt-BR")}
            </p>
          </div>
          <div>
            <span className="font-medium text-gray-900">ID</span>
            <p className="mt-0.5 text-gray-600">#{edital.id}</p>
          </div>
        </div>
      </div>
    </main>
  );
}
