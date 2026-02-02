FROM node:20-alpine

WORKDIR /app

COPY package.json package-lock.json* ./

# Increase npm network robustness
RUN npm config set fetch-retries 5 \
    && npm config set fetch-retry-factor 2 \
    && npm config set fetch-retry-mintimeout 10000 \
    && npm config set fetch-retry-maxtimeout 60000

# Use 'npm ci' if lockfile exists, otherwise 'npm install'
CMD ["sh", "-c", "if [ -f package-lock.json ]; then npm ci; else npm install; fi && npm run dev"]
