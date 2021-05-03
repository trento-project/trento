FROM opensuse/leap:15.2
MAINTAINER Rubén Torrero Marijnissen <rtorreromarijnissen@suse.com>

COPY . /app

RUN zypper -n in go1.16 nodejs14 make
RUN cd /app && make build

EXPOSE 8080/tcp

ENTRYPOINT ["/app/trento"]
