# Create a minimal container to run a Golang static binary

FROM alpine:3.3
EXPOSE 80

COPY image/ /
RUN apk add --update ca-certificates

ENV PORT=80
CMD ["/gateway", "-debug"]
