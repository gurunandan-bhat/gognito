create table user (
    id varchar(40) not null primary key,
    name varchar(64)
    email varchar(64) not null
    roles tnytext
    lastSeen timestamp
    createdAt timestamp not null default current_timestamp
);

