# Resource: bitbucket_repository
Represents a repository within a project
## Example usage
```hcl
data "bitbucket_project" "Proj" {
  key = "MyProjectKey"
}
resource "bitbucket_repository" "example" {
  project_id = data.bitbucket_project.Proj.id
  name = "My Repo"
  is_private = true
}
```
## Argument Reference
* `project_id` - **(Required, String)** The id of the project.
* `name` - **(Required, String)** The name of the repository.
* `is_private` - **(Optional, Boolean)** Whether the repository is private. Default: `true`
* `use_existing` - **(Optional, Boolean, IgnoreDiffs)** During a CREATE only, look for an existing repository with the same `name`.  Prevents the need for an import. Default: `false`
## Attribute Reference
* `id` - **(String)** The UUID of the repository.
## Import
Repositories can be imported using a proper value of `id` as described above
