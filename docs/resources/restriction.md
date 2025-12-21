# Resource: bitbucket_restriction
Represents a branch restriction on a repository.
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
resource "bitbucket_restriction" "example" {
  repository_id = bitbucket_repository.Repo.id
  kind = "require_passing_builds_to_merge"
  branch_match_kind = "branching_model"
  branch_type = "production"
  value = 1
}
```
## Argument Reference
* `repository_id` - **(Required, ForceNew, String)** The id of the repository.
* `kind` - **(Required, String)** The kind of the restriction. Allowed values: `require_approvals_to_merge`, `require_default_reviewer_approvals_to_merge`, `require_passing_builds_to_merge`, `require_commits_behind`
* `branch_match_kind` - **(Required, String)** The method to match the branch for the restriction. Allowed values: `branching_model`, `glob`
* `branch_type` - **(Optional, String)** When using `branching_model` matching, the model of the branch for the restriction. Allowed values: `feature`, `bugfix`, `release`, `hotfix`, `development`, `production`
* `pattern` - **(Optional, String)** When using `glob` matching, the wildcarded name of the branch for the restriction.
* `value` - **(Optional, Integer)** When using `kind` equal to one of: `require_approvals_to_merge`, `require_default_reviewer_approvals_to_merge`, `require_passing_builds_to_merge`, `require_commits_behind`, the numerical value of the restriction.
* `use_existing` - **(Optional, Boolean, IgnoreDiffs)** During a CREATE only, look for an existing restriction with the same `kind`, `branch_match_kind`, `branch_type`, and `pattern`.  Prevents the need for an import. Default: `false`
## Attribute Reference
* `id` - **(String)** Same as `repository_id`:`restriction_id`
* `restriction_id` - **(String)** Id of the restriction alone
## Import
Restrictions can be imported using a proper value of `id` as described above
