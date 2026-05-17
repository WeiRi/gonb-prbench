FROM inp-dns-459
# The inplace test needs RunLocalUDPServer which was deleted from the image.
# Copy a bundled helper file that provides it.
COPY dns_helper.go /go/src/github.com/miekg/dns/dns_helper.go
