create table if not exists CommandLog
(
    ID       int auto_increment
        primary key,
    Command  varchar(50)                        null,
    Count    int      default 1                 null,
    UserID   varchar(50)                        null,
    LastUsed datetime default CURRENT_TIMESTAMP null
);