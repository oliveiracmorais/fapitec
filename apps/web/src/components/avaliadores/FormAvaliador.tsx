"use client";

import { useState, FormEvent } from "react";
import Link from "next/link";
import CampoFormulario from "../campo-formulario";
import type { AvaliadorSaida, AtualizarAvaliadorEntrada } from "../../lib/api-avaliadores";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Props = { mode: "criar" | "editar"; initialData?: AvaliadorSaida; onSubmit: (data: any) => Promise<void>; onCancel: () => void };

const TITULACAO_OPTIONS = [
  { value: "Graduado", label: "Graduado" },
  { value: "Especialista", label: "Especialista" },
  { value: "Mestre", label: "Mestre" },
  { value: "Doutor", label: "Doutor" },
  { value: "Pos-Doutor", label: "Pós-Doutor" },
];

export default function FormAvaliador({ mode, initialData, onSubmit, onCancel }: Props) {
  const [form, setForm] = useState({
    nome: initialData?.nome || "",
    cpf: initialData?.cpf || "",
    email: initialData?.email || "",
    titulacao_maxima: initialData?.titulacao_maxima || "",
    area_conhecimento: initialData?.area_conhecimento || "",
    instituicao: initialData?.instituicao || "",
    curriculo_resumido: initialData?.curriculo_resumido || "",
    usuario_id: initialData?.usuario_id || 0,
  });
  const [erro, setErro] = useState("");
  const [carregando, setCarregando] = useState(false);

  function atualizar(chave: string, valor: string) {
    setForm((prev) => ({ ...prev, [chave]: valor }));
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setErro("");

    if (!form.nome || !form.cpf || !form.email) {
      setErro("Nome, CPF e Email são obrigatórios.");
      return;
    }

    setCarregando(true);
    try {
      if (mode === "criar") {
        await onSubmit({
          usuario_id: form.usuario_id,
          nome: form.nome,
          cpf: form.cpf,
          email: form.email,
          titulacao_maxima: form.titulacao_maxima,
          area_conhecimento: form.area_conhecimento,
          instituicao: form.instituicao,
          curriculo_resumido: form.curriculo_resumido,
        });
      } else {
        const payload: AtualizarAvaliadorEntrada = {};
        if (form.nome !== initialData?.nome) payload.nome = form.nome;
        if (form.cpf !== initialData?.cpf) payload.cpf = form.cpf;
        if (form.email !== initialData?.email) payload.email = form.email;
        if (form.titulacao_maxima !== initialData?.titulacao_maxima) payload.titulacao_maxima = form.titulacao_maxima;
        if (form.area_conhecimento !== initialData?.area_conhecimento) payload.area_conhecimento = form.area_conhecimento;
        if (form.instituicao !== initialData?.instituicao) payload.instituicao = form.instituicao;
        if (form.curriculo_resumido !== initialData?.curriculo_resumido) payload.curriculo_resumido = form.curriculo_resumido;
        await onSubmit(payload);
      }
    } catch (err) {
      setErro(err instanceof Error ? err.message : "Erro ao salvar");
    } finally {
      setCarregando(false);
    }
  }

  return (
    <div className="rounded-xl border border-gray-200 bg-white p-8 shadow-sm">
      <h1 className="text-2xl font-bold text-gray-900">
        {mode === "criar" ? "Novo Avaliador" : "Editar Avaliador"}
      </h1>
      <p className="mt-1 text-sm text-gray-600">
        {mode === "criar" ? "Cadastre um novo avaliador no sistema" : "Altere os dados do avaliador"}
      </p>

      <form onSubmit={handleSubmit} className="mt-6 space-y-6">
        <div className="border-b border-gray-100 pb-6">
          <h2 className="mb-4 text-lg font-semibold text-gray-800">Dados Pessoais</h2>

          <div className="grid grid-cols-2 gap-4">
            <CampoFormulario
              id="nome"
              label="Nome Completo"
              placeholder="Nome do avaliador"
              value={form.nome}
              onChange={(e) => atualizar("nome", e.target.value)}
              required
            />
            <CampoFormulario
              id="cpf"
              label="CPF"
              placeholder="000.000.000-00"
              value={form.cpf}
              onChange={(e) => atualizar("cpf", e.target.value)}
              required
            />
          </div>

          <div className="mt-4">
            <CampoFormulario
              id="email"
              label="E-mail"
              type="email"
              placeholder="avaliador@email.com"
              value={form.email}
              onChange={(e) => atualizar("email", e.target.value)}
              required
            />
          </div>
        </div>

        <div className="border-b border-gray-100 pb-6">
          <h2 className="mb-4 text-lg font-semibold text-gray-800">Perfil Técnico</h2>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label htmlFor="titulacao_maxima" className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700">
                Titulação Máxima
              </label>
              <select
                id="titulacao_maxima"
                value={form.titulacao_maxima}
                onChange={(e) => atualizar("titulacao_maxima", e.target.value)}
                className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
              >
                <option value="">Selecione...</option>
                {TITULACAO_OPTIONS.map((opt) => (
                  <option key={opt.value} value={opt.value}>
                    {opt.label}
                  </option>
                ))}
              </select>
            </div>

            <CampoFormulario
              id="area_conhecimento"
              label="Área de Conhecimento"
              placeholder="Ex: Ciência da Computação"
              value={form.area_conhecimento}
              onChange={(e) => atualizar("area_conhecimento", e.target.value)}
            />
          </div>

          <div className="mt-4">
            <CampoFormulario
              id="instituicao"
              label="Instituição"
              placeholder="Ex: Universidade Federal"
              value={form.instituicao}
              onChange={(e) => atualizar("instituicao", e.target.value)}
            />
          </div>

          <div className="mt-4">
            <label htmlFor="curriculo_resumido" className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700">
              Currículo Resumido
            </label>
            <textarea
              id="curriculo_resumido"
              value={form.curriculo_resumido}
              onChange={(e) => atualizar("curriculo_resumido", e.target.value)}
              placeholder="Resumo da experiência e qualificações"
              rows={4}
              className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
            />
          </div>
        </div>

        {erro && (
          <div className="rounded-lg bg-red-50 p-3 text-sm text-red-700">{erro}</div>
        )}

        <div className="flex gap-3">
          <Link
            href={{ pathname: "/avaliadores" }}
            onClick={(e) => { e.preventDefault(); onCancel(); }}
            className="rounded-lg border border-gray-300 px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50"
          >
            Cancelar
          </Link>
          <button
            type="submit"
            disabled={carregando}
            className="flex-1 rounded-lg bg-brand-600 px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-brand-700 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {carregando ? "Salvando..." : mode === "criar" ? "Cadastrar Avaliador" : "Salvar Alterações"}
          </button>
        </div>
      </form>
    </div>
  );
}
