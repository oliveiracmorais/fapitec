import catalogo from "../../../../packages/config/src/catalogo-plataforma.json";

export type PerfilId = (typeof catalogo.perfis)[number]["id"];
export type Modulo = (typeof catalogo.modulos)[number];
export type IdentidadeDemonstrativa = (typeof catalogo.identidadesDemonstrativas)[number];

export type SessaoDemonstrativa = {
  id: string;
  ator: string;
  identidadeId: string;
  nome: string;
  perfil: string;
  perfilNome: string;
  demonstrativa: true;
};

export type EventoAuditoria = {
  id: string;
  ator: string;
  perfil?: string;
  acao: string;
  resultado: "sucesso" | "falha" | "negado" | "erro";
  modulo?: string;
  contexto: Record<string, string>;
  dataHora: string;
};

export const perfis = catalogo.perfis;
export const modulos = catalogo.modulos;
export const identidadesDemonstrativas = catalogo.identidadesDemonstrativas;

export function listarModulosVisiveis(perfil: string) {
  return modulos.filter((modulo) => modulo.perfisPermitidos.includes(perfil));
}

export function buscarModulo(id: string) {
  return modulos.find((modulo) => modulo.id === id);
}

export function buscarPerfil(id: string) {
  return perfis.find((perfil) => perfil.id === id);
}

export function buscarIdentidadeDemonstrativa(id: string) {
  return identidadesDemonstrativas.find((identidade) => identidade.id === id);
}

export function criarSessaoDemonstrativa(identidadeId: string): SessaoDemonstrativa | undefined {
  const identidade = buscarIdentidadeDemonstrativa(identidadeId);

  if (!identidade) {
    return undefined;
  }

  const perfil = buscarPerfil(identidade.perfil);

  return {
    id: `sessao-${identidade.id}`,
    ator: `identidade-demonstrativa:${identidade.id}`,
    identidadeId: identidade.id,
    nome: identidade.nome,
    perfil: perfil?.id ?? identidade.perfil,
    perfilNome: perfil?.nome ?? identidade.perfil,
    demonstrativa: true
  };
}

export function obterSessaoDemonstrativa(sessaoId?: string): SessaoDemonstrativa | undefined {
  if (!sessaoId) {
    return undefined;
  }

  const identidadeId = sessaoId.startsWith("sessao-") ? sessaoId.slice("sessao-".length) : sessaoId;
  return criarSessaoDemonstrativa(identidadeId);
}

export function podeAcessarModulo(perfil: string, moduloId: string) {
  const modulo = buscarModulo(moduloId);

  return Boolean(modulo?.perfisPermitidos.includes(perfil));
}

export function criarEventoAuditoria({
  ator = "sistema",
  perfil,
  acao,
  moduloId,
  resultado,
  contexto = {}
}: {
  ator?: string;
  perfil?: string;
  acao: string;
  moduloId?: string;
  resultado: EventoAuditoria["resultado"];
  contexto?: Record<string, string>;
}): EventoAuditoria {
  const modulo = moduloId ? buscarModulo(moduloId) : undefined;

  return {
    id: `evt-${acao}-${resultado}`,
    ator,
    perfil,
    acao,
    resultado,
    modulo: modulo?.id,
    contexto: {
      ...(modulo ? { nomeModulo: modulo.nome, tipoModulo: modulo.tipo } : {}),
      ...contexto
    },
    dataHora: new Date().toISOString()
  };
}

export function formatarDataHoraBrasil(dataHora: string) {
  return new Intl.DateTimeFormat("pt-BR", {
    dateStyle: "short",
    timeStyle: "medium",
    timeZone: "America/Maceio"
  }).format(new Date(dataHora));
}
