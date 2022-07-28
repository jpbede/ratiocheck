FROM chromedp/headless-shell:latest
COPY ratiocheck /
EXPOSE 3000
ENTRYPOINT ["/ratiocheck", "l"]