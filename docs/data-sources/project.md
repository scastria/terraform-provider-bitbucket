# Data Source: bitbucket_project
Represents a project
## Example usage
```hcl
data "bitbucket_project" "example" {
  key = "MyProjectKey"
}
```
## Argument Reference
* `key` - **(Optional, String)** The key of the project.
* `contains_repository_name` - **(Optional, String)** The name of a repository that is contained within the project.
## Attribute Reference
* `id` - **(String)** The UUID of the project.
