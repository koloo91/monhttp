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
WORKDIR /app
COPY backend/migrations ./migrations/
COPY --from=backend /app/ .
COPY --from=frontend /app/dist/monhttp/ ./public/
ENTRYPOINT ["./monhttp"]
