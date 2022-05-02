ARG BASE_GO_IMAGE
ARG BASE_TARGET_IMAGE

FROM ${BASE_GO_IMAGE} AS build
ARG APP_NAME
WORKDIR /usr/local/go/src/${ALIAS}
COPY ./ ./
RUN go mod download
ARG TARGETOS
ARG TARGETARCH
ARG CGO_ENABLED
RUN env CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/${APP_NAME} ./

# FROM basego AS unit-test
# RUN --mount=target=. \
#     --mount=type=cache,target=/root/.cache/go-build \
#     go test -v .

FROM ${BASE_TARGET_IMAGE}
ARG APP_NAME
COPY --from=build /out/${APP_NAME} /bin/${APP_NAME}
COPY config /etc/${APP_NAME}/
RUN mkdir -p /${APP_NAME} && \
    chown -R nobody:nobody /etc/${APP_NAME} /${APP_NAME}

USER nobody
WORKDIR /${APP_NAME}
