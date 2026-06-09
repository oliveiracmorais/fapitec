"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import PropostaPainel from "../../../components/proposta-painel";
import { listarPropostas, type PropostaResumo } from "../../../lib/api-propostas";

export default function MinhasPropostasPage() {
  const [propostas, setPropostas] = useState<PropostaResumo[]>([]);
  const [carregando, setCarregando] = useState(true);

  useEffect(() => {
    listarPropostas()
      .then(setPropostas)
      .catch(() => setPropostas([]))
      .finally(() => setCarregando(false));
  }, []);

  return (
    <main className="mx-auto max-w-6xl px-4 py-6">
      <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Minhas Propostas</h1>
          <p className="mt-1 text-sm text-gray-600">
            Acompanhe o status das suas propostas submetidas.
          </p>
        </div>
        <Link
          href="/editais"
          className="inline-flex items-center gap-1.5 rounded-lg bg-brand-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-700"
        >
          Ver Editais Disponíveis
        </Link>
      </div>

      <PropostaPainel propostas={propostas} carregando={carregando} />
    </main>
  );
}
