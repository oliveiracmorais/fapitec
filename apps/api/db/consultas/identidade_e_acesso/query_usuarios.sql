-- name: InserirUsuario :one
INSERT INTO usuarios (nome, cpf, email, senha_hash, eh_estrangeiro)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, nome, cpf, email, senha_hash, eh_estrangeiro, tentativas, bloqueado_ate, criado_em;

-- name: BuscarUsuarioPorCPF :one
SELECT id, nome, cpf, email, senha_hash, eh_estrangeiro, tentativas, bloqueado_ate, criado_em
FROM usuarios
WHERE cpf = $1;

-- name: BuscarUsuarioPorEmail :one
SELECT id, nome, cpf, email, senha_hash, eh_estrangeiro, tentativas, bloqueado_ate, criado_em
FROM usuarios
WHERE email = $1;

-- name: AtualizarTentativasUsuario :exec
UPDATE usuarios
SET tentativas = $2, bloqueado_ate = $3
WHERE id = $1;

-- name: AtualizarSenhaUsuario :exec
UPDATE usuarios
SET senha_hash = $2
WHERE id = $1;
