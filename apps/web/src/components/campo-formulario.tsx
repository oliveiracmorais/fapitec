"use client";

import { useState, InputHTMLAttributes } from "react";

type CampoFormularioProps = InputHTMLAttributes<HTMLInputElement> & {
  id: string;
  label: string;
  erro?: string;
  touched?: boolean;
  icone?: string;
  tipoSenha?: boolean;
};

export default function CampoFormulario({
  id,
  label,
  erro,
  touched,
  icone,
  tipoSenha,
  className,
  ...props
}: CampoFormularioProps) {
  const [mostrarSenha, setMostrarSenha] = useState(false);
  const temErro = touched && erro;
  const valido = touched && !erro && props.value && String(props.value).length > 0;

  const borda = temErro
    ? "border-red-400 focus:border-red-500 focus:ring-red-500"
    : valido
      ? "border-green-400 focus:border-green-600 focus:ring-green-600"
      : "border-gray-300 focus:border-blue-600 focus:ring-blue-600";

  return (
    <div>
      <label
        htmlFor={id}
        className="flex items-center gap-1.5 text-sm font-medium text-gray-700"
      >
        {icone && <span className="text-base">{icone}</span>}
        {label}
      </label>

      <div className="relative mt-1">
        <input
          id={id}
          {...props}
          type={tipoSenha && !mostrarSenha ? "password" : props.type || "text"}
          className={`block w-full rounded-lg border px-4 py-2.5 pr-10 text-sm outline-none transition-colors ${borda} ${className ?? ""}`}
        />

        {tipoSenha && props.value && String(props.value).length > 0 && (
          <button
            type="button"
            onClick={() => setMostrarSenha((v) => !v)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-600 hover:text-gray-700"
            aria-label={mostrarSenha ? "Esconder senha" : "Mostrar senha"}
          >
            {mostrarSenha ? (
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94" />
                <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
                <line x1="1" y1="1" x2="23" y2="23" />
              </svg>
            ) : (
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                <circle cx="12" cy="12" r="3" />
              </svg>
            )}
          </button>
        )}
      </div>

      {temErro && (
        <p className="mt-1 text-xs text-red-700">{erro}</p>
      )}
    </div>
  );
}
