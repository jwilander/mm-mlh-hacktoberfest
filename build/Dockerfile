ARG DOCKER_BUILD_IMAGE=golang:1.14
ARG DOCKER_BASE_IMAGE=alpine:3.12

FROM ${DOCKER_BUILD_IMAGE} AS build
WORKDIR /mm-mlh-hacktoberfest/
COPY . /mm-mlh-hacktoberfest/
RUN make build

# Final Image
FROM ${DOCKER_BASE_IMAGE}

COPY --from=build /mm-mlh-hacktoberfest/build/_output/bin/main /mm-mlh-hacktoberfest/main
WORKDIR /mm-mlh-hacktoberfest/

CMD ./main
