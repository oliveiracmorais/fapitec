"use client";

import type { ItemOrcamentario } from "../lib/api-propostas";
import { formatarBRL } from "../lib/api-propostas";

type Props = {
  item: ItemOrcamentario;
  indice: number;
  onChange: (indice: number, campo: keyof ItemOrcamentario, valor: string | number) => void;
  onRemover: (indice: number) => void;
};

export default function PropostaItemOrcamento({ item, indice, onChange, onRemover }: Props) {
  return (
    <div className="flex flex-wrap items-end gap-2 rounded-lg border border-gray-200 bg-gray-50 p-3">
      <div className="flex-1">
        <label className="block text-xs font-medium text-gray-600">Descrição</label>
        <input
          type="text"
          value={item.descricao}
          onChange={(e) => onChange(indice, "descricao", e.target.value)}
          className="mt-0.5 w-full rounded border border-gray-300 px-2 py-1.5 text-sm outline-none focus:border-brand-600"
          placeholder="Ex: Materiais de laboratório"
        />
      </div>
      <div>
        <label className="block text-xs font-medium text-gray-600">Tipo</label>
        <select
          value={item.tipo}
          onChange={(e) => onChange(indice, "tipo", e.target.value)}
          className="mt-0.5 rounded border border-gray-300 px-2 py-1.5 text-sm outline-none focus:border-brand-600"
        >
          <option value="consumo">Consumo</option>
          <option value="permanente">Permanente</option>
        </select>
      </div>
      <div>
        <label className="block text-xs font-medium text-gray-600">Qtd</label>
        <input
          type="number"
          min="1"
          value={item.quantidade || ""}
          onChange={(e) => onChange(indice, "quantidade", Number(e.target.value))}
          className="mt-0.5 w-16 rounded border border-gray-300 px-2 py-1.5 text-sm outline-none focus:border-brand-600"
        />
      </div>
      <div>
        <label className="block text-xs font-medium text-gray-600">Valor Unit.</label>
        <input
          type="number"
          min="0"
          step="0.01"
          value={item.valor_unitario > 0 ? item.valor_unitario / 100 : ""}
          onChange={(e) => onChange(indice, "valor_unitario", Math.round(Number(e.target.value) * 100))}
          className="mt-0.5 w-24 rounded border border-gray-300 px-2 py-1.5 text-sm outline-none focus:border-brand-600"
        />
      </div>
      <div className="text-right">
        <label className="block text-xs font-medium text-gray-600">Total</label>
        <p className="mt-0.5 text-sm font-semibold text-gray-900">
          {formatarBRL(item.valor_total)}
        </p>
      </div>
      <button
        type="button"
        onClick={() => onRemover(indice)}
        className="mb-0.5 rounded p-1 text-red-500 hover:bg-red-50"
        aria-label="Remover item"
      >
        <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
      </button>
    </div>
  );
}
