create table if not exists WowCount
(
    ID         int auto_increment
        primary key,
    UserID     varchar(50)                        null,
    TotalCount int                                null,
    MaxWow     int                                null,
    LastUsed   datetime default CURRENT_TIMESTAMP null
);

