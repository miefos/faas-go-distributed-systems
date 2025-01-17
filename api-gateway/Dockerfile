FROM apache/apisix

WORKDIR /app
# Copy the setup script and secret file into the container
COPY setup.sh ./setup.sh
COPY secret.txt ./secret.txt

# Run the setup script at startup
CMD ["sh", "setup.sh"]