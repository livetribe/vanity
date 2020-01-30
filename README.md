# vanity
Go vanity server

[![Build Status](https://travis-ci.org/livetribe/vanity.svg?branch=master)](https://travis-ci.org/livetribe/vanity) 
[![Go Report Card](https://goreportcard.com/badge/github.com/livetribe/vanity)](https://goreportcard.com/report/github.com/livetribe/vanity) 
[![Documentation](https://godoc.org/github.com/livetribe/vanity?status.svg)](http://godoc.org/github.com/livetribe/vanity) 
[![Coverage Status](https://coveralls.io/repos/github/livetribe/vanity/badge.svg)](https://coveralls.io/github/livetribe/vanity)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/livetribe/vanity.svg?style=social)](https://github.com/livetribe/vanity/tags)

## Concepts<sup>[1](#inspiration)</sup>
- VCS is Version Control System (such as 'git')
- Repo root is the root path the source code repository (such as 'https://github.com/livetribe')
- Domain is the internet address where the Go Vanity server is hosted (such as
  m4o.io or l7e.io). Domain is deduced from HTTP request.
- Path is the path component of the Go package (such as /cmd/tcpproxy in
  kkn.fi/cmd/tcpproxy)

## Features
- Redirects browsers to pkg.go.dev, configurable to godoc.org
- Redirects Go tool to VCS
- Redirects HTTP to HTTPS
- Configurable logger which is fully compatible with standard log package.
  Stdout is default.

## Installation
```
go get l7e.io/vanity
```

## Specification
- [Go 1.4 Custom Import Path Checking](https://docs.google.com/document/d/1jVFkZTcYbNLaTxXD9OcGfn7vYv5hWtPx9--lTx1gPMs/edit)

___
<a name="inspiration">1</a>: Inspired by [kkn.fi/vanity](https://github.com/kare/vanity).
