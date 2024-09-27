run:
	go run cmd/main.go

up:
	sql-migrate new up

down:
	sql-migrate down

swag:
	swag init -d internal/handlers/ -g router.go --parseDependency --parseDepth 3