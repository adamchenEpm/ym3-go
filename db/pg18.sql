
create user ym3_user with password 'ym3_pwd';

create database ym3_main
ENCODING = 'UTF8'
LC_COLLATE = 'zh_CN.utf8'
LC_CTYPE = 'zh_CN.utf8'
TEMPLATE = template0
OWNER = ym3_user;

\c ym3_main

create extension vector;

alter database ym3_main SET timezone TO 'Asia/Shanghai';

-- 部门表
CREATE TABLE IF NOT EXISTS sys_dept (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    pid          BIGSERIAL DEFAULT NULL,
    remark      VARCHAR(100) DEFAULT NULL,
    create_time TIMESTAMP DEFAULT NULL,
    update_time TIMESTAMP DEFAULT NULL,
    status      SMALLINT DEFAULT 2
);

COMMENT ON TABLE sys_dept IS '部门';
COMMENT ON COLUMN sys_dept.id IS '部门Id';
COMMENT ON COLUMN sys_dept.name IS '部门名称';
COMMENT ON COLUMN sys_dept.remark IS '备注';
COMMENT ON COLUMN sys_dept.create_time IS '创建时间';
COMMENT ON COLUMN sys_dept.update_time IS '修改时间';
COMMENT ON COLUMN sys_role.status IS '状态(1:停用,2:启用)';

-- 部门表.索引
CREATE INDEX idx_sys_dept_name ON sys_dept(name);
CREATE INDEX idx_sys_dept_status ON sys_dept(status);

-- 部门表.基础数据
insert into sys_dept (id, name, pid, remark, status)
values (1, 'xx公司', null, '', 2),
(101, '营销中心', 1, '', 2),
(102, '设计中心', 1, '', 2),
(104, '采购中心', 1, '', 2),
(105, '财务中心', 1, '', 2),
(106, '行政中心', 1, '', 2),
(107, '信息中心', 1, '', 2),
(109, '产品中心', 1, '', 2),
(110, '客户中心', 1,, '', 2);

-- 角色表
CREATE TABLE IF NOT EXISTS sys_role (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    remark      VARCHAR(100) DEFAULT NULL,
    create_time TIMESTAMP DEFAULT NULL,
    update_time TIMESTAMP DEFAULT NULL,
    status      SMALLINT DEFAULT 2
);

COMMENT ON TABLE sys_role IS '角色';
COMMENT ON COLUMN sys_role.id IS '角色Id';
COMMENT ON COLUMN sys_role.name IS '角色名称';
COMMENT ON COLUMN sys_role.remark IS '备注';
COMMENT ON COLUMN sys_role.create_time IS '创建时间';
COMMENT ON COLUMN sys_role.update_time IS '修改时间';
COMMENT ON COLUMN sys_role.status IS '状态(1:停用,2:启用)';

-- 角色表.索引
CREATE INDEX idx_sys_role_name ON sys_role(name);
CREATE INDEX idx_sys_role_status ON sys_role(status);

-- 角色表.基础数据
insert into sys_role (id, name, remark, status)
values (1, '管理员', '', 2);

-- 用户表
CREATE TABLE IF NOT EXISTS sys_user (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    nickname    VARCHAR(100) DEFAULT NULL,
    dept_id     BIGSERIAL DEFAULT NULL,
    mobile      VARCHAR(20) DEFAULT '',
    email       VARCHAR(100) DEFAULT '',
    password    VARCHAR(500) NOT NULL,
    remark      VARCHAR(100) DEFAULT NULL,
    create_time TIMESTAMP DEFAULT NULL,
    update_time TIMESTAMP DEFAULT NULL,
    status      SMALLINT DEFAULT 2
);

