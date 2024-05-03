bin.build:
	mkdir -p build
	rm -f build/blockd
	go build -ldflags="-s -w" -o build/blockd cmd/main.go

d.build:
	sudo docker buildx build . -t blockd:latest

.PHONY: run.local
run.local: bin.build
	./build/blockd \
		-log-level=debug \
		-log-local=true \
		-log-add-source=true \
		-rest-address=localhost:8080 \
		-db-host=localhost:8432 \
		-db-database=blockd \
		-db-user=blockd \
		-db-secret=blockd \
		-db-enable-tls=false

.PHONY: run.debug
run.debug: bin.build
	./build/blockd \
		-log-level=debug \
		-log-local=false \
		-log-add-source=true \
		-rest-address=localhost:8080 \
		-db-host=localhost:8432 \
		-db-database=blockd \
		-db-user=blockd \
		-db-secret=blockd \
		-db-enable-tls=false