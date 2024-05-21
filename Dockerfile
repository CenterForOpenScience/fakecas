FROM golang:1.22-alpine as build


COPY ./ /go/src/github.com/CenterForOpenScience/fakecas

ARG GIT_COMMIT=
ENV GIT_COMMIT ${GIT_COMMIT}

ARG GIT_TAG=
ENV GIT_TAG ${GIT_TAG}

RUN cd /go/src/github.com/CenterForOpenScience/fakecas \
    && go mod download \
    && VERSION=${GIT_TAG} go build -o fakecas

FROM alpine:3.19 as runtime

COPY --from=build /go/src/github.com/CenterForOpenScience/fakecas/fakecas /usr/local/bin/

CMD ["fakecas"]
