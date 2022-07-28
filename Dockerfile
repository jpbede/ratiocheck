FROM alpine:3.16.1
COPY ratiocheck /
EXPOSE 3000
CMD ["/ratiocheck"]