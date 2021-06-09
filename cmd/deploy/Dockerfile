FROM golang AS build
WORKDIR /build
COPY . ./
RUN CGO_ENABLED=0 go build arch-repo/cmd/deploy

FROM alpine
COPY --from=build /build/deploy /bin/deploy
ENTRYPOINT ["/bin/deploy"]
