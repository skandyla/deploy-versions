FROM golang:1.17 as build
RUN update-ca-certificates

FROM scratch
COPY main /main
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
ENV PORT 8080
EXPOSE 8080
USER nobody
ENTRYPOINT ["/main"]
