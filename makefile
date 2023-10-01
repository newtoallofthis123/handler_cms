build:
	@go build -o bin/cms.exe

run: build
	@bin/cms.exe