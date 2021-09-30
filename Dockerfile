FROM node:16 AS node-build
WORKDIR /build
ADD . /build
RUN make web-assets

FROM golang:1.16 as go-build
COPY --from=node-build /build /build
WORKDIR /build
RUN make build

FROM python:3.7-slim AS trento-runner
RUN ln -s /usr/local/bin/python /usr/bin/python \
    && /usr/bin/python -m venv /venv \
    && /venv/bin/pip install ansible ara \
    && apt-get update && apt-get install -y --no-install-recommends \
      ssh \
    && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false \
    && rm -rf /var/lib/apt/lists/*

ENV PATH="/venv/bin:$PATH"
ENV PYTHONPATH=/venv/lib/python3.7/site-packages
COPY --from=go-build /build/trento /app/trento
LABEL org.opencontainers.image.source="https://github.com/trento-project/trento"
ENTRYPOINT ["/app/trento"]

FROM gcr.io/distroless/base:debug AS trento-web
COPY --from=go-build /build/trento /app/trento
LABEL org.opencontainers.image.source="https://github.com/trento-project/trento"
EXPOSE 8080/tcp
ENTRYPOINT ["/app/trento"]
