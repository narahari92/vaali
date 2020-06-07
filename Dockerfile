FROM golang:1.14-alpine AS build

WORKDIR /go/src/
RUN mkdir -p github.com/narahari92/vaali
COPY ./ /go/src/github.com/narahari92/vaali

ARG GOARCH

RUN go install github.com/narahari92/vaali/cmd/vaali

FROM scratch

COPY --from=build /go/bin/vaali /bin/vaali

ENTRYPOINT ["/bin/vaali"]
