import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Counter } from 'k6/metrics';

// ====== CONFIGURATION ======
export const options = {
  stages: [
    { duration: '10s', target: 100 },  // ramp-up to 100 users
    { duration: '20s', target: 500 },   // spike to 500 users
    { duration: '30s', target: 1000 },   // spike to 1000 users
    { duration: '10s', target: 0 },   // ramp-down
  ],
    thresholds: {
        http_req_duration: ['p(95)<10000'], // 95% requests < 10s
        http_req_failed: ['rate<0.01'],
    },

};

// ====== TOKENS ======
const CUSTOMER_TOKEN = 'eyJhbGciOiJSUzI1NiIsImtpZCI6IjM4MDI5MzRmZTBlZWM0NmE1ZWQwMDA2ZDE0YTFiYWIwMWUzNDUwODMiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoidGVzdCIsInJvbGUiOiJjdXN0b21lciIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS9ob3JvLWQ0N2IxIiwiYXVkIjoiaG9yby1kNDdiMSIsImF1dGhfdGltZSI6MTc2Mjg4NjM0MywidXNlcl9pZCI6Inl4MmRid3RPVFdRUnQxTzBaSXRVOUVaSUtNSjMiLCJzdWIiOiJ5eDJkYnd0T1RXUVJ0MU8wWkl0VTlFWklLTUozIiwiaWF0IjoxNzYyODg2MzQzLCJleHAiOjE3NjI4ODk5NDMsImVtYWlsIjoiY3VzdG9tZXJ0ZXN0QGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJlbWFpbCI6WyJjdXN0b21lcnRlc3RAZ21haWwuY29tIl19LCJzaWduX2luX3Byb3ZpZGVyIjoicGFzc3dvcmQifX0.NubaLv2DgifKXoGqA66rhOQ6SErSQ-3qtOu5S1lWN0ZhAn3CmzDCnUs6F95K2iDCraTe6oXlKk1x-07wlr4zlFPuP0N4wfcVs00Nv2KC6s94FzikdGdJkwvmVQy9-zWtJya_ShfQXqnWIR5dS1CwtaYf1Lc3ONOc_jrkx3EVUzhrJtWUHbY1tPCVL8lk_WzxLT7a3A88rWWyQCDyx1p_LsWRguHR6FasQjUEaBXxIiRq4pW1ayg6WBO5eL3vRsK2F0_WBCABkJcyQ4Rx47ZgK6Xj5HeLNdkCM4xjG6UDorUU8vMfcwviQRYgca3cDC1iai4SC35_O__pdJlzEhRiCg';
const PROPHET_TOKEN = 'eyJhbGciOiJSUzI1NiIsImtpZCI6IjM4MDI5MzRmZTBlZWM0NmE1ZWQwMDA2ZDE0YTFiYWIwMWUzNDUwODMiLCJ0eXAiOiJKV1QifQ.eyJyb2xlIjoicHJvcGhldCIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS9ob3JvLWQ0N2IxIiwiYXVkIjoiaG9yby1kNDdiMSIsImF1dGhfdGltZSI6MTc2Mjg4NjQwMiwidXNlcl9pZCI6InNtOFNhWmw1UHhPSlNhcHV3bFM1ZmNlNEtTNTMiLCJzdWIiOiJzbThTYVpsNVB4T0pTYXB1d2xTNWZjZTRLUzUzIiwiaWF0IjoxNzYyODg2NDAyLCJleHAiOjE3NjI4OTAwMDIsImVtYWlsIjoicHJvcGhldHRlc3RAZ21haWwuY29tIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7ImVtYWlsIjpbInByb3BoZXR0ZXN0QGdtYWlsLmNvbSJdfSwic2lnbl9pbl9wcm92aWRlciI6InBhc3N3b3JkIn19.okcFNQrd65pk-VQdsrAHJhPRy0O_mL0JJoKg-iWS_D_5wvT-zsfLN-P8_yLInsWyI5ztUluk_ZBH2OEoJ-vd-f2ILsXwD0T-o4norvu_JtrHSMlzKlMpLyYznN7CHA8zJSCQPev1j1zPjIP89DUl1ngZiQs0D36z2Q_HnMDwnVq8INLDfYpSoXBCqbT0nclvtOntfsLYr2K_khru7fIoik6OiLHnnx4j3cSTTBa2MXLsfX-21A02AmWOUH4IRYW7xusoXyjcMN3VhAmqXIIgpya5nsSKwbpglmGRuTCkTnw2SMqxLtHwKJkrMZIrCbP6qrPI-4NglewdiA824Jenyg';

// ====== TARGET BASE URL ======
const BASE_URL = 'http://192.168.0.97:8080';


// ====== METRICS ======
let registerTrend = new Trend('register_duration');
let courseTrend = new Trend('get_courses_duration');
let orderTrend = new Trend('create_order_duration');
let successCount = new Counter('successful_requests');
let failCount = new Counter('failed_requests');

// ====== TEST LOGIC ======
export default function () {
  // Randomly act as customer or prophet
  const isProphet = Math.random() < 0.3;
  const token = isProphet ? PROPHET_TOKEN : CUSTOMER_TOKEN;
  const role = isProphet ? 'prophet' : 'customer';

  // GET courses
  let courseRes = http.get(`${BASE_URL}/api/courses`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  check(courseRes, { 'get courses status 200': (r) => r.status === 200 });
  courseTrend.add(courseRes.timings.duration);

  // CREATE order (customer)
  if (!isProphet) {
    const payload = JSON.stringify({
      courseId: 'test-course-id',
      roomId: 'test-room-id',
    });
    let orderRes = http.post(`${BASE_URL}/api/orders`, payload, {
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
    });
    check(orderRes, {
      'order create ok': (r) => r.status === 201 || r.status === 200,
    })
      ? successCount.add(1)
      : failCount.add(1);
    orderTrend.add(orderRes.timings.duration);
  }

  // GET prophet balance
  if (isProphet) {
    let balRes = http.get(`${BASE_URL}/api/payments/balance`, {
      headers: {
        Authorization: `Bearer ${token}`,
        'X-User-Id': 'sm8SaZl5PxOJSapuwlS5fce4KS53',
        'X-User-Role': 'prophet',
      },
    });
    check(balRes, { 'get balance ok': (r) => r.status === 200 });
  }

  sleep(1); // wait between actions
}
