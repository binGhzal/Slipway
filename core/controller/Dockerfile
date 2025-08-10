# SPDX-License-Identifier: Apache-2.0
FROM golang:1.22 as build
WORKDIR /src
COPY ../../ go.mod go.sum ./
RUN go mod download
COPY ../../ .
RUN go build -o controller ./core/controller

FROM gcr.io/distroless/base
COPY --from=build /src/controller /controller
ENTRYPOINT ["/controller"]
