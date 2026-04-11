
-- 创建数据库
create database if not exists ym3_main;

-- 启用 pgvector 扩展（向量记忆需要）
create extension if not exists vector;

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
