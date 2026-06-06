-- name: InserirEdital :one
INSERT INTO editais (nome, descricao, data_inicio, data_fim, status, tipo_chamada, nota_de_corte, valor_global, modelo_formulario, relatorios_exigidos, titulo_minimo_elegibilidade, exige_empresa, porte_empresa, enquadramento_empresa, documentos_obrigatorios)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING id, nome, descricao, data_inicio, data_fim, status, tipo_chamada, nota_de_corte, valor_global, modelo_formulario, relatorios_exigidos, titulo_minimo_elegibilidade, exige_empresa, porte_empresa, enquadramento_empresa, documentos_obrigatorios, criado_em;

-- name: BuscarEditalPorID :one
SELECT id, nome, descricao, data_inicio, data_fim, status, tipo_chamada, nota_de_corte, valor_global, modelo_formulario, relatorios_exigidos, titulo_minimo_elegibilidade, exige_empresa, porte_empresa, enquadramento_empresa, documentos_obrigatorios, criado_em
FROM editais
WHERE id = $1;

-- name: ListarEditais :many
SELECT id, nome, descricao, data_inicio, data_fim, status, tipo_chamada, nota_de_corte, valor_global, modelo_formulario, relatorios_exigidos, titulo_minimo_elegibilidade, exige_empresa, porte_empresa, enquadramento_empresa, documentos_obrigatorios, criado_em
FROM editais
WHERE ($1 = '' OR nome ILIKE '%' || $1 || '%')
  AND ($2 = '' OR status = $2)
  AND ($3 = '' OR tipo_chamada = $3)
ORDER BY criado_em DESC;

-- name: AtualizarEdital :exec
UPDATE editais
SET nome = $2, descricao = $3, data_inicio = $4, data_fim = $5, status = $6, tipo_chamada = $7, nota_de_corte = $8, valor_global = $9,
    modelo_formulario = $10, relatorios_exigidos = $11, titulo_minimo_elegibilidade = $12, exige_empresa = $13, porte_empresa = $14, enquadramento_empresa = $15, documentos_obrigatorios = $16
WHERE id = $1;

-- name: DeletarEdital :exec
DELETE FROM editais WHERE id = $1;
