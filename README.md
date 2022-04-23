# img-site

## Diagram
```
┌────► Uploader             ┌───────────────┐  Access  ┌────────────┐
│          │                │               │  logs    │            │
│          │                │ LoggingBucket │◄─────────┤ CloudFront │◄─────── Viewer
│          │                │               │          │            │
│          │ Image          └───────────────┘          └────────────┘
│          │ upload                                           ▲
│          │                                                  │
│          │                                                  │
│          ▼                                                  ▼
│  ┌──────────────┐        ┌────────────────┐         ┌───────────────┐
│  │              │        │                │         │               │
│  │ UploadBucket ├───────►│ IngestFunction ├────────►│ HostingBucket │
│  │              │        │                ├───┐     │               │
│  └──────────────┘        └────────────────┘   │     └───────────────┘
│                           Generate HTML       └────┐
│                           Upload to HostingBucket  │
│                           Send SNS notification    │
│                                                    │
└────────────────────────────────────────────────────┘
      SNS email notification with generated link
```

## Parameters

| Parameter                        | Description                                                                                                             | Example                                                     |
|----------------------------------|-------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------|
| `UploadBucketName`               | Name of the upload S3 bucket to create                                                                                  | `img-bucket-name`                                           |
| `EmailNotificationEnabled`       | Whether to enable email notification on image processing completion - `true` or `false`                                 | `true`                                                      |
| `NotificationEmail`              | Email address to notify on processing completion                                                                        | `someone@example.com`                                       |
| `SiteName`                       | Name for site - used in HTML template for title                                                                         | `CoolImages`                                                |
| `NoDomain`                       | If `true`, a custom domain and certificate is not attached to the CloudFront distribution - useful for testing purposes | `false`                                                     |
| `Domain`                         | Domain name that the image site will be hosted on (ignored if `NoDomain` is `true`)                                     | `img.example.com`                                           |
| `CertArn`                        | ARN of the ACM certificate to use with the CloudFront distribution (ignored if `NoDomain` is `true`)                    | `arn:aws:acm:us-east-1:0123456789:certificate/abc-1234-def` |
| `CreateUploadIamUser`            | Whether to create an IAM user for upload to the ingest bucket - user created if `true`                                  | `false`                                                     |
| `CloudFrontRedirectFunctionName` | Name for the CloudFront redirect function - can leave as default unless you're deploying multiple copies of this stack  | `img_site_folder_index_redirect`                            |
| `CloudFrontCachePolicyName`      | Name for the CloudFront cache policy - can leave as default unless you're deploying multiple copies of this stack       | `img-site-cache-policy`                                     |


## Deploy
```shell
# Build and deploy SAM stack
sam build --use-container
sam deploy --guided

# Now upload error.html and robots.txt to the deployed HostingBucket
aws s3 cp error.html robots.txt s3://<name of HostingBucket>/
```