import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { registerFirebaseUser } from './firebase-auth.js';

// Custom metrics
const errorRate = new Rate('errors');
const firebaseRegistrationDuration = new Trend('firebase_registration_duration');
const registrationDuration = new Trend('registration_duration');
const coursesFetchDuration = new Trend('courses_fetch_duration');
const chatRoomCreationDuration = new Trend('chat_room_creation_duration');
const orderCreationDuration = new Trend('order_creation_duration');
const paymentCompletionDuration = new Trend('payment_completion_duration');

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const FIREBASE_API_KEY = __ENV.FIREBASE_API_KEY;

export const options = {
  stages: [
    { duration: '20s', target: 50 },   // Ramp up to 50 users
    { duration: '30s', target: 100 },  // Ramp up to 100 users
    { duration: '40s', target: 100 },  // Stay at 100 users
    { duration: '30s', target: 0 },    // Ramp down to 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<3000'], // 95% of requests must complete within 3s
    errors: ['rate<0.1'],              // Error rate must be below 10%
    'http_req_duration{name:registration}': ['p(95)<3000'],
    'http_req_duration{name:get_courses}': ['p(95)<3000'],
    'http_req_duration{name:create_chat_room}': ['p(95)<3000'],
    'http_req_duration{name:create_order}': ['p(95)<3000'],
    'http_req_duration{name:complete_payment}': ['p(95)<3000'],
  },
};

// Helper function to generate random email
function generateRandomEmail() {
  const timestamp = Date.now();
  const random = Math.random().toString(36).substring(7);
  const vuId = __VU; // Virtual User ID (unique per VU)
  return `loadtest_${vuId}_${timestamp}_${random}@example.com`;
}

// Helper function to generate random password
function generateRandomPassword() {
  return `LoadTest${Math.random().toString(36).substring(2, 15)}!`;
}