COMMENT ON TABLE sys_user IS '用户';
COMMENT ON COLUMN sys_user.id IS '用户Id';
COMMENT ON COLUMN sys_user.name IS '用户名称';
COMMENT ON COLUMN sys_user.nickname IS '用户昵称';
COMMENT ON COLUMN sys_user.dept_id IS '部门Id';
COMMENT ON COLUMN sys_user.mobile IS '用户手机';
COMMENT ON COLUMN sys_user.email IS '用户邮箱';
COMMENT ON COLUMN sys_user.password IS '用户密码';
COMMENT ON COLUMN sys_user.remark IS '备注';
COMMENT ON COLUMN sys_user.create_time IS '创建时间';
COMMENT ON COLUMN sys_user.update_time IS '修改时间';
COMMENT ON COLUMN sys_user.status IS '状态(1:停用,2:启用)';

-- 用户表.索引
CREATE INDEX idx_sys_user_mobile ON sys_user(mobile);
CREATE INDEX idx_sys_user_status ON sys_user(status);

-- 用户表.基础数据
insert into sys_user (id, name, nickname, dept_id, mobile, email, password, remark, status)
values (1, 'admin', '管理员', 107, '13922110987', '232028819@qq.com',
    'AES-GCM:VkAO3Xpls98BOQ1bILjBcGzC7Ip7zydsgyvd2VxHzA==', '', 2);


-- 大模型表
create table sys_llm (
    id bigserial primary key,
    name varchar(100),
    type varchar(100),
    base_url text,
    api_key varchar(100),
    model_id varchar(100),
    model_name varchar(100),
    remark text,
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp,
    status smallint default 1
);

comment on table sys_llm is '大模型表';
comment on column sys_llm.id is 'ID';
comment on column sys_llm.name is '名称';
comment on column sys_llm.type is '类型';
comment on column sys_llm.base_url is '基本URL';
comment on column sys_llm.api_key is 'API Key';
comment on column sys_llm.model_id is '模型ID';
comment on column sys_llm.model_name is '模型名称';
comment on column sys_llm.remark is '备注';
comment on column sys_llm.create_time is '创建时间';
comment on column sys_llm.update_time is '更新时间';
comment on column sys_llm.status is '状态(1:无效,2:有效)';

insert into sys_llm (id,name, type, base_url, api_key, model_id, model_name, remark,status) values
(1, 'Aliyun' , 'aliyun',  'https://dashscope.aliyuncs.com/compatible-mode/v1', 'sk-16ab6965525b4bd4bd245d9e8a3a693c', 'glm-5', 'glm-5', '智普',   2);

insert into sys_llm (id,name, type, base_url, api_key, model_id, model_name, remark,status) values
(2, 'OpenAI' , 'openai', '', 'GPT-6', 'GPT-6', '未来',  '',2);










-- 1. 角色表
create table if not exists sys_role (
    id serial primary key,
    name varchar(50) not null unique,
    description text,
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp,
    status smallint default 1
);

comment on table sys_role is '用户角色表';
comment on column sys_role.id is '角色ID';
comment on column sys_role.name is '角色名称：admin, user, viewer';
comment on column sys_role.description is '角色描述';
comment on column sys_role.create_time is '创建时间';
comment on column sys_role.update_time is '更新时间';
comment on column sys_role.status is '状态(1:有效,0:无效)';

-- 2. 用户表
create table if not exists sys_user (
    id serial primary key,
    username varchar(100),
    password_hash varchar(255),
    email varchar(255),
    phone varchar(50) not null unique,
    avatar_url text,
    role_id int references sys_role(id) on delete set null,
    last_login_at timestamp,
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp,
    status smallint default 1
    );
comment on table sys_user is '系统内部用户表';
comment on column sys_user.id is '用户ID';
comment on column sys_user.username is '用户名';
comment on column sys_user.password_hash is '密码哈希';
comment on column sys_user.email is '邮箱';
comment on column sys_user.phone is '手机号';
comment on column sys_user.avatar_url is '头像URL';
comment on column sys_user.role_id is '角色ID';
comment on column sys_user.last_login_at is '最后登录时间';
comment on column sys_user.create_time is '创建时间';
comment on column sys_user.update_time is '更新时间';
comment on column sys_user.status is '状态(1:有效,0:无效)';


