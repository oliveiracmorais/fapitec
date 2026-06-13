"use client";

import { useEffect, useState, useCallback } from "react";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import PropostaFormMultiEtapas from "../../../../../components/proposta-form-multi-etapas";
import { criarProposta, editarProposta, visualizarProposta, type DadosProponente, type DadosAcademicos, type ItemOrcamentario } from "../../../../../lib/api-propostas";
import { useAuth } from "../../../../../context/auth-context";

type DocumentoEntry = {
  tipo: string;
  arquivo?: File;
  nome_arquivo?: string;
};

const PROPOSTA_PADRAO: DadosProponente = {
  nome: "", cpf: "", rg: "", genero: "", etnia: "",
  data_nascimento: "", endereco: "", cep: "", logradouro: "",
  numero: "", complemento: "", bairro: "", cidade: "", uf: "",
  telefone: "", email: "",
};

const ACADEMICO_PADRAO: DadosAcademicos = {
  maior_titulacao: "", curso: "", instituicao: "",
  ano_conclusao: 0, area_conhecimento: "",
};

type EditalInfo = {
  id: number;
  nome: string;
  status: string;
  valor_global: number;
  documentos_obrigatorios: string[];
};

export default function InscreverPage() {
  const params = useParams();
  const router = useRouter();
  const searchParams = useSearchParams();
  const editarId = searchParams.get("editar");
  const { usuario } = useAuth();

  const [carregando, setCarregando] = useState(true);
  const [enviando, setEnviando] = useState(false);
  const [erroGeral, setErroGeral] = useState("");
  const [edital, setEdital] = useState<EditalInfo | null>(null);

  const [dadosProponente, setDadosProponente] = useState<DadosProponente>(PROPOSTA_PADRAO);
  const [dadosAcademicos, setDadosAcademicos] = useState<DadosAcademicos>(ACADEMICO_PADRAO);
  const [itensOrcamentarios, setItensOrcamentarios] = useState<ItemOrcamentario[]>([]);
  const [documentos, setDocumentos] = useState<DocumentoEntry[]>([]);

  useEffect(() => {
    if (!params?.id) return;
    const editalId = Number(params.id);
    if (!editalId) { router.push("/editais"); return; }

    Promise.all([
      fetch(`/api/v1/editais/${editalId}`).then((r) => r.json()),
      editarId ? visualizarProposta(Number(editarId)).catch(() => null) : Promise.resolve(null),
    ])
      .then(([editalData, propostaExistente]) => {
        if (editalData.status !== "ativo") {
          setErroGeral("Este edital não está ativo para inscrições.");
          setCarregando(false);
          return;
        }
        setEdital({
          id: editalData.id,
          nome: editalData.nome,
          status: editalData.status,
          valor_global: editalData.valor_global || 0,
          documentos_obrigatorios: editalData.documentos_obrigatorios || [],
        });

        if (propostaExistente) {
          setDadosProponente(propostaExistente.dados_proponente || PROPOSTA_PADRAO);
          setDadosAcademicos(propostaExistente.dados_academicos || ACADEMICO_PADRAO);
          setItensOrcamentarios(propostaExistente.itens_orcamentarios || []);
        } else if (usuario) {
          setDadosProponente({
            nome: usuario.nome,
            cpf: usuario.documento,
            email: usuario.email,
            rg: "",
            genero: "",
            etnia: "",
            data_nascimento: "",
            endereco: "",
            cep: "",
            logradouro: "",
            numero: "",
            complemento: "",
            bairro: "",
            cidade: "",
            uf: "",
            telefone: "",
          });
        }

        setCarregando(false);
      })
      .catch(() => {
        setErroGeral("Erro ao carregar dados do edital.");
        setCarregando(false);
      });
  }, [params, router, editarId, usuario]);

  const handleProponenteChange = useCallback((campo: keyof DadosProponente, valor: string) => {
    setDadosProponente((prev) => ({ ...prev, [campo]: valor }));
  }, []);

  const handleAcademicoChange = useCallback((campo: keyof DadosAcademicos, valor: string | number) => {
    setDadosAcademicos((prev) => ({ ...prev, [campo]: valor }));
  }, []);

  const handleOrcamentoChange = useCallback(
    (indice: number, campo: keyof ItemOrcamentario, valor: string | number) => {
      setItensOrcamentarios((prev) => {
        const novos = [...prev];
        novos[indice] = { ...novos[indice], [campo]: valor };
        const qtd = campo === "quantidade" ? Number(valor) : novos[indice].quantidade;
        const vu = campo === "valor_unitario" ? Number(valor) : novos[indice].valor_unitario;
        novos[indice].valor_total = qtd * vu;
        return novos;
      });
    },
    []
  );

  const handleOrcamentoAdicionar = useCallback(() => {
    setItensOrcamentarios((prev) => [
      ...prev,
      { descricao: "", tipo: "consumo", quantidade: 1, valor_unitario: 0, valor_total: 0 },
    ]);
  }, []);

  const handleOrcamentoRemover = useCallback((indice: number) => {
    setItensOrcamentarios((prev) => prev.filter((_, i) => i !== indice));
  }, []);

  const handleDocumentoAdicionar = useCallback((tipo: string, arquivo: File) => {
    setDocumentos((prev) => {
      const existente = prev.findIndex((d) => d.tipo === tipo);
      if (existente >= 0) {
        const novos = [...prev];
        novos[existente] = { tipo, arquivo, nome_arquivo: arquivo.name };
        return novos;
      }
      return [...prev, { tipo, arquivo, nome_arquivo: arquivo.name }];
    });
  }, []);

  const handleDocumentoRemover = useCallback((tipo: string) => {
    setDocumentos((prev) => prev.filter((d) => d.tipo !== tipo));
  }, []);

  async function handleSubmit() {
    if (!edital) return;
    setEnviando(true);
    setErroGeral("");

    try {
      const enderecoCompleto = [
        dadosProponente.logradouro,
        dadosProponente.numero,
        dadosProponente.complemento,
        dadosProponente.bairro,
        dadosProponente.cidade,
        dadosProponente.uf,
        dadosProponente.cep,
      ]
        .filter(Boolean)
        .join(", ");

      const payload = {
        edital_id: edital.id,
        proponente_id: Number(usuario?.id) || 0,
        dados_proponente: { ...dadosProponente, endereco: enderecoCompleto },
        dados_academicos: dadosAcademicos,
        itens_orcamentarios: itensOrcamentarios.map((item) => ({
          descricao: item.descricao,
          tipo: item.tipo,
          quantidade: item.quantidade,
          valor_unitario: item.valor_unitario,
        })),
      };

      if (editarId) {
        await editarProposta(Number(editarId), payload);
      } else {
        const proposta = await criarProposta(payload);
        await submeterEDirecionar(proposta.id);
        return;
      }

      router.push(`/minhas-propostas/${editarId}` as never);
    } catch (e) {
      setErroGeral(e instanceof Error ? e.message : "Erro ao salvar proposta");
    } finally {
      setEnviando(false);
    }
  }

  async function submeterEDirecionar(propostaId: number) {
    try {
      const res = await fetch(`/api/v1/propostas/${propostaId}/submeter`, { method: "POST" });
      if (!res.ok) throw new Error("Erro ao submeter");
      router.push(`/minhas-propostas/${propostaId}?sucesso=true` as never);
    } catch {
      router.push(`/minhas-propostas/${propostaId}` as never);
    }
  }

  if (carregando) {
    return (
      <main className="mx-auto max-w-3xl px-4 py-6">
        <div className="flex items-center justify-center py-20 text-gray-600">
          Carregando...
        </div>
      </main>
    );
  }

  if (erroGeral && !edital) {
    return (
      <main className="mx-auto max-w-3xl px-4 py-6">
        <div className="rounded-xl border border-red-200 bg-red-50 p-8 text-center">
          <p className="text-red-700">{erroGeral}</p>
          <Link
            href="/editais"
            className="mt-4 inline-block text-sm font-medium text-brand-600 hover:text-brand-700"
          >
            ← Voltar para editais
          </Link>
        </div>
      </main>
    );
  }

  return (
    <main className="mx-auto max-w-7xl px-4 py-6">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">
          {editarId ? "Editar Proposta" : "Nova Proposta"}
        </h1>
        <p className="mt-1 text-sm text-gray-600">
          Edital: <span className="font-medium">{edital?.nome}</span>
        </p>
      </div>

      {erroGeral && (
        <div className="mb-4 rounded-lg bg-red-50 p-3 text-sm text-red-700">
          {erroGeral}
        </div>
      )}

      <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <PropostaFormMultiEtapas
          dadosProponente={dadosProponente}
          dadosAcademicos={dadosAcademicos}
          itensOrcamentarios={itensOrcamentarios}
          documentos={documentos}
          tiposDocumentosExigidos={edital?.documentos_obrigatorios || []}
          valorGlobalEdital={edital?.valor_global || 0}
          carregando={enviando}
          onProponenteChange={handleProponenteChange}
          onAcademicoChange={handleAcademicoChange}
          onOrcamentoChange={handleOrcamentoChange}
          onOrcamentoAdicionar={handleOrcamentoAdicionar}
          onOrcamentoRemover={handleOrcamentoRemover}
          onDocumentoAdicionar={handleDocumentoAdicionar}
          onDocumentoRemover={handleDocumentoRemover}
          onSubmit={handleSubmit}
        />
      </div>
    </main>
  );
}
