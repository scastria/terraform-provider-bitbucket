# Resource: bitbucket_environment
Represents an environment in a repository.
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
resource "bitbucket_environment" "example" {
  repository_id = bitbucket_repository.Repo.id
  name = "My Environment"
  type = "Production"
  is_admin_only = true
}
```
## Argument Reference
* `repository_id` - **(Required, ForceNew, String)** The id of the repository.
* `name` - **(Required, String)** The name of the environment.
* `type` - **(Required, ForceNew, String)** The type of the environment. Allowed values: `Test`, `Staging`, `Production`
* `is_admin_only` - **(Optional, Boolean)** Whether deployment to the environment requires admin permissions. Default: `false`
* `use_existing` - **(Optional, Boolean, IgnoreDiffs)** During a CREATE only, look for an existing environment with the same `name`.  Prevents the need for an import. Default: `false`
## Attribute Reference
* `id` - **(String)** Same as `repository_id`:`uuid`
* `uuid` - **(String)** Uuid of the environment
## Import
Environments can be imported using a proper value of `id` as described above
