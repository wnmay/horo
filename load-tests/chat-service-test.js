import http from "k6/http";
import { check, sleep } from "k6";
import { Trend } from "k6/metrics";
import { getFirebaseToken } from "./firebase-auth.js"; 

export const message_delivery_time = new Trend("message_delivery_time");

export const options = {
  vus: 5,
  iterations: 20,
  thresholds: {
    "http_req_failed": ["rate<0.01"],
    "http_req_duration": ["p(95)<2000"],
    "message_delivery_time": ["p(95)<1000"],
  },
};

const BASE_URL = __ENV.BASE_URL || "http://localhost:3000/api";

// === setup() จะรันก่อน VU ทุกตัว ===
export function setup() {
  const auth = getFirebaseToken();
  return { token: auth.idToken };
}

export default function (data) {
  const start = Date.now();

  const params = {
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${data.token}`, 
    },
  };

  // === 1 Create Room ===
  const createPayload = JSON.stringify({
    courseID: "course-999",
  });
  const resCreate = http.post(`${BASE_URL}/chats/rooms`, createPayload, params);
  check(resCreate, {
    "create room success": (r) => r.status === 201,
  });

  const body = JSON.parse(resCreate.body || "{}");
  const roomId = body.roomID || "test-room-123";

  // === 2 Get Messages by Room ===
  const resMessages = http.get(`${BASE_URL}/chats/${roomId}/messages`, params);
  check(resMessages, {
    "get messages success": (r) => r.status === 200,
  });

  // === 3️⃣ Record latency ===
  const end = Date.now();
  const total = end - start;
  message_delivery_time.add(total);
  console.log(`💬 Chat round trip took ${total} ms`);

  sleep(1);
}
