ARG POSTGRES_VERSION=12

FROM postgres:${POSTGRES_VERSION}

ADD init-user-db.sh /docker-entrypoint-initdb.d/