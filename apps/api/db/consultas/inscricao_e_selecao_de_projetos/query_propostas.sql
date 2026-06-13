-- name: InserirProposta :one
INSERT INTO propostas (
    edital_id, proponente_id, protocolo, status, versao,
    proponente_nome, proponente_cpf, proponente_rg, proponente_genero, proponente_etnia,
    proponente_data_nascimento,
    proponente_cep, proponente_logradouro, proponente_numero, proponente_complemento,
    proponente_bairro, proponente_cidade, proponente_uf,
    proponente_telefone, proponente_email,
    academico_maior_titulacao, academico_curso, academico_instituicao, academico_ano_conclusao,
    academico_area_conhecimento, empresa_vinculada, valor_total_solicitado,
    data_submissao, data_atualizacao
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16,
    $17, $18, $19, $20,
    $21, $22, $23, $24,
    $25, $26, $27,
    $28, NOW()
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
    proponente_data_nascimento = $9,
    proponente_cep = $10, proponente_logradouro = $11, proponente_numero = $12,
    proponente_complemento = $13, proponente_bairro = $14,
    proponente_cidade = $15, proponente_uf = $16,
    proponente_telefone = $17, proponente_email = $18,
    academico_maior_titulacao = $19, academico_curso = $20,
    academico_instituicao = $21, academico_ano_conclusao = $22,
    academico_area_conhecimento = $23, empresa_vinculada = $24,
    valor_total_solicitado = $25, protocolo = $26,
    data_submissao = $27, data_atualizacao = NOW()
WHERE id = $1;

-- name: DeletarProposta :exec
DELETE FROM propostas WHERE id = $1;

-- name: ContarPropostasPorEdital :one
SELECT COUNT(*) FROM propostas WHERE edital_id = $1;

-- name: InserirVersaoProposta :exec
INSERT INTO versoes_proposta (
    proposta_id, versao, status,
    proponente_nome, proponente_cpf, proponente_rg, proponente_genero, proponente_etnia,
    proponente_data_nascimento,
    proponente_cep, proponente_logradouro, proponente_numero, proponente_complemento,
    proponente_bairro, proponente_cidade, proponente_uf,
    proponente_telefone, proponente_email,
    academico_maior_titulacao, academico_curso, academico_instituicao, academico_ano_conclusao,
    academico_area_conhecimento, empresa_vinculada, valor_total_solicitado,
    protocolo, data_submissao
) VALUES (
    $1, $2, $3,
    $4, $5, $6, $7, $8,
    $9, $10, $11, $12, $13, $14, $15,
    $16, $17, $18, $19, $20,
    $21, $22, $23, $24,
    $25, $26, $27
);
