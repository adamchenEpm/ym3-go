
-- 创建数据库
create database if not exists ym3_main character set utf8mb4 collate utf8mb4_unicode_ci;

use ym3_main;

-- 菜单
create table if not exists sys_menu
(
    id          bigint  not null auto_increment comment '菜单Id',
    up_id       bigint  not null comment '上级菜单Id(一级菜单为0)',
    name        varchar(50)  default null comment '菜单名称',
    title       varchar(50)  default '' comment '菜单标题',
    perm        varchar(100) default '' comment '授权(user:list)',
    type        tinyint not null comment '菜单类型(0:主目录,1:目录,2:菜单,3:按钮)',
    style       tinyint      default null comment '菜单风格(0:default,1:tabView,2:drawer,3:dialog,4:url)',
    icon        varchar(50)  default null comment '菜单图标',
    sort        smallint     default null comment '排序',
    remark      varchar(1000) default null comment '备注',

    create_time datetime     default null comment '创建时间',
    update_time datetime     default null comment '修改时间',
    status      tinyint      default '1' comment '状态(0:停用,1:启用)',
    key (type),
    primary key (id)
) auto_increment = 101 engine=InnoDB comment='菜单';

-- 角色
create table if not exists sys_role
(
    id          bigint      not null auto_increment comment '角色Id',
    name        varchar(50) not null comment '角色名称',
    remark      varchar(200) default null comment '备注',

    create_time datetime     default null comment '创建时间',
    update_time datetime     default null comment '修改时间',
    status      tinyint      default '0' comment '状态(0:停用,1:启用)',
    primary key (id),
    key (status)
) engine=InnoDB comment='角色';

-- 用户
create table if not exists sys_user
(
    id                      bigint not null auto_increment          comment '用户Id',
    name                    varchar(100) default null               comment '用户名称',
    nickname                varchar(100) default null               comment '用户昵称',
    mobile                  varchar(20) default ''                  comment '用户手机',
    email                   varchar(100) default ''                 comment '用户邮箱',
    password                varchar(500) not null                   comment '用户密码',

    remark                  varchar(100) default null               comment '备注',

    create_time             datetime default null                   comment '创建时间',
    update_time             datetime default null                   comment '修改时间',
    status                  tinyint default '0'                     comment '状态(0:停用,1:启用)',
    primary key (id),
    key (mobile),
    key (status)
) engine=InnoDB comment='用户';


-- 通道
create table if not exists sys_channel
(
    id                          bigint not null auto_increment          comment 'ID',
    name                        varchar(255) not null                   comment '名称',
    type                        varchar(255) not null                   comment '类型(feishu,dingtalk,wechat,qq)',
    app_id                      varchar(255) default ''                 comment '应用ID',
    app_secret                  varchar(255) default ''                 comment '应用密钥',
    link_way                    varchar(255) default ''                 comment '链接方式(websocket,webhook,http)',
    sort                        int default 101                         comment '排序',
    remark                      varchar(255) default ''                 comment '备注',

    create_time                 datetime default current_timestamp      comment '创建时间',
    update_time                 datetime default current_timestamp on update current_timestamp comment '修改时间',
    status                      tinyint default 1                       comment '状态(1:无效,2:有效)',

    primary key (id)
) comment='通道';



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

CREATE TABLE IF NOT EXISTS sys_user_binding (
                                                id INT AUTO_INCREMENT PRIMARY KEY COMMENT '绑定ID',
                                                internal_user_id INT NOT NULL COMMENT '内部用户ID',
                                                channel_instance_id INT NOT NULL COMMENT '通道实例ID',
                                                external_user_id VARCHAR(100) NOT NULL COMMENT '平台侧用户ID（飞书open_id、微信openid等）',
    external_chat_id VARCHAR(100) COMMENT '会话ID（私聊时等于external_user_id，群聊时为群ID）',
    external_name VARCHAR(100) COMMENT '昵称（冗余）',
    bind_data JSON COMMENT '扩展信息（头像、手机号等）',
    bind_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '绑定时间',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    status TINYINT DEFAULT 1 COMMENT '状态：1正常，0禁用',
    UNIQUE KEY uk_binding (channel_instance_id, external_user_id),
    FOREIGN KEY (internal_user_id) REFERENCES sys_user(id) ON DELETE CASCADE,
    FOREIGN KEY (channel_instance_id) REFERENCES sys_channel_instance(id) ON DELETE CASCADE
    ) COMMENT='内部用户与外部通道账户绑定关系表';



-- 角色表
CREATE TABLE IF NOT EXISTS sys_role (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '角色ID',
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '角色名称：admin, user, viewer',
    description TEXT COMMENT '角色描述',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    status TINYINT DEFAULT 1 COMMENT '状态：1正常，0禁用'
) COMMENT='用户角色表';

