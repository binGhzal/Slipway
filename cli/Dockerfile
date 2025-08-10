# SPDX-License-Identifier: Apache-2.0
FROM golang:1.22 as build
WORKDIR /src
COPY ../ go.mod go.sum ./
RUN go mod download
COPY ../ .
RUN go build -o slipway ./cli

FROM alpine:3.19
COPY --from=build /src/slipway /usr/local/bin/slipway
ENTRYPOINT ["slipway"]
