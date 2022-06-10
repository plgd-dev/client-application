FROM node:16 AS build-web
COPY . /web
RUN cd /web && npm install && npm run build
