#!/usr/bin/env python3

import sys

import boto3

if len(sys.argv) != 2:
    print("Missing / invalid arguments\nRequires bucket name as 1st and only script argument.")
bucket_name = sys.argv[1]
s3 = boto3.resource('s3')
bucket = s3.Bucket(bucket_name)
bucket.object_versions.delete()
