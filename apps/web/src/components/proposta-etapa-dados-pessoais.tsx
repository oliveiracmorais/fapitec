"use client";

import type { DadosProponente } from "../lib/api-propostas";
import { formatarCPF } from "../lib/validacao";
import CampoFormulario from "./campo-formulario";

type Props = {
  dados: DadosProponente;
  onChange: (campo: keyof DadosProponente, valor: string) => void;
  erros: Record<string, string>;
  touched: Record<string, boolean>;
  onBlur: (campo: string) => void;
};

export default function PropostaEtapaDadosPessoais({ dados, onChange, erros, touched, onBlur }: Props) {
  function handleChange(campo: keyof DadosProponente, valor: string) {
    if (campo === "cpf") valor = formatarCPF(valor);
    onChange(campo, valor);
  }

  return (
    <div className="space-y-4">
      <h2 className="text-lg font-semibold text-gray-900">Dados do Proponente</h2>
      <CampoFormulario
        id="nome"
        label="Nome Completo"
        value={dados.nome}
        onChange={(e) => handleChange("nome", e.target.value)}
        erro={erros.nome}
        touched={touched.nome}
        onBlur={() => onBlur("nome")}
        placeholder="Seu nome completo"
        icone="👤"
      />
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <CampoFormulario
          id="cpf"
          label="CPF"
          value={dados.cpf}
          onChange={(e) => handleChange("cpf", e.target.value)}
          erro={erros.cpf}
          touched={touched.cpf}
          onBlur={() => onBlur("cpf")}
          placeholder="000.000.000-00"
          icone="📋"
        />
        <CampoFormulario
          id="rg"
          label="RG"
          value={dados.rg}
          onChange={(e) => handleChange("rg", e.target.value)}
          erro={erros.rg}
          touched={touched.rg}
          onBlur={() => onBlur("rg")}
          placeholder="RG"
          icone="🆔"
        />
      </div>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <CampoFormulario
          id="genero"
          label="Gênero"
          value={dados.genero}
          onChange={(e) => handleChange("genero", e.target.value)}
          erro={erros.genero}
          touched={touched.genero}
          onBlur={() => onBlur("genero")}
          placeholder="Feminino / Masculino / Outro"
          icone="⚧"
        />
        <CampoFormulario
          id="etnia"
          label="Etnia"
          value={dados.etnia}
          onChange={(e) => handleChange("etnia", e.target.value)}
          erro={erros.etnia}
          touched={touched.etnia}
          onBlur={() => onBlur("etnia")}
          placeholder="Branca / Parda / Preta / Amarela / Indígena"
          icone="🌍"
        />
      </div>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <CampoFormulario
          id="data_nascimento"
          label="Data de Nascimento"
          type="date"
          value={dados.data_nascimento}
          onChange={(e) => handleChange("data_nascimento", e.target.value)}
          erro={erros.data_nascimento}
          touched={touched.data_nascimento}
          onBlur={() => onBlur("data_nascimento")}
          icone="🎂"
        />
        <CampoFormulario
          id="telefone"
          label="Telefone"
          value={dados.telefone}
          onChange={(e) => handleChange("telefone", e.target.value)}
          erro={erros.telefone}
          touched={touched.telefone}
          onBlur={() => onBlur("telefone")}
          placeholder="(79) 99999-9999"
          icone="📱"
        />
      </div>
      <CampoFormulario
        id="email"
        label="E-mail"
        type="email"
        value={dados.email}
        onChange={(e) => handleChange("email", e.target.value)}
        erro={erros.email}
        touched={touched.email}
        onBlur={() => onBlur("email")}
        placeholder="email@exemplo.com"
        icone="📧"
      />
      <div>
        <label htmlFor="endereco" className="flex items-center gap-1.5 text-sm font-medium text-gray-700">
          <span className="text-base">📍</span> Endereço
        </label>
        <textarea
          id="endereco"
          value={dados.endereco}
          onChange={(e) => handleChange("endereco", e.target.value)}
          onBlur={() => onBlur("endereco")}
          placeholder="Rua, número, bairro, cidade, estado, CEP"
          rows={2}
          className={`mt-1 block w-full rounded-lg border px-4 py-2.5 text-sm outline-none transition-colors ${
            touched.endereco && erros.endereco
              ? "border-red-400 focus:border-red-500 focus:ring-red-500"
              : touched.endereco && !erros.endereco && dados.endereco
                ? "border-green-400"
                : "border-gray-300 focus:border-blue-600 focus:ring-blue-600"
          }`}
        />
        {touched.endereco && erros.endereco && (
          <p className="mt-1 text-xs text-red-700">{erros.endereco}</p>
        )}
      </div>
    </div>
  );
}
