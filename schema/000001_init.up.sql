create table users
(
    id  bigserial not null unique ,
    login varchar unique not null ,
    password varchar not null
);
