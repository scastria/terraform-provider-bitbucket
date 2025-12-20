terraform {
  required_providers {
    bitbucket = {
      source = "github.com/scastria/bitbucket"
    }
  }
}

provider "bitbucket" {
}

data "bitbucket_project" "Proj" {
  key = "DATA"
}

resource "bitbucket_repository" "Repo" {
  project_id = data.bitbucket_project.Proj.id
  key         = "shawn-test"
  name = "Shawn-Test"
  is_private   = true
  # use_existing = true
}

resource "bitbucket_pipelines_config" "PipeConfig" {
  repository_id = bitbucket_repository.Repo.id
  is_enabled       = true
}

# resource "bitbucket_webhook" "Hook" {
#   repository_id = bitbucket_repository.Repo.id
#   url = "XXXX"
#   title = "PR_CLEANUP"
#   events = [
#     "pullrequest:fulfilled",
#     "pullrequest:rejected"
#   ]
#   use_existing = true
# }

# resource "bitbucket_dynamic_pipelines_provider" "DynoProvider" {
#   repository_id = bitbucket_repository.Repo.id
#   provider_id = "XXXX"
# }

# output "Test" {
#   value = data.bitbucket_project.Proj.id
# }
