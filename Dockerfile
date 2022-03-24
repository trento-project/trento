FROM registry.suse.com/bci/nodejs:16 AS node-build
WORKDIR /build
# skip the hack/get_version_from_git.sh execution in the frontend build
ENV VERSION=""
# first we only add what's needed to run the makefile
ADD Makefile /build/
# then we add what's needed just to install node packages so that dependencies can be cached in a dedicate layer
ADD web/frontend/package.json web/frontend/package-lock.json /build/web/frontend/
RUN zypper -n in make && make web-deps
# only as last step we add the web frontend sources, this way changing these doesn't retrigger the slow npm install
ADD web/frontend /build/web/frontend/
RUN make web-assets

FROM registry.suse.com/bci/golang:1.16 as go-build
WORKDIR /build
# we add what's needed to download go modules so that dependencies can be cached in a dedicate layer
ADD go.mod go.sum /build/
RUN go mod download
ADD . /build
COPY --from=node-build /build /build
RUN zypper -n in git-core && make build

FROM registry.suse.com/bci/python:3.9 AS trento-runner
RUN /usr/bin/python3 -m venv /venv \
    && /venv/bin/pip install 'ansible~=4.6.0' 'requests~=2.26.0' 'rpm==0.0.2' 'pyparsing~=2.0' \
    && zypper -n ref && zypper -n in --no-recommends openssh \
    && zypper -n clean

ENV PATH="/venv/bin:$PATH"
ENV PYTHONPATH=/venv/lib/python3.9/site-packages

# Add Tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

COPY --from=go-build /build/trento /usr/bin/trento
LABEL org.opencontainers.image.source="https://github.com/trento-project/trento"
ENTRYPOINT ["/tini", "--", "/usr/bin/trento"]

FROM registry.suse.com/bci/bci-micro:15.3 AS trento-web
COPY --from=go-build /build/trento /usr/bin/trento
LABEL org.opencontainers.image.source="https://github.com/trento-project/trento"
EXPOSE 8080/tcp
ENTRYPOINT ["/usr/bin/trento"]
