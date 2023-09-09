FROM nginx:latest

# COPY ./docker/nginx.prod/config /etc/nginx/conf.d
COPY ./config/nginx.conf /etc/nginx

RUN service nginx start