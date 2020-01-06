provider "gcp" {
  // Provider would pick the gcp's default credentials if not specified.
  //credentials = "/path/to/credentials.json"
  project = "${var.project_id}"
  zone    = "${var.zone}"
  version = "0.0.3"
}
