"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import PropostaResumoComponent from "../../../../components/proposta-resumo";
import PropostaStatusBadge from "../../../../components/proposta-status-badge";
import { visualizarProposta, deletarProposta, submeterProposta, formatarBRL, type Proposta } from "../../../../lib/api-propostas";

export default function PropostaDetalhePage() {
  const params = useParams();
  const router = useRouter();
  const [proposta, setProposta] = useState<Proposta | null>(null);
  const [carregando, setCarregando] = useState(true);
  const [acao, setAcao] = useState("");

  function irParaMinhasPropostas() {
    router.push("/minhas-propostas" as never);
  }

  useEffect(() => {
    if (!params?.id) return;
    const id = Number(params.id);
    if (!id) { irParaMinhasPropostas(); return; }

    visualizarProposta(id)
      .then(setProposta)
      .catch(() => irParaMinhasPropostas())
      .finally(() => setCarregando(false));
  }, [params, router]);

  async function handleSubmeter() {
    if (!proposta) return;
    setAcao("submeter");
    try {
      const atualizada = await submeterProposta(proposta.id);
      setProposta(atualizada);
    } catch (e) {
      alert(e instanceof Error ? e.message : "Erro ao submeter proposta");
    } finally {
      setAcao("");
    }
  }

  async function handleDeletar() {
    if (!proposta) return;
    if (!confirm("Tem certeza que deseja excluir esta proposta?")) return;
    setAcao("deletar");
    try {
      await deletarProposta(proposta.id);
      irParaMinhasPropostas();
    } catch (e) {
      alert(e instanceof Error ? e.message : "Erro ao deletar proposta");
    } finally {
      setAcao("");
    }
  }

  if (carregando) {
    return (
      <main className="mx-auto max-w-3xl px-4 py-6">
        <div className="flex items-center justify-center py-20 text-gray-600">
          Carregando...
        </div>
      </main>
    );
  }

  if (!proposta) return null;

  const isRascunho = proposta.status === "rascunho";

  return (
    <main className="mx-auto max-w-3xl px-4 py-6">
      <div className="rounded-xl border border-gray-200 bg-white p-8 shadow-sm">
        <div className="mb-6 flex flex-wrap items-start justify-between gap-4">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              Proposta #{proposta.id}
            </h1>
            <p className="mt-1 text-sm text-gray-600">
              Protocolo: <span className="font-mono font-medium">{proposta.protocolo || "—"}</span>
            </p>
            <p className="text-xs text-gray-500">
              Versão: {proposta.versao} · {new Date(proposta.data_atualizacao).toLocaleDateString("pt-BR")}
            </p>
          </div>
          <PropostaStatusBadge status={proposta.status} />
        </div>

        <div className="mb-6 flex flex-wrap gap-4">
          {isRascunho && (
            <>
          <Link
            href={{ pathname: `/editais/${proposta.edital_id}/inscrever`, query: { editar: proposta.id } }}
                className="rounded-lg bg-brand-600 px-4 py-2 text-sm font-medium text-white hover:bg-brand-700"
              >
                Editar Proposta
              </Link>
              <button
                type="button"
                onClick={handleSubmeter}
                disabled={acao === "submeter"}
                className="rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-50"
              >
                {acao === "submeter" ? "Enviando..." : "Submeter para Avaliação"}
              </button>
              <button
                type="button"
                onClick={handleDeletar}
                disabled={acao === "deletar"}
                className="rounded-lg border border-red-300 px-4 py-2 text-sm font-medium text-red-700 hover:bg-red-50 disabled:opacity-50"
              >
                {acao === "deletar" ? "Excluindo..." : "Excluir Rascunho"}
              </button>
            </>
          )}
        </div>

        <div className="space-y-6">
          <div className="rounded-lg border border-gray-100 bg-gray-50 p-4">
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="font-medium text-gray-900">Edital</span>
                <p className="mt-0.5 text-gray-600">#{proposta.edital_id}</p>
              </div>
              <div>
                <span className="font-medium text-gray-900">Valor Solicitado</span>
                <p className="mt-0.5 text-gray-600">{formatarBRL(proposta.valor_total_solicitado)}</p>
              </div>
              {proposta.data_submissao && (
                <div>
                  <span className="font-medium text-gray-900">Submetida em</span>
                  <p className="mt-0.5 text-gray-600">
                    {new Date(proposta.data_submissao).toLocaleDateString("pt-BR")}
                  </p>
                </div>
              )}
            </div>
          </div>

          <PropostaResumoComponent
            dadosProponente={proposta.dados_proponente}
            dadosAcademicos={proposta.dados_academicos}
            itensOrcamentarios={proposta.itens_orcamentarios}
            valorGlobalEdital={0}
            status={proposta.status}
            protocolo={proposta.protocolo}
          />
        </div>

        <div className="mt-6 border-t border-gray-100 pt-6">
          <Link
            href={{ pathname: "/minhas-propostas" }}
            className="text-sm font-medium text-brand-600 hover:text-brand-700"
          >
            ← Voltar para minhas propostas
          </Link>
        </div>
      </div>
    </main>
  );
}
