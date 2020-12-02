/*

Create our Cloudfront Distrobution

*/

resource "aws_cloudfront_origin_access_identity" "origin_access_identity" {
  comment = "OriginAccessIdentityServerlessDemo"
}

resource "aws_cloudfront_distribution" "s3_distribution" {

  origin {
    domain_name = aws_s3_bucket.bucket.bucket_regional_domain_name
    origin_id   = var.bucket_origin_id

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.origin_access_identity.cloudfront_access_identity_path
    }
  }

  enabled             = true
  is_ipv6_enabled     = true
  comment             = "Serverless Demo"
  default_root_object = "index.html"

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = var.bucket_origin_id

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  price_class = "PriceClass_200"

  restrictions {
    geo_restriction {
      restriction_type = "whitelist"
      locations        = ["US", "CA"]
    }
  }

  tags = var.tags

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}

data "aws_iam_policy_document" "web_distribution" {
  statement {
    actions = ["s3:GetObject"]
    principals {
      type        = "AWS"
      identifiers = [aws_cloudfront_origin_access_identity.origin_access_identity.iam_arn]
    }
    resources = ["${aws_s3_bucket.bucket.arn}/*"]
  }
}

resource "aws_s3_bucket_policy" "web_distribution" {
  bucket = aws_s3_bucket.bucket.id
  policy = data.aws_iam_policy_document.web_distribution.json
}
