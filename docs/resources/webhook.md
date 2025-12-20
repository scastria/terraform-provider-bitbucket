# Resource: bitbucket_webhook
Represents a webhook on a repository.
## Example usage
```hcl
data "bitbucket_project" "Proj" {
  key = "MyProjectKey"
}
resource "bitbucket_repository" "Repo" {
  project_id = data.bitbucket_project.Proj.id
  key = "MyRepoKey"
  name = "My-Repo"
  is_private = true
}
resource "bitbucket_webhook" "example" {
  repository_id = bitbucket_repository.Repo.id
  url = "https://example.com/webhook"
  title = "My Webhook"
  events = [
    "pullrequest:fulfilled",
    "pullrequest:rejected"
  ]
}
```
## Argument Reference
* `repository_id` - **(Required, ForceNew, String)** The id of the repository.
* `url` - **(Required, String)** The url of the webhook.
* `title` - **(Optional, String)** The title of the webhook.
* `events` - **(Required, List of String)** The events that should cause the webhook to be triggered.
* `is_active` - **(Optional, Boolean)** Whether the webhook is active. Default: `true`
* `use_existing` - **(Optional, Boolean, IgnoreDiffs)** During a CREATE only, look for an existing webhook with the same url.  Prevents the need for an import. Default: `false`
## Attribute Reference
* `id` - **(String)** Same as `repository_id`:`uuid`
* `uuid` - **(String)** Uuid of the webhook
## Import
Webhooks can be imported using a proper value of `id` as described above
