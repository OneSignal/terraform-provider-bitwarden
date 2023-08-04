# Bitwarden Terraform Provider

_This template repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)._

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

This provider is intended to manage users and groups in Bitwarden for an organisation.

```hcl
provider "bitwarden" {
  # Configure the authentication values needed to authenticate with Bitwarden API.
  # More information can be found here: https://bitwarden.com/help/public-api/#authentication
  client_id     = "client_id"
  client_secret = "client_api_secret"
}

resource "bitwarden_group" "example" {
  name = "example"
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.
If you like to test it out, you will need to create a `~/.terraformrc` file for the `dev_overrides`
```hcl
provider_installation {

  dev_overrides {
      "onesignal/bitwarden" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
export BITWARDEN_CLIENT_ID="organization.xxxx"
export BITWARDEN_CLIENT_SECRET="xxx"
make testacc
```


## TODOs

- [ ] implement group members
- [ ] gracefully handle 404 when group is manually deleted? How do we do that?
- [ ] ensure urls don't end on trailing /, use validators?
- [ ] Add docs to group resource
