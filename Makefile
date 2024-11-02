.SILENT:

run:
	go run cmd/sso/main.go --config config/local.yml

# просто "build" почему-то не работал
build-app:
	go build -o build/main -a cmd/sso/main.go && echo "Build completed!"

build-app-with-logs:
	go build -o build/main -a -v -x cmd/sso/main.go && echo "Build completed!"

create-migration:
	# Создаст два файла миграции с именем указанным в параметре name или по имени текущей ветки git
	# make create-migration
	# make create-migration name="<migration name>"
	go run cmd/migrator/create/main.go --migration_name=$(or $(name),$(shell git rev-parse --abbrev-ref HEAD))

migrate:
	CONFIG_PATH="config/local.yml" go run cmd/migrator/up/main.go

migrate-down-to:
	if [ -z "$(version)" ]; then \
		echo "Error: 'version' parameter is required."; \
		exit 1; \
	fi
	CONFIG_PATH="config/local.yml" go run cmd/migrator/down-to/main.go $(version)

# Обновит пакет до последней версии
update-protos:
	go get github.com/Felya-a/chat-app-protos

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=tmp/coverage.txt ./...; \
	go tool cover -html=tmp/coverage.txt -o tmp/coverage.html;

ginkgo:
	ginkgo -v ./...

ginkgo-watch:
	ginkgo watch -v ./...