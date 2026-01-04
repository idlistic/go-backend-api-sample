INSERT INTO branches (name)
VALUES ('Chiang Mai - Branch 1')
ON CONFLICT (name) DO NOTHING;

-- Create a few demo timeslots for today and tomorrow
WITH b AS (
  SELECT id AS branch_id FROM branches WHERE name = 'Chiang Mai - Branch 1'
)
INSERT INTO timeslots (branch_id, service_date, start_time, end_time, capacity, reserved, is_active)
SELECT b.branch_id, CURRENT_DATE, TIME '10:00', TIME '11:00', 3, 0, TRUE FROM b
UNION ALL
SELECT b.branch_id, CURRENT_DATE, TIME '11:00', TIME '12:00', 3, 0, TRUE FROM b
UNION ALL
SELECT b.branch_id, CURRENT_DATE + 1, TIME '10:00', TIME '11:00', 2, 0, TRUE FROM b
ON CONFLICT (branch_id, service_date, start_time, end_time) DO NOTHING;
