# syntax=docker/dockerfile:1
ARG K6_VERSION=2.3.0

# Stage 1: typecheck and compile TypeScript tests to JavaScript.
FROM node:22-alpine AS build
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run typecheck && npm run build

# Stage 2: k6 runtime with only the compiled scripts.
FROM grafana/k6:${K6_VERSION}
COPY --from=build /src/dist /scripts
