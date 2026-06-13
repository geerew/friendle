FROM node:22-alpine AS ui

WORKDIR /src/ui

COPY ui/package.json ui/pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile

COPY ui/ ./
RUN pnpm run build

FROM golang:1.25-alpine AS backend

ARG VERSION
ARG COMMIT

RUN apk add --no-cache gcc git musl-dev sqlite-dev

ENV CGO_ENABLED=1

WORKDIR /src/

COPY . .
RUN rm -rf ui

COPY --from=ui /src/ui/build ./ui/build
COPY --from=ui /src/ui/embed.go ./ui/embed.go

RUN go mod download
RUN set -eux; \
 SET_VERSION="${VERSION:-dev}"; \
 SET_COMMIT="${COMMIT:-unknown}"; \
 go build -ldflags "-X github.com/geerew/friendle/version.Version=${SET_VERSION} -X github.com/geerew/friendle/version.Commit=${SET_COMMIT}" -o friendle .

FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite-libs

COPY --from=backend /src/friendle /usr/local/bin/friendle

RUN chmod +x /usr/local/bin/friendle

EXPOSE 80

VOLUME [ "/friendle_data" ]

ENTRYPOINT [ "/usr/local/bin/friendle", "serve", "--data-dir", "/friendle_data", "--http", "0.0.0.0:80", "--enable-signup" ]
