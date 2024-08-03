create table if not exists BotFeatures
(
    ID           int auto_increment
        primary key,
    Version      varchar(30)                        null,
    Feature      varchar(500)                       null,
    DateReleased datetime default CURRENT_TIMESTAMP null
);

