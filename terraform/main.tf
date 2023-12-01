variable "project-id" {
  type = string
  default = "festive-antenna-402105"
}
 
provider "google" {
  # credentials = file("$HOME/.config/gcloud/application_default_credentials.json")
  project     = var.project-id
  region      = "asia-southeast2"
}

# // create a new bucket
# resource "google_storage_bucket" "hijalearn-bucket" {
#   name          = "hijalearn-bucket-402105"
#   location      = "asia-southeast2"
#   storage_class = "STANDARD"
#   force_destroy = true
# }
#
# // firestore native mode
# resource "google_project_service" "firestore" {
#   project = var.project-id
#   service = "firestore.googleapis.com"
# }
#
# resource "google_firestore_database" "database" {
#   name        = "hijalearn-db"
#   project     = var.project-id
#   location_id = "asia-southeast2"
#   type        = "FIRESTORE_NATIVE"
#
#   depends_on = [google_project_service.firestore]
# }

// create cloud run service
resource "google_cloud_run_v2_service" "default" {
  name     = "hijalearn-service"
  location = "asia-southeast2"
  template {
      containers {
        image = "gcr.io/${var.project-id}/hijalearn-testbuild:latest"
      }
  }
}

resource "google_cloud_run_v2_service_iam_member" "noauth" {
  location = google_cloud_run_v2_service.default.location
  name     = google_cloud_run_v2_service.default.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