-- 2. 用户表
CREATE TABLE IF NOT EXISTS sys_user (
                                        id INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
                                        username VARCHAR(100) NOT NULL UNIQUE COMMENT '用户名',
    password_hash VARCHAR(255) COMMENT '密码哈希（支持SSO时可空）',
    email VARCHAR(255) COMMENT '邮箱',
    phone VARCHAR(50) COMMENT '手机号',
    avatar_url TEXT COMMENT '头像URL',
    role_id INT COMMENT '角色ID',
    last_login_at TIMESTAMP NULL COMMENT '最后登录时间',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    status TINYINT DEFAULT 1 COMMENT '状态：1正常，0禁用',
    FOREIGN KEY (role_id) REFERENCES sys_role(id) ON DELETE SET NULL
    ) COMMENT='系统内部用户表';

-- 3. 通道实例表（每个机器人/应用）
CREATE TABLE IF NOT EXISTS sys_channel_instance (
                                                    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '通道实例ID',
                                                    channel_type VARCHAR(20) NOT NULL COMMENT '通道类型：feishu, wechat_mp, wechat_work, dingtalk, qq, web',
    name VARCHAR(100) NOT NULL COMMENT '实例名称，如“客服机器人-生产”',
    config JSON NOT NULL COMMENT '平台配置JSON，如 {"app_id":"xxx","app_secret":"xxx"}，敏感字段建议加密',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    status TINYINT DEFAULT 1 COMMENT '状态：1正常，0禁用'
    ) COMMENT='通道实例表';

-- 4. 用户绑定表（内部用户 ↔ 外部通道身份）
CREATE TABLE IF NOT EXISTS sys_user_binding (
                                                id INT AUTO_INCREMENT PRIMARY KEY COMMENT '绑定ID',
                                                internal_user_id INT NOT NULL COMMENT '内部用户ID',
                                                channel_instance_id INT NOT NULL COMMENT '通道实例ID',
                                                external_user_id VARCHAR(100) NOT NULL COMMENT '平台侧用户ID（飞书open_id、微信openid等）',
    external_chat_id VARCHAR(100) COMMENT '会话ID（私聊时等于external_user_id，群聊时为群ID）',
    external_name VARCHAR(100) COMMENT '昵称（冗余）',
    bind_data JSON COMMENT '扩展信息（头像、手机号等）',
    bind_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '绑定时间',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    status TINYINT DEFAULT 1 COMMENT '状态：1正常，0禁用',
    UNIQUE KEY uk_binding (channel_instance_id, external_user_id),
    FOREIGN KEY (internal_user_id) REFERENCES sys_user(id) ON DELETE CASCADE,
    FOREIGN KEY (channel_instance_id) REFERENCES sys_channel_instance(id) ON DELETE CASCADE
    ) COMMENT='内部用户与外部通道账户绑定关系表';

-- 5. 群组配置表（群级别设置）
CREATE TABLE IF NOT EXISTS sys_group_config (
                                                id INT AUTO_INCREMENT PRIMARY KEY COMMENT '配置ID',
                                                channel_instance_id INT NOT NULL COMMENT '通道实例ID',
                                                external_chat_id VARCHAR(100) NOT NULL COMMENT '群ID（飞书chat_id、微信群ID等）',
    config JSON NOT NULL COMMENT '群配置JSON，如 {"only_reply_when_at": true, "default_model": "gpt-4"}',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    status TINYINT DEFAULT 1 COMMENT '状态：1正常，0禁用',
    UNIQUE KEY uk_group (channel_instance_id, external_chat_id),
    FOREIGN KEY (channel_instance_id) REFERENCES sys_channel_instance(id) ON DELETE CASCADE
    ) COMMENT='群组配置表';

-- 6. 模型规则表（按维度匹配AI模型）
CREATE TABLE IF NOT EXISTS sys_model_rule (
                                              id INT AUTO_INCREMENT PRIMARY KEY COMMENT '规则ID',
                                              channel_instance_id INT COMMENT '通道实例ID（NULL表示全局规则）',
                                              match_type VARCHAR(20) NOT NULL COMMENT '匹配类型：command, group, user_role, user_id, global',
    match_value TEXT COMMENT '匹配值：如 "/help"、"admin"、群ID、用户ID等',
    priority INT DEFAULT 0 COMMENT '优先级（数值越小越高）',
    model_provider VARCHAR(50) NOT NULL COMMENT '模型提供商：openai, azure, local, anthropic',
    model_name VARCHAR(100) NOT NULL COMMENT '模型名称',
    model_params JSON COMMENT '模型参数，如 {"temperature":0.7, "max_tokens":2000}',
    enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    status TINYINT DEFAULT 1 COMMENT '状态：1正常，0禁用',
    FOREIGN KEY (channel_instance_id) REFERENCES sys_channel_instance(id) ON DELETE CASCADE
    ) COMMENT='AI模型调用规则表';

