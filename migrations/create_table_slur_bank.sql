create table if not exists SlurBank
(
    ID              int auto_increment
        primary key,
    Slur            varchar(150)  null,
    SlurTarget      varchar(150)  null,
    SlurDescription varchar(1000) null
);

