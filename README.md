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

| Parameter             | Description                                                                            | Example                                                     |
|-----------------------|----------------------------------------------------------------------------------------|-------------------------------------------------------------|
| `UploadBucketName`    | Name of the upload S3 bucket to create                                                 | `img-bucket-name`                                           |
| `NotificationEmail`   | Email address to notify on processing completion                                       | `someone@example.com`                                       |
| `SiteName`            | Name for site - used in HTML template for title                                        | `CoolImages`                                                |
| `Domain`              | Domain name that the image site will be hosted on                                      | `img.example.com`                                           |
| `CertArn`             | ARN of the ACM certificate to use with the CloudFront distribution                     | `arn:aws:acm:us-east-1:0123456789:certificate/abc-1234-def` |
| `CreateUploadIamUser` | Whether to create an IAM user for upload to the ingest bucket - user created if `true` | `false`                                                     |


## Deploy
```shell
# Build and deploy SAM stack
sam build --use-container
sam deploy --guided

# Now upload error.html and robots.txt to the deployed HostingBucket
aws s3 cp error.html robots.txt s3://<name of HostingBucket>/
```