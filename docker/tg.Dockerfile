# builder
FROM node:24-alpine AS builder
WORKDIR /app
COPY ./tg/package.json ./tg/package-lock.json* ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci --no-audit --no-fund
COPY ./tg .
RUN npm run build

# runner
FROM node:24-alpine AS runner
WORKDIR /app
COPY --from=builder /app/package.json ./
COPY --from=builder /app/package-lock.json* ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci --only=production --no-audit --no-fund
COPY --from=builder /app/dist ./dist
CMD ["npm", "run", "run"]
