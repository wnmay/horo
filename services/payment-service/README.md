# Services

# API design

## complete payment (Testing for Powershell)

```
$paymentId = "PUT_PAYMENT_ID_HERE"
Invoke-RestMethod -Uri "http://localhost:3001/api/v1/payments/$paymentId/complete" -Method PUT

```
