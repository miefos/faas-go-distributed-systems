# This docker only runs the python setup script
FROM python:3.9-alpine

WORKDIR /app

# Install curl
RUN apk add --no-cache curl

# Copy the setup script and secret file into the container
COPY setup.py ./setup.py

# Install required Python packages
RUN pip install requests

# Run the setup script
CMD ["python", "setup.py"]