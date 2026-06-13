export type DadosProponente = {
  nome: string;
  cpf: string;
  rg: string;
  genero: string;
  etnia: string;
  data_nascimento: string;
  endereco: string;
  cep: string;
  logradouro: string;
  numero: string;
  complemento: string;
  bairro: string;
  cidade: string;
  uf: string;
  telefone: string;
  email: string;
};

export type DadosAcademicos = {
  maior_titulacao: string;
  curso: string;
  instituicao: string;
  ano_conclusao: number;
  area_conhecimento: string;
};

export type EmpresaVinculada = {
  nome: string;
  cnpj: string;
  porte: string;
  enquadramento: string;
};

export type ItemOrcamentario = {
  id?: number;
  descricao: string;
  tipo: "consumo" | "permanente";
  quantidade: number;
  valor_unitario: number;
  valor_total: number;
};

export type DocumentoProposta = {
  id?: number;
  tipo: string;
  nome_arquivo: string;
  data_upload: string;
};

export type Proposta = {
  id: number;
  edital_id: number;
  proponente_id: string;
  protocolo: string;
  status: string;
  dados_proponente: DadosProponente;
  dados_academicos: DadosAcademicos;
  empresa_vinculada?: EmpresaVinculada;
  itens_orcamentarios: ItemOrcamentario[];
  documentos: DocumentoProposta[];
  valor_total_solicitado: number;
  data_submissao?: string;
  data_atualizacao: string;
  versao: number;
};

export type PropostaResumo = {
  id: number;
  edital_id: number;
  edital_nome: string;
  protocolo: string;
  status: string;
  valor_total_solicitado: number;
  data_submissao?: string;
  data_atualizacao: string;
  versao: number;
};

export type CriarPropostaPayload = {
  edital_id: number;
  proponente_id: number;
  dados_proponente: DadosProponente;
  dados_academicos: DadosAcademicos;
  empresa_vinculada?: EmpresaVinculada;
  itens_orcamentarios: Omit<ItemOrcamentario, "valor_total">[];
};

async function apiFetch<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(url, {
    headers: { "Content-Type": "application/json", ...options?.headers },
    ...options,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.erro || `Erro ${res.status}: ${res.statusText}`);
  }
  return res.json();
}

export async function criarProposta(payload: CriarPropostaPayload): Promise<Proposta> {
  return apiFetch<Proposta>("/api/v1/propostas", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function listarPropostas(params?: { edital_id?: number }): Promise<PropostaResumo[]> {
  const query = params?.edital_id ? `?edital_id=${params.edital_id}` : "";
  const data = await apiFetch<{ propostas: PropostaResumo[] }>(`/api/v1/propostas${query}`);
  return data.propostas ?? [];
}

export async function visualizarProposta(id: number): Promise<Proposta> {
  return apiFetch<Proposta>(`/api/v1/propostas/${id}`);
}

export async function editarProposta(id: number, payload: Partial<CriarPropostaPayload>): Promise<Proposta> {
  return apiFetch<Proposta>(`/api/v1/propostas/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export async function deletarProposta(id: number): Promise<void> {
  await apiFetch<void>(`/api/v1/propostas/${id}`, { method: "DELETE" });
}

export async function submeterProposta(id: number): Promise<Proposta> {
  return apiFetch<Proposta>(`/api/v1/propostas/${id}/submeter`, { method: "POST" });
}

export async function uploadDocumento(propostaId: number, arquivo: File, tipo: string): Promise<DocumentoProposta> {
  const formData = new FormData();
  formData.append("arquivo", arquivo);
  formData.append("tipo", tipo);
  const res = await fetch(`/api/v1/propostas/${propostaId}/documentos`, {
    method: "POST",
    body: formData,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.erro || "Erro ao fazer upload");
  }
  return res.json();
}

const STATUS_LABEL: Record<string, string> = {
  rascunho: "Rascunho",
  submetida: "Submetida",
  em_avaliacao: "Em Avaliação",
  avaliada: "Avaliada",
  aprovada: "Aprovada",
  rejeitada: "Rejeitada",
  recursada: "Recursada",
  finalizada: "Finalizada",
};

const STATUS_COLOR: Record<string, string> = {
  rascunho: "bg-gray-100 text-gray-600",
  submetida: "bg-blue-100 text-blue-800",
  em_avaliacao: "bg-yellow-100 text-yellow-800",
  avaliada: "bg-purple-100 text-purple-800",
  aprovada: "bg-green-100 text-green-800",
  rejeitada: "bg-red-100 text-red-800",
  recursada: "bg-orange-100 text-orange-800",
  finalizada: "bg-gray-200 text-gray-700",
};

export function formatarStatus(status: string): string {
  return STATUS_LABEL[status] || status;
}

export function corStatus(status: string): string {
  return STATUS_COLOR[status] || "bg-gray-100 text-gray-600";
}

export function formatarBRL(valor: number): string {
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(valor / 100);
}
