"use client";

import { useEffect, useState } from "react";
import { Trigger } from "@/types/contracts";
import {
  ChatMessage,
  OrderCompletedNotification,
  OrderPaidNotification,
  OrderPaymentBoundNotification,
} from "@/types/ws-message";

export function NotificationMessage({ msg }: { msg: ChatMessage }) {
  const [time, setTime] = useState("");

  useEffect(() => {
    const d = msg.createdAt ? new Date(msg.createdAt) : new Date();
    setTime(d.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", hour12: false }));
  }, [msg.createdAt]);

  if (msg.type !== "notification") return null;

  let icon = "ℹ️";
  let title = "Notification";
  let desc = "";
  let color = "gray";

  switch (msg.trigger) {
    case Trigger.OrderPaymentBound: {
      const data = (msg as OrderPaymentBoundNotification).messageDetail;
      icon = "💰";
      title = "Payment Started";
      desc = `Customer ${data?.customerId ?? "Unknown"} started payment for ${data?.courseName ?? "unknown"}.`;
      color = "blue";
      break;
    }
    case Trigger.OrderPaid: {
      const data = (msg as OrderPaidNotification).messageDetail;
      icon = "✅";
      title = "Payment Successful";
      desc = `Order ${data?.orderId ?? "Unknown"} (${data?.courseName ?? "Unknown"}) has been paid successfully.`;
      color = "green";
      break;
    }
    case Trigger.OrderCompleted: {
      const data = (msg as OrderCompletedNotification).messageDetail;
      icon = "🎉";
      title = "Order Completed";
      desc = `Course "${data?.courseName?? "Unknown"}" has been marked completed.`;
      color = "purple";
      break;
    }
  }

  const bg =
    color === "green"
      ? "bg-green-50 border-green-300"
      : color === "blue"
      ? "bg-blue-50 border-blue-300"
      : color === "purple"
      ? "bg-purple-50 border-purple-300"
      : "bg-gray-50 border-gray-300";

  return (
    <div className="flex flex-col items-center gap-2 p-3 w-full">
      <p className="text-gray-400 text-sm">{time}</p>
      <div className={`flex items-start gap-3 p-3 border rounded-lg w-full ${bg}`}>
        <span className="text-2xl">{icon}</span>
        <div>
          <p className="font-semibold">{title}</p>
          <p className="text-sm text-gray-600">{desc}</p>
        </div>
      </div>
    </div>
  );
}