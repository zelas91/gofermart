create table users
(
    id  bigserial not null unique primary key ,
    login varchar unique not null ,
    password varchar not null
);
