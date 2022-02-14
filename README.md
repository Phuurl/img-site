# img-site

## Diagram
```
┌────► Uploader             ┌───────────────┐          ┌────────────┐
│          │                │               │  Access  │            │
│          │                │ LoggingBucket │◄─────────┤ CloudFront │◄──────────── Viewer
│          │                │               │  Logs    │            │
│          │ File           └───────────────┘          └────────────┘
│          │ Upload                                           ▲
│          │                                                  │
│          │                                                  │
│          ▼                                                  ▼
│  ┌──────────────┐        ┌────────────────┐         ┌───────────────┐
│  │              │        │                │         │               │
│  │ UploadBucket ├───────►│ IngestFunction ├─────────┤ HostingBucket │
│  │              │        │                ├───┐     │               │
│  └──────────────┘        └────────────────┘   │     └───────────────┘
│                           Generates HTML      └────┐
│                           Copies to HostingBucket  │
│                           Sends SNS notification   │
│                                                    │
└────────────────────────────────────────────────────┘
      SNS email notification with generated link
```

## Deploy
```shell
# Build and deploy SAM stack
sam build --use-container
sam deploy --guided

# Now upload error.html and robots.txt to the deployed HostingBucket
aws s3 cp error.html robots.txt s3://<name of HostingBucket>/
```