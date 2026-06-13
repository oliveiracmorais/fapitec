"use client";

import { useState } from "react";

type FormParecerProps = {
  propostaId: number;
  avaliadorId: number;
  onSubmit: (data: { nota: number; parecer_texto: string }) => Promise<void>;
  onCancel: () => void;
};

export default function FormParecer({
  propostaId,
  avaliadorId,
  onSubmit,
  onCancel,
}: FormParecerProps) {
  const [nota, setNota] = useState(0);
  const [parecerTexto, setParecerTexto] = useState("");
  const [enviando, setEnviando] = useState(false);
  const [erro, setErro] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setErro("");

    if (nota < 0 || nota > 100) {
      setErro("A nota deve estar entre 0 e 100.");
      return;
    }
    if (!parecerTexto.trim()) {
      setErro("O parecer textual é obrigatório.");
      return;
    }

    setEnviando(true);
    try {
      await onSubmit({ nota, parecer_texto: parecerTexto });
    } catch (err) {
      setErro(err instanceof Error ? err.message : "Erro ao salvar parecer.");
    } finally {
      setEnviando(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <h2 className="text-lg font-semibold text-gray-900">Emitir Parecer</h2>

      <div>
        <label className="mb-1 block text-sm font-medium text-gray-900">
          Nota (0–100)
        </label>
        <div className="flex items-center gap-4">
          <input
            type="range"
            min={0}
            max={100}
            value={nota}
            onChange={(e) => setNota(Number(e.target.value))}
            className="h-2 w-full cursor-pointer appearance-none rounded-lg bg-gray-200 accent-brand-600"
          />
          <span className="min-w-[3ch] text-right text-lg font-bold text-brand-600">
            {nota}
          </span>
        </div>
        <div className="mt-1 flex justify-between text-xs text-gray-400">
          <span>0</span>
          <span>100</span>
        </div>
      </div>

      <div>
        <label className="mb-1 block text-sm font-medium text-gray-900">
          Parecer Técnico
        </label>
        <textarea
          value={parecerTexto}
          onChange={(e) => setParecerTexto(e.target.value)}
          rows={6}
          placeholder="Descreva sua avaliação técnica da proposta..."
          className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm outline-none focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
        />
      </div>

      {erro && (
        <div className="rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">
          {erro}
        </div>
      )}

      <div className="flex items-center gap-3">
        <button
          type="submit"
          disabled={enviando}
          className="rounded-lg bg-brand-600 px-6 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-700 disabled:opacity-50"
        >
          {enviando ? "Salvando..." : "Salvar Parecer"}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="rounded-lg border border-gray-300 px-6 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
        >
          Cancelar
        </button>
      </div>
    </form>
  );
}
