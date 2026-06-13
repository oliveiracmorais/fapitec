"use client";

import { useState, useEffect, FormEvent } from "react";
import { atribuirEdital } from "../../lib/api-avaliadores";

type EditalResumo = {
  id: number;
  nome: string;
};

type Props = {
  isOpen: boolean;
  onClose: () => void;
  avaliadorId: number;
  avaliadorNome: string;
  onAtribuido: () => void;
};

export default function AtribuirEditalDialog({
  isOpen,
  onClose,
  avaliadorId,
  avaliadorNome,
  onAtribuido,
}: Props) {
  const [editais, setEditais] = useState<EditalResumo[]>([]);
  const [editalId, setEditalId] = useState("");
  const [dataInicio, setDataInicio] = useState("");
  const [dataFim, setDataFim] = useState("");
  const [carregando, setCarregando] = useState(false);
  const [carregandoEditais, setCarregandoEditais] = useState(true);
  const [erro, setErro] = useState("");

  useEffect(() => {
    if (!isOpen) return;

    setCarregandoEditais(true);
    fetch("/api/v1/editais")
      .then((res) => res.json())
      .then((data) => setEditais(data.editais ?? []))
      .catch(() => setEditais([]))
      .finally(() => setCarregandoEditais(false));
  }, [isOpen]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setErro("");

    if (!editalId || !dataInicio || !dataFim) {
      setErro("Preencha todos os campos.");
      return;
    }

    if (dataInicio > dataFim) {
      setErro("Data de início não pode ser posterior à data de fim.");
      return;
    }

    setCarregando(true);
    try {
      await atribuirEdital(avaliadorId, {
        edital_id: Number(editalId),
        data_inicio: dataInicio,
        data_fim: dataFim,
      });
      onAtribuido();
      onClose();
    } catch (err) {
      setErro(err instanceof Error ? err.message : "Erro ao atribuir edital");
    } finally {
      setCarregando(false);
    }
  }

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="w-full max-w-lg rounded-xl bg-white p-6 shadow-xl">
        <h2 className="text-lg font-bold text-gray-900">Atribuir Edital</h2>
        <p className="mt-1 text-sm text-gray-600">
          Atribuir edital para <span className="font-medium">{avaliadorNome}</span>
        </p>

        <form onSubmit={handleSubmit} className="mt-4 space-y-4">
          <div>
            <label htmlFor="edital" className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700">
              Edital
            </label>
            {carregandoEditais ? (
              <p className="text-sm text-gray-500">Carregando editais...</p>
            ) : (
              <select
                id="edital"
                value={editalId}
                onChange={(e) => setEditalId(e.target.value)}
                required
                className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
              >
                <option value="">Selecione um edital</option>
                {editais.map((edital) => (
                  <option key={edital.id} value={edital.id}>
                    {edital.nome}
                  </option>
                ))}
              </select>
            )}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label htmlFor="dataInicio" className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700">
                Data de Início
              </label>
              <input
                id="dataInicio"
                type="date"
                value={dataInicio}
                onChange={(e) => setDataInicio(e.target.value)}
                required
                className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
              />
            </div>
            <div>
              <label htmlFor="dataFim" className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700">
                Data de Fim
              </label>
              <input
                id="dataFim"
                type="date"
                value={dataFim}
                onChange={(e) => setDataFim(e.target.value)}
                required
                className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
              />
            </div>
          </div>

          {erro && (
            <div className="rounded-lg bg-red-50 p-3 text-sm text-red-700">{erro}</div>
          )}

          <div className="flex gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="rounded-lg border border-gray-300 px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50"
            >
              Cancelar
            </button>
            <button
              type="submit"
              disabled={carregando || carregandoEditais}
              className="flex-1 rounded-lg bg-brand-600 px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-brand-700 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {carregando ? "Atribuindo..." : "Atribuir"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
