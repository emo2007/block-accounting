PROJECT_DIR = $(CURDIR)
PROJECT_BIN = ${PROJECT_DIR}/bin
TOOLS_BIN = ${PROJECT_BIN}/tools

bin.build:
	mkdir -p ${PROJECT_DIR}/build
	rm -f ${PROJECT_DIR}/build/blockd
	go build -ldflags="-s -w" -o ${PROJECT_DIR}/build/blockd ${PROJECT_DIR}/cmd/main.go

d.build:
	sudo docker buildx build . -t blockd:latest

d.net:
	sudo docker network create --driver bridge --subnet=192.168.2.0/24 --attachable blockd-net

d.drop-net:
	sudo docker network rm blockd-net

up: d.build d.net
	sudo docker compose up -d

.PHONY: run.local
run.local: bin.build
	${PROJECT_DIR}/build/blockd \
		-log-level=debug \
		-log-local=true \
		-log-add-source=true \
		-rest-address=localhost:8081 \
		-db-host=localhost:8432 \
		-db-database=blockd \
		-db-user=blockd \
		-db-secret=blockd \
		-db-enable-tls=false \
		-jwt-secret=local_jwt_secret

.PHONY: run.debug
run.debug: bin.build
	${PROJECT_DIR}/build/blockd \
		-log-level=debug \
		-log-local=false \
		-log-add-source=true \
		-rest-address=localhost:8081 \
		-db-host=localhost:8432 \
		-db-database=blockd \
		-db-user=blockd \
		-db-secret=blockd \
		-db-enable-tls=false

start.d:
	sudo systemctl start docker