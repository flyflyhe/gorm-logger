create table if not exists `user`
(
    `id`       bigint(20) unsigned not null auto_increment primary key ,
    `username` varchar(255) default ''     not null,
    `password` varchar(255) default ''     not null,
    `status`   tinyint default 0 not null,
    `created_at`  datetime           not null,
    `updated_at`  datetime           not null
)