CREATE TABLE IF NOT EXISTS branches (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- optional: prevent duplicate branch names
CREATE UNIQUE INDEX IF NOT EXISTS ux_branches_name ON branches (name);
