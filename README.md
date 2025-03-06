# Terraform Provider for Hetzner Robot

This is a Terraform provider for interacting with the Hetzner Robot API. It allows management of Hetzner resources through Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
```sh
$ git clone https://github.com/lazureykis/terraform-provider-hetzner-robot
```

2. Enter the repository directory
```sh
$ cd terraform-provider-hetzner-robot
```

3. Build the provider
```sh
$ go build -o terraform-provider-hetznerrobot
```

## Using the provider

```hcl
terraform {
  required_providers {
    hetznerrobot = {
      source = "lazureykis/hetznerrobot"
    }
  }
}

provider "hetznerrobot" {
  username = "your-robot-username"
  password = "your-robot-password"
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go build`. This will build the provider and put the provider binary in the current directory.

```sh
$ go build
```

To generate or update documentation, run `go generate`.

```sh
$ go generate
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
