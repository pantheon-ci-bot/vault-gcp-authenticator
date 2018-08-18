FROM scratch
ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY ./vault-gcp-authenticator /
CMD ["/vault-gcp-authenticator"]
