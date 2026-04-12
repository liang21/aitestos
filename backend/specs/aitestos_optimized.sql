-- ============================================================
-- 智能测试管理平台 - 优化版数据库 Schema
-- Version: 2.1 Optimized
-- Date: 2026-04-01
-- ============================================================

-- ============================================================
-- 1. 枚举类型定义
-- ============================================================

CREATE TYPE case_status_enum AS ENUM ('unexecuted', 'pass', 'block', 'fail');
CREATE TYPE case_type_enum AS ENUM ('functionality', 'performance', 'api', 'ui', 'security');
CREATE TYPE plan_status_enum AS ENUM ('draft', 'active', 'completed', 'archived');
CREATE TYPE priority_enum AS ENUM ('P0', 'P1', 'P2', 'P3');
CREATE TYPE result_status_enum AS ENUM ('pass', 'fail', 'block', 'skip');
CREATE TYPE user_role_enum AS ENUM ('super_admin', 'admin', 'normal');

-- 新增枚举类型
CREATE TYPE task_status_enum AS ENUM ('pending', 'processing', 'completed', 'failed');
CREATE TYPE draft_status_enum AS ENUM ('pending', 'confirmed', 'rejected');
CREATE TYPE document_type_enum AS ENUM ('prd', 'figma', 'api_spec');

-- ============================================================
-- 2. 核心业务表
-- ============================================================

-- -----------------------------------------------------------
-- project 项目表
-- -----------------------------------------------------------
CREATE TABLE project (
    id          uuid                        DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name        varchar(255)                                            NOT NULL UNIQUE,
    prefix      varchar(4)                                             NOT NULL UNIQUE,
    description text,
    config      jsonb                       DEFAULT '{}'::jsonb,
    created_at  timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN project.prefix IS '项目前缀，用于用例编号生成，2-4位大写字母';

-- -----------------------------------------------------------
-- project_config 项目配置表 (新增)
-- -----------------------------------------------------------
CREATE TABLE project_config (
    id          uuid                        DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    project_id  uuid                        NOT NULL REFERENCES project ON DELETE CASCADE,
    key         varchar(255)                NOT NULL,
    value       jsonb                       DEFAULT '{}'::jsonb,
    description text,
    created_at  timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, key)
);

-- -----------------------------------------------------------
-- module 模块表
-- -----------------------------------------------------------
CREATE TABLE module (
    id            uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    project_id    uuid                            NOT NULL REFERENCES project ON DELETE CASCADE,
    name          varchar(255)                    NOT NULL,
    abbreviation  varchar(4)                      NOT NULL,
    description   text,
    UNIQUE(project_id, name),
    UNIQUE(project_id, abbreviation)
);

COMMENT ON COLUMN module.abbreviation IS '模块缩写，用于用例编号生成，2-4位大写字母，项目内唯一';

-- -----------------------------------------------------------
-- users 用户表
-- -----------------------------------------------------------
CREATE TABLE users (
    id         uuid                        DEFAULT gen_random_uuid()       NOT NULL PRIMARY KEY,
    username   varchar(32)                                                  NOT NULL,
    email      varchar(255)                                                 NOT NULL UNIQUE,
    password   varchar(255)                                                 NOT NULL,
    role       user_role_enum              DEFAULT 'normal'::user_role_enum NOT NULL,
    created_at timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp(3) with time zone
);

