---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bitwarden_member Resource - terraform-provider-bitwarden"
subcategory: ""
description: |-
  The Bitwarden member resources manage the members aka users within an Bitwarden organization. We leverage the public [Bitwarden API](https://bitwarden.com/help/api/]
---

# bitwarden_member (Resource)

The Bitwarden member resources manage the members aka users within an Bitwarden organization. We leverage the public [Bitwarden API](https://bitwarden.com/help/api/]

## Example Usage

```terraform
resource "bitwarden_member" "example" {
  type        = 2
  email       = "niels@fake.com"
  external_id = "external-niels-id"
  access_all  = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `email` (String) The member's email address.
- `type` (Number) Defines Member type, needs to be one of:
    Owner = 0,
    Admin = 1,
    User = 2,
    Manager = 3,
    Custom = 4.
    See https://github.com/bitwarden/server/blob/master/src/Core/Enums/OrganizationUserType.cs

### Optional

- `access_all` (Boolean) Determines if this member can access all collections within the organization, or only the associated collections. If set to {true}, this option overrides any collection assignments
- `external_id` (String) External identifier for reference or linking this member to another system, such as a user directory

### Read-Only

- `id` (String) The member's unique identifier within the organization
- `last_updated` (String)
- `name` (String) The member's name, set from their user account profile
- `status` (Number) The member's status within the organisation, is one of the following:
    Invited = 0,
    Accepted = 1,
    Confirmed = 2,
    Revoked = -1.
    See https://github.com/bitwarden/server/blob/master/src/Core/Enums/OrganizationUserStatusType.cs