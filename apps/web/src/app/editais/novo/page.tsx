"use client";

import { useState, FormEvent, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Image from "next/image";
import CampoFormulario from "../../../components/campo-formulario";

export default function NovoEditalPage() {
  const router = useRouter();
  const [form, setForm] = useState({
    nome: "",
    descricao: "",
    dataInicio: "",
    dataFim: "",
    tipoChamada: "",
  });
  const [erro, setErro] = useState("");
  const [carregando, setCarregando] = useState(false);

  useEffect(() => {
    const sessao =
      typeof window !== "undefined" && localStorage.getItem("sessao");
    if (!sessao) router.replace("/");
  }, [router]);

  function atualizar(chave: string, valor: string) {
    setForm((prev) => ({ ...prev, [chave]: valor }));
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setErro("");

    if (form.dataInicio && form.dataFim && form.dataInicio > form.dataFim) {
      setErro("Data de início não pode ser posterior à data de fim.");
      return;
    }

    setCarregando(true);

    try {
      const res = await fetch("/api/v1/editais", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          nome: form.nome,
          descricao: form.descricao,
          data_inicio: form.dataInicio,
          data_fim: form.dataFim,
          tipo_chamada: form.tipoChamada,
        }),
      });

      const data = await res.json();

      if (!res.ok) {
        setErro(data.erro || "Erro ao criar edital");
        return;
      }

      router.push("/editais");
    } catch {
      setErro("Erro de conexão com o servidor");
    } finally {
      setCarregando(false);
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="border-b border-gray-200 bg-white shadow-sm">
        <div className="mx-auto flex max-w-6xl items-center px-4 py-3">
          <Link href="/dashboard">
            <Image
              src="/logo-2.png"
              alt="FAPITEC-SE"
              width={140}
              height={40}
              className="h-9 w-auto"
            />
          </Link>
        </div>
      </header>

      <main className="mx-auto max-w-2xl px-4 py-8">
        <Link
          href="/editais"
          className="mb-4 inline-block text-sm text-gray-600 hover:text-gray-700"
        >
          ← Voltar para editais
        </Link>

        <div className="rounded-xl border border-gray-200 bg-white p-8 shadow-sm">
          <h1 className="text-2xl font-bold text-gray-900">Novo Edital</h1>
          <p className="mt-1 text-sm text-gray-600">
            Preencha os dados para criar um novo edital
          </p>

          <form onSubmit={handleSubmit} className="mt-6 space-y-5">
            <CampoFormulario
              id="nome"
              label="Nome do Edital"
              placeholder="Ex: Edital APQ 2026"
              value={form.nome}
              onChange={(e) => atualizar("nome", e.target.value)}
              icone="📋"
              required
            />

            <div>
              <label
                htmlFor="descricao"
                className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700"
              >
                <span>📝</span> Descrição
              </label>
              <textarea
                id="descricao"
                value={form.descricao}
                onChange={(e) => atualizar("descricao", e.target.value)}
                placeholder="Descrição do edital"
                required
                rows={3}
                className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label
                  htmlFor="dataInicio"
                  className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700"
                >
                  <span>📅</span> Data de Início
                </label>
                <input
                  id="dataInicio"
                  type="date"
                  value={form.dataInicio}
                  onChange={(e) => atualizar("dataInicio", e.target.value)}
                  required
                  className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
                />
              </div>
              <div>
                <label
                  htmlFor="dataFim"
                  className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700"
                >
                  <span>📅</span> Data de Fim
                </label>
                <input
                  id="dataFim"
                  type="date"
                  value={form.dataFim}
                  onChange={(e) => atualizar("dataFim", e.target.value)}
                  required
                  className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
                />
              </div>
            </div>

            <CampoFormulario
              id="tipoChamada"
              label="Tipo de Chamada"
              placeholder="Ex: APQ, ARC"
              value={form.tipoChamada}
              onChange={(e) => atualizar("tipoChamada", e.target.value)}
              icone="🏷️"
            />

            {erro && (
              <div className="rounded-lg bg-red-50 p-3 text-sm text-red-700">
                {erro}
              </div>
            )}

            <div className="flex gap-3">
              <Link
                href="/editais"
                className="rounded-lg border border-gray-300 px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50"
              >
                Cancelar
              </Link>
              <button
                type="submit"
                disabled={carregando}
                className="flex-1 rounded-lg bg-brand-600 px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-brand-700 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {carregando ? "Criando..." : "Criar Edital"}
              </button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}
