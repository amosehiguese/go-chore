make.PHONY: clean
clean:
	rm -rf ./bin

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
	protoc proto/v1/*.proto --go_out=. --go_opt=paths=source_relative --proto_path=.