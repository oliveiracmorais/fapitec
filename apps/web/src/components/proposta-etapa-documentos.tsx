"use client";

import PropostaUpload from "./proposta-upload";

type DocumentoEntry = {
  tipo: string;
  arquivo?: File;
  nome_arquivo?: string;
};

type Props = {
  documentos: DocumentoEntry[];
  tiposExigidos: string[];
  onAdicionar: (tipo: string, arquivo: File) => void;
  onRemover: (tipo: string) => void;
  erros?: Record<string, string>;
  touched?: Record<string, boolean>;
};

const TIPOS_ACEITOS = ["pdf", "doc", "docx", "jpg", "jpeg", "png"];

const ROTULOS_DOCUMENTOS: Record<string, string> = {
  comprovante_conclusao: "Comprovante de Conclusão",
  comprovante_residencia: "Comprovante de Residência",
  comprovante_inscricao: "Comprovante de Inscrição",
  identidade: "Documento de Identidade",
  cpf: "Cadastro de Pessoa Física (CPF)",
  currículo: "Currículo Lattes",
  plano_trabalho: "Plano de Trabalho",
  declaracao_vinculo: "Declaração de Vínculo Institucional",
};

export default function PropostaEtapaDocumentos({ documentos, tiposExigidos, onAdicionar, onRemover, erros, touched }: Props) {
  if (tiposExigidos.length === 0) {
    return (
      <div className="space-y-4">
        <h2 className="text-lg font-semibold text-gray-900">Documentação</h2>
        <p className="text-sm text-gray-500">
          Nenhum documento obrigatório exigido para este edital.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <h2 className="text-lg font-semibold text-gray-900">Documentação</h2>
      <p className="text-sm text-gray-500">
        Envie os documentos obrigatórios exigidos pelo edital.
      </p>

      {tiposExigidos.map((tipo) => {
        const doc = documentos.find((d) => d.tipo === tipo);
        return (
          <div key={tipo}>
            <label className="mb-1 block text-sm font-medium text-gray-700">
              {ROTULOS_DOCUMENTOS[tipo] || tipo}
            </label>
            <PropostaUpload
              tiposAceitos={TIPOS_ACEITOS}
              onUpload={(arquivo) => onAdicionar(tipo, arquivo)}
              onRemover={() => onRemover(tipo)}
              arquivo={doc?.arquivo}
              nomeArquivo={doc?.nome_arquivo}
            />
            {touched?.[`doc_${tipo}`] && erros?.[`doc_${tipo}`] && (
              <p className="mt-1 text-xs text-red-700">{erros[`doc_${tipo}`]}</p>
            )}
          </div>
        );
      })}
    </div>
  );
}
