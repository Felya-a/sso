.SILENT:

run:
	go run cmd/sso/main.go --config config/local.yml

create-migration:
	# Создаст два фала миграции с именем указанным в параметре name или по имени текущей ветки git
	go run cmd/migrator/create/main.go --migration_name=$(or $(name),$(shell git rev-parse --abbrev-ref HEAD))

migrate:
	CONFIG_PATH="config/local.yml" go run cmd/migrator/up/main.go

migrate-down-to:
	if [ -z "$(version)" ]; then \
		echo "Error: 'version' parameter is required."; \
		exit 1; \
	fi
	CONFIG_PATH="config/local.yml" go run cmd/migrator/down-to/main.go $(version)