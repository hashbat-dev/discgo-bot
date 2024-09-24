create table if not exists TrackWords
(
    ID          int auto_increment
        primary key,
    Phrase      varchar(50)                        null,
    UserID      varchar(50)                        null,
    Count       bigint                             null,
    LastTracked datetime default CURRENT_TIMESTAMP null
);

