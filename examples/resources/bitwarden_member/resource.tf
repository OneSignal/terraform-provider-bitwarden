resource "bitwarden_member" "example" {
  type        = 2
  email       = "niels@fake.com"
  external_id = "external-niels-id"
  access_all  = false
}
