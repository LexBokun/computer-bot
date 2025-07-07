generate:
	buf dep update
	buf dep prune
	buf lint
	#buf breaking --against ".git#subdir=."
	buf generate

linter:
	golangci-lint --config golangci.yaml