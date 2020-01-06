 output "maximum_persistent_disks" {
   value = "${gcp_instance_type.instance[0].maximum_persistent_disks}"
 }