-- 4. 通道表（飞书机器人/微信公众号等）
create table if not exists sys_channel (
    id bigserial primary key,
    name varchar(255) not null,
    type varchar(255) not null,
    app_id varchar(255) default '',
    app_secret varchar(255) default '',
    link_way varchar(255) default '',
    sort int default 101,
    remark varchar(255) default '',
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp,
    status smallint default 1
    );

comment on table sys_channel is '通道';
comment on column sys_channel.id is 'ID';
comment on column sys_channel.name is '名称';
comment on column sys_channel.type is '类型(feishu,dingtalk,wechat_personal,wechat_work,qq)';
comment on column sys_channel.app_id is '应用ID';
comment on column sys_channel.app_secret is '应用密钥';
comment on column sys_channel.link_way is '链接方式(websocket,webhook,http)';
comment on column sys_channel.sort is '排序';
comment on column sys_channel.remark is '备注';
comment on column sys_channel.create_time is '创建时间';
comment on column sys_channel.update_time is '修改时间';
comment on column sys_channel.status is '状态(1:无效,2:有效)';



-- 3. 通道实例表（对应每个机器人/应用）
create table if not exists sys_channel_instance (
    id serial primary key,
    channel_type varchar(20) not null,
    name varchar(100) not null,
    config jsonb not null,
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp,
    status smallint default 1
);
comment on table sys_channel_instance is '通道实例表（飞书机器人/微信公众号等）';
comment on column sys_channel_instance.id is '实例ID';
comment on column sys_channel_instance.channel_type is '通道类型(feishu,dingtalk,wechat_personal,wechat_work,qq)';
comment on column sys_channel_instance.name is '实例名称';
comment on column sys_channel_instance.config is '平台配置JSON，如{"app_id":"xxx","app_secret":"xxx"}';
comment on column sys_channel_instance.create_time is '创建时间';
comment on column sys_channel_instance.update_time is '更新时间';
comment on column sys_channel_instance.status is '状态(1:有效,0:无效)';

-- 4. 用户绑定表（内部用户 ↔ 外部通道身份）
create table if not exists sys_user_binding (
    id serial primary key,
    internal_user_id int not null references sys_user(id) on delete cascade,
    channel_instance_id int not null references sys_channel_instance(id) on delete cascade,
    external_user_id varchar(100) not null,
    external_chat_id varchar(100),
    external_name varchar(100),
    bind_data jsonb,
    bind_time timestamp default current_timestamp,
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp,
    status smallint default 1,
    unique(channel_instance_id, external_user_id)
    );
comment on table sys_user_binding is '内部用户与外部通道账户绑定关系表';
comment on column sys_user_binding.id is '绑定ID';
comment on column sys_user_binding.internal_user_id is '内部用户ID';
comment on column sys_user_binding.channel_instance_id is '通道实例ID';
comment on column sys_user_binding.external_user_id is '平台侧用户ID（飞书open_id、微信openid等）';
comment on column sys_user_binding.external_chat_id is '会话ID（私聊时等于external_user_id，群聊时为群ID）';
comment on column sys_user_binding.external_name is '昵称（冗余）';
comment on column sys_user_binding.bind_data is '扩展信息';
comment on column sys_user_binding.bind_time is '绑定时间';
comment on column sys_user_binding.create_time is '创建时间';
comment on column sys_user_binding.update_time is '更新时间';
comment on column sys_user_binding.status is '状态(1:有效,0:无效)';

-- 5. 群组配置表
create table if not exists sys_group_config (
    id serial primary key,
    channel_instance_id int not null references sys_channel_instance(id) on delete cascade,
    external_chat_id varchar(100) not null,
    config jsonb not null,
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp,
    status smallint default 1,
    unique(channel_instance_id, external_chat_id)
    );
