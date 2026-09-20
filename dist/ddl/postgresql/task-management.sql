DROP SCHEMA IF EXISTS task CASCADE;
CREATE SCHEMA task;

-- TODO: PK,FK,UNIQUE,INDEX
-- TODO: FUNCTION/TRIGGER(BUILT IN)
-- TODO: COMMENT
-- TODO: FUNCTION/TRIGGER(CUSTOM)
-- TODO: CONSTRAINT(CUSTOM)

CREATE TABLE task.tasks (
  task_id varchar(26) NOT NULL, --TODO: CUSTOM FUNCTIO
  CHECK (LENGTH(task_id) = 26),
  CHECK (task_id ~* '^[0-9A-HJKMNP-TV-Z]{26}$'),
    -- elementRef: $taskId

  title varchar(50) NOT NULL,
    -- elementRef: $title

  description text,
  CHECK (LENGTH(description) >= 10),
    -- elementRef: $description

  create_user varchar(5),
    -- elementRef: $userId
    -- override:
    --   name: 登録ユーザー

  task_pic varchar(5),
  CHECK (LENGTH(task_pic) = 5),
  CHECK (task_pic ~* '^U[0-9]{4}$'),
    -- elementRef: $userId
    -- override:
    --   name: 担当者

  priority bigint NOT NULL DEFAULT 3,
  CHECK (priority >= 1),
  CHECK (priority <= 5),
    -- elementRef: $priority

  status_code varchar(15) NOT NULL,
    -- elementRef: $taskStatus
    -- override:
    --   name: ステータスコード

  created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,

  updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);


CREATE TABLE task.statuses (
  status_code varchar(15) NOT NULL,
    -- elementRef: $taskStatus
    -- override:
    --   name: ステータスコード

  description text,
  CHECK (LENGTH(description) >= 10),
    -- elementRef: $description

  created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,

  updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp
);
