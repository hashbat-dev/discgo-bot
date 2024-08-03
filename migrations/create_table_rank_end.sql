create table if not exists RankEnd
(
    ID           int auto_increment
        primary key,
    RankCount    int         null,
    RankCategory varchar(50) null,
    RankEnd      varchar(50) null
);

