FROM ubuntu:xenial

ADD https://github.com/mconftec/fsbeat/releases/download/v1.0.1-alpha/fsbeat-1.0.1-alpha-amd64.deb /

RUN dpkg -i /fsbeat-1.0.1-alpha-amd64.deb

CMD ["/usr/share/fsbeat/bin/fsbeat", "-e", "-c", "/usr/share/fsbeat/fsbeat.yml"]
