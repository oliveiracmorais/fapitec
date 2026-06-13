-- Migration 007: Substitui proponente_endereco (TEXT) por colunas estruturadas
-- Depende: 005_criar_propostas.sql
-- Nota: Projeto nao esta em producao — dados existentes sao descartados.

ALTER TABLE propostas
    DROP COLUMN IF EXISTS proponente_endereco,
    ADD COLUMN IF NOT EXISTS proponente_cep VARCHAR(8) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS proponente_logradouro VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS proponente_numero VARCHAR(20) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS proponente_complemento VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS proponente_bairro VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS proponente_cidade VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS proponente_uf VARCHAR(2) NOT NULL DEFAULT '';

ALTER TABLE versoes_proposta
    DROP COLUMN IF EXISTS proponente_endereco,
    ADD COLUMN IF NOT EXISTS proponente_cep VARCHAR(8) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS proponente_logradouro VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS proponente_numero VARCHAR(20) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS proponente_complemento VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS proponente_bairro VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS proponente_cidade VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS proponente_uf VARCHAR(2) NOT NULL DEFAULT '';
