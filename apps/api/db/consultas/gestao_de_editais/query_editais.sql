-- name: InserirEdital :one
INSERT INTO editais (nome, descricao, data_inicio, data_fim, status, tipo_chamada, nota_de_corte, valor_global)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, nome, descricao, data_inicio, data_fim, status, tipo_chamada, nota_de_corte, valor_global, criado_em;

-- name: BuscarEditalPorID :one
SELECT id, nome, descricao, data_inicio, data_fim, status, tipo_chamada, nota_de_corte, valor_global, criado_em
FROM editais
WHERE id = $1;

-- name: ListarEditais :many
SELECT id, nome, descricao, data_inicio, data_fim, status, tipo_chamada, nota_de_corte, valor_global, criado_em
FROM editais
WHERE ($1 = '' OR nome ILIKE '%' || $1 || '%')
  AND ($2 = '' OR status = $2)
  AND ($3 = '' OR tipo_chamada = $3)
ORDER BY criado_em DESC;

-- name: AtualizarEdital :exec
UPDATE editais
SET nome = $2, descricao = $3, data_inicio = $4, data_fim = $5, status = $6, tipo_chamada = $7, nota_de_corte = $8, valor_global = $9
WHERE id = $1;

-- name: DeletarEdital :exec
DELETE FROM editais WHERE id = $1;
