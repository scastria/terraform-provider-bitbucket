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

# output "Test" {
#   value = data.bitbucket_project.Proj.id
# }
