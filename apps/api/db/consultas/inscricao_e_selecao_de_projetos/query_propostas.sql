-- name: InserirProposta :one
INSERT INTO propostas (
    edital_id, proponente_id, protocolo, status, versao,
    proponente_nome, proponente_cpf, proponente_rg, proponente_genero, proponente_etnia,
    proponente_data_nascimento, proponente_endereco, proponente_telefone, proponente_email,
    academico_maior_titulacao, academico_curso, academico_instituicao, academico_ano_conclusao,
    academico_area_conhecimento, empresa_vinculada, valor_total_solicitado,
    data_submissao, data_atualizacao
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10,
    $11, $12, $13, $14,
    $15, $16, $17, $18,
    $19, $20, $21,
    $22, NOW()
) RETURNING id, criado_em;

-- name: BuscarPropostaPorID :one
SELECT * FROM propostas WHERE id = $1;

-- name: BuscarPropostaPorProtocolo :one
SELECT * FROM propostas WHERE protocolo = $1;

-- name: ListarPropostas :many
SELECT * FROM propostas
WHERE (edital_id = $1 OR $1 = 0)
  AND (proponente_id = $2 OR $2 = 0)
  AND (status = $3 OR $3 = '')
ORDER BY criado_em DESC;

-- name: AtualizarProposta :exec
UPDATE propostas SET
    versao = $2, status = $3,
    proponente_nome = $4, proponente_cpf = $5, proponente_rg = $6,
    proponente_genero = $7, proponente_etnia = $8,
    proponente_data_nascimento = $9, proponente_endereco = $10,
    proponente_telefone = $11, proponente_email = $12,
    academico_maior_titulacao = $13, academico_curso = $14,
    academico_instituicao = $15, academico_ano_conclusao = $16,
    academico_area_conhecimento = $17, empresa_vinculada = $18,
    valor_total_solicitado = $19, protocolo = $20,
    data_submissao = $21, data_atualizacao = NOW()
WHERE id = $1;

-- name: DeletarProposta :exec
DELETE FROM propostas WHERE id = $1;

-- name: ContarPropostasPorEdital :one
SELECT COUNT(*) FROM propostas WHERE edital_id = $1;

-- name: InserirVersaoProposta :exec
INSERT INTO versoes_proposta (
    proposta_id, versao, status,
    proponente_nome, proponente_cpf, proponente_rg, proponente_genero, proponente_etnia,
    proponente_data_nascimento, proponente_endereco, proponente_telefone, proponente_email,
    academico_maior_titulacao, academico_curso, academico_instituicao, academico_ano_conclusao,
    academico_area_conhecimento, empresa_vinculada, valor_total_solicitado,
    protocolo, data_submissao
) VALUES (
    $1, $2, $3,
    $4, $5, $6, $7, $8,
    $9, $10, $11, $12,
    $13, $14, $15, $16,
    $17, $18, $19,
    $20, $21
);
