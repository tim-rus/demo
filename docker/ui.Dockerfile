# builder
FROM node:24-alpine AS builder
WORKDIR /app
ARG VITE_SERVER_BASE_URL
COPY ./ui/package.json ./ui/package-lock.json* ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci --no-audit --no-fund
COPY ./ui .
RUN VITE_SERVER_BASE_URL=$VITE_SERVER_BASE_URL npm run build-only

# runner
FROM sebp/lighttpd:latest
COPY --from=builder /app/dist /app/www
COPY ./ui/lighttpd.conf /etc/lighttpd/lighttpd.conf
EXPOSE $DOCKER_PORT

CMD ["start.sh"]
