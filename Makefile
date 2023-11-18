.PHONY: build
build: 
	go build -C cmd -o ../bin/chore

.PHONY: run
run: build
	./bin/chore $(filter-out $@, $(MAKECMDGOALS))
%:
	@true

.PHONY: compile
compile:
	protoc proto/v1/*.proto \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	--proto_path=.

.PHONY: gencert
gencert:
	cfssl gencert -initca test/ca-csr.json | cfssljson -bare ca

	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test/ca-config.json -profile=server test/server-csr.json | cfssljson -bare server

.PHONY: client
client:
	go run rosie_client/rosie_client.go -ca-cert ca.pem $(filter-out $@, $(MAKECMDGOALS))
%:
	@true

.PHONY: server
server:
	go run rosie/rosie.go -cert server.pem -key server-key.pem