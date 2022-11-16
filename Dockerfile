FROM chromedp/headless-shell:latest
RUN apt install dumb-init
COPY ratiocheck /
EXPOSE 3000
ENTRYPOINT ["dumb-init", "--", "/ratiocheck", "l"]