#!/bin/sh
# wait-for-db.sh

HOST=$1
shift
CMD="$@"

echo "Waiting for database at $HOST..."

until pg_isready -h "$HOST" -U postgres; do
  echo "Waiting for database..."
  sleep 1
done

echo "Database is ready!"
exec $CMD
