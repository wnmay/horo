# Load Testing with K6

## Quick Start

### Option 1: Use .env File (Recommended) 🎯

1. The `.env` file is already configured with Firebase API key
2. Just run:

```powershell
.\load-tests\run-load-test.ps1
```

That's it! The script loads environment variables from `.env` and runs the test.

### Option 2: Manual Setup

```powershell
# Set environment variables
$env:FIREBASE_API_KEY="AIzaSyBnN693MGJZ2p1yRICP6IzoOgEKj3ZBsKI"
$env:BASE_URL="http://localhost:8080"

# Run the test
k6 run load-tests/full-flow-load-test.js
```

### Step 1: Get Your Firebase Web API Key

1. Go to [Firebase Console](https://console.firebase.google.com/project/horo-d47b1/settings/general)
2. Scroll to "Your apps" section
3. Copy the **"Web API Key"** (starts with `AIzaSy...`)

### Step 2: Run the Test

```powershell
# Set the API key
$env:FIREBASE_API_KEY="AIzaSyXXXXXXXXXXXXXXXXXXXXXXXXXX"

# Run the test (100 users, 2 minutes)
k6 run load-tests/full-flow-load-test.js

# Or test with fewer users first
$env:FIREBASE_API_KEY="AIzaSy..."; k6 run --vus 10 --duration 30s load-tests/full-flow-load-test.js
```

### Quick One-Liner

````powershell
$env:FIREBASE_API_KEY="AIzaSy..."; k6 run load-tests/full-flow-load-test.js
```## What the Test Does

Tests complete user journey with **100 unique users**:

1. ✅ **Register Firebase user** - Create unique Firebase account
2. ✅ **Register in backend** - Register with backend using Firebase token
3. ✅ **Get all courses** using the token
4. ✅ **Create chat room** with a course
5. ✅ **Create order** with the chat room
6. ✅ **Complete payment** for the order

**Load Pattern:**

- 0-20s: Ramp to 50 users
- 20-50s: Ramp to 100 users
- 50-90s: Hold at 100 users
- 90-120s: Ramp down to 0

**Requirements:**

- ✅ P95 response time < 3 seconds
- ✅ Error rate < 10%

## Prerequisites

1. **Install K6**: https://k6.io/docs/get-started/installation/

   ```powershell
   choco install k6
````

2. **Firebase Web API Key**:

   - Go to [Firebase Console](https://console.firebase.google.com/)
   - Select your project
   - Click ⚙️ → Project settings → General tab
   - Copy "Web API Key" (starts with `AIzaSy...`)

3. **Enable Firebase Email/Password Auth**:

   - Firebase Console → Authentication → Sign-in method
   - Enable Email/Password authentication
   - ⚠️ Free tier: 100 signups/hour

4. **Start all services**: API Gateway (8080), User Management, Order, Payment, Chat, Course

5. **Test data**: Ensure at least one course exists in the database

## How It Works

Each Virtual User (VU):

- Generates unique credentials: `loadtest_<VU>_<timestamp>_<random>@example.com`
- **Step 1**: Registers Firebase account directly (gets `idToken`)
- **Step 2**: Calls `/api/users/register` with `idToken`, `fullName`, `role`
- Uses the token for all subsequent authenticated requests
- Tests **100 concurrent different users** (realistic scenario!)

## Understanding Results

Example output:

```
✓ registration status is 201
✓ registration returns token
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

## Troubleshooting

### ❌ "API key not valid. Please pass a valid API key"

**Most common error!**

1. **Get the correct key**:

   - Go to Firebase Console → ⚙️ Project settings → General
   - Find **"Web API Key"** (format: `AIzaSy...`)
   - NOT "Server key" or Service Account key

2. **Set it correctly**:
   ```powershell
   $env:FIREBASE_API_KEY="AIzaSy..."
   echo $env:FIREBASE_API_KEY  # Verify it's set
   k6 run load-tests/full-flow-load-test.js
   ```

### ❌ "FIREBASE_API_KEY environment variable not set"

```powershell
# Set the variable before running
$env:FIREBASE_API_KEY="your-actual-key"
k6 run load-tests/full-flow-load-test.js
```

### ❌ Registration fails with 500 error

- Check user-management service logs
- Verify Firebase Admin SDK is configured correctly
- Ensure firebase-key.json is loaded in backend

### ❌ "courses list is empty"

- Create at least one course in the database
- Verify course-service is running

### ❌ High error rate (>10%)

- Check service logs for errors
- Verify database connection pool size
- Reduce concurrent users (use `--vus 10`)

### ❌ Slow responses (P95 >3s)

- Add database indexes
- Increase service resources (CPU/Memory)
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

## Files

- **`full-flow-load-test.js`**: Main load test (100 users, 5 steps)
- **`firebase-auth.js`**: Firebase helper (not used anymore, kept for reference)
- **`run-load-test.ps1`**: Interactive script (optional)
- **`README.md`**: This documentation
