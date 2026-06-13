export type PropostaParaAvaliarSaida = {
  id: number;
  edital_id: number;
  protocolo: string;
  status: string;
  dados_proponente: {
    nome: string;
    cpf: string;
    email: string;
  };
  dados_academicos: {
    maior_titulacao: string;
    area_conhecimento: string;
  };
  valor_total: number;
  pareceres: ParecerSaida[];
  data_submissao: string;
};

export type ParecerSaida = {
  id: number;
  proposta_id: number;
  etapa: string;
  nota: number;
  parecer_texto: string;
  data: string;
};

export type ParecerAnonimizadoSaida = {
  id: number;
  proposta_id: number;
  etapa: string;
  hash_avaliador: string;
  nota: number;
  parecer_texto: string;
  data: string;
};

export type EmitirParecerEntrada = {
  proposta_id: number;
  etapa: string;
  avaliador_id: number;
  nota: number;
  parecer_texto: string;
};

export type ClassificacaoSaida = {
  proposta_id: number;
  protocolo: string;
  nota_final: number;
  status: string;
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

export async function listarPropostasParaAvaliar(avaliadorId: number): Promise<PropostaParaAvaliarSaida[]> {
  return apiFetch<PropostaParaAvaliarSaida[]>(
    `/api/v1/avaliadores/me/propostas?avaliador_id=${avaliadorId}`
  );
}

export async function emitirParecer(propostaId: number, payload: EmitirParecerEntrada): Promise<ParecerSaida> {
  return apiFetch<ParecerSaida>(`/api/v1/propostas/${propostaId}/pareceres`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function listarPareceres(propostaId: number): Promise<ParecerAnonimizadoSaida[]> {
  return apiFetch<ParecerAnonimizadoSaida[]>(`/api/v1/propostas/${propostaId}/pareceres`);
}

export async function finalizarAvaliacao(
  editalId: number,
  payload: { nota_de_corte: number }
): Promise<ClassificacaoSaida[]> {
  return apiFetch<ClassificacaoSaida[]>(`/api/v1/editais/${editalId}/finalizar-avaliacao`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

const STATUS_LABEL: Record<string, string> = {
  submetida: "Submetida",
  em_avaliacao: "Em Avaliação",
  aprovada: "Aprovada",
  reprovada: "Reprovada",
};

const STATUS_COLOR: Record<string, string> = {
  submetida: "bg-blue-100 text-blue-800",
  em_avaliacao: "bg-yellow-100 text-yellow-800",
  aprovada: "bg-green-100 text-green-800",
  reprovada: "bg-red-100 text-red-800",
};

export function formatarStatusProposta(status: string): string {
  return STATUS_LABEL[status] || status;
}

export function corStatusProposta(status: string): string {
  return STATUS_COLOR[status] || "bg-gray-100 text-gray-600";
}

export function formatarDataISO(data: string): string {
  return new Date(data).toLocaleDateString("pt-BR");
}

export function formatarMoeda(valor: number): string {
  return (valor / 100).toLocaleString("pt-BR", {
    style: "currency",
    currency: "BRL",
  });
}
