import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { getFirebaseToken, refreshFirebaseToken } from './firebase-auth.js';

export const options = {
  stages: [
    { duration: '1m', target: 50 },
    { duration: '2m', target: 100 },
    { duration: '3m', target: 100 },  // Main test period
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    'http_req_duration': ['p(95)<3000'],
    'http_req_duration{operation:create_order}': ['p(95)<3000'],
    'http_req_duration{operation:get_order}': ['p(95)<2000'],
    'http_req_failed': ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000/api';
const COURSE_ID = __ENV.COURSE_ID || '126e7890-e89b-12d3-a456-426614174001';

// Global token management
let authData = null;
let tokenExpireTime = 0;

// Setup function - runs once per VU
export function setup() {
  // Get fresh token at the start
  console.log('🚀 Setting up test - getting fresh Firebase token...');
  const auth = getFirebaseToken();
  
  return {
    idToken: auth.idToken,
    refreshToken: auth.refreshToken,
    expiresIn: auth.expiresIn,
    startTime: Date.now(),
  };
}

export default function (data) {
  // Check if token needs refresh (refresh 5 minutes before expiry)
  const currentTime = Date.now();
  const tokenAge = (currentTime - data.startTime) / 1000; // seconds
  const timeUntilExpiry = data.expiresIn - tokenAge;

  let token = data.idToken;

  // Refresh token if it's about to expire (less than 5 minutes left)
  if (timeUntilExpiry < 300 && __VU === 1 && __ITER % 10 === 0) {
    console.log(`⚠️  Token expiring soon (${Math.floor(timeUntilExpiry / 60)} min left), refreshing...`);
    try {
      const refreshed = refreshFirebaseToken(data.refreshToken);
      token = refreshed.idToken;
      data.idToken = refreshed.idToken;
      data.refreshToken = refreshed.refreshToken;
      data.expiresIn = refreshed.expiresIn;
      data.startTime = Date.now();
    } catch (e) {
      console.error('Failed to refresh token:', e.message);
    }
  }

  const params = {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  };

  let orderId;

  // Scenario 1: User creates an order
  group('Create Order', () => {
    const payload = JSON.stringify({
      courseId: COURSE_ID,
    });
    
    const res = http.post(`${BASE_URL}/orders`, payload, {
      ...params,
      tags: { operation: 'create_order' },
    });
    
    const success = check(res, {
      'order created': (r) => r.status === 201,
      'has order_id': (r) => {
        try {
          const body = JSON.parse(r.body);
          orderId = body.data?.order_id;
          return orderId !== undefined;
        } catch (e) {
          return false;
        }
      },
      'create time < 3s': (r) => r.timings.duration < 3000,
    });

    if (!success) {
      console.error(`Create failed: ${res.status} - ${res.body}`);
      return;
    }
  });

  sleep(2); // User thinks/views order

  // Scenario 2: User views their order
  if (orderId) {
    group('Get Order Details', () => {
      const res = http.get(`${BASE_URL}/orders/${orderId}`, {
        ...params,
        tags: { operation: 'get_order' },
      });
      
      check(res, {
        'order retrieved': (r) => r.status === 200,
        'correct order returned': (r) => {
          try {
            const body = JSON.parse(r.body);
            return body.data?.order_id === orderId;
          } catch (e) {
            return false;
          }
        },
        'get time < 2s': (r) => r.timings.duration < 2000,
      });
    });
  }

  sleep(1);

  // Scenario 3: User views all their orders
  group('Get All Orders', () => {
    const res = http.get(`${BASE_URL}/orders`, params);
    
    check(res, {
      'orders list retrieved': (r) => r.status === 200,
      'list time < 3s': (r) => r.timings.duration < 3000,
    });
  });

  sleep(2); // User browsing time
}