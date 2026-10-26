# builder
FROM node:24-alpine AS builder
WORKDIR /app
COPY ./mcp/package.json ./mcp/package-lock.json* ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci --no-audit --no-fund
COPY ./mcp .
RUN npm run build

# runner
FROM node:24-alpine AS runner
WORKDIR /app
USER node
COPY ./mcp/package.json ./mcp/package-lock.json* ./
COPY --from=builder /app/dist ./dist
CMD ["npm", "run", "run"]
