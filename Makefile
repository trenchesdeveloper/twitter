mock:
	mockery --all --keeptree


migrate-up:
	migrate -source file://postgres/migrations -database "postgres://postgres:postgres@localhost:5432/twitter_clone_dev?sslmode=disable" up

migrate-down:
	migrate -source file://postgres/migrations -database "postgres://postgres:postgres@localhost:5432/twitter_clone_dev?sslmode=disable" down

migrate-drop:
	migrate -source file://postgres/migrations -database "postgres://postgres:postgres@localhost:5432/twitter_clone_dev?sslmode=disable" drop

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir postgres/migrations $$name

run:
	go run cmd/graphqlserver/main.go


generate:
	go run github.com/99designs/gqlgen generate --config graph/gqlgen.yml