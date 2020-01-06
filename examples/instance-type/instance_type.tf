resource "gcp_instance_type" "instance" {
  count        = "${length(var.instance_types)}"
  machine_type = "${element(var.instance_types, count.index)}"
}