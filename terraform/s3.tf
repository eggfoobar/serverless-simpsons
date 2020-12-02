
/*

Create our bucket

*/

resource "aws_s3_bucket" "bucket" {
  bucket = lower(format("%s-%s", var.bucket_name, substr(data.aws_caller_identity.current.user_id, 0, 8)))
  acl    = "private"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  tags = var.tags
}

/*

Lock it down!

*/

resource "aws_s3_bucket_public_access_block" "block_public" {
  bucket = aws_s3_bucket.bucket.id
  block_public_acls   = true
  block_public_policy = true
  restrict_public_buckets = true
  ignore_public_acls = true
}

/*

Upload the goods

*/

resource "aws_s3_bucket_object" "static_assets" {
  for_each = fileset(path.module, "../cosmic-ballet/frontend/**")
  bucket = aws_s3_bucket.bucket.id
  key    = trimprefix(each.value, "../cosmic-ballet/frontend/")
  source = "${path.module}/${each.value}"
  etag                   = filemd5("${path.module}/${each.value}")
  tags                   = var.tags
  server_side_encryption = "AES256"
  content_type           = "text/html"
}
