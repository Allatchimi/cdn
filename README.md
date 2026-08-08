# Overview

This project leverages Docker and Make to streamline development and deployment.

# Prerequisites

- Docker installed and running

- Make installed

- A GitHub account with a personal access token with the necessary permissions

# Usage

1. Clone, build and run docker image

   Clone the repository:

   ```
       git clone https://github.com/4kpros/cdn.git
       cd cdn
   ```

   Add env file:

   ```
       cp .env.example app.env
   ```

   Uncomment this line on this file docker/cdn/Dockerfile:

   ```
       # COPY --from=builder /app/app.env ./app.env
   ```

   To build and run locally using Docker, run the following command:

   ```
       make docker-cdn
   ```

   This will build and start the Docker container.

# Update GitHub Action Secrets for continuous integration(build and package)

Go to this link: [GitHub Action Secrets](https://github.com/4kpros/cdn/settings/secrets/actions)

- ------------- On your GitHub Action Secrets page -------------

  - Set Secrets `GHCR_USERNAME` `GHCR_PASSWORD` with value your GitHub credentials. `GHCR_PASSWORD` is your personal access token with `write package` permission enabled

# Makefile Targets

- `docker-cdn`: Builds and starts the Docker container for local development.

- `docker-ghcr-push`: Builds the Docker image and pushes it to the GitHub Container Registry.

- `docker-ghcr-pull`: Pulls a specific image from the GitHub Container Registry.

# Additional Notes

By following these steps and customizing the Makefile to fit your specific needs, you can effectively manage your project using Docker and Make.
