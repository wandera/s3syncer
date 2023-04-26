# s3syncer

[![Test](https://github.com/wandera/s3syncer/actions/workflows/test.yml/badge.svg)](https://github.com/wandera/s3syncer/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/wandera/s3syncer)](https://goreportcard.com/report/github.com/wandera/s3syncer)
[![GitHub release](https://img.shields.io/github/release/wandera/s3syncer.svg)](https://github.com/wandera/s3syncer/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/wandera/s3syncer/blob/master/LICENSE)

S3 syncer - Easy way how to sync local directory with S3.

## Usage

```text

Usage:
  s3syncer [flags]

Flags:
  -f, --folder string      folder to watch
  -h, --help               help for s3syncer
  -l, --log-level string   command log level (options: [panic fatal error warning info debug trace]) (default "info")
  -p, --s3-path string     S3 path (s3://<bucket name>/<path>)
  -r, --s3-region string   S3 region (default "eu-west-1")
```