comment on table sys_group_config is '群组配置表';
comment on column sys_group_config.id is '配置ID';
comment on column sys_group_config.channel_instance_id is '通道实例ID';
comment on column sys_group_config.external_chat_id is '群ID（飞书chat_id、微信群ID等）';
comment on column sys_group_config.config is '群配置JSON，如{"only_reply_when_at":true,"default_model":"gpt-4"}';
comment on column sys_group_config.create_time is '创建时间';
comment on column sys_group_config.update_time is '更新时间';
comment on column sys_group_config.status is '状态(1:有效,0:无效)';

-- 6. 模型规则表
create table if not exists sys_model_rule (
    id serial primary key,
    channel_instance_id int references sys_channel_instance(id) on delete cascade,
    match_type varchar(20) not null,
    match_value text,
    priority int default 0,
    model_provider varchar(50) not null,
    model_name varchar(100) not null,
    model_params jsonb,
    enabled boolean default true,
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp,
    status smallint default 1
    );
comment on table sys_model_rule is 'AI模型调用规则表';
comment on column sys_model_rule.id is '规则ID';
comment on column sys_model_rule.channel_instance_id is '通道实例ID（null表示全局规则）';
comment on column sys_model_rule.match_type is '匹配类型(command,group,user_role,user_id,global)';
comment on column sys_model_rule.match_value is '匹配值（如/help、admin、群ID、用户ID）';
comment on column sys_model_rule.priority is '优先级（数值越小越高）';
comment on column sys_model_rule.model_provider is '模型提供商(openai,azure,local,anthropic)';
comment on column sys_model_rule.model_name is '模型名称';
comment on column sys_model_rule.model_params is '模型参数，如{"temperature":0.7,"max_tokens":2000}';
comment on column sys_model_rule.enabled is '是否启用';
comment on column sys_model_rule.create_time is '创建时间';
comment on column sys_model_rule.update_time is '更新时间';
comment on column sys_model_rule.status is '状态(1:有效,0:无效)';

-- 7. 菜单表
create table if not exists sys_menu (
    id serial primary key,
    parent_id int default 0,
    name varchar(100) not null,
    path varchar(200),
    icon varchar(50),
    permission_key varchar(100) unique,
    sort_order int default 0,
    role_ids jsonb default '[]',
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp,
    status smallint default 1
    );
comment on table sys_menu is '前端菜单配置表';
comment on column sys_menu.id is '菜单ID';
comment on column sys_menu.parent_id is '父菜单ID';
comment on column sys_menu.name is '菜单名称';
comment on column sys_menu.path is '路由路径';
comment on column sys_menu.icon is '图标';
comment on column sys_menu.permission_key is '权限标识';
comment on column sys_menu.sort_order is '排序号';
comment on column sys_menu.role_ids is '允许访问的角色ID数组，如[1,2]';
comment on column sys_menu.create_time is '创建时间';
comment on column sys_menu.update_time is '更新时间';
comment on column sys_menu.status is '状态(1:有效,0:无效)';

-- 8. 审计日志表
create table if not exists sys_audit_log (
    id bigserial primary key,
    user_id int references sys_user(id) on delete set null,
    channel_instance_id int references sys_channel_instance(id) on delete set null,
    external_user_id varchar(100),
    action varchar(50) not null,
    request_data jsonb,
    response_data jsonb,
    error_message text,
    duration_ms int,
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp,
    status smallint default 1
    );
comment on table sys_audit_log is '审计日志表';
comment on column sys_audit_log.id is '日志ID';
comment on column sys_audit_log.user_id is '内部用户ID';
comment on column sys_audit_log.channel_instance_id is '通道实例ID';
comment on column sys_audit_log.external_user_id is '外部用户ID';
comment on column sys_audit_log.action is '操作类型(send_message,receive_message,model_call,login,config_change)';
comment on column sys_audit_log.request_data is '请求数据';
comment on column sys_audit_log.response_data is '响应数据';
comment on column sys_audit_log.error_message is '错误信息';
comment on column sys_audit_log.duration_ms is '耗时(毫秒)';
comment on column sys_audit_log.create_time is '创建时间';
comment on column sys_audit_log.update_time is '更新时间';
comment on column sys_audit_log.status is '状态(1:有效,0:无效)';

