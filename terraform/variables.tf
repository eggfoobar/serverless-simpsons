variable "tags" {
  type        = object({})
  description = "Tags to use in your AWS resourcers"
}

variable "bucket_name" {
  type        = string
  description = "Bucket name"
}

variable "bucket_origin_id" {
  type        = string
  description = "Bucket origin for cloudfront"
  default     = "my.awesome.bucket"
}
