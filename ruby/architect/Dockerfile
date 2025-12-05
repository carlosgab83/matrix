# ruby/Dockerfile
FROM ruby:3.3-slim AS base

# System deps
RUN apt-get update -y && \
apt-get install -y --no-install-recommends build-essential git curl && \
rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

FROM base

# Copy Gemfile and Gemfile.lock for better caching
COPY ./Gemfile ./Gemfile.lock ./
RUN bundle install

# Copy source code
# COPY . .