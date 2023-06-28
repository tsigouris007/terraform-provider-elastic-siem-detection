default: testacc doc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

doc:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	terraform fmt -recursive ./examples/