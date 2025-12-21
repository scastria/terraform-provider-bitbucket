# Resource: bitbucket_dynamic_pipelines_provider
Represents the provider for dynamic pipelines.
## Example usage
```hcl
data "bitbucket_project" "Proj" {
  key = "MyProjectKey"
}
resource "bitbucket_repository" "Repo" {
  project_id = data.bitbucket_project.Proj.id
  name = "My Repo"
  is_private = true
}
resource "bitbucket_dynamic_pipelines_provider" "example" {
  repository_id = bitbucket_repository.Repo.id
  provider_id = "XXXX"
}
```
## Argument Reference
* `repository_id` - **(Required, ForceNew, String)** The id of the repository.
* `provider_id` - **(Required, String)** The id of the provider.
## Attribute Reference
* `id` - **(String)** The UUID of the repository.
## Import
Dynamic pipeline providers can be imported using a proper value of `id` as described above
