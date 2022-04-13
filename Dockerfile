FROM golang:latest AS build

# NOTE! You need to find etcd or download it. See README.md
COPY ./etcd/bin/etcd /bin/etcd



FROM gcr.io/distroless/static
COPY --from=build /bin/etcd /bin/etcd
ENTRYPOINT ["/bin/etcd"]
# Args to project
CMD []

# docker build --no-cache -t us-central1-docker.pkg.dev/mchirico/public/etcd:v0.0.1 -f Dockerfile .
# docker push us-central1-docker.pkg.dev/mchirico/public/etcd:v0.0.1
