FROM node:16 AS node-build
WORKDIR /build
ADD . /build
RUN make web-assets

FROM golang:1.16 as go-build
COPY --from=node-build /build /build
WORKDIR /build
RUN make build

FROM python:3.7-slim AS python-build
RUN ln -s /usr/local/bin/python /usr/bin/python
RUN /usr/bin/python -m venv /venv
RUN /venv/bin/pip install ansible ara

FROM gcr.io/distroless/python3:debug
COPY --from=python-build /venv /venv
ENV PATH="/venv/bin:$PATH"
ENV PYTHONPATH=/venv/lib/python3.7/site-packages
COPY --from=go-build /build/trento /app/trento

EXPOSE 8080/tcp
ENTRYPOINT ["/app/trento"]
