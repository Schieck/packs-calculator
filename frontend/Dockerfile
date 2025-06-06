# Build stage
FROM node:24.1.0 AS builder

WORKDIR /app

# Accept build arguments
ARG AUTH_SECRET
ARG API_BASE_URL=https://packs-backend.fly.dev/api/v1
ENV VITE_AUTH_SECRET=${AUTH_SECRET}
ENV VITE_API_BASE_URL=${API_BASE_URL}

# Copy package files
COPY package.json pnpm-lock.yaml ./

# Install pnpm and dependencies
RUN npm install -g pnpm && pnpm install --frozen-lockfile

# Copy source code
COPY . .

# Build the application
RUN pnpm run build

# Production stage
FROM nginx:alpine

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"] 