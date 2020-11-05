docker cp ./create.sql psql-main:/tmp/create.sql
docker exec -u postgres psql-main psql -U postgres -f /tmp/create.sql
