export type AvaliadorSaida = {
  id: number;
  usuario_id: number;
  nome: string;
  cpf: string;
  email: string;
  titulacao_maxima: string;
  area_conhecimento: string;
  instituicao: string;
  curriculo_resumido: string;
  estado: string;
  data_cadastro: string;
  data_atualizacao: string;
  total_propostas?: number;
  atribuicoes_ativas?: number;
};

export type CadastrarAvaliadorEntrada = {
  usuario_id: number;
  nome: string;
  cpf: string;
  email: string;
  titulacao_maxima: string;
  area_conhecimento: string;
  instituicao: string;
  curriculo_resumido: string;
};

export type AtualizarAvaliadorEntrada = {
  nome?: string;
  cpf?: string;
  email?: string;
  titulacao_maxima?: string;
  area_conhecimento?: string;
  instituicao?: string;
  curriculo_resumido?: string;
  estado?: string;
};

export type AtribuicaoSaida = {
  id: number;
  avaliador_id: number;
  edital_id: number;
  data_inicio: string;
  data_fim: string;
  status_convite: string;
  hash_anonimizacao: string;
  criado_em: string;
};

export type AtribuirEditalEntrada = {
  edital_id: number;
  data_inicio: string;
  data_fim: string;
};

export type GerenciarConviteEntrada = {
  acao: "aceitar" | "recusar";
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

export async function listarAvaliadores(params?: {
  nome?: string;
  cpf?: string;
  area_conhecimento?: string;
  estado?: string;
}): Promise<AvaliadorSaida[]> {
  const query = new URLSearchParams();
  if (params?.nome) query.set("nome", params.nome);
  if (params?.cpf) query.set("cpf", params.cpf);
  if (params?.area_conhecimento) query.set("area_conhecimento", params.area_conhecimento);
  if (params?.estado) query.set("estado", params.estado);
  const qs = query.toString();
  return apiFetch<AvaliadorSaida[]>(`/api/v1/avaliadores${qs ? `?${qs}` : ""}`);
}

export async function cadastrarAvaliador(payload: CadastrarAvaliadorEntrada): Promise<AvaliadorSaida> {
  return apiFetch<AvaliadorSaida>("/api/v1/avaliadores", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function visualizarAvaliador(id: number): Promise<AvaliadorSaida> {
  return apiFetch<AvaliadorSaida>(`/api/v1/avaliadores/${id}`);
}

export async function editarAvaliador(id: number, payload: AtualizarAvaliadorEntrada): Promise<AvaliadorSaida> {
  return apiFetch<AvaliadorSaida>(`/api/v1/avaliadores/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export async function atribuirEdital(avaliadorId: number, payload: AtribuirEditalEntrada): Promise<AtribuicaoSaida> {
  return apiFetch<AtribuicaoSaida>(`/api/v1/avaliadores/${avaliadorId}/atribuir`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function listarAtribuicoesPorAvaliador(avaliadorId: number): Promise<AtribuicaoSaida[]> {
  return apiFetch<AtribuicaoSaida[]>(`/api/v1/avaliadores/${avaliadorId}/atribuicoes`);
}

export async function listarAtribuicoesPorEdital(editalId: number): Promise<AtribuicaoSaida[]> {
  return apiFetch<AtribuicaoSaida[]>(`/api/v1/editais/${editalId}/avaliadores`);
}

export async function gerenciarConvite(atribuicaoId: number, payload: GerenciarConviteEntrada): Promise<AtribuicaoSaida> {
  return apiFetch<AtribuicaoSaida>(`/api/v1/atribuicoes/${atribuicaoId}/convite`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

const ESTADO_LABEL: Record<string, string> = {
  ativo: "Ativo",
  inativo: "Inativo",
};

const ESTADO_COLOR: Record<string, string> = {
  ativo: "bg-green-100 text-green-800",
  inativo: "bg-gray-100 text-gray-600",
};

const STATUS_CONVITE_LABEL: Record<string, string> = {
  pendente: "Pendente",
  aceito: "Aceito",
  recusado: "Recusado",
};

const STATUS_CONVITE_COLOR: Record<string, string> = {
  pendente: "bg-yellow-100 text-yellow-800",
  aceito: "bg-green-100 text-green-800",
  recusado: "bg-red-100 text-red-800",
};

export function formatarEstado(estado: string): string {
  return ESTADO_LABEL[estado] || estado;
}

export function corEstado(estado: string): string {
  return ESTADO_COLOR[estado] || "bg-gray-100 text-gray-600";
}

export function formatarStatusConvite(status: string): string {
  return STATUS_CONVITE_LABEL[status] || status;
}

export function corStatusConvite(status: string): string {
  return STATUS_CONVITE_COLOR[status] || "bg-gray-100 text-gray-600";
}

export function formatarDataISO(data: string): string {
  return new Date(data).toLocaleDateString("pt-BR");
}
