-- name: InserirAtribuicao :one
INSERT INTO atribuicoes_editais (avaliador_id, edital_id, data_inicio, data_fim, status_convite, hash_anonimizacao, criado_em)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: BuscarAtribuicaoPorID :one
SELECT * FROM atribuicoes_editais WHERE id = $1;

-- name: ListarAtribuicoesPorAvaliador :many
SELECT * FROM atribuicoes_editais WHERE avaliador_id = $1 ORDER BY criado_em DESC;

-- name: ListarAtribuicoesPorEdital :many
SELECT * FROM atribuicoes_editais WHERE edital_id = $1 ORDER BY criado_em DESC;

-- name: ListarAtribuicoes :many
SELECT * FROM atribuicoes_editais
WHERE ($1 = 0 OR avaliador_id = $1)
  AND ($2 = 0 OR edital_id = $2)
  AND ($3 = '' OR status_convite = $3)
ORDER BY criado_em DESC;

-- name: AtualizarAtribuicaoStatus :exec
UPDATE atribuicoes_editais SET status_convite = $2 WHERE id = $1;

-- name: ContarAtribuicoesAtivasPorAvaliador :one
SELECT COUNT(*) FROM atribuicoes_editais
WHERE avaliador_id = $1
  AND status_convite = 'aceito'
  AND data_inicio <= NOW()
  AND data_fim >= NOW();