-- 7. 菜单表（前端动态菜单）
CREATE TABLE IF NOT EXISTS sys_menu (
                                        id INT AUTO_INCREMENT PRIMARY KEY COMMENT '菜单ID',
                                        parent_id INT DEFAULT 0 COMMENT '父菜单ID',
                                        name VARCHAR(100) NOT NULL COMMENT '菜单名称',
    path VARCHAR(200) COMMENT '路由路径',
    icon VARCHAR(50) COMMENT '图标',
    permission_key VARCHAR(100) UNIQUE COMMENT '权限标识，如 "menu:admin:users"',
    sort_order INT DEFAULT 0 COMMENT '排序号',
    role_ids JSON DEFAULT NULL COMMENT '允许访问的角色ID数组，如 [1,2]，NULL或空数组表示公开',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    status TINYINT DEFAULT 1 COMMENT '状态：1正常，0禁用'
    ) COMMENT='前端菜单配置表';

-- 8. 审计日志表（记录所有消息和操作）
CREATE TABLE IF NOT EXISTS sys_audit_log (
                                             id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '日志ID',
                                             user_id INT COMMENT '内部用户ID',
                                             channel_instance_id INT COMMENT '通道实例ID',
                                             external_user_id VARCHAR(100) COMMENT '外部用户ID',
    action VARCHAR(50) NOT NULL COMMENT '操作类型：send_message, receive_message, model_call, login, config_change',
    request_data JSON COMMENT '请求数据',
    response_data JSON COMMENT '响应数据',
    error_message TEXT COMMENT '错误信息',
    duration_ms INT COMMENT '耗时（毫秒）',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    status TINYINT DEFAULT 1 COMMENT '状态：1正常，0禁用',
    FOREIGN KEY (user_id) REFERENCES sys_user(id) ON DELETE SET NULL,
    FOREIGN KEY (channel_instance_id) REFERENCES sys_channel_instance(id) ON DELETE SET NULL
    ) COMMENT='审计日志表';

-- 9. 系统配置表（全局配置）
CREATE TABLE IF NOT EXISTS sys_config (
                                          config_key VARCHAR(100) PRIMARY KEY COMMENT '配置键',
    value TEXT COMMENT '配置值',
    description TEXT COMMENT '配置描述',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    status TINYINT DEFAULT 1 COMMENT '状态：1正常，0禁用'
    ) COMMENT='系统配置表';

-- ============================================
-- 初始化数据
-- ============================================

INSERT INTO sys_role (name, description) VALUES
                                             ('admin', '系统管理员，所有权限'),
                                             ('user', '普通用户，可使用AI功能'),
                                             ('viewer', '只读用户，仅查看日志')
    ON DUPLICATE KEY UPDATE name = VALUES(name);

INSERT INTO sys_menu (parent_id, name, path, icon, permission_key, sort_order, role_ids) VALUES
                                                                                             (0, '控制台', '/dashboard', 'Dashboard', 'menu:dashboard', 1, NULL),
                                                                                             (0, '对话记录', '/conversations', 'Message', 'menu:conversations', 2, NULL),
                                                                                             (0, '模型管理', '/models', 'Robot', 'menu:models', 3, '[1]'),
                                                                                             (0, '通道管理', '/channels', 'Api', 'menu:channels', 4, '[1]'),
                                                                                             (0, '用户管理', '/users', 'User', 'menu:users', 5, '[1]'),
                                                                                             (0, '系统设置', '/settings', 'Setting', 'menu:settings', 6, '[1]')
    ON DUPLICATE KEY UPDATE name = VALUES(name);

INSERT INTO sys_config (config_key, value, description) VALUES
                                                            ('auto_create_user', 'true', '是否自动为外部用户创建内部账户'),
                                                            ('default_model_provider', 'openai', '默认模型提供商'),
                                                            ('default_model_name', 'gpt-3.5-turbo', '默认模型名称'),
                                                            ('ws_ping_interval', '30', 'WebSocket心跳间隔（秒）')
    ON DUPLICATE KEY UPDATE value = VALUES(value);












create table if not exists `sys_user` (
    `id` int auto_increment primary key,
    `name` varchar(255),
    `age` int,
    `create_time` datetime default current_timestamp
    `update_time` datetime default current_timestamp on update current_timestamp
    key(age_idx(age))
)

