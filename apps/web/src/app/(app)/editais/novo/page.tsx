"use client";

import { useState, FormEvent } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import CampoFormulario from "../../../../components/campo-formulario";

type CheckboxOption = {
  value: string;
  label: string;
};

const porteOptions: CheckboxOption[] = [
  { value: "MEI", label: "MEI" },
  { value: "ME", label: "ME" },
  { value: "EPP", label: "EPP" },
];

const enquadramentoOptions: CheckboxOption[] = [
  { value: "Simples Nacional", label: "Simples Nacional" },
  { value: "Lucro Presumido", label: "Lucro Presumido" },
];

const documentoOptions: CheckboxOption[] = [
  { value: "RG", label: "RG" },
  { value: "CPF", label: "CPF" },
  { value: "comprovante_conclusao", label: "Comprovante de Conclusão" },
  { value: "comprovante_residencia", label: "Comprovante de Residência" },
];

const modeloOptions = [
  { value: 0, label: "Nenhum" },
  { value: 1, label: "Modelo 1" },
  { value: 2, label: "Modelo 2" },
  { value: 3, label: "Modelo 3" },
  { value: 4, label: "Modelo 4" },
  { value: 5, label: "Modelo 5" },
  { value: 6, label: "Modelo 6" },
];

type MultiCheckboxProps = {
  label: string;
  options: CheckboxOption[];
  selected: string[];
  onChange: (selected: string[]) => void;
};

function MultiCheckbox({ label, options, selected, onChange }: MultiCheckboxProps) {
  function toggle(val: string) {
    if (selected.includes(val)) {
      onChange(selected.filter((v) => v !== val));
    } else {
      onChange([...selected, val]);
    }
  }

  return (
    <div>
      <label className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700">
        {label}
      </label>
      <div className="mt-1 flex flex-wrap gap-2">
        {options.map((opt) => (
          <label
            key={opt.value}
            className={`flex cursor-pointer items-center gap-1.5 rounded-lg border px-3 py-1.5 text-sm transition-colors ${
              selected.includes(opt.value)
                ? "border-brand-600 bg-brand-50 text-brand-700"
                : "border-gray-300 text-gray-600 hover:border-gray-400"
            }`}
          >
            <input
              type="checkbox"
              checked={selected.includes(opt.value)}
              onChange={() => toggle(opt.value)}
              className="sr-only"
            />
            {opt.label}
          </label>
        ))}
      </div>
    </div>
  );
}

