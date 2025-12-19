create table if not exists users (
    id char(36) primary key,
    name varchar(255) not null,
    email varchar(100) not null unique,
    email_verified_at timestamp null,
    password varchar(100) not null,
    avatar varchar(255) null,
    phone varchar(20) null,
    is_active bool default true,
    last_login_at timestamp null default current_timestamp,
    role enum('Admin', 'User') default 'User',
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    deleted_at timestamp null
)engine = InnoDB;