-- 9. 系统配置表
create table if not exists sys_config (
    config_key varchar(100) primary key,
    value text,
    description text,
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp,
    status smallint default 1
    );
comment on table sys_config is '系统配置表';
comment on column sys_config.config_key is '配置键';
comment on column sys_config.value is '配置值';
comment on column sys_config.description is '配置描述';
comment on column sys_config.create_time is '创建时间';
comment on column sys_config.update_time is '更新时间';
comment on column sys_config.status is '状态(1:有效,0:无效)';

-- 10. 对话记忆表（向量记忆，用于龙虾记忆）
create table if not exists sys_conversation_memory (
    id bigserial primary key,
    user_id int not null references sys_user(id) on delete cascade,
    channel_instance_id int references sys_channel_instance(id) on delete set null,
    role varchar(20) not null,
    content text not null,
    embedding vector(1536),
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp,
    status smallint default 1
    );
comment on table sys_conversation_memory is '对话记忆表，存储向量化的对话片段（龙虾记忆）';
comment on column sys_conversation_memory.id is '记忆ID';
comment on column sys_conversation_memory.user_id is '用户ID';
comment on column sys_conversation_memory.channel_instance_id is '来源通道实例ID';
comment on column sys_conversation_memory.role is '角色(user/assistant)';
comment on column sys_conversation_memory.content is '原始对话内容';
comment on column sys_conversation_memory.embedding is '向量嵌入(1536维)';
comment on column sys_conversation_memory.create_time is '创建时间';
comment on column sys_conversation_memory.update_time is '更新时间';
comment on column sys_conversation_memory.status is '状态(1:有效,0:无效)';

-- 可选：创建向量检索索引（数据量大时再建）
-- create index idx_memory_embedding on sys_conversation_memory using hnsw (embedding vector_cosine_ops);
create index idx_memory_user_time on sys_conversation_memory(user_id, create_time);

-- ============================================
-- 初始化数据
-- ============================================

insert into sys_role (name, description) values
                                             ('admin', '系统管理员，所有权限'),
                                             ('user', '普通用户，可使用AI功能'),
                                             ('viewer', '只读用户，仅查看日志')
    on conflict (name) do nothing;

insert into sys_menu (parent_id, name, path, icon, permission_key, sort_order, role_ids) values
                                                                                             (0, '控制台', '/dashboard', 'Dashboard', 'menu:dashboard', 1, '[]'),
                                                                                             (0, '对话记录', '/conversations', 'Message', 'menu:conversations', 2, '[]'),
                                                                                             (0, '模型管理', '/models', 'Robot', 'menu:models', 3, '[1]'),
                                                                                             (0, '通道管理', '/channels', 'Api', 'menu:channels', 4, '[1]'),
                                                                                             (0, '用户管理', '/users', 'User', 'menu:users', 5, '[1]'),
                                                                                             (0, '系统设置', '/settings', 'Setting', 'menu:settings', 6, '[1]')
    on conflict (permission_key) do nothing;

insert into sys_config (config_key, value, description) values
                                                            ('auto_create_user', 'true', '是否自动为外部用户创建内部账户'),
                                                            ('default_model_provider', 'openai', '默认模型提供商'),
                                                            ('default_model_name', 'gpt-3.5-turbo', '默认模型名称'),
                                                            ('ws_ping_interval', '30', 'WebSocket心跳间隔（秒）')
    on conflict (config_key) do nothing;

-- 可选：创建自动更新 update_time 的触发器（如果需要）
-- 通用触发器函数
create or replace function update_updated_at_column()
returns trigger as $$
begin
    new.update_time = current_timestamp;
return new;
end;
$$ language plpgsql;

-- 为需要自动更新 update_time 的表创建触发器（以 sys_channel 为例）
-- 其他表按需添加
-- create trigger trigger_sys_channel_update_time
--     before update on sys_channel
--     for each row
--     execute function update_updated_at_column();
