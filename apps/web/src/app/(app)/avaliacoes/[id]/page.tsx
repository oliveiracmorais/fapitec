"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import Link from "next/link";
import { useAuth } from "../../../../context/auth-context";
import FormParecer from "../../../../components/avaliacoes/FormParecer";
import {
  listarPropostasParaAvaliar,
  emitirParecer,
  listarPareceres,
  formatarStatusProposta,
  corStatusProposta,
  formatarDataISO,
} from "../../../../lib/api-avaliacoes";
import type {
  PropostaParaAvaliarSaida,
  ParecerAnonimizadoSaida,
} from "../../../../lib/api-avaliacoes";

export default function AvaliarPropostaPage() {
  const router = useRouter();
  const params = useParams();
  const { usuario } = useAuth();
  const [proposta, setProposta] = useState<PropostaParaAvaliarSaida | null>(null);
  const [pareceres, setPareceres] = useState<ParecerAnonimizadoSaida[]>([]);
  const [carregando, setCarregando] = useState(true);
  const [parecerEnviado, setParecerEnviado] = useState(false);

  useEffect(() => {
    if (!params?.id || !usuario?.id) return;
    const id = Number(params.id);

    Promise.all([
      listarPropostasParaAvaliar(Number(usuario.id)),
      listarPareceres(id),
    ])
      .then(([propostas, p]) => {
        const found = propostas.find((p) => p.id === id);
        if (!found) {
          router.push("/avaliacoes" as any);
          return;
        }
        setProposta(found);
        setPareceres(p);
      })
      .catch(() => router.push("/avaliacoes" as any))
      .finally(() => setCarregando(false));
  }, [params, usuario, router]);

  async function handleParecer(data: { nota: number; parecer_texto: string }) {
    if (!proposta || !usuario?.id) return;

    await emitirParecer(proposta.id, {
      proposta_id: proposta.id,
      etapa: "unica",
      avaliador_id: Number(usuario.id),
      nota: data.nota,
      parecer_texto: data.parecer_texto,
    });

    setParecerEnviado(true);
    const atualizados = await listarPareceres(proposta.id);
    setPareceres(atualizados);
  }

  if (carregando) {
    return (
      <div className="flex items-center justify-center py-20 text-gray-600">
        Carregando...
      </div>
    );
  }

  if (!proposta) return null;

  return (
    <main className="mx-auto max-w-4xl px-4 py-6">
      <div className="rounded-xl border border-gray-200 bg-white p-8 shadow-sm">
        <div className="mb-6 flex items-start justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              {proposta.dados_proponente.nome}
            </h1>
            <p className="mt-1 text-sm text-gray-600">
              Protocolo: {proposta.protocolo} — Edital #{proposta.edital_id}
            </p>
          </div>
          <span
            className={`inline-block rounded-full px-3 py-1 text-xs font-medium ${corStatusProposta(proposta.status)}`}
          >
            {formatarStatusProposta(proposta.status)}
          </span>
        </div>

        <div className="mb-6 grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="font-medium text-gray-900">CPF</span>
            <p className="mt-0.5 text-gray-600">{proposta.dados_proponente.cpf}</p>
          </div>
          <div>
            <span className="font-medium text-gray-900">E-mail</span>
            <p className="mt-0.5 text-gray-600">{proposta.dados_proponente.email}</p>
          </div>
          <div>
            <span className="font-medium text-gray-900">Maior Titulação</span>
            <p className="mt-0.5 text-gray-600">
              {proposta.dados_academicos.maior_titulacao}
            </p>
          </div>
          <div>
            <span className="font-medium text-gray-900">Área de Conhecimento</span>
            <p className="mt-0.5 text-gray-600">
              {proposta.dados_academicos.area_conhecimento}
            </p>
          </div>
          <div>
            <span className="font-medium text-gray-900">Data de Submissão</span>
            <p className="mt-0.5 text-gray-600">
              {formatarDataISO(proposta.data_submissao)}
            </p>
          </div>
        </div>

        {pareceres.length > 0 && (
          <div className="mb-6 rounded-lg bg-gray-50 p-4">
            <h3 className="mb-3 text-sm font-semibold text-gray-700">
              Pareceres já emitidos ({pareceres.length})
            </h3>
            <div className="space-y-3">
              {pareceres.map((p) => (
                <div
                  key={p.id}
                  className="rounded-lg border border-gray-200 bg-white p-4"
                >
                  <div className="mb-2 flex items-center justify-between text-sm">
                    <span className="font-medium text-gray-900">
                      Nota: {p.nota}
                    </span>
                    <span className="text-xs text-gray-500">
                      {formatarDataISO(p.data)} — {p.etapa}
                    </span>
                  </div>
                  <p className="whitespace-pre-wrap text-sm text-gray-600">
                    {p.parecer_texto}
                  </p>
                  <p className="mt-2 text-xs text-gray-400">
                    Avaliador: {p.hash_avaliador}
                  </p>
                </div>
              ))}
            </div>
          </div>
        )}

        {parecerEnviado ? (
          <div className="rounded-lg bg-green-50 p-6 text-center">
            <p className="font-medium text-green-800">
              Parecer salvo com sucesso!
            </p>
            <Link
              href={{ pathname: "/avaliacoes" }}
              className="mt-3 inline-block text-sm font-medium text-brand-600 hover:text-brand-700"
            >
              ← Voltar para lista de avaliações
            </Link>
          </div>
        ) : (
          <FormParecer
            propostaId={proposta.id}
            avaliadorId={Number(usuario?.id) || 0}
            onSubmit={handleParecer}
            onCancel={() => router.push("/avaliacoes" as any)}
          />
        )}
      </div>
    </main>
  );
}
