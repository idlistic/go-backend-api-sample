DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
    CREATE TYPE order_status AS ENUM ('created', 'cancelled');
  END IF;
END$$;

CREATE TABLE IF NOT EXISTS orders (
  id BIGSERIAL PRIMARY KEY,

  branch_id BIGINT NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
  timeslot_id BIGINT NOT NULL REFERENCES timeslots(id) ON DELETE RESTRICT,

  customer_name TEXT NOT NULL,
  status order_status NOT NULL DEFAULT 'created',

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS ix_orders_branch_created_at
  ON orders (branch_id, created_at DESC);

CREATE INDEX IF NOT EXISTS ix_orders_timeslot_id
  ON orders (timeslot_id);

CREATE INDEX IF NOT EXISTS ix_orders_status
  ON orders (status);
