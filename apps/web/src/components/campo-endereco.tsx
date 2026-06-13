"use client";

import { useState, useCallback } from "react";

type CampoEnderecoProps = {
  cep: string;
  logradouro: string;
  numero: string;
  complemento: string;
  bairro: string;
  cidade: string;
  uf: string;
  onChange: (campo: string, valor: string) => void;
  onBlur: (campo: string) => void;
  erros: Record<string, string>;
  touched: Record<string, boolean>;
};

const UFS = [
  "AC", "AL", "AP", "AM", "BA", "CE", "DF", "ES", "GO",
  "MA", "MT", "MS", "MG", "PA", "PB", "PR", "PE", "PI",
  "RJ", "RN", "RS", "RO", "RR", "SC", "SP", "SE", "TO",
];

function apenasDigitos(v: string) {
  return v.replace(/\D/g, "").slice(0, 8);
}

export default function CampoEndereco({
  cep, logradouro, numero, complemento, bairro, cidade, uf,
  onChange, onBlur, erros, touched,
}: CampoEnderecoProps) {
  const [buscandoCEP, setBuscandoCEP] = useState(false);

  const buscarCEP = useCallback(async (cepCompleto: string) => {
    const digitos = apenasDigitos(cepCompleto);
    if (digitos.length !== 8) return;
    setBuscandoCEP(true);
    try {
      const res = await fetch(`https://viacep.com.br/ws/${digitos}/json/`);
      if (!res.ok) return;
      const data = await res.json();
      if (data.erro) return;
      onChange("logradouro", data.logradouro || "");
      onChange("bairro", data.bairro || "");
      onChange("cidade", data.localidade || "");
      onChange("uf", data.uf || "");
      if (!data.logradouro) {
        onChange("logradouro", "");
      }
    } catch {
      // Silencia erro de rede
    } finally {
      setBuscandoCEP(false);
    }
  }, [onChange]);

  function handleCEPChange(valor: string) {
    const digitos = apenasDigitos(valor);
    onChange("cep", digitos);
    if (digitos.length === 8) {
      buscarCEP(digitos);
    }
  }

  function borda(campo: string) {
    if (touched[campo] && erros[campo])
      return "border-red-400 focus:border-red-500 focus:ring-red-500";
    if (touched[campo] && !erros[campo])
      return "border-green-400 focus:border-green-600 focus:ring-green-600";
    return "border-gray-300 focus:border-blue-600 focus:ring-blue-600";
  }

  function campo(label: string, icone: string, campo: string, valor: string, placeholder = "",小型 = false) {
    return (
      <div className={小型 ? "" : "sm:col-span-2"}>
        <label htmlFor={campo} className="flex items-center gap-1.5 text-sm font-medium text-gray-700">
          <span className="text-base">{icone}</span> {label}
        </label>
        <input
          id={campo}
          value={valor}
          onChange={(e) => onChange(campo, e.target.value)}
          onBlur={() => onBlur(campo)}
          placeholder={placeholder}
          className={`mt-1 block w-full rounded-lg border px-4 py-2.5 text-sm outline-none transition-colors ${borda(campo)}`}
        />
        {touched[campo] && erros[campo] && (
          <p className="mt-1 text-xs text-red-700">{erros[campo]}</p>
        )}
      </div>
    );
  }

  return (
    <div>
      <label className="flex items-center gap-1.5 text-sm font-medium text-gray-700 mb-2">
        <span className="text-base">📍</span> Endereço
      </label>
      <div className="space-y-3 rounded-lg border border-gray-200 bg-gray-50 p-4">
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
          <div>
            <label htmlFor="cep" className="flex items-center gap-1.5 text-sm font-medium text-gray-700">
              <span className="text-base">📮</span> CEP
            </label>
            <div className="relative mt-1">
              <input
                id="cep"
                value={cep}
                onChange={(e) => handleCEPChange(e.target.value)}
                onBlur={() => onBlur("cep")}
                placeholder="49000-000"
                className={`block w-full rounded-lg border px-4 py-2.5 pr-10 text-sm outline-none transition-colors ${borda("cep")}`}
              />
              {buscandoCEP && (
                <span className="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-gray-400">
                  buscando...
                </span>
              )}
            </div>
            {touched.cep && erros.cep && (
              <p className="mt-1 text-xs text-red-700">{erros.cep}</p>
            )}
          </div>
          <div className="sm:col-span-2">
            {campo("Logradouro", "🏛️", "logradouro", logradouro, "Rua, Avenida...")}
          </div>
        </div>

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
          <div>
            <label htmlFor="numero" className="flex items-center gap-1.5 text-sm font-medium text-gray-700">
              <span className="text-base">🔢</span> Número
            </label>
            <input
              id="numero"
              value={numero}
              onChange={(e) => onChange("numero", e.target.value)}
              onBlur={() => onBlur("numero")}
              placeholder="S/N"
              className={`mt-1 block w-full rounded-lg border px-4 py-2.5 text-sm outline-none transition-colors ${borda("numero")}`}
            />
            {touched.numero && erros.numero && (
              <p className="mt-1 text-xs text-red-700">{erros.numero}</p>
            )}
          </div>
          <div className="sm:col-span-2">
            <label htmlFor="complemento" className="flex items-center gap-1.5 text-sm font-medium text-gray-700">
              <span className="text-base">➕</span> Complemento
            </label>
            <input
              id="complemento"
              value={complemento}
              onChange={(e) => onChange("complemento", e.target.value)}
              onBlur={() => onBlur("complemento")}
              placeholder="Apto, Bloco, Sala..."
              className={`mt-1 block w-full rounded-lg border px-4 py-2.5 text-sm outline-none transition-colors ${borda("complemento")}`}
            />
            {touched.complemento && erros.complemento && (
              <p className="mt-1 text-xs text-red-700">{erros.complemento}</p>
            )}
          </div>
        </div>

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
          <div>
            <label htmlFor="bairro" className="flex items-center gap-1.5 text-sm font-medium text-gray-700">
              <span className="text-base">🏘️</span> Bairro
            </label>
            <input
              id="bairro"
              value={bairro}
              onChange={(e) => onChange("bairro", e.target.value)}
              onBlur={() => onBlur("bairro")}
              placeholder="Centro"
              className={`mt-1 block w-full rounded-lg border px-4 py-2.5 text-sm outline-none transition-colors ${borda("bairro")}`}
            />
            {touched.bairro && erros.bairro && (
              <p className="mt-1 text-xs text-red-700">{erros.bairro}</p>
            )}
          </div>
          <div>
            <label htmlFor="cidade" className="flex items-center gap-1.5 text-sm font-medium text-gray-700">
              <span className="text-base">🌆</span> Cidade
            </label>
            <input
              id="cidade"
              value={cidade}
              onChange={(e) => onChange("cidade", e.target.value)}
              onBlur={() => onBlur("cidade")}
              placeholder="Aracaju"
              className={`mt-1 block w-full rounded-lg border px-4 py-2.5 text-sm outline-none transition-colors ${borda("cidade")}`}
            />
            {touched.cidade && erros.cidade && (
              <p className="mt-1 text-xs text-red-700">{erros.cidade}</p>
            )}
          </div>
          <div>
            <label htmlFor="uf" className="flex items-center gap-1.5 text-sm font-medium text-gray-700">
              <span className="text-base">🗺️</span> UF
            </label>
            <select
              id="uf"
              value={uf}
              onChange={(e) => onChange("uf", e.target.value)}
              onBlur={() => onBlur("uf")}
              className={`mt-1 block w-full rounded-lg border px-4 py-2.5 text-sm outline-none transition-colors ${borda("uf")}`}
            >
              <option value="">Selecione...</option>
              {UFS.map((s) => (
                <option key={s} value={s}>{s}</option>
              ))}
            </select>
            {touched.uf && erros.uf && (
              <p className="mt-1 text-xs text-red-700">{erros.uf}</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
