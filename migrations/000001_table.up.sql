CREATE table if not exists users(
    id serial primary key,
    uuid VARCHAR,
    name VARCHAR,
    age int
);