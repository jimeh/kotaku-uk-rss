FROM scratch
ADD bin/kotaku-uk-rss_linux_amd64 /kotaku-uk-rss
EXPOSE 8080
WORKDIR /
CMD ["/kotaku-uk-rss", "--port", "8080"]
