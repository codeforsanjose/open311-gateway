# Create a minimal container to run a Golang static binary

FROM tianon/true
EXPOSE 8080

COPY image/ /

ENV PORT=8080
CMD ["/gateway", "-debug"]
