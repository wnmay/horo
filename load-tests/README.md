# Load Testing with K6

## Quick Start

```powershell
# Set Firebase API Key (ONLY credential needed!)
$env:FIREBASE_API_KEY="your-firebase-api-key"

# Run full flow test (100 unique users, 2 minutes)
k6 run load-tests/full-flow-load-test.js
```

## Prerequisites

### 1. Install K6

https://k6.io/docs/getting-started/installation/

```powershell
# Windows (using Chocolatey)
choco install k6
```

### 2. Firebase Configuration

- Go to Firebase Console → Authentication → Sign-in method
- **Enable Email/Password Authentication**
- ⚠️ **Important**: Firebase free tier allows **100 signups per hour**

### 3. Firebase API Key

- Go to Firebase Console → Project Settings
- Find "Web API Key" under General tab
- This is the ONLY credential you need!

### 4. Test Data

- Ensure at least one course exists in the database
- All services must be running (API Gateway port 8080)

## Test Scenarios

### Full Flow Load Test (`full-flow-load-test.js`)

Tests complete user journey with **100 unique real Firebase users**:

0. ✅ **Register Firebase user** (creates unique Firebase account per VU)
1. ✅ Register user in backend with Firebase token
2. ✅ Get all courses
3. ✅ Create chat room with a course
4. ✅ Create order with the chat room
5. ✅ Complete payment for the order

**Load Pattern**:

- 0-20s: Ramp to 50 users
- 20-50s: Ramp to 100 users
- 50-90s: Hold at 100 users
- 90-120s: Ramp down to 0

**Requirements**:

- ✅ P95 response time < 3 seconds
- ✅ Error rate < 10%

**How it works** (TRUE load testing):

- Each VU creates a **unique Firebase account**
- Email format: `loadtest_<VU>_<timestamp>_<random>@example.com`
- Each VU gets their **own Firebase token**
- Tests **100 concurrent different users** (realistic scenario!)

## Running Tests

### Minimal Command

```powershell
$env:FIREBASE_API_KEY="AIza..."; k6 run load-tests/full-flow-load-test.js
```

### Minimal Command

```powershell
$env:FIREBASE_API_KEY="AIza..."; k6 run load-tests/full-flow-load-test.js
```

### With Custom Base URL

```powershell
$env:FIREBASE_API_KEY="your-key"
$env:BASE_URL="http://your-api-gateway:8080"

k6 run load-tests/full-flow-load-test.js
```

### Test with Fewer Users (for testing)

```powershell
# Test with 10 users to verify setup
$env:FIREBASE_API_KEY="your-key"; k6 run --vus 10 --duration 30s load-tests/full-flow-load-test.js
```

## Why 100 Unique Users?

### ❌ Wrong Approach (Shared Token):

- 1 Firebase user → 1 token → 100 VUs using same token
- All VUs have the same identity
- Not realistic - doesn't test concurrent different users
- Can't properly test user isolation, permissions, data segregation

### ✅ Correct Approach (Unique Users):

- 100 Firebase users → 100 tokens → 100 VUs with different identities
- Each VU has unique email, token, and backend user ID
- **TRUE concurrent user simulation**
- Tests realistic load: 100 different customers using the system
- Validates proper user isolation and data separation

⚠️ **Firebase Rate Limits**: Free tier = 100 signups/hour. Paid tier = unlimited.

## Understanding Results

Example output:

```
🔐 Registering new Firebase user: loadtest_1_1699824123456_abc@example.com
✅ Successfully registered Firebase user
✓ registration status is 201
✓ get courses completed within 3s
✓ create chat room status is 200 or 201
✓ create order completed within 3s
✓ complete payment status is 200

============================================
  LOAD TEST SUMMARY
============================================

HTTP Requests:
  Total: 5000
  Rate: 41.67/s
  Failed: 2.5%

Response Times:
  Avg: 850ms
  P95: 2450ms

REQUIREMENT CHECK:
  ✓ 95% of requests < 3s: PASS
  ✓ Error rate < 10%: PASS
============================================
```

Key metrics:

- **http_req_duration P95**: Should be < 3000ms
- **Error Rate**: Should be < 10%
- **Step timings**: Individual operation performance

## Troubleshooting

### ❌ "Firebase authentication failed"

- Verify `FIREBASE_API_KEY` is correct
- Confirm user exists in Firebase Authentication
- Check password is correct

### ❌ "courses list is empty"

- Create at least one course in the database
- Verify course-service is running

### ❌ High error rate (>10%)

- Check service logs for errors
- Verify database connection pool size
- Reduce concurrent users (modify stages)
- Check token hasn't expired

### ❌ Slow responses (P95 >3s)

- Add database indexes
- Increase service resources (CPU/Memory)
- Optimize slow queries
- Scale services horizontally

## Cleanup Test Data

After testing, remove test users:

### MongoDB

```javascript
use horo_users;
db.users.deleteMany({ email: /^loadtest_/ });
db.chatrooms.deleteMany({ "members": /^loadtest_/ });
```

### PostgreSQL

```sql
DELETE FROM payments WHERE order_id IN (
  SELECT order_id FROM orders WHERE user_id LIKE 'loadtest_%'
);
DELETE FROM orders WHERE user_id LIKE 'loadtest_%';
```

## Security Notes

⚠️ **Important**:

- Use test Firebase project (not production)
- Use separate test database
- Never commit credentials to git
- Be aware of Firebase free tier limits (100 signups/hour)

## Files

- **`full-flow-load-test.js`**: Main load test (100 users, 5 steps)
- **`firebase-auth.js`**: Firebase authentication helper
- **`README.md`**: This documentation
