.PHONY: build run eval

build:
	@echo "Building..."
	@go build -o bin/vrp cmd/vrp/*

run:
	@echo "Running..."
	@./bin/vrp $(FILE_NAME)

eval:
	python3 evaluateShared.py --cmd ./bin/vrp --problemDir ./trainingProblems