import Link from "next/link";
import {
  criarEventoAuditoria,
  formatarDataHoraBrasil,
  identidadesDemonstrativas,
  listarModulosVisiveis,
  obterSessaoDemonstrativa
} from "../lib/catalogo-plataforma";

type PageProps = {
  searchParams?: Promise<Record<string, string | string[] | undefined>>;
};

function obterParametro(valor: string | string[] | undefined) {
  return Array.isArray(valor) ? valor[0] : valor;
}

export default async function DashboardPage({ searchParams }: PageProps) {
  const parametros = (await searchParams) ?? {};
  const sessaoId = obterParametro(parametros.sessao);
  const sessaoEncerradaId = obterParametro(parametros.sair);
  const sessaoEncerrada = obterSessaoDemonstrativa(sessaoEncerradaId);
  const sessao = obterSessaoDemonstrativa(sessaoId);

  if (!sessao) {
    const evento = sessaoEncerrada
      ? criarEventoAuditoria({
          ator: sessaoEncerrada.ator,
          perfil: sessaoEncerrada.perfil,
          acao: "logout_demonstrativo",
          resultado: "sucesso",
          contexto: {
            sessao: sessaoEncerrada.id,
            identidadeDemonstrativa: sessaoEncerrada.identidadeId
          }
        })
      : criarEventoAuditoria({
          acao: "acessar_dashboard",
          resultado: "negado",
          contexto: { motivo: "sessao_ausente" }
        });

    return (
      <main className="shell">
        <header className="topbar">
          <div className="brand">
            <strong>FAPITEC-SE</strong>
            <span>Validação de acesso à plataforma</span>
          </div>
          <div className="userbox">Sessão demonstrativa ausente</div>
        </header>

        <section className="content">
          <div className="summary">
            <span className={`status ${sessaoEncerrada ? "" : "status-alerta"}`}>
              {sessaoEncerrada ? "Sessão encerrada" : "Acesso restrito"}
            </span>
            <h1>Entrar com identidade demonstrativa</h1>
            <p>
              Este fluxo usa identidades controladas apenas para validação institucional. Não
              colete nem informe dados reais de usuários da FAPITEC-SE nesta etapa.
            </p>
          </div>

          <div className="grid" aria-label="Identidades demonstrativas">
            {identidadesDemonstrativas.map((identidade) => (
              <Link className="card" href={`/?sessao=sessao-${identidade.id}`} key={identidade.id}>
                <h2>{identidade.nome}</h2>
                <p>Perfil demonstrativo: {identidade.perfil}</p>
                <div className="card-footer">
                  <span className="status">Validação</span>
                  <span className="link">Entrar</span>
                </div>
              </Link>
            ))}
          </div>

          <section className="audit-panel" aria-label="Evento de auditoria">
            <h2>Evento auditável</h2>
            <dl>
              <div>
                <dt>Ação</dt>
                <dd>{evento.acao}</dd>
              </div>
              <div>
                <dt>Resultado</dt>
                <dd>{evento.resultado}</dd>
              </div>
              <div>
                <dt>Data e hora</dt>
                <dd>{formatarDataHoraBrasil(evento.dataHora)}</dd>
              </div>
            </dl>
          </section>
        </section>
      </main>
    );
  }

  const modulosVisiveis = listarModulosVisiveis(sessao.perfil);
  const eventoLogin = criarEventoAuditoria({
    ator: sessao.ator,
    perfil: sessao.perfil,
    acao: "login_demonstrativo",
    resultado: "sucesso",
    contexto: { sessao: sessao.id, identidadeDemonstrativa: sessao.identidadeId }
  });

  return (
    <main className="shell">
      <header className="topbar">
        <div className="brand">
          <strong>FAPITEC-SE</strong>
          <span>Plataforma integrada de gestão institucional</span>
        </div>
        <div className="userbox">
          {sessao.nome}
          <span>{sessao.perfilNome}</span>
          <Link href={`/?sair=${sessao.id}`}>Sair</Link>
        </div>
      </header>

      <section className="content">
        <div className="summary">
          <h1>Dashboard modular</h1>
          <p>
            Macroestrutura inicial para validação dos módulos previstos. Os cards abaixo são
            placeholders navegáveis e não detalham regras de negócio ainda não confirmadas.
          </p>
        </div>

        <div className="grid" aria-label="Módulos da plataforma">
          {modulosVisiveis.map((modulo) => (
            <Link className="card" href={`/modulos/${modulo.id}?sessao=${sessao.id}`} key={modulo.id}>
              <h2>{modulo.nome}</h2>
              <p>{modulo.descricao}</p>
              <div className="card-footer">
                <span className="status">{modulo.status}</span>
                <span className="link">Acessar</span>
              </div>
            </Link>
          ))}
        </div>

        <section className="audit-panel" aria-label="Trilha auditável mínima">
          <h2>Trilha auditável mínima</h2>
          <dl>
            <div>
              <dt>Ator</dt>
              <dd>{eventoLogin.ator}</dd>
            </div>
            <div>
              <dt>Ação</dt>
              <dd>{eventoLogin.acao}</dd>
            </div>
            <div>
              <dt>Resultado</dt>
              <dd>{eventoLogin.resultado}</dd>
            </div>
            <div>
              <dt>Data e hora</dt>
              <dd>{formatarDataHoraBrasil(eventoLogin.dataHora)}</dd>
            </div>
          </dl>
        </section>
      </section>
    </main>
  );
}
