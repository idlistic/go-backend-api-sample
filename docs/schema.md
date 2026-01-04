# Database Schema (Postgres)

## branches
- id (PK)
- name (unique)
- created_at

## timeslots
- id (PK)
- branch_id (FK -> branches.id)
- service_date, start_time, end_time
- capacity, reserved, is_active
- created_at, updated_at

Unique:
- (branch_id, service_date, start_time, end_time)

Indexes:
- (branch_id, service_date)

## orders
- id (PK)
- branch_id (FK -> branches.id)
- timeslot_id (FK -> timeslots.id)
- customer_name
- status: created | cancelled
- created_at, updated_at
