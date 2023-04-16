from alpine:latest
LABEL org.opencontainers.image.source="https://github.com/AnestLarry/Simserver"
copy bin/* /Simserver/
CMD /Simserver/ftpserver