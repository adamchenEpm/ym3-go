
create database if not exists `ym3_main`;
use `ym3_main`;


create table if not exists `sys_user` (
                                          `id` int auto_increment primary key,
                                          `name` varchar(255),
    `age` int,
    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    `status` int default 1 comment '状态(1:正常,0:禁用)',
    key(age_idx(age))
    )

create table if not exists `sys_role` (
                                          `id` int auto_increment primary key,
                                          `name` varchar(255),
    `role_type` varchar(255) default '' comment '角色类型(user,admin)',
    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
    )

create table if not exists `sys_user_role` (
                                          `id` int auto_increment primary key,
                                          `name` varchar(255),
    `age` int,
    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
    )

create table if not exists `sys_department` (
                                          `id` int auto_increment primary key,
                                          `name` varchar(255),
    `role_type` varchar(255) default '' comment '角色类型(user,admin)',
    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
    )


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






create table if not exists `sys_channel_user` (
                                          `id` int auto_increment primary key,
                                          `name` varchar(255),

    `user_id` int,
    `channel_id` int,
    `channel_name` varchar(255) default '' comment '渠道名称',
    `channel_user_id` varchar(255) default '' comment '渠道用户ID',
    `channel_chat_id` varchar(255) default '' comment '渠道聊天ID',

    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
)

create table if not exists `sys_agent` (
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

create table if not exists `sys__user_agent` (
                                                  `id` int auto_increment primary key,

    `user_id` int,
    `channel_id` int,
    `channel_name` varchar(255) default '' comment '渠道名称',
    `channel_user_id` varchar(255) default '' comment '渠道用户ID',
    `channel_chat_id` varchar(255) default '' comment '渠道聊天ID',

    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
    )


create table if not exists `sys_app` (
                                           `id` int auto_increment primary key,
                                           `name` varchar(255) default '' comment '应用名称',
    `type` varchar(255) default '' comment '应用类型',
    `app_id` varchar(255) default '' comment '应用ID',
    `app_secret` varchar(255) default '' comment '应用密钥',
    `link_way` varchar(255) default '' comment '链接方式(websocket,webhook,http)',

    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
    ) comment '应用表'



create table if not exists `sys_user_app` (
                                         `id` int auto_increment primary key,
                                         `name` varchar(255) default '' comment '渠道名称',
    `type` varchar(255) default '' comment '渠道类型',
    `app_id` varchar(255) default '' comment '应用ID',
    `app_secret` varchar(255) default '' comment '应用密钥',
    `link_way` varchar(255) default '' comment '链接方式(websocket,webhook,http)',
    `base_url` varchar(255) default '' comment '基础URL',
    `outer_app_id` varchar(255) default '' comment '应用聊天ID',
    `outer_app_secret` varchar(255) default '' comment '应用聊天密钥',

    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
    ) comment '应用表'

create table if not exists `sys_solu` (
                                         `id` int auto_increment primary key,
                                         `name` varchar(255) default '' comment '渠道名称',
    `type` varchar(255) default '' comment '渠道类型',
    `app_id` varchar(255) default '' comment '应用ID',
    `app_secret` varchar(255) default '' comment '应用密钥',
    `link_way` varchar(255) default '' comment '链接方式(websocket,webhook,http)',

    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
    ) comment '解决方案表'


create table if not exists `sys_` (
                                         `id` int auto_increment primary key,
                                         `name` varchar(255) default '' comment '渠道名称',
    `type` varchar(255) default '' comment '渠道类型',
    `app_id` varchar(255) default '' comment '应用ID',
    `app_secret` varchar(255) default '' comment '应用密钥',
    `link_way` varchar(255) default '' comment '链接方式(websocket,webhook,http)',

    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
    ) comment '应用表'


create table if not exists `sys_blacklist` (
                                      `id` int auto_increment primary key,
                                      `name` varchar(255) default '' comment '渠道名称',
    `type` varchar(255) default '' comment '渠道类型',
    `app_id` varchar(255) default '' comment '应用ID',
    `app_secret` varchar(255) default '' comment '应用密钥',
    `link_way` varchar(255) default '' comment '链接方式(websocket,webhook,http)',

    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
    ) comment '应用表'


create table if not exists `sys_role` (
                                      `id` int auto_increment primary key,
                                      `name` varchar(255) default '' comment '渠道名称',
    `type` varchar(255) default '' comment '渠道类型',
    `app_id` varchar(255) default '' comment '应用ID',
    `app_secret` varchar(255) default '' comment '应用密钥',
    `link_way` varchar(255) default '' comment '链接方式(websocket,webhook,http)',

    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
    ) comment '应用表'

