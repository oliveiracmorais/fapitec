"use client";

import type { DadosAcademicos } from "../lib/api-propostas";
import CampoFormulario from "./campo-formulario";

type Props = {
  dados: DadosAcademicos;
  onChange: (campo: keyof DadosAcademicos, valor: string | number) => void;
  erros: Record<string, string>;
  touched: Record<string, boolean>;
  onBlur: (campo: string) => void;
};

const TITULACOES = [
  "Ensino Médio",
  "Graduação",
  "Especialização",
  "Mestrado",
  "Doutorado",
  "Pós-Doutorado",
];

export default function PropostaEtapaFormacaoAcademica({ dados, onChange, erros, touched, onBlur }: Props) {
  return (
    <div className="space-y-4">
      <h2 className="text-lg font-semibold text-gray-900">Formação Acadêmica</h2>

      <div>
        <label htmlFor="maior_titulacao" className="flex items-center gap-1.5 text-sm font-medium text-gray-700">
          <span className="text-base">🎓</span> Maior Titulação
        </label>
        <select
          id="maior_titulacao"
          value={dados.maior_titulacao}
          onChange={(e) => onChange("maior_titulacao", e.target.value)}
          onBlur={() => onBlur("maior_titulacao")}
          className={`mt-1 block w-full rounded-lg border px-4 py-2.5 text-sm outline-none transition-colors ${
            touched.maior_titulacao && erros.maior_titulacao
              ? "border-red-400 focus:border-red-500"
              : "border-gray-300 focus:border-blue-600 focus:ring-blue-600"
          }`}
        >
          <option value="">Selecione...</option>
          {TITULACOES.map((t) => (
            <option key={t} value={t}>{t}</option>
          ))}
        </select>
        {touched.maior_titulacao && erros.maior_titulacao && (
          <p className="mt-1 text-xs text-red-700">{erros.maior_titulacao}</p>
        )}
      </div>

      <CampoFormulario
        id="curso"
        label="Curso"
        value={dados.curso}
        onChange={(e) => onChange("curso", e.target.value)}
        erro={erros.curso}
        touched={touched.curso}
        onBlur={() => onBlur("curso")}
        placeholder="Nome do curso"
        icone="📚"
      />

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <CampoFormulario
          id="instituicao"
          label="Instituição"
          value={dados.instituicao}
          onChange={(e) => onChange("instituicao", e.target.value)}
          erro={erros.instituicao}
          touched={touched.instituicao}
          onBlur={() => onBlur("instituicao")}
          placeholder="Nome da IES"
          icone="🏛️"
        />
        <CampoFormulario
          id="ano_conclusao"
          label="Ano de Conclusão"
          type="number"
          value={dados.ano_conclusao ? String(dados.ano_conclusao) : ""}
          onChange={(e) => onChange("ano_conclusao", Number(e.target.value))}
          erro={erros.ano_conclusao}
          touched={touched.ano_conclusao}
          onBlur={() => onBlur("ano_conclusao")}
          placeholder="2024"
          icone="📅"
        />
      </div>

      <CampoFormulario
        id="area_conhecimento"
        label="Área de Conhecimento"
        value={dados.area_conhecimento}
        onChange={(e) => onChange("area_conhecimento", e.target.value)}
        erro={erros.area_conhecimento}
        touched={touched.area_conhecimento}
        onBlur={() => onBlur("area_conhecimento")}
        placeholder="Ex: Ciências Biológicas"
        icone="🔬"
      />
    </div>
  );
}