-- -----------------------------------------------------------
-- test_case 测试用例表
-- -----------------------------------------------------------
CREATE TABLE test_case (
    id            uuid                        DEFAULT gen_random_uuid()              NOT NULL PRIMARY KEY,
    module_id     uuid                                                                NOT NULL REFERENCES module ON DELETE CASCADE,
    user_id       uuid                                                                NOT NULL REFERENCES users ON DELETE NO ACTION,
    number        varchar(32)                                                         NOT NULL UNIQUE,
    title         varchar(255)                                                        NOT NULL,
    preconditions jsonb                       DEFAULT '[]'::jsonb,
    steps         jsonb                       DEFAULT '[]'::jsonb                     NOT NULL,
    expected      jsonb                       DEFAULT '{}'::jsonb                     NOT NULL,
    ai_metadata   jsonb                       DEFAULT '{}'::jsonb,
    case_type     case_type_enum              DEFAULT 'functionality'::case_type_enum NOT NULL,
    priority      priority_enum               DEFAULT 'P2'::priority_enum             NOT NULL,
    status        case_status_enum            DEFAULT 'unexecuted'::case_status_enum  NOT NULL,
    created_at    timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------
-- test_plan 测试计划表
-- -----------------------------------------------------------
CREATE TABLE test_plan (
    id           uuid                        DEFAULT gen_random_uuid()        NOT NULL PRIMARY KEY,
    project_id   uuid                                                          NOT NULL REFERENCES project ON DELETE CASCADE,
    user_id      uuid                                                          NOT NULL REFERENCES users ON DELETE NO ACTION,
    name         varchar(255)                                                  NOT NULL,
    status       plan_status_enum            DEFAULT 'draft'::plan_status_enum NOT NULL,
    extra_config jsonb                       DEFAULT '{}'::jsonb,
    created_at   timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at   timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------
-- test_result 测试结果表
-- -----------------------------------------------------------
CREATE TABLE test_result (
    id             uuid                        DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    case_id        uuid                                                   NOT NULL REFERENCES test_case ON DELETE CASCADE,
    plan_id        uuid                                                   NOT NULL REFERENCES test_plan ON DELETE CASCADE,
    executor_id    uuid                                                   NOT NULL REFERENCES users ON DELETE NO ACTION,
    execute        result_status_enum                                     NOT NULL,
    result_details jsonb                       DEFAULT '{}'::jsonb,
    executed_at    timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================
-- 3. 知识库与 AI 相关表
-- ============================================================

-- -----------------------------------------------------------
-- document 文档表
-- -----------------------------------------------------------
CREATE TABLE document (
    id           uuid                        DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    project_id   uuid                                                   NOT NULL REFERENCES project ON DELETE CASCADE,
    name         varchar(255)                                           NOT NULL,
    type         document_type_enum                                     NOT NULL,
    status       task_status_enum            DEFAULT 'pending'::task_status_enum NOT NULL,
    url          text,
    content_text text,
    metadata     jsonb                       DEFAULT '{}'::jsonb,
    created_at   timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at   timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN document.status IS '文档处理状态：pending/processing/completed/failed';

-- -----------------------------------------------------------
-- document_chunk 文档分块表
-- -----------------------------------------------------------
CREATE TABLE document_chunk (
    id          uuid                        DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    document_id uuid                                                   NOT NULL REFERENCES document ON DELETE CASCADE,
    chunk_index integer                                                NOT NULL,
    content     text                                                   NOT NULL,
    metadata    jsonb                       DEFAULT '{}'::jsonb,
    created_at  timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------
-- generation_task 生成任务表
-- -----------------------------------------------------------
CREATE TABLE generation_task (
    id             uuid                        DEFAULT gen_random_uuid()           NOT NULL PRIMARY KEY,
    project_id     uuid                                                             NOT NULL REFERENCES project ON DELETE CASCADE,
    user_id        uuid                                                             NOT NULL REFERENCES users ON DELETE NO ACTION,
    status         task_status_enum            DEFAULT 'pending'::task_status_enum  NOT NULL,
    prompt         text,
    result_summary jsonb                       DEFAULT '{}'::jsonb,
    error_msg      text,
    created_at     timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at     timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------
-- generated_case_draft 生成用例草稿表
-- -----------------------------------------------------------
CREATE TABLE generated_case_draft (
    id            uuid                        DEFAULT gen_random_uuid()           NOT NULL PRIMARY KEY,
    task_id       uuid                                                             NOT NULL REFERENCES generation_task ON DELETE CASCADE,
    module_id     uuid                                                             REFERENCES module ON DELETE SET NULL,
    title         varchar(255)                                                     NOT NULL,
    preconditions jsonb                       DEFAULT '[]'::jsonb,
    steps         jsonb                       DEFAULT '[]'::jsonb                  NOT NULL,
    expected      jsonb                       DEFAULT '{}'::jsonb                  NOT NULL,
    case_type     case_type_enum              DEFAULT 'functionality'::case_type_enum,
    priority      priority_enum               DEFAULT 'P2'::priority_enum,
    ai_metadata   jsonb                       DEFAULT '{}'::jsonb,
    status        draft_status_enum           DEFAULT 'pending'::draft_status_enum NOT NULL,
    feedback      text,
    created_at    timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================
-- 4. 索引定义
-- ============================================================

-- project_config 索引
CREATE INDEX idx_project_config_project_id ON project_config(project_id);
CREATE INDEX idx_project_config_key ON project_config(key);

-- users 索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NOT NULL;

-- module 索引
CREATE INDEX idx_module_project_id ON module(project_id);

-- test_case 索引
CREATE INDEX idx_test_case_module_id ON test_case(module_id);
CREATE INDEX idx_test_case_user_id ON test_case(user_id);
CREATE INDEX idx_test_case_number ON test_case(number);
CREATE INDEX idx_test_case_status ON test_case(status);
CREATE INDEX idx_test_case_ai ON test_case USING gin (ai_metadata);
CREATE INDEX idx_test_case_steps ON test_case USING gin (steps);

-- test_plan 索引
CREATE INDEX idx_test_plan_project_id ON test_plan(project_id);
CREATE INDEX idx_test_plan_user_id ON test_plan(user_id);

-- test_result 索引
CREATE INDEX idx_test_result_case_id ON test_result(case_id);
CREATE INDEX idx_test_result_plan_id ON test_result(plan_id);
CREATE INDEX idx_test_result_executor_id ON test_result(executor_id);
CREATE INDEX idx_test_result_executed_at ON test_result(executed_at);

-- document 索引
CREATE INDEX idx_document_project_id ON document(project_id);
CREATE INDEX idx_document_type ON document(type);
CREATE INDEX idx_document_status ON document(status);

-- document_chunk 索引
CREATE INDEX idx_document_chunk_document_id ON document_chunk(document_id);

-- generation_task 索引
CREATE INDEX idx_generation_task_project_id ON generation_task(project_id);
CREATE INDEX idx_generation_task_user_id ON generation_task(user_id);
CREATE INDEX idx_generation_task_status ON generation_task(status);

-- generated_case_draft 索引
CREATE INDEX idx_draft_task_id ON generated_case_draft(task_id);
CREATE INDEX idx_draft_module_id ON generated_case_draft(module_id);
CREATE INDEX idx_draft_status ON generated_case_draft(status);

-- 复合索引（性能优化）
CREATE INDEX idx_test_result_plan_executed ON test_result(plan_id, executed_at DESC);
CREATE INDEX idx_generation_task_project_status ON generation_task(project_id, status);

-- ============================================================
-- 5. 触发器
-- ============================================================

-- updated_at 自动更新函数
CREATE FUNCTION update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;

-- 为各表添加触发器
CREATE TRIGGER update_project_updated_at
    BEFORE UPDATE ON project FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_project_config_updated_at
    BEFORE UPDATE ON project_config FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_test_case_updated_at
    BEFORE UPDATE ON test_case FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_test_plan_updated_at
    BEFORE UPDATE ON test_plan FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_document_updated_at
    BEFORE UPDATE ON document FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_generation_task_updated_at
    BEFORE UPDATE ON generation_task FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_generated_case_draft_updated_at
    BEFORE UPDATE ON generated_case_draft FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================
-- 6. 注释 (可选，提升可读性)
-- ============================================================

COMMENT ON TABLE project IS '项目表';
COMMENT ON TABLE project_config IS '项目配置表，支持键值对扩展';
COMMENT ON TABLE module IS '模块表，项目下的功能模块';
COMMENT ON TABLE users IS '用户表';
COMMENT ON TABLE test_case IS '测试用例表';
COMMENT ON TABLE test_plan IS '测试计划表';
COMMENT ON TABLE test_result IS '测试执行结果表';
COMMENT ON TABLE document IS '文档表（PRD/Figma/API Spec）';
COMMENT ON TABLE document_chunk IS '文档分块表，用于向量检索';
COMMENT ON TABLE generation_task IS 'AI 用例生成任务表';
COMMENT ON TABLE generated_case_draft IS '生成的用例草稿表';

COMMENT ON COLUMN project.prefix IS '项目前缀，用于用例编号生成，2-4位大写字母';
COMMENT ON COLUMN module.abbreviation IS '模块缩写，用于用例编号生成，2-4位大写字母，项目内唯一';
COMMENT ON COLUMN document.status IS '文档处理状态：pending/processing/completed/failed';
COMMENT ON COLUMN test_case.ai_metadata IS 'AI 元数据，包含生成任务ID、置信度、引用文档块等信息';
COMMENT ON COLUMN generated_case_draft.ai_metadata IS 'AI 元数据，包含置信度、引用文档块等信息';

COMMENT ON TYPE task_status_enum IS '任务状态枚举';
COMMENT ON TYPE draft_status_enum IS '草稿状态枚举';
COMMENT ON TYPE document_type_enum IS '文档类型枚举';
