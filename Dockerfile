FROM node:lts-alpine AS frontend
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build:prod

FROM golang:alpine AS backend
WORKDIR /app
COPY backend/ .
RUN go build -o monhttp

FROM alpine
WORKDIR /monhttp
RUN mkdir -p /monhttp/config
COPY backend/migrations ./migrations/
COPY --from=backend /app/monhttp ./monhttp
COPY --from=frontend /app/dist/monhttp/ ./public/
ENTRYPOINT ["./monhttp"]
