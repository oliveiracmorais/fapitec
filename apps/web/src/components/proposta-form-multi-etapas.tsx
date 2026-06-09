"use client";

import { useState, useCallback } from "react";
import type { DadosProponente, DadosAcademicos, ItemOrcamentario } from "../lib/api-propostas";
import PropostaEtapaDadosPessoais from "./proposta-etapa-dados-pessoais";
import PropostaEtapaFormacaoAcademica from "./proposta-etapa-formacao-academica";
import PropostaEtapaDocumentos from "./proposta-etapa-documentos";
import PropostaEtapaOrcamento from "./proposta-etapa-orcamento";
import PropostaResumo from "./proposta-resumo";
import { validarNome, validarCPF, validarEmail } from "../lib/validacao";

type DocumentoEntry = {
  tipo: string;
  arquivo?: File;
  nome_arquivo?: string;
};

type Props = {
  dadosProponente: DadosProponente;
  dadosAcademicos: DadosAcademicos;
  itensOrcamentarios: ItemOrcamentario[];
  documentos: DocumentoEntry[];
  tiposDocumentosExigidos: string[];
  valorGlobalEdital: number;
  carregando: boolean;
  onProponenteChange: (campo: keyof DadosProponente, valor: string) => void;
  onAcademicoChange: (campo: keyof DadosAcademicos, valor: string | number) => void;
  onOrcamentoChange: (indice: number, campo: keyof ItemOrcamentario, valor: string | number) => void;
  onOrcamentoAdicionar: () => void;
  onOrcamentoRemover: (indice: number) => void;
  onDocumentoAdicionar: (tipo: string, arquivo: File) => void;
  onDocumentoRemover: (tipo: string) => void;
  onSubmit: () => void;
};

const ETAPAS = [
  "Dados Pessoais",
  "Formação Acadêmica",
  "Documentação",
  "Planilha Orçamentária",
  "Revisão e Envio",
];

const STATUS_ETAPA = {
  pendente: "bg-gray-200 text-gray-500",
  atual: "bg-brand-600 text-white",
  concluida: "bg-green-500 text-white",
} as const;

