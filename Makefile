SSH_HOST ?= ubuntu@137.131.166.115
SSH_KEY ?= ~/.ssh/ssh-key-2026-05-26.key
COMPOSE_FILE ?= docker-compose.coolify.yml
ENV_FILE ?= .env
PROJECT_DIR ?= /home/ubuntu/apps/fapitec

ssh:
	ssh -i $(SSH_KEY) $(SSH_HOST)

status:
	ssh -i $(SSH_KEY) $(SSH_HOST) "sudo docker ps --filter 'name=fapitec' --format 'table {{.Names}}\t{{.Status}}'"

logs-api:
	ssh -i $(SSH_KEY) $(SSH_HOST) "sudo docker logs fapitec-api --tail 30"

logs-web:
	ssh -i $(SSH_KEY) $(SSH_HOST) "sudo docker logs fapitec-web --tail 30"

logs-casdoor:
	ssh -i $(SSH_KEY) $(SSH_HOST) "sudo docker logs fapitec-casdoor --tail 30"

logs-postgres:
	ssh -i $(SSH_KEY) $(SSH_HOST) "sudo docker logs fapitec-postgres --tail 30"

restart-api:
	ssh -i $(SSH_KEY) $(SSH_HOST) "sudo docker restart fapitec-api"

restart-web:
	ssh -i $(SSH_KEY) $(SSH_HOST) "sudo docker restart fapitec-web"

restart-casdoor:
	ssh -i $(SSH_KEY) $(SSH_HOST) "sudo docker restart fapitec-casdoor"

restart-all: restart-api restart-web restart-casdoor

deploy:
	ssh -i $(SSH_KEY) $(SSH_HOST) "cd $(PROJECT_DIR) && sudo docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) up -d --build"

backup:
	ssh -i $(SSH_KEY) $(SSH_HOST) "sudo docker exec fapitec-postgres pg_dumpall -U fapitec -c > ~/backup-$$$$(date +%Y%m%d-%H%M%S).sql && echo 'Backup criado em ~/' && ls -la ~/backup-*.sql | tail -1"

psql:
	ssh -i $(SSH_KEY) $(SSH_HOST) "sudo docker exec -it fapitec-postgres psql -U fapitec -d fapitec"

down:
	ssh -i $(SSH_KEY) $(SSH_HOST) "cd $(PROJECT_DIR) && sudo docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) down"

help:
	@echo "Comandos disponiveis:"
	@echo "  make ssh            - Conectar via SSH na VM"
	@echo "  make status         - Status dos containers FAPITEC"
	@echo "  make logs-api       - Logs da API"
	@echo "  make logs-web       - Logs do frontend"
	@echo "  make logs-casdoor   - Logs do Casdoor"
	@echo "  make logs-postgres  - Logs do PostgreSQL"
	@echo "  make restart-api    - Reiniciar API"
	@echo "  make restart-web    - Reiniciar web"
	@echo "  make restart-all    - Reiniciar todos os servicos"
	@echo "  make deploy         - Rebuild e deploy do stack"
	@echo "  make backup         - Backup completo do banco"
	@echo "  make psql           - Conectar no PostgreSQL (CLI)"
	@echo "  make down           - Parar todos os servicos"
	@echo ""
	@echo "Variaveis (podem ser sobrescritas):"
	@echo "  SSH_HOST=    (padrao: $(SSH_HOST))"
	@echo "  SSH_KEY=     (padrao: $(SSH_KEY))"
