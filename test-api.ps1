# Test Health
Write-Host "Testing Health Endpoint..." -ForegroundColor Green
Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get | ConvertTo-Json

# Test Chat
Write-Host "`nTesting Chat Endpoint..." -ForegroundColor Green
$body = @{
    session_id = "test-session-1"
    message = "Ищу смартфон"
    country = "CH"
    language = "de"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/chat" -Method Post -Body $body -ContentType "application/json" | ConvertTo-Json