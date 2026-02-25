FROM golang:alpine AS base
RUN mkdir -p /opt/app
WORKDIR /opt/app
RUN apk add build-base

COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY . .

FROM base AS build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags musl -ldflags="-w -s" \
    -o user-service cmd/api/main.go

FROM base AS test-exec

ARG _outputdir="/tmp/coverage"
ARG COVERAGE_EXCLUDE

RUN mkdir -p ${_outputdir} && \
    go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1 && \
	grep -v -E "${COVERAGE_EXCLUDE}" coverage.tmp > ${_outputdir}/coverage.out && \
    go tool cover -html=${_outputdir}/coverage.out -o ${_outputdir}/coverage.html

FROM scratch AS test
ARG _outputdir="/tmp/coverage"
COPY --from=test-exec ${_outputdir}/coverage.out /
COPY --from=test-exec ${_outputdir}/coverage.html /

FROM alpine AS final

ARG app_name=app
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

COPY --from=build /opt/app/user-service /app/user-service
COPY --from=build /opt/app/docs /app/docs
COPY --from=build /opt/app/migrations /app/migrations

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

CMD ["/app/user-service"]
 