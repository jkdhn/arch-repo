FROM golang AS build
WORKDIR /build
COPY . ./
RUN CGO_ENABLED=0 go build arch-repo/cmd/server

FROM alpine
COPY --from=build /build/server /bin/server
ENTRYPOINT ["/bin/server"]
