
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


-- 组织表
CREATE TABLE IF NOT EXISTS sys_org (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    invite_code    VARCHAR(100) NOT NULL,
    user_id     BIGSERIAL DEFAULT NULL,
    remark      VARCHAR(100) DEFAULT NULL,
    create_time TIMESTAMP DEFAULT NULL,
    update_time TIMESTAMP DEFAULT NULL,
    status      SMALLINT DEFAULT 2
);

COMMENT ON TABLE sys_org IS '组织';
COMMENT ON COLUMN sys_org.id IS '组织Id';
COMMENT ON COLUMN sys_org.name IS '组织名称';
COMMENT ON COLUMN sys_org.invite_code IS '邀请编码';
COMMENT ON COLUMN sys_org.user_id IS '管理员用户Id';
COMMENT ON COLUMN sys_org.remark IS '备注';
COMMENT ON COLUMN sys_org.create_time IS '创建时间';
COMMENT ON COLUMN sys_org.update_time IS '修改时间';
COMMENT ON COLUMN sys_org.status IS '状态(1:停用,2:启用)';

-- 组织表.索引
CREATE INDEX idx_sys_org_name ON sys_org(name);
CREATE INDEX idx_sys_org_invite_code ON sys_org(invite_code);
CREATE INDEX idx_sys_org_status ON sys_org(status);

-- 组织表.基础数据
insert into sys_org (id, name, invite_code, user_id, remark, status)
values (1, 'xx公司邀请编码', 'xGoxNodeXpy',1, '', 2);


