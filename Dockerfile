FROM golang:1.22-alpine


COPY ./ /go/src/github.com/CenterForOpenScience/fakecas

ARG GIT_COMMIT=
ENV GIT_COMMIT ${GIT_COMMIT}

ARG GIT_TAG=
ENV GIT_TAG ${GIT_TAG}

RUN cd /go/src/github.com/CenterForOpenScience/fakecas \
    && go mod download \
    && VERSION=${GIT_TAG} go build -o fakecas \
    && mv /go/src/github.com/CenterForOpenScience/fakecas/fakecas /usr/local/bin/

CMD ["fakecas"]
