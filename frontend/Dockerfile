# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# Runtime stage
FROM node:18-alpine

RUN npm install -g http-server@0.12.3

WORKDIR /app

COPY --from=builder /app/dist ./dist

EXPOSE 5173

CMD ["http-server", "dist", "-p", "5173", "--cors"]
