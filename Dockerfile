FROM node:16 as web

ARG TARGETARCH
# ARG VERSION="$(git describe --tags --abbre)"
ARG VERSION="$(cat VERSION)"
ARG GOFLAGS="'-ldflags=-w -s' -X 'dashboard/version.Version=$VERSION' -extldflags '-static'"

WORKDIR /build
COPY ./web .
RUN npm install
ENV REACT_APP_VERSION=$VERSION
RUN npm run build

FROM golang AS builder

ENV GOARCH=$TARGETARCH
ENV GOFLAGS=$GOFLAGS
ENV GO111MODULE=on
ENV CGO_ENABLED=1
WORKDIR /build
COPY . .
COPY --from=web /build/build ./web/build
COPY --from=web /build/web.go ./web/web.go
RUN go mod download
RUN go build -o actgpt

FROM alpine

RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true

COPY --from=builder /build/actgpt /

EXPOSE 3000
WORKDIR /data
ENTRYPOINT ["/actgpt"]
