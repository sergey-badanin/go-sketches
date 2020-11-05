CREATE DATABASE composed;

\connect composed;

CREATE TABLE message(
    created_at TIMESTAMP,
    text varchar
);