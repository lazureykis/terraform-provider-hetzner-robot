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
$ go build -o terraform-provider-hetzner-robot
```

## Using the provider

```hcl
terraform {
  required_providers {
    hetzner_robot = {
      source = "lazureykis/hetzner_robot"
    }
  }
}

provider "hetzner_robot" {
  username = "your-robot-username"
  password = "your-robot-password"
}
```
