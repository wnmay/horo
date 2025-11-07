"use client";

import { use, useEffect, useState } from "react";
import { useRouter } from "next/navigation";

interface Order {
  id: string;
  customerName: string;
  status: string;
  date: string;
}

interface Course {
  id: string;
  title: string;
  description: string;
  prophet: string;
}

export default function CoursePage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  // ✅ unwrap params once using React’s `use`
  const { id } = use(params);
  const router = useRouter();

  const [course, setCourse] = useState<Course | null>(null);
  const [orders, setOrders] = useState<Order[]>([]);
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    if (loaded) return; // prevent duplicate calls (React Strict Mode)
    console.log("🟢 Fetching mock course + orders once...");

    // 🧩 mock course data
    const mockCourse: Course = {
      id,
      title: `Course ${id}`,
      description: "This is a mock course description for demonstration.",
      prophet: "Prophet Dev1",
    };

    // 🧩 mock order data
    const mockOrders: Order[] = [
      {
        id: "order1",
        customerName: "Alice",
        status: "PENDING",
        date: "2025-11-01",
      },
      {
        id: "order2",
        customerName: "Bob",
        status: "DONE",
        date: "2025-11-03",
      },
      {
        id: "order3",
        customerName: "Charlie",
        status: "CANCELLED",
        date: "2025-11-05",
      },
    ];

    // ✅ simulate async delay
    setTimeout(() => {
      setCourse(mockCourse);
      setOrders(mockOrders);
      setLoaded(true);
    }, 500);
  }, [id, loaded]);

  if (!course) {
    return (
      <div className="p-6 text-center text-gray-600">
        Loading course information...
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* ✅ Course Info */}
      <div className="border rounded-xl p-6 bg-white shadow">
        <h1 className="text-2xl font-bold mb-2">{course.title}</h1>
        <p className="text-gray-600 mb-3">{course.description}</p>
        <p className="text-sm text-gray-500">
          Taught by <span className="font-medium">{course.prophet}</span>
        </p>
      </div>

      {/* ✅ Orders List */}
      <div className="border rounded-xl p-6 bg-white shadow">
        <h2 className="text-xl font-semibold mb-4">Orders</h2>

        {orders.length === 0 ? (
          <p className="text-gray-500 text-sm">No orders yet.</p>
        ) : (
          <ul className="space-y-3">
            {orders.map((order) => (
              <li
                key={order.id}
                onClick={() =>
                  router.push(`/prophet/chat/${order.customerName}`)
                }
                className="border p-3 rounded-lg flex justify-between items-center hover:bg-gray-50 transition cursor-pointer"
              >
                <div>
                  <p className="font-medium text-blue-600 hover:underline">
                    {order.customerName}
                  </p>
                  <p className="text-xs text-gray-500">{order.date}</p>
                </div>
                <span
                  className={`px-3 py-1 text-sm rounded-full ${
                    order.status === "DONE"
                      ? "bg-green-100 text-green-700"
                      : order.status === "PENDING"
                      ? "bg-yellow-100 text-yellow-700"
                      : "bg-red-100 text-red-700"
                  }`}
                >
                  {order.status}
                </span>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