export default function PropostaFormMultiEtapas({
  dadosProponente, dadosAcademicos, itensOrcamentarios, documentos,
  tiposDocumentosExigidos, valorGlobalEdital, carregando,
  onProponenteChange, onAcademicoChange,
  onOrcamentoChange, onOrcamentoAdicionar, onOrcamentoRemover,
  onDocumentoAdicionar, onDocumentoRemover, onSubmit,
}: Props) {
  const [etapa, setEtapa] = useState(0);
  const [erros, setErros] = useState<Record<string, string>>({});
  const [touched, setTouched] = useState<Record<string, boolean>>({});

  const marcarTouched = useCallback((campo: string) => {
    setTouched((prev) => ({ ...prev, [campo]: true }));
  }, []);

  function validarEtapaAtual(): boolean {
    const novosErros: Record<string, string> = {};
    const novosTouched: Record<string, boolean> = {};

    if (etapa === 0) {
      novosErros.nome = validarNome(dadosProponente.nome);
      novosErros.cpf = validarCPF(dadosProponente.cpf);
      novosErros.email = validarEmail(dadosProponente.email);
      novosErros.rg = dadosProponente.rg.trim() ? "" : "RG é obrigatório";
      novosErros.data_nascimento = dadosProponente.data_nascimento ? "" : "Data de nascimento é obrigatória";
      novosErros.telefone = dadosProponente.telefone.trim() ? "" : "Telefone é obrigatório";
      novosErros.endereco = dadosProponente.endereco.trim() ? "" : "Endereço é obrigatório";
      novosErros.genero = dadosProponente.genero.trim() ? "" : "Gênero é obrigatório";
      novosErros.etnia = dadosProponente.etnia.trim() ? "" : "Etnia é obrigatória";
      ["nome", "cpf", "email", "rg", "data_nascimento", "telefone", "endereco", "genero", "etnia"].forEach(
        (c) => (novosTouched[c] = true)
      );
    } else if (etapa === 1) {
      novosErros.maior_titulacao = dadosAcademicos.maior_titulacao ? "" : "Maior titulação é obrigatória";
      novosErros.curso = dadosAcademicos.curso.trim() ? "" : "Curso é obrigatório";
      novosErros.instituicao = dadosAcademicos.instituicao.trim() ? "" : "Instituição é obrigatória";
      novosErros.ano_conclusao = dadosAcademicos.ano_conclusao > 0 ? "" : "Ano de conclusão é obrigatório";
      novosErros.area_conhecimento = dadosAcademicos.area_conhecimento.trim() ? "" : "Área de conhecimento é obrigatória";
      ["maior_titulacao", "curso", "instituicao", "ano_conclusao", "area_conhecimento"].forEach(
        (c) => (novosTouched[c] = true)
      );
    } else if (etapa === 2) {
      if (tiposDocumentosExigidos.length > 0) {
        for (const tipo of tiposDocumentosExigidos) {
          const doc = documentos.find((d) => d.tipo === tipo);
          if (!doc?.arquivo && !doc?.nome_arquivo) {
            novosErros[`doc_${tipo}`] = `Documento "${tipo}" é obrigatório`;
          }
        }
      }
    } else if (etapa === 3) {
      if (itensOrcamentarios.length === 0) {
        novosErros.orcamento = "Adicione ao menos um item orçamentário";
      }
    }

    setErros(novosErros);
    setTouched((prev) => ({ ...prev, ...novosTouched }));

    return Object.values(novosErros).every((e) => !e);
  }

  function avancar() {
    if (!validarEtapaAtual()) return;
    setErros({});
    setEtapa((e) => Math.min(e + 1, ETAPAS.length - 1));
  }

  function voltar() {
    setErros({});
    setEtapa((e) => Math.max(e - 1, 0));
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2">
        {ETAPAS.map((nome, i) => (
          <div key={nome} className="flex items-center gap-2">
            <span
              className={`flex h-8 w-8 items-center justify-center rounded-full text-xs font-bold ${
                i < etapa ? STATUS_ETAPA.concluida : i === etapa ? STATUS_ETAPA.atual : STATUS_ETAPA.pendente
              }`}
            >
              {i < etapa ? "✓" : i + 1}
            </span>
            <span className={`hidden text-xs sm:inline ${i === etapa ? "font-medium text-gray-900" : "text-gray-500"}`}>
              {nome}
            </span>
            {i < ETAPAS.length - 1 && <div className="h-px w-6 bg-gray-300 sm:w-12" />}
          </div>
        ))}
      </div>

      {etapa === 0 && (
        <PropostaEtapaDadosPessoais
          dados={dadosProponente}
          onChange={onProponenteChange}
          erros={erros}
          touched={touched}
          onBlur={marcarTouched}
        />
      )}
      {etapa === 1 && (
        <PropostaEtapaFormacaoAcademica
          dados={dadosAcademicos}
          onChange={onAcademicoChange}
          erros={erros}
          touched={touched}
          onBlur={marcarTouched}
        />
      )}
      {etapa === 2 && (
        <PropostaEtapaDocumentos
          documentos={documentos}
          tiposExigidos={tiposDocumentosExigidos}
          onAdicionar={onDocumentoAdicionar}
          onRemover={onDocumentoRemover}
        />
      )}
      {etapa === 3 && (
        <PropostaEtapaOrcamento
          itens={itensOrcamentarios}
          valorGlobalEdital={valorGlobalEdital}
          onChangeItem={onOrcamentoChange}
          onAdicionar={onOrcamentoAdicionar}
          onRemover={onOrcamentoRemover}
        />
      )}
      {etapa === 4 && (
        <PropostaResumo
          dadosProponente={dadosProponente}
          dadosAcademicos={dadosAcademicos}
          itensOrcamentarios={itensOrcamentarios}
          valorGlobalEdital={valorGlobalEdital}
        />
      )}

      {erros.orcamento && (
        <p className="text-xs text-red-700">{erros.orcamento}</p>
      )}

      <div className="flex items-center justify-between border-t border-gray-200 pt-4">
        <button
          type="button"
          onClick={voltar}
          disabled={etapa === 0}
          className="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:opacity-50"
        >
          ← Voltar
        </button>

        {etapa < ETAPAS.length - 1 ? (
          <button
            type="button"
            onClick={avancar}
            className="rounded-lg bg-brand-600 px-6 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-700"
          >
            Avançar →
          </button>
        ) : (
          <button
            type="button"
            onClick={onSubmit}
            disabled={carregando}
            className="rounded-lg bg-green-600 px-6 py-2 text-sm font-medium text-white transition-colors hover:bg-green-700 disabled:opacity-50"
          >
            {carregando ? "Enviando..." : "Submeter Proposta"}
          </button>
        )}
      </div>
    </div>
  );
}
