FROM chromedp/headless-shell:latest
RUN apt update && apt install dumb-init
RUN rm -rf /var/lib/{apt,dpkg,cache,log}/
COPY ratiocheck /
EXPOSE 3000
ENTRYPOINT ["dumb-init", "--", "/ratiocheck", "l"]