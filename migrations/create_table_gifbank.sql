create table if not exists GifBank
(
    ID            int auto_increment
        primary key,
    Category      varchar(40)                        null,
    GifURL        varchar(500)                       null,
    DateTimeAdded datetime default CURRENT_TIMESTAMP null,
    AddedByID     varchar(50)                        null
);