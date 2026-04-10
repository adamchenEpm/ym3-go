
create database if not exists `ym3_main`;
use `ym3_main`;

create table if not exists `sys_channel` (
    `id` int auto_increment primary key,
    `name` varchar(255) default '' comment '渠道名称',
    `type` varchar(255) default '' comment '渠道类型',
    `app_id` varchar(255) default '' comment '应用ID',
    `app_secret` varchar(255) default '' comment '应用密钥',
    `link_way` varchar(255) default '' comment '链接方式(websocket,webhook,http)',

    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
)




create table if not exists `sys_user` (
    `id` int auto_increment primary key,
    `name` varchar(255),
    `age` int,
    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
)

