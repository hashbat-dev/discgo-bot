create table if not exists TTSModels
(
    ID              int auto_increment
        primary key,
    Command         varchar(50)                        null,
    Description     varchar(150)                       null,
    Model           varchar(250)                       null,
    UpdatedBy       varchar(50)                        null,
    UpdatedDateTime datetime default CURRENT_TIMESTAMP null
);

