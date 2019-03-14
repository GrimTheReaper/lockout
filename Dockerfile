FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD lockout /
ADD pb /
CMD ["/lockout"]
