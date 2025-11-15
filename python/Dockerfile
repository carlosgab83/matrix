# python/Dockerfile
FROM python:3.12-slim

# Install pipenv
RUN pip install --no-cache-dir pipenv

# Set working directory
WORKDIR /app

# Development target (default)
FROM python:3.12-slim AS dev

RUN pip install --no-cache-dir pipenv

WORKDIR /app

# Install system dependencies for Kafka
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

# Copy Pipfile first for dependency caching
COPY oracle/Pipfile* ./

# Install dependencies
RUN pipenv install --dev --system --deploy || pipenv install --dev --system

CMD ["python", "main.py"]