ARG NODEJS_VERSION=14
ARG POSTGRES_VERSION=12

FROM node:${NODEJS_VERSION}-alpine

RUN apk add postgresql-client

# Latest version of graphile-migrate
RUN npm install -g graphile-migrate

# Default working directory. Map your migrations folder in here with `docker -v`
WORKDIR /repo

ENTRYPOINT ["/usr/local/bin/graphile-migrate"]