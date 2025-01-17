FROM apache/apisix

WORKDIR /app
# Copy the setup script and secret file into the container
COPY setup.sh ./setup.sh
COPY secret.txt ./secret.txt

USER root

RUN apt-get update
RUN apt-get install -y curl

# Run the setup script
CMD [ "./setup.sh" ]