export default function NovoEditalPage() {
  const router = useRouter();
  const [form, setForm] = useState({
    nome: "",
    descricao: "",
    dataInicio: "",
    dataFim: "",
    tipoChamada: "",
    modeloFormulario: 0,
    tituloMinimoElegibilidade: "",
  });
  const [exigeEmpresa, setExigeEmpresa] = useState(false);
  const [porteEmpresa, setPorteEmpresa] = useState<string[]>([]);
  const [enquadramentoEmpresa, setEnquadramentoEmpresa] = useState<string[]>([]);
  const [documentosObrigatorios, setDocumentosObrigatorios] = useState<string[]>([]);
  const [erro, setErro] = useState("");
  const [carregando, setCarregando] = useState(false);

  function atualizar(chave: string, valor: string) {
    setForm((prev) => ({ ...prev, [chave]: valor }));
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setErro("");

    if (form.dataInicio && form.dataFim && form.dataInicio > form.dataFim) {
      setErro("Data de início não pode ser posterior à data de fim.");
      return;
    }

    setCarregando(true);

    try {
      const body: Record<string, unknown> = {
        nome: form.nome,
        descricao: form.descricao,
        data_inicio: form.dataInicio,
        data_fim: form.dataFim,
        tipo_chamada: form.tipoChamada,
        modelo_formulario: form.modeloFormulario,
        titulo_minimo_elegibilidade: form.tituloMinimoElegibilidade,
        exige_empresa: exigeEmpresa,
        porte_empresa: porteEmpresa,
        enquadramento_empresa: enquadramentoEmpresa,
        documentos_obrigatorios: documentosObrigatorios,
      };

      const res = await fetch("/api/v1/editais", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });

      const data = await res.json();

      if (!res.ok) {
        setErro(data.erro || "Erro ao criar edital");
        return;
      }

      router.push("/editais");
    } catch {
      setErro("Erro de conexão com o servidor");
    } finally {
      setCarregando(false);
    }
  }

  return (
    <main className="mx-auto max-w-3xl px-4 py-6">
      <div className="rounded-xl border border-gray-200 bg-white p-8 shadow-sm">
        <h1 className="text-2xl font-bold text-gray-900">Novo Edital</h1>
        <p className="mt-1 text-sm text-gray-600">
          Preencha os dados para criar um novo edital
        </p>

        <form onSubmit={handleSubmit} className="mt-6 space-y-6">
          <div className="border-b border-gray-100 pb-6">
            <h2 className="mb-4 text-lg font-semibold text-gray-800">Formulário Base</h2>

            <CampoFormulario
              id="nome"
              label="Nome do Edital"
              placeholder="Ex: Edital APQ 2026"
              value={form.nome}
              onChange={(e) => atualizar("nome", e.target.value)}
              required
            />

            <div className="mt-4">
              <label htmlFor="descricao" className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700">
                Descrição
              </label>
              <textarea
                id="descricao"
                value={form.descricao}
                onChange={(e) => atualizar("descricao", e.target.value)}
                placeholder="Descrição do edital"
                required
                rows={3}
                className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
              />
            </div>

            <div className="mt-4 grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="dataInicio" className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700">
                  Data de Início
                </label>
                <input
                  id="dataInicio"
                  type="date"
                  value={form.dataInicio}
                  onChange={(e) => atualizar("dataInicio", e.target.value)}
                  required
                  className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
                />
              </div>
              <div>
                <label htmlFor="dataFim" className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700">
                  Data de Fim
                </label>
                <input
                  id="dataFim"
                  type="date"
                  value={form.dataFim}
                  onChange={(e) => atualizar("dataFim", e.target.value)}
                  required
                  className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
                />
              </div>
            </div>

            <div className="mt-4 grid grid-cols-2 gap-4">
              <CampoFormulario
                id="tipoChamada"
                label="Tipo de Chamada"
                placeholder="Ex: APQ, ARC"
                value={form.tipoChamada}
                onChange={(e) => atualizar("tipoChamada", e.target.value)}
              />
              <div>
                <label htmlFor="modeloFormulario" className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700">
                  Modelo de Formulário
                </label>
                <select
                  id="modeloFormulario"
                  value={form.modeloFormulario}
                  onChange={(e) => setForm((prev) => ({ ...prev, modeloFormulario: Number(e.target.value) }))}
                  className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
                >
                  {modeloOptions.map((opt) => (
                    <option key={opt.value} value={opt.value}>
                      {opt.label}
                    </option>
                  ))}
                </select>
              </div>
            </div>
          </div>

          <div className="border-b border-gray-100 pb-6">
            <h2 className="mb-4 text-lg font-semibold text-gray-800">Proponente</h2>

            <div>
              <label htmlFor="tituloMinimo" className="mb-1 flex items-center gap-1.5 text-sm font-medium text-gray-700">
                Titulação Mínima para Elegibilidade
              </label>
              <select
                id="tituloMinimo"
                value={form.tituloMinimoElegibilidade}
                onChange={(e) => atualizar("tituloMinimoElegibilidade", e.target.value)}
                className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
              >
                <option value="">Não exigido</option>
                <option value="Graduado">Graduado</option>
                <option value="Mestre">Mestre</option>
                <option value="Doutor">Doutor</option>
              </select>
            </div>
          </div>

          <div className="border-b border-gray-100 pb-6">
            <h2 className="mb-4 text-lg font-semibold text-gray-800">Empresa</h2>

            <label className="flex cursor-pointer items-center gap-2 text-sm">
              <input
                type="checkbox"
                checked={exigeEmpresa}
                onChange={(e) => {
                  setExigeEmpresa(e.target.checked);
                  if (!e.target.checked) {
                    setPorteEmpresa([]);
                    setEnquadramentoEmpresa([]);
                  }
                }}
                className="h-4 w-4 rounded border-gray-300 text-brand-600 focus:ring-brand-600"
              />
              Exigir empresa vinculada
            </label>

            {exigeEmpresa && (
              <div className="mt-4 space-y-4 pl-4 border-l-2 border-brand-200">
                <MultiCheckbox
                  label="Porte da Empresa"
                  options={porteOptions}
                  selected={porteEmpresa}
                  onChange={setPorteEmpresa}
                />
                <MultiCheckbox
                  label="Enquadramento"
                  options={enquadramentoOptions}
                  selected={enquadramentoEmpresa}
                  onChange={setEnquadramentoEmpresa}
                />
              </div>
            )}
          </div>

          <div className="border-b border-gray-100 pb-6">
            <h2 className="mb-4 text-lg font-semibold text-gray-800">Documentações Complementares</h2>

            <MultiCheckbox
              label="Documentos Obrigatórios"
              options={documentoOptions}
              selected={documentosObrigatorios}
              onChange={setDocumentosObrigatorios}
            />
          </div>

          {erro && (
            <div className="rounded-lg bg-red-50 p-3 text-sm text-red-700">
              {erro}
            </div>
          )}

          <div className="flex gap-3">
            <Link
              href="/editais"
              className="rounded-lg border border-gray-300 px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50"
            >
              Cancelar
            </Link>
            <button
              type="submit"
              disabled={carregando}
              className="flex-1 rounded-lg bg-brand-600 px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-brand-700 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {carregando ? "Criando..." : "Criar Edital"}
            </button>
          </div>
        </form>
      </div>
    </main>
  );
}
