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
};

const TIPOS_ACEITOS = ["pdf", "doc", "docx", "jpg", "jpeg", "png"];

export default function PropostaEtapaDocumentos({ documentos, tiposExigidos, onAdicionar, onRemover }: Props) {
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
              {tipo}
            </label>
            <PropostaUpload
              tiposAceitos={TIPOS_ACEITOS}
              onUpload={(arquivo) => onAdicionar(tipo, arquivo)}
              onRemover={() => onRemover(tipo)}
              arquivo={doc?.arquivo}
              nomeArquivo={doc?.nome_arquivo}
            />
          </div>
        );
      })}
    </div>
  );
}
