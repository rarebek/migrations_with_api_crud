CREATE TABLE posts(
    id serial primary key,
    content VARCHAR,
    user_id INT,
    foreign key (user_id) references users(id)
);