-- name: InserirAvaliador :one
INSERT INTO avaliadores (usuario_id, nome, cpf, email, titulacao_maxima, area_conhecimento, instituicao, curriculo_resumido, estado, data_cadastro, data_atualizacao)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id;

-- name: BuscarAvaliadorPorID :one
SELECT * FROM avaliadores WHERE id = $1;

-- name: BuscarAvaliadorPorCPF :one
SELECT * FROM avaliadores WHERE cpf = $1;

-- name: BuscarAvaliadorPorUsuarioID :one
SELECT * FROM avaliadores WHERE usuario_id = $1;

-- name: ListarAvaliadores :many
SELECT * FROM avaliadores
WHERE ($1 = '' OR nome ILIKE '%' || $1 || '%')
  AND ($2 = '' OR cpf ILIKE '%' || $2 || '%')
  AND ($3 = '' OR area_conhecimento ILIKE '%' || $3 || '%')
  AND ($4 = '' OR estado = $4)
ORDER BY nome;

-- name: AtualizarAvaliador :exec
UPDATE avaliadores SET
    nome = $2,
    cpf = $3,
    email = $4,
    titulacao_maxima = $5,
    area_conhecimento = $6,
    instituicao = $7,
    curriculo_resumido = $8,
    estado = $9,
    data_atualizacao = NOW()
WHERE id = $1;

-- name: ContarPropostasPorAvaliador :one
SELECT COUNT(*) FROM pareceres WHERE avaliador_id = $1;
