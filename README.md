# terraform-provider-elastic-siem-detection
A complete Elastic SIEM rules / exception containers / exceptions terraform provider

This repository is a provider for [Elastic SIEM](https://www.elastic.co/security/siem).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.18

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Publishing The Provider

Follow the instructions here: https://developer.hashicorp.com/terraform/tutorials/providers/provider-release-publish

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider locally

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To use it locally copy the compiled provider from `$GOPATH/bin/<COMPILED_PROVIDER>` to `/home/$USER/.terraform.d/plugins/local/elastic-siem-detection/elastic-siem-detection/1.0.0/linux_amd64/<COMPILED_PROVIDER>`.
- Replace `local` with any path of your choice.
- Replace 1st occurence of `elastic-siem-detection` with any path of your choice.
- Replace 2nd occurence of `elastic-siem-detection` with any path of your choice.
- Replace `1.0.0` with any version of your choice.
- Replace `linux_amd64` with the corresponding OS platform.
- Replace `<COMPILED_PROVIDER>` with the proper compiled binary name. Suggested to use `terraform-provider-elastic-siem-detection`.

To use it in your terraform:
```terraform
terraform {
  required_version = ">= 0.13.0"
  required_providers {
    elastic-siem-detection = {
      source = "local/elastic-siem-detection/elastic-siem-detection"
      version = "1.0.0"
    }
  }
}
```

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.
```shell
make testacc
```

## Notes

**Update**

Now supported:
- `match` and `match_any` clause in exception items.
- `wildcard` clause in exception items.

## Usage

You can find a recommended way to use this provider under the `./usage` directory.

## Credits

The original idea was taken from [https://github.com/fabiogermann/terraform-provider-elastic-siem](https://github.com/fabiogermann/terraform-provider-elastic-siem)

Special thanks to [https://github.com/gpapakyriakopoulos](https://github.com/gpapakyriakopoulos) for the debugging sessions and emotional support.
