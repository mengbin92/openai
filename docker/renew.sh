#!/bin/bash

docker run -it --rm \
-v ./certbot/etc/letsencrypt:/etc/letsencrypt \
-v ./certbot/var/lib/letsencrypt:/var/lib/letsencrypt \
-v ./certbot/var/log/letsencrypt:/var/log/letsencrypt \
-v ./site:/data/letsencrypt \
certbot/certbot \
renew --webroot -w /data/letsencrypt --quiet && docker kill --signal=HUP frontend