export default function () {
  // Check if Firebase API Key is set
  if (!FIREBASE_API_KEY) {
    console.error('❌ FIREBASE_API_KEY environment variable not set!');
    console.error('   Set it with: $env:FIREBASE_API_KEY="your-api-key"');
    errorRate.add(1);
    return;
  }

  // Generate unique credentials for this virtual user
  const email = generateRandomEmail();
  const password = generateRandomPassword();
  const fullName = `Load Test User ${__VU}`;
  const role = Math.random() > 0.5 ? 'customer' : 'prophet';
  
  let idToken, courseId, chatRoomId, orderId, paymentId;

  // ===== STEP 1: Register Firebase user =====
  const startFirebaseReg = Date.now();
  let firebaseAuth;
  
  try {
    firebaseAuth = registerFirebaseUser(email, password);
    idToken = firebaseAuth.idToken;
    firebaseRegistrationDuration.add(Date.now() - startFirebaseReg);
  } catch (error) {
    console.error(`❌ Failed to register Firebase user ${email}: ${error}`);
    errorRate.add(1);
    return;
  }
  
  const bearerToken = `Bearer ${idToken}`;
  
  sleep(0.5);

  // ===== STEP 2: Register user in backend =====
  const registerPayload = JSON.stringify({
    idToken: idToken,
    fullName: fullName,
    role: role,
  });

  const registerParams = {
    headers: {
      'Content-Type': 'application/json',
    },
    tags: { name: 'registration' },
  };

  const registerRes = http.post(
    `${BASE_URL}/api/users/register`,
    registerPayload,
    registerParams
  );

  const registerSuccess = check(registerRes, {
    'registration status is 201': (r) => r.status === 201,
    'registration completed within 3s': (r) => r.timings.duration < 3000,
  });

  registrationDuration.add(registerRes.timings.duration);
  errorRate.add(!registerSuccess);

  if (!registerSuccess) {
    console.error(`Registration failed for ${email}: ${registerRes.status} - ${registerRes.body}`);
    return;
  }

  sleep(0.5);

  // ===== STEP 3: Get all courses =====
  const coursesParams = {
    headers: {
      'Authorization': bearerToken,
      'Content-Type': 'application/json',
    },
    tags: { name: 'get_courses' },
  };

  const coursesRes = http.get(
    `${BASE_URL}/api/courses`,
    coursesParams
  );

  const coursesSuccess = check(coursesRes, {
    'get courses status is 200': (r) => r.status === 200,
    'get courses completed within 3s': (r) => r.timings.duration < 3000,
    'courses list is not empty': (r) => {
      try {
        const body = JSON.parse(r.body);
        const courses = body.data || body;
        return Array.isArray(courses) && courses.length > 0;
      } catch (e) {
        return false;
      }
    },
  });

  coursesFetchDuration.add(coursesRes.timings.duration);
  errorRate.add(!coursesSuccess);

  if (!coursesSuccess) {
    console.error(`Get courses failed: ${coursesRes.status} - ${coursesRes.body}`);
    return;
  }

  // Extract a course ID
  try {
    const coursesBody = JSON.parse(coursesRes.body);
    const courses = coursesBody.data || coursesBody;
    if (courses && courses.length > 0) {
      courseId = courses[0].course_id || courses[0].courseId || courses[0].id || courses[0]._id;
    }
  } catch (e) {
    console.error(`Failed to parse courses response: ${e}`);
    return;
  }

  if (!courseId) {
    console.error('No course ID found');
    return;
  }

  sleep(0.5);

  // ===== STEP 3: Create chat room =====
  const chatRoomPayload = JSON.stringify({
    courseId: courseId,
    customerId: `customer_${__VU}`, // Mock customer ID
  });

  const chatRoomParams = {
    headers: {
      'Authorization': bearerToken,
      'Content-Type': 'application/json',
    },
    tags: { name: 'create_chat_room' },
  };

  const chatRoomRes = http.post(
    `${BASE_URL}/api/chat/rooms`,
    chatRoomPayload,
    chatRoomParams
  );

  const chatRoomSuccess = check(chatRoomRes, {
    'create chat room status is 200 or 201': (r) => r.status === 200 || r.status === 201,
    'create chat room completed within 3s': (r) => r.timings.duration < 3000,
  });

  chatRoomCreationDuration.add(chatRoomRes.timings.duration);
  errorRate.add(!chatRoomSuccess);

  if (!chatRoomSuccess) {
    console.error(`Create chat room failed: ${chatRoomRes.status} - ${chatRoomRes.body}`);
    return;
  }

  // Extract room ID
  try {
    const chatRoomBody = JSON.parse(chatRoomRes.body);
    chatRoomId = chatRoomBody.roomId || chatRoomBody.room_id || chatRoomBody.id || chatRoomBody.data?.id;
  } catch (e) {
    console.error(`Failed to parse chat room response: ${e}`);
    return;
  }

  if (!chatRoomId) {
    console.error('No chat room ID found');
    return;
  }

  sleep(0.5);

  // ===== STEP 4: Create order =====
  const orderPayload = JSON.stringify({
    courseId: courseId,
    roomId: chatRoomId,
  });

  const orderParams = {
    headers: {
      'Authorization': bearerToken,
      'Content-Type': 'application/json',
    },
    tags: { name: 'create_order' },
  };

  const orderRes = http.post(
    `${BASE_URL}/api/orders`,
    orderPayload,
    orderParams
  );

  const orderSuccess = check(orderRes, {
    'create order status is 200 or 201': (r) => r.status === 200 || r.status === 201,
    'create order completed within 3s': (r) => r.timings.duration < 3000,
  });

  orderCreationDuration.add(orderRes.timings.duration);
  errorRate.add(!orderSuccess);

  if (!orderSuccess) {
    console.error(`Create order failed: ${orderRes.status} - ${orderRes.body}`);
    return;
  }

  // Extract order ID and payment ID
  try {
    const orderBody = JSON.parse(orderRes.body);
    orderId = orderBody.orderId || orderBody.order_id || orderBody.id || orderBody.data?.order_id || orderBody.data?.id;
    paymentId = orderBody.paymentId || orderBody.payment_id || orderBody.data?.payment_id;
  } catch (e) {
    console.error(`Failed to parse order response: ${e}`);
    return;
  }

  if (!orderId) {
    console.error('No order ID found');
    return;
  }

  // If payment ID not returned, try to fetch it
  if (!paymentId) {
    const paymentByOrderRes = http.get(
      `${BASE_URL}/api/payments/order/${orderId}`,
      {
        headers: {
          'Authorization': bearerToken,
          'Content-Type': 'application/json',
        },
      }
    );

    if (paymentByOrderRes.status === 200) {
      try {
        const paymentBody = JSON.parse(paymentByOrderRes.body);
        paymentId = paymentBody.id || paymentBody.payment_id || paymentBody.data?.id || paymentBody.data?.payment_id;
      } catch (e) {
        console.error(`Failed to parse payment response: ${e}`);
      }
    }
  }

  if (!paymentId) {
    console.error('No payment ID found');
    return;
  }

  sleep(0.5);

  // ===== STEP 5: Complete payment =====
  const paymentParams = {
    headers: {
      'Authorization': bearerToken,
      'Content-Type': 'application/json',
    },
    tags: { name: 'complete_payment' },
  };

  const paymentRes = http.put(
    `${BASE_URL}/api/payments/${paymentId}/complete`,
    '{}',
    paymentParams
  );

  const paymentSuccess = check(paymentRes, {
    'complete payment status is 200': (r) => r.status === 200,
    'complete payment completed within 3s': (r) => r.timings.duration < 3000,
  });

  paymentCompletionDuration.add(paymentRes.timings.duration);
  errorRate.add(!paymentSuccess);

  if (!paymentSuccess) {
    console.error(`Complete payment failed: ${paymentRes.status} - ${paymentRes.body}`);
    return;
  }

  // Small delay between iterations
  sleep(1);
}

