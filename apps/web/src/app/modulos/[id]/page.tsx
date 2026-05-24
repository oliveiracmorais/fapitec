import Link from "next/link";
import { notFound } from "next/navigation";
import {
  buscarModulo,
  criarEventoAuditoria,
  formatarDataHoraBrasil,
  modulos,
  obterSessaoDemonstrativa,
  podeAcessarModulo
} from "../../../lib/catalogo-plataforma";

export function generateStaticParams() {
  return modulos.map((modulo) => ({ id: modulo.id }));
}

type PageProps = {
  params: Promise<{ id: string }>;
  searchParams?: Promise<Record<string, string | string[] | undefined>>;
};

function obterParametro(valor: string | string[] | undefined) {
  return Array.isArray(valor) ? valor[0] : valor;
}

export default async function ModuloPlaceholderPage({ params, searchParams }: PageProps) {
  const { id } = await params;
  const parametros = (await searchParams) ?? {};
  const modulo = buscarModulo(id);

  if (!modulo) {
    notFound();
  }

  const sessaoId = obterParametro(parametros.sessao);
  const sessao = obterSessaoDemonstrativa(sessaoId);
  const permitido = sessao ? podeAcessarModulo(sessao.perfil, modulo.id) : false;

  if (!sessao || !permitido) {
    const eventoNegado = criarEventoAuditoria({
      ator: sessao?.ator,
      perfil: sessao?.perfil,
      acao: "acessar_modulo_dashboard",
      moduloId: modulo.id,
      resultado: "negado",
      contexto: {
        motivo: sessao ? "perfil_sem_permissao" : "sessao_ausente",
        ...(sessao ? { sessao: sessao.id, identidadeDemonstrativa: sessao.identidadeId } : {})
      }
    });

    return (
      <main className="shell">
        <header className="topbar">
          <div className="brand">
            <strong>FAPITEC-SE</strong>
            <span>{modulo.nome}</span>
          </div>
          <Link href={sessao ? `/?sessao=${sessao.id}` : "/"}>Voltar ao dashboard</Link>
        </header>

        <section className="content">
          <div className="placeholder">
            <span className="status status-alerta">Acesso negado</span>
            <h1>{modulo.nome}</h1>
            <p>
              A sessão demonstrativa atual não possui permissão para acessar este módulo.
            </p>
          </div>

          <section className="audit-panel" aria-label="Evento de auditoria">
            <h2>Evento auditável</h2>
            <dl>
              <div>
                <dt>Ação</dt>
                <dd>{eventoNegado.acao}</dd>
              </div>
              <div>
                <dt>Resultado</dt>
                <dd>{eventoNegado.resultado}</dd>
              </div>
              <div>
                <dt>Módulo</dt>
                <dd>{modulo.nome}</dd>
              </div>
              <div>
                <dt>Data e hora</dt>
                <dd>{formatarDataHoraBrasil(eventoNegado.dataHora)}</dd>
              </div>
            </dl>
          </section>
        </section>
      </main>
    );
  }

  const eventoSucesso = criarEventoAuditoria({
    ator: sessao.ator,
    perfil: sessao.perfil,
    acao: "acessar_modulo_dashboard",
    moduloId: modulo.id,
    resultado: "sucesso",
    contexto: { sessao: sessao.id, identidadeDemonstrativa: sessao.identidadeId }
  });

  return (
    <main className="shell">
      <header className="topbar">
        <div className="brand">
          <strong>FAPITEC-SE</strong>
          <span>{modulo.nome}</span>
        </div>
        <Link href={`/?sessao=${sessao.id}`}>Voltar ao dashboard</Link>
      </header>

      <section className="content">
        <div className="placeholder">
          <span className="status">{modulo.status}</span>
          <h1>{modulo.nome}</h1>
          <p>{modulo.descricao}</p>
          <p>
            Este módulo está reservado para detalhamento em story própria após validação da
            macroestrutura com a equipe da FAPITEC-SE.
          </p>
        </div>

        <section className="audit-panel" aria-label="Evento de auditoria">
          <h2>Evento auditável</h2>
          <dl>
            <div>
              <dt>Ator</dt>
              <dd>{eventoSucesso.ator}</dd>
            </div>
            <div>
              <dt>Ação</dt>
              <dd>{eventoSucesso.acao}</dd>
            </div>
            <div>
              <dt>Resultado</dt>
              <dd>{eventoSucesso.resultado}</dd>
            </div>
            <div>
              <dt>Data e hora</dt>
              <dd>{formatarDataHoraBrasil(eventoSucesso.dataHora)}</dd>
            </div>
          </dl>
        </section>
      </section>
    </main>
  );
}
