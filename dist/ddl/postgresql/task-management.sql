-- INFO: Create Schema
DROP SCHEMA IF EXISTS task CASCADE;
CREATE SCHEMA task;

-- INFO: Create Table(tasks)
CREATE TABLE task.tasks (
  task_id varchar(26) NOT NULL,
  CHECK (LENGTH(task_id) = 26),
  CHECK (task_id ~* '^[0-9A-HJKMNP-TV-Z]{26}$'),

  title varchar(50) NOT NULL,

  description text,
  CHECK (LENGTH(description) >= 10),

  create_user varchar(5),

  task_pic varchar(5),
  CHECK (LENGTH(task_pic) = 5),
  CHECK (task_pic ~* '^U[0-9]{4}$'),

  priority bigint NOT NULL DEFAULT 3,
  CHECK (priority >= 1),
  CHECK (priority <= 5),

  status_code varchar(15) NOT NULL,

  created_at timestamp with time zone NOT NULL,

  updated_at timestamp with time zone NOT NULL
);

-- INFO: Create Table(statuses)
CREATE TABLE task.statuses (
  status_code varchar(15) NOT NULL,

  description text NOT NULL,
  CHECK (LENGTH(description) >= 10),

  created_at timestamp with time zone NOT NULL,

  updated_at timestamp with time zone NOT NULL
);

-- INFO: Set PK Constraint
ALTER TABLE task.tasks
  ADD CONSTRAINT pk_tasks
  PRIMARY KEY (task_id);

ALTER TABLE task.statuses
  ADD CONSTRAINT pk_statuses
  PRIMARY KEY (status_code);

-- INFO: Set FK constraint
ALTER TABLE task.tasks
  ADD CONSTRAINT fk_tasks_7a3910e6
  FOREIGN KEY (status_code)
  REFERENCES task.statuses (status_code)
  ON DELETE NO ACTION ON UPDATE NO ACTION;

-- INFO: Create Index
CREATE INDEX idx_tasks_fb4a92c9 ON task.tasks (
  status_code ASC,
  priority DESC,
  task_pic ASC NULLS FIRST,
  task_id ASC
);

CREATE INDEX idx_tasks_e0739e6d ON task.tasks (
  create_user ASC,
  task_id ASC
);

-- INFO: Set Comment
COMMENT ON TABLE task.tasks IS 'タスク';
COMMENT ON COLUMN task.tasks.task_id IS 'タスクID [Element: taskId]';
COMMENT ON COLUMN task.tasks.title IS 'タイトル [Element: title]';
COMMENT ON COLUMN task.tasks.description IS '説明 [Element: description]';
COMMENT ON COLUMN task.tasks.create_user IS '登録ユーザー [Element: userId]';
COMMENT ON COLUMN task.tasks.task_pic IS '担当者 [Element: userId]';
COMMENT ON COLUMN task.tasks.priority IS '優先度 [Element: priority]';
COMMENT ON COLUMN task.tasks.status_code IS 'ステータスコード [Element: taskStatus]';
COMMENT ON COLUMN task.tasks.created_at IS '作成日時 [BuiltIn]';
COMMENT ON COLUMN task.tasks.updated_at IS '更新日時 [BuiltIn]';

COMMENT ON TABLE task.statuses IS 'ステータス';
COMMENT ON COLUMN task.statuses.status_code IS 'ステータスコード [Element: taskStatus]';
COMMENT ON COLUMN task.statuses.description IS '説明 [Element: description]';
COMMENT ON COLUMN task.statuses.created_at IS '作成日時 [BuiltIn]';
COMMENT ON COLUMN task.statuses.updated_at IS '更新日時 [BuiltIn]';

-- INFO: ARIADNE generated BuiltIn Function
CREATE FUNCTION task.ariadne_builtin_row_metadata_insert()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  NEW.created_at := current_timestamp;
  NEW.updated_at := current_timestamp;
  RETURN NEW;
END;
$$;

CREATE FUNCTION task.ariadne_builtin_row_metadata_update()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  NEW.updated_at := current_timestamp;
  RETURN NEW;
END;
$$;

CREATE FUNCTION task.ariadne_builtin_ulid()
RETURNS varchar
LANGUAGE plpgsql
AS $$
DECLARE
  v_alphabet CONSTANT varchar := '0123456789ABCDEFGHJKMNPQRSTVWXYZ';
  v_timestamp bigint;
  v_timestamp_part varchar := '';
  v_random_part varchar := '';
  v_value bigint;
  v_index integer;
BEGIN
  -- Unix time in milliseconds (48-bit timestamp component)
  v_timestamp := floor(extract(epoch FROM clock_timestamp()) * 1000)::bigint;
  v_value := v_timestamp;

  -- Encode timestamp as 10-character Crockford Base32.
  -- Prepending each digit preserves the most-significant-digit-first order.
  FOR i IN 1..10 LOOP
    v_index := (v_value % 32)::integer;
    v_timestamp_part :=
      substr(v_alphabet, v_index + 1, 1) || v_timestamp_part;
    v_value := v_value / 32;
  END LOOP;

  -- Generate 16-character / 80-bit randomness component.
  FOR i IN 1..16 LOOP
    v_index := floor(random() * 32)::integer;
    v_random_part :=
      v_random_part || substr(v_alphabet, v_index + 1, 1);
  END LOOP;

  RETURN v_timestamp_part || v_random_part;
END;
$$;

-- INFO: ARIADNE generated BuiltIn Trigger
CREATE TRIGGER ariadne_builtin_tasks_before_insert
BEFORE INSERT ON task.tasks
FOR EACH ROW
EXECUTE FUNCTION task.ariadne_builtin_row_metadata_insert();

CREATE TRIGGER ariadne_builtin_tasks_before_update
BEFORE UPDATE ON task.tasks
FOR EACH ROW
EXECUTE FUNCTION task.ariadne_builtin_row_metadata_update();

CREATE TRIGGER ariadne_builtin_statuses_before_insert
BEFORE INSERT ON task.statuses
FOR EACH ROW
EXECUTE FUNCTION task.ariadne_builtin_row_metadata_insert();

CREATE TRIGGER ariadne_builtin_statuses_before_update
BEFORE UPDATE ON task.statuses
FOR EACH ROW
EXECUTE FUNCTION task.ariadne_builtin_row_metadata_update();

-- INFO: ARIADNE generated Generation Trigger Function
CREATE OR REPLACE FUNCTION task.ariadne_generation_tasks()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  NEW.task_id := task.ariadne_builtin_ulid();

  RETURN NEW;
END;
$$;

CREATE TRIGGER ariadne_generation_tasks
BEFORE INSERT
ON task.tasks
FOR EACH ROW
EXECUTE FUNCTION task.ariadne_generation_tasks();