export function handleSummary(data) {
  return {
    'stdout': textSummary(data, { indent: ' ', enableColors: true }),
    'load-test-results.json': JSON.stringify(data),
  };
}

function textSummary(data, options) {
  const indent = options?.indent || '';
  const enableColors = options?.enableColors || false;
  
  let summary = `
${indent}============================================
${indent}  LOAD TEST SUMMARY
${indent}============================================
${indent}
${indent}Duration: ${data.state.testRunDurationMs / 1000}s
${indent}VUs: ${data.metrics.vus?.values?.max || 0} max
${indent}
${indent}HTTP Requests:
${indent}  Total: ${data.metrics.http_reqs?.values?.count || 0}
${indent}  Rate: ${data.metrics.http_reqs?.values?.rate?.toFixed(2) || 0}/s
${indent}  Failed: ${data.metrics.http_req_failed?.values?.rate ? (data.metrics.http_req_failed.values.rate * 100).toFixed(2) : 0}%
${indent}
${indent}Response Times:
${indent}  Avg: ${data.metrics.http_req_duration?.values?.avg?.toFixed(2) || 0}ms
${indent}  Min: ${data.metrics.http_req_duration?.values?.min?.toFixed(2) || 0}ms
${indent}  Max: ${data.metrics.http_req_duration?.values?.max?.toFixed(2) || 0}ms
${indent}  P95: ${data.metrics.http_req_duration?.values?.['p(95)']?.toFixed(2) || 0}ms
${indent}  P99: ${data.metrics.http_req_duration?.values?.['p(99)']?.toFixed(2) || 0}ms
${indent}
${indent}Step-by-Step Performance:
${indent}  Registration:
${indent}    Avg: ${data.metrics.registration_duration?.values?.avg?.toFixed(2) || 0}ms
${indent}    P95: ${data.metrics.registration_duration?.values?.['p(95)']?.toFixed(2) || 0}ms
${indent}
${indent}  Get Courses:
${indent}    Avg: ${data.metrics.courses_fetch_duration?.values?.avg?.toFixed(2) || 0}ms
${indent}    P95: ${data.metrics.courses_fetch_duration?.values?.['p(95)']?.toFixed(2) || 0}ms
${indent}
${indent}  Create Chat Room:
${indent}    Avg: ${data.metrics.chat_room_creation_duration?.values?.avg?.toFixed(2) || 0}ms
${indent}    P95: ${data.metrics.chat_room_creation_duration?.values?.['p(95)']?.toFixed(2) || 0}ms
${indent}
${indent}  Create Order:
${indent}    Avg: ${data.metrics.order_creation_duration?.values?.avg?.toFixed(2) || 0}ms
${indent}    P95: ${data.metrics.order_creation_duration?.values?.['p(95)']?.toFixed(2) || 0}ms
${indent}
${indent}  Complete Payment:
${indent}    Avg: ${data.metrics.payment_completion_duration?.values?.avg?.toFixed(2) || 0}ms
${indent}    P95: ${data.metrics.payment_completion_duration?.values?.['p(95)']?.toFixed(2) || 0}ms
${indent}
${indent}Error Rate: ${data.metrics.errors?.values?.rate ? (data.metrics.errors.values.rate * 100).toFixed(2) : 0}%
${indent}
${indent}REQUIREMENT CHECK:
${indent}  ✓ 95% of requests < 3s: ${(data.metrics.http_req_duration?.values?.['p(95)'] || 0) < 3000 ? 'PASS' : 'FAIL'}
${indent}  ✓ Error rate < 10%: ${(data.metrics.errors?.values?.rate || 0) < 0.1 ? 'PASS' : 'FAIL'}
${indent}============================================
`;

  return summary;
}
