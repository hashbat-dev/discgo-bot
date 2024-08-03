create table if not exists ToDoList
(
    ID               int auto_increment
        primary key,
    Category         varchar(10) null,
    AssignedID       int         null,
    ToDoText         text        null,
    CreatedBy        varchar(30) null,
    CreatedDateTime  datetime    null,
    StartedBy        varchar(30) null,
    StartedDateTime  datetime    null,
    FinishedBy       varchar(30) null,
    FinishedDateTime datetime    null,
    Version          varchar(10) null,
    DeletedDateTime  datetime    null
);

