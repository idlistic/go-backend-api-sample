CREATE TABLE IF NOT EXISTS timeslots (
  id BIGSERIAL PRIMARY KEY,

  branch_id BIGINT NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,

  service_date DATE NOT NULL,
  start_time TIME NOT NULL,
  end_time TIME NOT NULL,

  capacity INT NOT NULL DEFAULT 1 CHECK (capacity > 0),
  reserved INT NOT NULL DEFAULT 0 CHECK (reserved >= 0),
  is_active BOOLEAN NOT NULL DEFAULT TRUE,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  CHECK (end_time > start_time),
  CHECK (reserved <= capacity)
);

-- index for common query: branch + date
CREATE INDEX IF NOT EXISTS ix_timeslots_branch_date
  ON timeslots (branch_id, service_date);

-- prevent exact duplicate slot in same branch/date/time range
CREATE UNIQUE INDEX IF NOT EXISTS ux_timeslots_unique_slot
  ON timeslots (branch_id, service_date, start_time, end_time);

-- useful when filtering active slots
CREATE INDEX IF NOT EXISTS ix_timeslots_active
  ON timeslots (is_active);
