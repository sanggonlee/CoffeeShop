.PHONY: all db_clean

db_clean:
	@docker rm coffeeshop_data && echo "Data cleanup complete." || echo "No data container to cleanup"
	@docker stop coffeeshop_db && docker rm coffeeshop_db && echo "DB cleanup complete." || echo "No DB container. Nothing to do."

db_init: db_clean
	@docker create -v /coffeeshop_data --name coffeeshop_data postgres
	@docker run --name coffeeshop_db --volumes-from coffeeshop_data -e POSTGRES_HOST=database -e POSTGRES_USER=coffee -e POSTGRES_PASSWORD=needcaffeine -e POSTGRES_DB=coffeeshop -p 5432:5432 -d postgres
	@echo "DB container initialized"

db_migrate:
	@goose -dir db/migrations postgres "user=coffee password=needcaffeine host=localhost port=5432 dbname=coffeeshop sslmode=disable" up

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o coffeeshop .

run:
	@go run main.go