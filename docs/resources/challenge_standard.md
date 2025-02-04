---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ctfd_challenge_standard Resource - terraform-provider-ctfd"
subcategory: ""
description: |-
  CTFd is built around the Challenge resource, which contains all the attributes to define a part of the Capture The Flag event.
  It is the first historic implementation of its kind, with basic functionalities.
---

# ctfd_challenge_standard (Resource)

CTFd is built around the Challenge resource, which contains all the attributes to define a part of the Capture The Flag event.

It is the first historic implementation of its kind, with basic functionalities.

## Example Usage

```terraform
resource "ctfd_challenge_standard" "http" {
  name        = "My Challenge"
  category    = "misc"
  description = "..."
  value       = 500

  topics = [
    "Misc"
  ]
  tags = [
    "misc",
    "basic"
  ]
}

resource "ctfd_flag" "http_flag" {
  challenge_id = ctfd_challenge_standard.http.id
  content      = "CTF{some_flag}"
}

resource "ctfd_hint" "http_hint_1" {
  challenge_id = ctfd_challenge_standard.http.id
  content      = "Some super-helpful hint"
  cost         = 50
}

resource "ctfd_hint" "http_hint_2" {
  challenge_id = ctfd_challenge_standard.http.id
  content      = "Even more helpful hint !"
  cost         = 50
  requirements = [ctfd_hint.http_hint_1.id]
}

resource "ctfd_file" "http_file" {
  challenge_id = ctfd_challenge_standard.http.id
  name         = "image.png"
  contentb64   = filebase64(".../image.png")
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `category` (String) Category of the challenge that CTFd groups by on the web UI.
- `description` (String) Description of the challenge, consider using multiline descriptions for better style.
- `name` (String) Name of the challenge, displayed as it.
- `value` (Number) The value (points) of the challenge once solved.

### Optional

- `attribution` (String) Attribution to the creator(s) of the challenge.
- `connection_info` (String) Connection Information to connect to the challenge instance, useful for pwn, web and infrastructure pentests.
- `max_attempts` (Number) Maximum amount of attempts before being unable to flag the challenge.
- `next` (Number) Suggestion for the end-user as next challenge to work on.
- `requirements` (Attributes) List of required challenges that needs to get flagged before this one being accessible. Useful for skill-trees-like strategy CTF. (see [below for nested schema](#nestedatt--requirements))
- `state` (String) State of the challenge, either hidden or visible.
- `tags` (List of String) List of challenge tags that will be displayed to the end-user. You could use them to give some quick insights of what a challenge involves.
- `topics` (List of String) List of challenge topics that are displayed to the administrators for maintenance and planification.

### Read-Only

- `id` (String) Identifier of the challenge.

<a id="nestedatt--requirements"></a>
### Nested Schema for `requirements`

Optional:

- `behavior` (String) Behavior if not unlocked, either hidden or anonymized.
- `prerequisites` (List of String) List of the challenges ID.
