create table users
(
    id  bigserial not null unique primary key ,
    login varchar unique not null ,
    password varchar not null
);

CREATE TABlE orders
(
    id      bigserial not null unique primary key,
    number  varchar(20) not null ,
    status  varchar(50) not null default 'NEW',
    user_id int references users (id) not null ,
    upload_at timestamp(0) with time zone default now()
)