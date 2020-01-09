# Go parameters
GOCMD=go

all: vet

# go vet:format check, bug check
vet:
	@$(GOCMD) vet `go list ./...`
