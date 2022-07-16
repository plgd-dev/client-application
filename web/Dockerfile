FROM node:18 AS build-web
COPY . /web
RUN cd /web && npm install && npm run build
