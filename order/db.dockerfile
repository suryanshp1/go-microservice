FROM postgres:15-alpine

COPY up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]