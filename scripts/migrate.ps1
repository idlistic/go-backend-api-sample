param (
  [string]$DbUser = "go_backend_api",
  [string]$DbName = "go_backend_api_db",
  [string]$Container = "go_backend_api_sample_db"
)

Write-Host "🚀 Running migrations..."

docker exec -i $Container psql -U $DbUser -d $DbName `
  -f /migrations/001_create_branches.sql `
  -f /migrations/002_create_timeslots.sql `
  -f /migrations/003_create_orders.sql `
  -f /seed/seed.sql

Write-Host "✅ Migration completed"
