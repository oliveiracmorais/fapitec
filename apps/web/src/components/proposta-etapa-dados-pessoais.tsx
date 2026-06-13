"use client";

import type { DadosProponente } from "../lib/api-propostas";
import { formatarCPF } from "../lib/validacao";
import CampoFormulario from "./campo-formulario";
import CampoEndereco from "./campo-endereco";

type Props = {
  dados: DadosProponente;
  onChange: (campo: keyof DadosProponente, valor: string) => void;
  erros: Record<string, string>;
  touched: Record<string, boolean>;
  onBlur: (campo: string) => void;
};

const GENEROS = ["Feminino", "Masculino", "Outro", "Prefiro não informar"] as const;
const ETNIAS = ["Branca", "Parda", "Preta", "Amarela", "Indígena", "Prefiro não informar"] as const;

export default function PropostaEtapaDadosPessoais({ dados, onChange, erros, touched, onBlur }: Props) {
  function handleChange(campo: string, valor: string) {
    if (campo === "cpf") valor = formatarCPF(valor);
    onChange(campo as keyof DadosProponente, valor);
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
        <div>
          <label htmlFor="genero" className="flex items-center gap-1.5 text-sm font-medium text-gray-700">
            <span className="text-base">⚧</span> Gênero
          </label>
          <select
            id="genero"
            value={dados.genero}
            onChange={(e) => handleChange("genero", e.target.value)}
            onBlur={() => onBlur("genero")}
            className={`mt-1 block w-full rounded-lg border px-4 py-2.5 text-sm outline-none transition-colors ${
              touched.genero && erros.genero
                ? "border-red-400 focus:border-red-500 focus:ring-red-500"
                : touched.genero && !erros.genero && dados.genero
                  ? "border-green-400 focus:border-green-600 focus:ring-green-600"
                  : "border-gray-300 focus:border-blue-600 focus:ring-blue-600"
            }`}
          >
            <option value="">Selecione...</option>
            {GENEROS.map((g) => (
              <option key={g} value={g}>{g}</option>
            ))}
          </select>
          {touched.genero && erros.genero && (
            <p className="mt-1 text-xs text-red-700">{erros.genero}</p>
          )}
        </div>
        <div>
          <label htmlFor="etnia" className="flex items-center gap-1.5 text-sm font-medium text-gray-700">
            <span className="text-base">🌍</span> Etnia
          </label>
          <select
            id="etnia"
            value={dados.etnia}
            onChange={(e) => handleChange("etnia", e.target.value)}
            onBlur={() => onBlur("etnia")}
            className={`mt-1 block w-full rounded-lg border px-4 py-2.5 text-sm outline-none transition-colors ${
              touched.etnia && erros.etnia
                ? "border-red-400 focus:border-red-500 focus:ring-red-500"
                : touched.etnia && !erros.etnia && dados.etnia
                  ? "border-green-400 focus:border-green-600 focus:ring-green-600"
                  : "border-gray-300 focus:border-blue-600 focus:ring-blue-600"
            }`}
          >
            <option value="">Selecione...</option>
            {ETNIAS.map((e) => (
              <option key={e} value={e}>{e}</option>
            ))}
          </select>
          {touched.etnia && erros.etnia && (
            <p className="mt-1 text-xs text-red-700">{erros.etnia}</p>
          )}
        </div>
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
      <CampoEndereco
        cep={dados.cep}
        logradouro={dados.logradouro}
        numero={dados.numero}
        complemento={dados.complemento}
        bairro={dados.bairro}
        cidade={dados.cidade}
        uf={dados.uf}
        onChange={handleChange}
        onBlur={onBlur}
        erros={erros}
        touched={touched}
      />
    </div>
  );
}
