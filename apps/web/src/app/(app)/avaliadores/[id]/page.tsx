"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import Link from "next/link";
import {
  visualizarAvaliador,
  listarAtribuicoesPorAvaliador,
  formatarEstado,
  corEstado,
  formatarStatusConvite,
  corStatusConvite,
  formatarDataISO,
} from "../../../../lib/api-avaliadores";
import type { AvaliadorSaida, AtribuicaoSaida } from "../../../../lib/api-avaliadores";
import AtribuirEditalDialog from "../../../../components/avaliadores/AtribuirEditalDialog";

type Aba = "dados" | "atribuicoes";

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="border-b border-gray-100 pb-6 last:border-b-0">
      <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-gray-500">{title}</h2>
      <div className="grid grid-cols-2 gap-4 text-sm">{children}</div>
    </div>
  );
}

function Field({ label, value }: { label: string; value: string | number | null | undefined }) {
  const display = value === null || value === undefined || value === "" ? "—" : String(value);
  return (
    <div>
      <span className="font-medium text-gray-900">{label}</span>
      <p className="mt-0.5 text-gray-600">{display}</p>
    </div>
  );
}

export default function AvaliadorDetalhePage() {
  const router = useRouter();
  const params = useParams();
  const [avaliador, setAvaliador] = useState<AvaliadorSaida | null>(null);
  const [atribuicoes, setAtribuicoes] = useState<AtribuicaoSaida[]>([]);
  const [carregando, setCarregando] = useState(true);
  const [aba, setAba] = useState<Aba>("dados");
  const [dialogOpen, setDialogOpen] = useState(false);

  useEffect(() => {
    if (!params?.id) return;
    const id = Number(params.id);

    Promise.all([
      visualizarAvaliador(id),
      listarAtribuicoesPorAvaliador(id),
    ])
      .then(([av, atr]) => {
        setAvaliador(av);
        setAtribuicoes(atr);
      })
      .catch(() => router.push("/avaliadores"))
      .finally(() => setCarregando(false));
  }, [params, router]);

  if (carregando) {
    return (
      <div className="flex items-center justify-center py-20 text-gray-600">Carregando...</div>
    );
  }

  if (!avaliador) return null;

  return (
    <main className="mx-auto max-w-4xl px-4 py-6">
      <div className="rounded-xl border border-gray-200 bg-white p-8 shadow-sm">
        <div className="mb-6 flex items-start justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{avaliador.nome}</h1>
            <p className="mt-1 text-sm text-gray-600">{avaliador.email}</p>
          </div>
          <span
            className={`rounded-full px-3 py-1 text-xs font-medium ${corEstado(avaliador.estado)}`}
          >
            {formatarEstado(avaliador.estado)}
          </span>
        </div>

        <div className="mb-6 border-b border-gray-200">
          <div className="flex gap-6">
            <button
              onClick={() => setAba("dados")}
              className={`border-b-2 pb-3 text-sm font-medium transition-colors ${
                aba === "dados"
                  ? "border-brand-600 text-brand-600"
                  : "border-transparent text-gray-600 hover:text-gray-900"
              }`}
            >
              Dados do Avaliador
            </button>
            <button
              onClick={() => setAba("atribuicoes")}
              className={`border-b-2 pb-3 text-sm font-medium transition-colors ${
                aba === "atribuicoes"
                  ? "border-brand-600 text-brand-600"
                  : "border-transparent text-gray-600 hover:text-gray-900"
              }`}
            >
              Atribuições ({atribuicoes.length})
            </button>
          </div>
        </div>

        {aba === "dados" && (
          <div className="space-y-6">
            <Section title="Informações Pessoais">
              <Field label="Nome" value={avaliador.nome} />
              <Field label="CPF" value={avaliador.cpf} />
              <Field label="E-mail" value={avaliador.email} />
              <Field label="Estado" value={formatarEstado(avaliador.estado)} />
            </Section>

            <Section title="Perfil Técnico">
              <Field label="Titulação Máxima" value={avaliador.titulacao_maxima} />
              <Field label="Área de Conhecimento" value={avaliador.area_conhecimento} />
              <Field label="Instituição" value={avaliador.instituicao} />
              <Field label="Total de Propostas" value={avaliador.total_propostas ?? 0} />
              <Field label="Atribuições Ativas" value={avaliador.atribuicoes_ativas ?? 0} />
            </Section>

            {avaliador.curriculo_resumido && (
              <Section title="Currículo Resumido">
                <div className="col-span-2">
                  <p className="text-gray-600 whitespace-pre-wrap">{avaliador.curriculo_resumido}</p>
                </div>
              </Section>
            )}

            <Section title="Datas">
              <Field label="Cadastro" value={formatarDataISO(avaliador.data_cadastro)} />
              <Field label="Atualização" value={formatarDataISO(avaliador.data_atualizacao)} />
            </Section>
          </div>
        )}

        {aba === "atribuicoes" && (
          <div>
            <div className="mb-4 flex items-center justify-between">
              <p className="text-sm text-gray-600">
                {atribuicoes.length === 0
                  ? "Nenhuma atribuição encontrada."
                  : `${atribuicoes.length} atribuição(ões)`}
              </p>
              <button
                onClick={() => setDialogOpen(true)}
                className="rounded-lg bg-brand-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-700"
              >
                + Atribuir Edital
              </button>
            </div>

            {atribuicoes.length > 0 && (
              <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b bg-gray-50 text-left text-gray-600">
                      <th className="px-4 py-3 font-medium">Edital</th>
                      <th className="px-4 py-3 font-medium">Vigência</th>
                      <th className="px-4 py-3 font-medium">Status Convite</th>
                      <th className="px-4 py-3 font-medium">Atribuído em</th>
                    </tr>
                  </thead>
                  <tbody>
                    {atribuicoes.map((atr) => (
                      <tr key={atr.id} className="border-b last:border-0 hover:bg-gray-50">
                        <td className="px-4 py-3 font-medium text-gray-900">
                          Edital #{atr.edital_id}
                        </td>
                        <td className="px-4 py-3 text-gray-600">
                          {formatarDataISO(atr.data_inicio)} a {formatarDataISO(atr.data_fim)}
                        </td>
                        <td className="px-4 py-3">
                          <span
                            className={`inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${corStatusConvite(atr.status_convite)}`}
                          >
                            {formatarStatusConvite(atr.status_convite)}
                          </span>
                        </td>
                        <td className="px-4 py-3 text-gray-600">
                          {formatarDataISO(atr.criado_em)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        <div className="mt-6 flex items-center gap-4 border-t border-gray-100 pt-6">
          <Link
            href={{ pathname: "/avaliadores" }}
            className="text-sm font-medium text-brand-600 hover:text-brand-700"
          >
            ← Voltar para lista
          </Link>
          <Link
            href={{ pathname: `/avaliadores/${avaliador.id}/editar` }}
            className="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
          >
            Editar
          </Link>
        </div>
      </div>

      <AtribuirEditalDialog
        isOpen={dialogOpen}
        onClose={() => setDialogOpen(false)}
        avaliadorId={avaliador.id}
        avaliadorNome={avaliador.nome}
        onAtribuido={() => {
          listarAtribuicoesPorAvaliador(avaliador.id).then(setAtribuicoes);
          visualizarAvaliador(avaliador.id).then(setAvaliador);
        }}
      />
    </main>
  );
}
