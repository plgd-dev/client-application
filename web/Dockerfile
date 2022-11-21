FROM node:16-bullseye-slim AS build-web
COPY . /web
RUN cd /web && npm install && npm run build