-- 用户表
CREATE TABLE IF NOT EXISTS sys_user (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    nickname    VARCHAR(100) DEFAULT NULL,
    org_id     BIGSERIAL DEFAULT NULL,
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
COMMENT ON COLUMN sys_user.org_id IS '组织Id';
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
insert into sys_user (id, name, nickname, org_id, mobile, email, password, remark, status)
values (1, 'admin', '管理员', 1, '13922110987', '232028819@qq.com',
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
    status smallint default 2
);

comment on table sys_llm is '用户大模型表';
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

insert into sys_llm (id, name, type, base_url, api_key, model_id, model_name, remark,status) values
    (1,'Aliyun' , 'aliyun',  'https://dashscope.aliyuncs.com/compatible-mode/v1', 'sk-16ab6965525b4bd4bd245d9e8a3a693c', 'glm-5', 'glm-5', '智普',   2);

insert into sys_llm (id, name, type, base_url, api_key, model_id, model_name, remark,status) values
    (2,'OpenAI' , 'openai', '', 'GPT-6', 'GPT-6', '未来',  '',2);

-- 用户大模型表
create table sys_user_llm (
    id bigserial primary key,
    user_id bigserial,
    llm_id bigserial,
    name varchar(100),
    type varchar(100),
    base_url text,
    api_key varchar(100),
    model_id varchar(100),
    model_name varchar(100),
    remark text,
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp,
    status smallint default 2
);

comment on table sys_user_llm is '用户大模型表';
comment on column sys_user_llm.user_id is '用户ID';
comment on column sys_user_llm.llm_id is '大模型ID';
comment on column sys_user_llm.name is '名称';
comment on column sys_user_llm.type is '类型';
comment on column sys_user_llm.base_url is '基本URL';
comment on column sys_user_llm.api_key is 'API Key';
comment on column sys_user_llm.model_id is '模型ID';
comment on column sys_user_llm.model_name is '模型名称';
comment on column sys_user_llm.remark is '备注';
comment on column sys_user_llm.create_time is '创建时间';
comment on column sys_user_llm.update_time is '更新时间';
comment on column sys_user_llm.status is '状态(1:无效,2:有效)';

insert into sys_user_llm (id,user_id, llm_id, name, type, base_url, api_key, model_id, model_name, remark,status) values
    (1,1,  1,'Aliyun' , 'aliyun',  'https://dashscope.aliyuncs.com/compatible-mode/v1', 'sk-16ab6965525b4bd4bd245d9e8a3a693c', 'glm-5', 'glm-5', '智普',   2);

insert into sys_user_llm (id,user_id, llm_id, name, type, base_url, api_key, model_id, model_name, remark,status) values
    (2,1,  2,'OpenAI' , 'openai', '', 'GPT-6', 'GPT-6', '未来',  '',2);


-- Agent表
CREATE TABLE IF NOT EXISTS sys_agent (
    id             BIGSERIAL PRIMARY KEY,
    name           VARCHAR(100) NOT NULL,
    nick_name   VARCHAR(200) DEFAULT NULL,
    system_prompt  TEXT NOT NULL,
    create_time    TIMESTAMP DEFAULT NULL,
    update_time    TIMESTAMP DEFAULT NULL,
    status         SMALLINT DEFAULT 2
    );

COMMENT ON TABLE sys_agent IS 'Agent定义表';
COMMENT ON COLUMN sys_agent.id IS 'Agent ID';
COMMENT ON COLUMN sys_agent.user_id IS '用户ID';
COMMENT ON COLUMN sys_agent.agent_id IS 'Agent ID';
COMMENT ON COLUMN sys_agent.name IS 'Agent唯一标识名';
COMMENT ON COLUMN sys_agent.nick_name IS '显示名称';
COMMENT ON COLUMN sys_agent.system_prompt IS '系统提示词';
COMMENT ON COLUMN sys_agent.create_time IS '创建时间';
COMMENT ON COLUMN sys_agent.update_time IS '修改时间';
COMMENT ON COLUMN sys_agent.status IS '状态(1:停用,2:启用)';

CREATE INDEX idx_sys_agent_name ON sys_agent(name);
CREATE INDEX idx_sys_agent_status ON sys_agent(status);

insert into sys_agent (id,name, nick_name, system_prompt, status) values
    (1,'weather_agent','天气助手','你是一个天气助手，你可以回答用户关于天气的问题。',2);

-- 用户Agent表
CREATE TABLE IF NOT EXISTS sys_user_agent (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGSERIAL,
    agent_id       BIGSERIAL,
    agent_name           VARCHAR(100) NOT NULL,
    agent_nick_name   VARCHAR(200) DEFAULT NULL,
    system_prompt  TEXT NOT NULL,
    create_time    TIMESTAMP DEFAULT NULL,
    update_time    TIMESTAMP DEFAULT NULL,
    status         SMALLINT DEFAULT 2
    );

COMMENT ON TABLE sys_user_agent IS '用户Agent定义表';
COMMENT ON COLUMN sys_user_agent.id IS '用户Agent ID';
COMMENT ON COLUMN sys_user_agent.user_id IS '用户ID';
COMMENT ON COLUMN sys_user_agent.agent_id IS 'Agent ID';
COMMENT ON COLUMN sys_user_agent.agent_name IS 'Agent唯一标识名';
COMMENT ON COLUMN sys_user_agent.agent_nick_name IS '显示名称';
COMMENT ON COLUMN sys_user_agent.system_prompt IS '系统提示词';
COMMENT ON COLUMN sys_user_agent.create_time IS '创建时间';
COMMENT ON COLUMN sys_user_agent.update_time IS '修改时间';
COMMENT ON COLUMN sys_user_agent.status IS '状态(1:停用,2:启用)';

CREATE INDEX idx_sys_user_agent_name ON sys_user_agent(name);
CREATE INDEX idx_sys_user_agent_status ON sys_user_agent(status);

insert into sys_user_agent (id,user_id, agent_id, agent_name, agent_nick_name, system_prompt, status) values
(1,1,1,'weather_agent','天气助手','你是一个天气助手，你可以回答用户关于天气的问题。',2);


-- skill表
CREATE TABLE IF NOT EXISTS sys_skill (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    description     TEXT NOT NULL,
    parameters      JSONB NOT NULL,
    executor_type   VARCHAR(20) NOT NULL,
    executor_config JSONB NOT NULL,
    create_time     TIMESTAMP DEFAULT NULL,
    update_time     TIMESTAMP DEFAULT NULL,
    status          SMALLINT DEFAULT 2
    );

COMMENT ON TABLE sys_skill IS '技能定义表';
COMMENT ON COLUMN sys_skill.id IS '技能ID';
COMMENT ON COLUMN sys_skill.name IS '技能名称';
COMMENT ON COLUMN sys_skill.description IS '技能描述';
COMMENT ON COLUMN sys_skill.parameters IS '参数JSON Schema';
COMMENT ON COLUMN sys_skill.executor_type IS '执行器类型';
COMMENT ON COLUMN sys_skill.executor_config IS '执行器配置';
COMMENT ON COLUMN sys_skill.create_time IS '创建时间';
COMMENT ON COLUMN sys_skill.update_time IS '修改时间';
COMMENT ON COLUMN sys_skill.status IS '状态(1:停用,2:启用)';

CREATE INDEX idx_sys_skill_name ON sys_skill(name);
CREATE INDEX idx_sys_skill_status ON sys_skill(status);

insert into sys_skill (id,name, description, parameters, executor_type, executor_config, status) values
    (1,'get_weather','获取天气信息','{"location": "string"}','python','{"script_path":"/path/to/get_weather.py"}',2);

-- 用户skill表
CREATE TABLE IF NOT EXISTS sys_user_skill (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGSERIAL,
    skill_id        BIGSERIAL,
    skill_name            VARCHAR(100) NOT NULL,
    description     TEXT NOT NULL,
    parameters      JSONB NOT NULL,
    executor_type   VARCHAR(20) NOT NULL,
    executor_config JSONB NOT NULL,
    create_time     TIMESTAMP DEFAULT NULL,
    update_time     TIMESTAMP DEFAULT NULL,
    status          SMALLINT DEFAULT 2
    );

COMMENT ON TABLE sys_user_skill IS '用户技能定义表';
COMMENT ON COLUMN sys_user_skill.id IS '技能ID';
COMMENT ON COLUMN sys_user_skill.skill_id IS '技能ID';
COMMENT ON COLUMN sys_user_skill.skill_name IS '技能名称';
COMMENT ON COLUMN sys_user_skill.description IS '技能描述';
COMMENT ON COLUMN sys_user_skill.parameters IS '参数JSON Schema';
COMMENT ON COLUMN sys_user_skill.executor_type IS '执行器类型';
COMMENT ON COLUMN sys_user_skill.executor_config IS '执行器配置';
COMMENT ON COLUMN sys_user_skill.create_time IS '创建时间';
COMMENT ON COLUMN sys_user_skill.update_time IS '修改时间';
COMMENT ON COLUMN sys_user_skill.status IS '状态(1:停用,2:启用)';

CREATE INDEX idx_sys_user_skill_name ON sys_user_skill(name);
CREATE INDEX idx_sys_user_skill_status ON sys_user_skill(status);

insert into sys_user_skill (id,user_id, skill_id, skill_name, description, parameters, executor_type, executor_config, status) values
    (1,1,1,'get_weather','获取天气信息','{"location": "string"}','python','{"script_path":"/path/to/get_weather.py"}',2);



-- 知识库定义表
CREATE TABLE IF NOT EXISTS sys_knowledge_base (
                                                  id          BIGSERIAL PRIMARY KEY,
                                                  org_id      BIGINT NOT NULL,
                                                  user_id     BIGINT,                 -- NULL 表示组织共有
                                                  name        VARCHAR(100) NOT NULL,
    scope       SMALLINT DEFAULT 3,     -- 1:系统, 2:组织, 3:个人
    is_active   BOOLEAN DEFAULT TRUE,
    remark      TEXT
    );

-- 知识库分片（带向量）
CREATE TABLE IF NOT EXISTS sys_knowledge_chunk (
                                                   id          BIGSERIAL PRIMARY KEY,
                                                   kb_id       BIGINT NOT NULL,
                                                   org_id      BIGINT NOT NULL,        -- 强制带上 org_id 以便物理隔离
                                                   content     TEXT,
                                                   embedding   VECTOR(1536),           -- 向量字段
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- 用户记忆表
CREATE TABLE IF NOT EXISTS sys_user_memory (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGSERIAL,
    description     TEXT NOT NULL,
    parameters      JSONB NOT NULL,
    create_time     TIMESTAMP DEFAULT NULL,
    update_time     TIMESTAMP DEFAULT NULL,
    status          SMALLINT DEFAULT 2
    );

COMMENT ON TABLE sys_user_memory IS '用户记忆表';
COMMENT ON COLUMN sys_user_memory.id IS '用户记忆ID';
COMMENT ON COLUMN sys_user_memory.user_id IS '用户ID';
COMMENT ON COLUMN sys_user_memory.description IS '记忆描述';
COMMENT ON COLUMN sys_user_memory.parameters IS '参数JSON Schema';
COMMENT ON COLUMN sys_user_memory.create_time IS '创建时间';
COMMENT ON COLUMN sys_user_memory.update_time IS '修改时间';
COMMENT ON COLUMN sys_user_memory.status IS '状态(1:停用,2:启用)';

CREATE INDEX idx_sys_user_memory_description ON sys_user_memory(description);
CREATE INDEX idx_sys_user_memory_status ON sys_user_memory(status);

insert into sys_user_memory (id,user_id, description, parameters) values
    (1,1, '你是一个技术专家','{"location": "string"}');


-- 用户通道表
create table if not exists sys_user_channel (
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
    status smallint default 2
    );

comment on table sys_user_channel is '用户通道表';
comment on column sys_user_channel.id is 'ID';
comment on column sys_user_channel.name is '名称';
comment on column sys_user_channel.type is '类型(feishu,dingtalk,wechat_personal,wechat_work,qq)';
comment on column sys_user_channel.app_id is '应用ID';
comment on column sys_user_channel.app_secret is '应用密钥';
comment on column sys_user_channel.link_way is '链接方式(websocket,webhook,http)';
comment on column sys_user_channel.sort is '排序';
comment on column sys_user_channel.remark is '备注';
comment on column sys_user_channel.create_time is '创建时间';
comment on column sys_user_channel.update_time is '修改时间';
comment on column sys_user_channel.status is '状态(1:无效,2:有效)';


insert into sys_user_channel (id,name,type,app_id,app_secret,link_way,sort,remark,create_time,update_time,status) values
    (1,'feishu','飞书','feishu_app_id','feishu_app_secret','feishu_link_way',101,'feishu_remark',current_timestamp,current_timestamp,2);



-- 会话表
CREATE TABLE IF NOT EXISTS sys_chat_session (
                                                id          BIGSERIAL PRIMARY KEY,
                                                org_id      BIGINT NOT NULL,
                                                user_id     BIGINT NOT NULL,
                                                agent_id    BIGINT NOT NULL,
                                                title       VARCHAR(200),
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- 消息明细表
CREATE TABLE IF NOT EXISTS sys_chat_message (
                                                id          BIGSERIAL PRIMARY KEY,
                                                session_id  BIGINT NOT NULL,
                                                role        VARCHAR(20) NOT NULL,      -- 'system', 'user', 'assistant', 'tool'
    content     TEXT,
    tool_calls  JSONB,                     -- 记录 AI 决定调用的工具参数
    token_usage INT,                       -- Token 消耗统计
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );