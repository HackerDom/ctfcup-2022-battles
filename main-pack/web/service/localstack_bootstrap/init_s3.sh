#!/usr/bin/env bash
set -x
awslocal s3 mb s3://users
awslocal s3 mb s3://notes
awslocal s3 mb s3://messages
set +x