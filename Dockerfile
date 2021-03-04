# Build:                   sudo docker build --no-cache . -t cacophony-api
# Run interactive session: sudo docker run -it cacophony-api

FROM golang:latest
RUN go version
RUN apt-get update
RUN apt-get install -y apt-utils
RUN apt-get install -y supervisor
RUN apt-get install -y dbus
RUN mkdir -p /var/run/dbus

RUN mkdir -p /var/log/supervisor
# server for automated testing
EXPOSE 2040
EXPOSE 80

COPY  thermal-recorder.conf /etc/supervisor/conf.d/thermal-recorder.conf
COPY  thermal-uploader.conf /etc/supervisor/conf.d/thermal-uploader.conf
COPY  event-reporter.conf /etc/supervisor/conf.d/event-reporter.conf
COPY  fake-lepton.conf /etc/supervisor/conf.d/fake-lepton.conf
COPY  management-interface.conf /etc/supervisor/conf.d/management-interface.conf


# WORKDIR /server
# RUN ls
# RUN go get ./...
# RUN go build ./...
# RUN testing-server
COPY  docker-entrypoint.sh /
RUN mkdir /etc/cacophony
RUN mkdir /var/spool/cptv

COPY recorder-config.toml /etc/cacophony/config.toml

ENTRYPOINT ["/docker-entrypoint.sh"]
