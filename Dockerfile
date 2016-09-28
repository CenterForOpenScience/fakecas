FROM golang:1

ENV GLIDE_VERSION 0.12.2
ENV GLIDE_SHA256 edd398b4e94116b289b9494d1c13ec2ea37386bad4ada91ecc9825f96b12143c
RUN apt-get update \
    && apt-get install -y \
        curl \
    && curl -o /tmp/glide.tar.gz -SL "https://github.com/Masterminds/glide/releases/download/v$GLIDE_VERSION/glide-v$GLIDE_VERSION-linux-$(dpkg --print-architecture).tar.gz" \
    && echo "$GLIDE_SHA256  /tmp/glide.tar.gz" | sha256sum -c - \
    && tar -xzf /tmp/glide.tar.gz -C /usr/local/bin --strip-components=1 \
    && chmod +x /usr/local/bin/glide \
    && rm /tmp/glide.tar.gz \
    && apt-get clean \
    && apt-get autoremove -y \
        curl \
    && rm -rf /var/lib/apt/lists/*

RUN go get -d github.com/CenterForOpenScience/fakecas \
    && cd /go/src/github.com/CenterForOpenScience/fakecas \
    && glide install \
    && go build \
    && mv /go/src/github.com/CenterForOpenScience/fakecas /usr/local/bin/

CMD ["fakecas"]
