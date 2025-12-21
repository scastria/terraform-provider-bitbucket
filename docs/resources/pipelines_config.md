# Resource: bitbucket_pipelines_config
Represents the configuration of pipelines for a repository.
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
resource "bitbucket_pipelines_config" "example" {
  repository_id = bitbucket_repository.Repo.id
  is_enabled = true
}
```
## Argument Reference
* `repository_id` - **(Required, ForceNew, String)** The id of the repository.
* `is_enabled` - **(Optional, Boolean)** Whether pipelines are enabled. Default: `false`
## Attribute Reference
* `id` - **(String)** The UUID of the repository.
## Import
Pipeline configs can be imported using a proper value of `id` as described above
