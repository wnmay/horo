import amqp from "amqplib";

// === CONFIG ===
const RABBIT_URL = process.env.RABBITMQ_URI || "amqp://guest:guest@localhost:5672/";
const OUT_QUEUE = process.env.OUT_Q || "chat.message.outgoing";
const IN_QUEUE = process.env.IN_Q || "chat.message.incoming"; // ถ้ายังไม่มี routing ให้ลองตั้ง IN_Q = OUT_Q
const TEST_MESSAGE = "Latency test message";

async function runLatencyTest() {
  console.log("🚀 Starting RabbitMQ latency test (Node.js)");

  // 1 Connect to RabbitMQ
  const conn = await amqp.connect(RABBIT_URL);
  const ch = await conn.createChannel();

  await ch.assertQueue(OUT_QUEUE, { durable: true });
  await ch.assertQueue(IN_QUEUE, { durable: true });

  // 2 Setup consumer to measure delivery time
  const gotMessage = new Promise((resolve, reject) => {
    const timeout = setTimeout(() => reject(new Error("Timeout: no message received")), 5000);

    ch.consume(
      IN_QUEUE,
      (msg) => {
        const now = Date.now();
        const data = JSON.parse(msg.content.toString());
        const sentAt = data.sentAt || now;
        const latency = now - sentAt;

        console.log(`💬 Message received from ${IN_QUEUE}`);
        console.log(`⏱ Latency: ${latency} ms`);

        if (latency <= 1000) {
          console.log("✅ PASSED: Delivered within 1 second");
        } else {
          console.log("❌ FAILED: Took longer than 1 second");
        }

        ch.ack(msg);
        clearTimeout(timeout);
        resolve();
      },
      { noAck: false }
    );
  });

  // 3 Publish a test message
  const start = Date.now();
  const message = {
    messageId: `test-${start}`,
    roomId: "room-latency",
    senderId: "tester",
    content: TEST_MESSAGE,
    sentAt: start,
  };

  ch.sendToQueue(OUT_QUEUE, Buffer.from(JSON.stringify(message)), {
    contentType: "application/json",
  });

  console.log(`📤 Published message to ${OUT_QUEUE} at ${start}`);

  // 4 Wait for message or timeout
  try {
    await gotMessage;
  } catch (err) {
    console.error("❌ No message received:", err.message);
    console.log("💡 If no routing exists from OUT → IN, set IN_Q=OUT_Q to test broker latency only.");
  }

  await ch.close();
  await conn.close();
  console.log("🧾 Test complete.");
}

runLatencyTest().catch((err) => {
  console.error("🔥 Error:", err);
});
