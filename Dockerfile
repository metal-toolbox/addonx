FROM gcr.io/distroless/static:nonroot

# `nonroot` coming from distroless
USER 65532:65532

# pass in name as --build-arg
ARG NAME

COPY ./bin/${NAME} /app

# Run the web service on container startup.
ENTRYPOINT ["/app"]
CMD ["serve"]
