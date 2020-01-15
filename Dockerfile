# Build:                   sudo docker build --no-cache . -t cacophony-api
# Run interactive session: sudo docker run -it cacophony-api

FROM golang:latest

RUN apt-get update
RUN apt-get install -y apt-utils

# server for automated testing
EXPOSE 2040

# WORKDIR /server
# RUN ls
# RUN go get ./...
# RUN go build ./...
# RUN testing-server
COPY  docker-entrypoint.sh /

RUN mkdir /etc/cacophony
RUN echo "" > /etc/cacophony/config.toml

ENTRYPOINT ["/docker-entrypoint.sh"]
