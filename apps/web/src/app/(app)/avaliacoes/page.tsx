"use client";

import { useEffect, useState } from "react";
import { useAuth } from "../../../context/auth-context";
import TabelaPropostasPendentes from "../../../components/avaliacoes/TabelaPropostasPendentes";
import { listarPropostasParaAvaliar } from "../../../lib/api-avaliacoes";
import type { PropostaParaAvaliarSaida } from "../../../lib/api-avaliacoes";

export default function AvaliacoesPage() {
  const { usuario } = useAuth();
  const [propostas, setPropostas] = useState<PropostaParaAvaliarSaida[]>([]);
  const [carregando, setCarregando] = useState(true);

  useEffect(() => {
    if (!usuario?.id) return;

    listarPropostasParaAvaliar(Number(usuario.id))
      .then(setPropostas)
      .catch(() => setPropostas([]))
      .finally(() => setCarregando(false));
  }, [usuario]);

  return (
    <main className="mx-auto max-w-4xl px-4 py-6">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Avaliações</h1>
        <p className="mt-1 text-sm text-gray-600">
          Propostas pendentes de avaliação nos editais aos quais você está atribuído.
        </p>
      </div>

      <TabelaPropostasPendentes propostas={propostas} carregando={carregando} />
    </main>
  );
}
