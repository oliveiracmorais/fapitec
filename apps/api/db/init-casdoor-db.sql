-- Cria database separada para o Casdoor (schema isolado no mesmo PostgreSQL)
SELECT 'CREATE DATABASE casdoor'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'casdoor');
