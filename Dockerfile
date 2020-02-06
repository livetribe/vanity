# Copyright 2020 the original author or authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM alpine:latest as certs

RUN apk --update add ca-certificates

FROM golang:1.13 as builder

RUN GO111MODULE=on go get github.com/ahmetb/govvv@master

WORKDIR /go/src/l7e.io/vanity

COPY . .

ARG VERSION="<unspecified>"

RUN CGO_ENABLED=0 GOBIN=/go/bin GOOS=linux GO111MODULE=on go install -ldflags="$(govvv -flags -version $VERSION -pkg l7e.io/vanity/cmd/vanity/cli)" l7e.io/vanity/cmd/vanity

FROM alpine:3.10

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/vanity /usr/local/bin

ENV VANITY_SPANNER_DATABASE=""
ENV VANITY_DATASTORE_PROJECT_ID=""
ENV VANITY_GOOGLE_API_KEY=""
ENV VANITY_GOOGLE_API_CREDENTIALS_FILE=""
ENV VANITY_GOOGLE_API_CREDENTIALS_JSON=""
ENV GOOGLE_APPLICATION_CREDENTIALS=""

EXPOSE 8080
EXPOSE 8081
EXPOSE 8082
EXPOSE 9100

ENTRYPOINT ["vanity"]
CMD ["--help"]
