variable "instance_types" {
  type        = "list"
  description = "list of instance types"
}

variable "project_id" {
  type        = "string"
  description = "Name/ID of the gcp project in which the operations has to be carried"
}

variable "zone" {
  type        = "string"
  description = "The zone from where the details has to be fetched"
}