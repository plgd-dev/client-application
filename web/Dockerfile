FROM node:16 AS build-web
COPY . /web
RUN cd /web && npm install --ignore-scripts && npm run build
