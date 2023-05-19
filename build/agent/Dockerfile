# syntax=docker/dockerfile:1
FROM golang:1.18
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0 
ENV GOOS=linux

RUN cd cmd/agent;go build -o /runner -mod mod -a .

FROM python:3.7-alpine
RUN pip install pytest
RUN apk --no-cache add ca-certificates git
WORKDIR /root/
COPY --from=0 /runner /bin/runner
ENV DEPENDENCY_MANAGER=pip
ENV CI=1
ENV PWTEST_SKIP_TEST_OUTPUT=1
ENTRYPOINT ["/bin/runner"]
