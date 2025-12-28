# Go Backend API Sample

[English Version](#english) | [เวอร์ชันภาษาไทย](#ภาษาไทย)

---

## English

### Project Overview
This project is a sample backend API built with Go and PostgreSQL.
It demonstrates real-world backend patterns such as transactional operations,
row-level locking, and basic business logic around timeslot reservations.

This repository is intended for job application and technical demonstration purposes.

---

### Features
- List branches
- List timeslots by branch and date
- Create order with timeslot reservation (transactional)
- Cancel order and release reserved timeslot

---

### Tech Stack
- Go (net/http)
- PostgreSQL
- Docker & Docker Compose

---

### API Endpoints
```http
GET    /branches
GET    /timeslots?branch_id=&date=
POST   /orders
PATCH  /orders/{id}/cancel
