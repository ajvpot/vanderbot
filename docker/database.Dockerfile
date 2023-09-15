ARG POSTGRES_VERSION=16

FROM postgres:${POSTGRES_VERSION}

ADD init-user-db.sh /docker-entrypoint-initdb.d/