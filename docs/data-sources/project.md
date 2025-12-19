# Data Source: bitbucket_project
Represents a project
## Example usage
```hcl
data "bitbucket_project" "example" {
  key = "MyProjectKey"
}
```
## Argument Reference
* `key` - **(Required, String)** The key of the project.
## Attribute Reference
* `id` - **(String)** The UUID of the project.
