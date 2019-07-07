src=./source/
bin=./bin/

.PHONY: build_plugin
build_plugin:
	@rm -rf $(bin)
	@mkdir $(bin)

	@echo "Building vintagemonster.go plugin ..."
	@go build -v -buildmode=plugin -o $(bin)vintagemonster.so $(src)vintagemonster.go
	@echo "All plugins are build successfully :)";


.PHONY: test
test: build_plugin
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "test report stored in coverage.html";

.PHONY: run
run: build_plugin
	@echo "Running main.go ...."
	@go run main.go