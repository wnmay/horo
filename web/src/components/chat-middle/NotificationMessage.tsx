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
      icon = "💰";
      title = "Payment Started";
      desc = msg.content;
      color = "blue";
      break;
    }
    case Trigger.OrderPaid: {
      icon = "✅";
      title = "Payment Successful";
      desc = msg.content;
      color = "green";
      break;
    }
    case Trigger.OrderCompleted: {
      icon = "🎉";
      title = "Order Completed";
      desc = msg.content;
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