# Python Narrative Detection API
FROM python:3.11-slim

WORKDIR /app

# Copy requirements first for layer caching
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application files
COPY . .

# Expose port
EXPOSE 5000

# Set environment variables
ENV DEMO_MODE=true
ENV PORT=5000

# Run the application
CMD ["python", "api.py"]
