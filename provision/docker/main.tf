locals {
  docker_context_path = "~/gosource/src/aceso/"
}

# Build docker image
resource "docker_image" "app" {
  name = "gcr.io/${var.gcp_project_id}/ascale:latest"
  build {
    context = "../../"
    dockerfile = "Dockerfile"
    build_args = {
    }
  }
  triggers = {
  }
}

# Push docker image
resource "docker_registry_image" "app" {
  name = docker_image.app.name
  triggers = {
    image_digest = docker_image.app.image_id
  }
}

