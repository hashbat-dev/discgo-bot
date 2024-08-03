create table if not exists FartCount
(
    ID          int auto_increment
        primary key,
    UserID      varchar(50)                        null,
    Count       int      default 1                 null,
    LastUpdated datetime default CURRENT_TIMESTAMP null
);

