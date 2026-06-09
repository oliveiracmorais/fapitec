"use client";

import PropostaItemOrcamento from "./proposta-item-orcamento";
import { formatarBRL } from "../lib/api-propostas";
import type { ItemOrcamentario } from "../lib/api-propostas";

type Props = {
  itens: ItemOrcamentario[];
  valorGlobalEdital: number;
  onChangeItem: (indice: number, campo: keyof ItemOrcamentario, valor: string | number) => void;
  onAdicionar: () => void;
  onRemover: (indice: number) => void;
};

export default function PropostaEtapaOrcamento({ itens, valorGlobalEdital, onChangeItem, onAdicionar, onRemover }: Props) {
  const valorTotal = itens.reduce((acc, item) => acc + item.valor_total, 0);
  const excede = valorGlobalEdital > 0 && valorTotal > valorGlobalEdital;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold text-gray-900">Planilha Orçamentária</h2>
        <button
          type="button"
          onClick={onAdicionar}
          className="rounded-lg bg-brand-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-brand-700"
        >
          + Adicionar Item
        </button>
      </div>

      {itens.length === 0 ? (
        <div className="rounded-lg border border-dashed border-gray-300 p-8 text-center">
          <p className="text-sm text-gray-500">
            Nenhum item orçamentário adicionado.
          </p>
          <p className="text-xs text-gray-400">
            Clique em "Adicionar Item" para começar.
          </p>
        </div>
      ) : (
        <div className="space-y-2">
          {itens.map((item, i) => (
            <PropostaItemOrcamento
              key={i}
              item={item}
              indice={i}
              onChange={onChangeItem}
              onRemover={onRemover}
            />
          ))}
        </div>
      )}

      <div className={`rounded-lg border p-4 ${excede ? "border-red-300 bg-red-50" : "border-gray-200 bg-gray-50"}`}>
        <div className="flex items-center justify-between text-sm">
          <span className="font-medium text-gray-700">Valor Total Solicitado</span>
          <span className={`text-lg font-bold ${excede ? "text-red-700" : "text-gray-900"}`}>
            {formatarBRL(valorTotal)}
          </span>
        </div>
        {valorGlobalEdital > 0 && (
          <div className="mt-1 flex items-center justify-between text-xs">
            <span className="text-gray-500">Valor Global do Edital</span>
            <span className="text-gray-600">{formatarBRL(valorGlobalEdital)}</span>
          </div>
        )}
        {excede && (
          <p className="mt-2 text-xs font-medium text-red-700">
            O valor total excede o valor global do edital. Reduza os itens ou ajuste os valores.
          </p>
        )}
      </div>
    </div>
  );
}
