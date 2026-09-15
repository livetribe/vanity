# vanity
Go vanity server

[![Build Status](https://github.com/livetribe/vanity/actions/workflows/ci.yml/badge.svg)](https://github.com/livetribe/vanity/actions/workflows/ci.yml) 
[![Go Report Card](https://goreportcard.com/badge/github.com/livetribe/vanity)](https://goreportcard.com/report/github.com/livetribe/vanity) 
[![Documentation](https://godoc.org/l7e.io/vanity?status.svg)](http://godoc.org/l7e.io/vanity) 
[![codecov](https://codecov.io/gh/livetribe/vanity/branch/master/graph/badge.svg)](https://codecov.io/gh/livetribe/vanity)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/livetribe/vanity.svg?style=social)](https://github.com/livetribe/vanity/tags)

## Concepts<sup>[1](#inspiration)</sup>
- VCS is Version Control System (such as 'git')
- Repo root is the root path the source code repository (such as 'https://github.com/livetribe')
- Domain is the internet address where the Go Vanity server is hosted (such as
  `m4o.io` or `l7e.io`). Domain is deduced from HTTP request.
- Path is the path component of the Go package (such as `/cmd/tcpproxy` in
  `kkn.fi/cmd/tcpproxy`)

## Features
- Redirects browsers to `pkg.go.dev`, configurable to `godoc.org`
- Redirects Go tool to VCS
- Redirects HTTP to HTTPS
- Configurable `log/slog` logger. The handler takes the logger from the request
  context, and falls back to `slog.Default()`. The `vanity` command writes
  Cloud Logging structured JSON to stderr, at the level that `--log-level` sets.

## Installation
```
go get l7e.io/vanity
```

## Specification
- [Go 1.4 Custom Import Path Checking](https://docs.google.com/document/d/1jVFkZTcYbNLaTxXD9OcGfn7vYv5hWtPx9--lTx1gPMs/edit)

___
<a name="inspiration">1</a>: Inspired by [kkn.fi/vanity](https://github.com/kare/vanity).
