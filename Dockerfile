FROM debian:latest

# Install networking tools
RUN apt-get update
RUN apt-get install net-tools

# Create log and binary directories
RUN mkdir -p /var/log/wololo
RUN mkdir /testdir

# Create config related files/directories
RUN mkdir -p /etc/wololo
RUN touch /etc/wololo/wololo.config

# Copy the application files
ADD wololo /testdir/wololo
ADD setup_cont_env.sh /testdir/setup_cont_env.sh

# Expose the application on port 5000
EXPOSE 5000

# Add a non-root user and set it as default to run applications
RUN useradd -ms /bin/bash testusr

# Set the entry point of the container to the application executable
ENTRYPOINT /testdir/setup_cont_env.sh
