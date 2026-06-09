"use client";

import { useState, useRef } from "react";
import { validarArquivo } from "../lib/validacao";

type Props = {
  tiposAceitos: string[];
  onUpload: (arquivo: File) => void;
  onRemover: () => void;
  arquivo?: File;
  nomeArquivo?: string;
};

export default function PropostaUpload({ tiposAceitos, onUpload, onRemover, arquivo, nomeArquivo }: Props) {
  const [erro, setErro] = useState("");
  const [arrastando, setArrastando] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  function handleFile(file: File) {
    const msg = validarArquivo(file, tiposAceitos);
    if (msg) {
      setErro(msg);
      return;
    }
    setErro("");
    onUpload(file);
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault();
    setArrastando(false);
    const file = e.dataTransfer.files[0];
    if (file) handleFile(file);
  }

  return (
    <div>
      <div
        onDragOver={(e) => { e.preventDefault(); setArrastando(true); }}
        onDragLeave={() => setArrastando(false)}
        onDrop={handleDrop}
        onClick={() => inputRef.current?.click()}
        className={`flex cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed p-6 transition-colors ${
          arrastando ? "border-brand-500 bg-brand-50" : "border-gray-300 hover:border-gray-400"
        }`}
      >
        {arquivo || nomeArquivo ? (
          <div className="flex items-center gap-2 text-sm">
            <svg className="h-5 w-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span className="text-gray-700">{nomeArquivo || arquivo?.name}</span>
            <button
              type="button"
              onClick={(e) => { e.stopPropagation(); onRemover(); }}
              className="text-red-500 hover:text-red-700"
              aria-label="Remover arquivo"
            >
              <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        ) : (
          <>
            <svg className="mb-2 h-8 w-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
            <p className="text-sm text-gray-600">
              <span className="font-medium text-brand-600">Clique</span> ou arraste o arquivo aqui
            </p>
            <p className="mt-0.5 text-xs text-gray-500">
              Formatos: {tiposAceitos.join(", ")} · Máx. 5MB
            </p>
          </>
        )}
      </div>
      <input
        ref={inputRef}
        type="file"
        accept={tiposAceitos.map((t) => `.${t}`).join(",")}
        className="hidden"
        onChange={(e) => { const f = e.target.files?.[0]; if (f) handleFile(f); }}
      />
      {erro && <p className="mt-1 text-xs text-red-700">{erro}</p>}
    </div>
  );
}
