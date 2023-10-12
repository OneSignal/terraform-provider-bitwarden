terraform {
  required_providers {
    bitwarden = {
      source = "onesignal/bitwarden"
    }
  }
}

# export BITWARDEN_CLIENT_ID=""
# export BITWARDEN_CLIENT_SECRET=""
provider "bitwarden" {
}

resource "bitwarden_group" "example" {
  name       = "example"
  access_all = true
}

resource "bitwarden_member" "example" {
  type        = 2
  email       = "niels@fake.com"
  external_id = "external-niels-id"
  access_all  = false
}

output "test_group" {
  value = bitwarden_group.example
}

output "test_member" {
  value = bitwarden_member.example
}
