CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS users_unique_name
    ON users
    ("name")
;
