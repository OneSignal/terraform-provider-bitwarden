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

output "test" {
  value = bitwarden_group.example
}
