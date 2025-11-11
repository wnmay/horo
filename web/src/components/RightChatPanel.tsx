"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useWebSocket } from "@/lib/ws/useWebSocket";
import api from "@/lib/api/api-client";
import { ChatMessage } from "@/types/ws-message";
import { Trigger } from "@/types/contracts";
import { auth } from "@/firebase/firebase";
import ReviewBox from "./ReviewBox";

type Role = "customer" | "prophet";

type OrderStatus =
  | "PENDING"
  | "CONFIRMED"
  | "PROPHET_DONE"
  | "CUSTOMER_DONE"
  | "COMPLETED";

type OrderSummary = {
  order_id: string;
  room_id: string;
  course_id: string;
  customer_id: string;
  is_customer_completed: boolean;
  is_prophet_completed: boolean;
  order_date: string;
  status: OrderStatus;
  amount: number;
};

type FetchOk =
  | { data: []; message: string }
  | { data: OrderSummary[]; message: string };

function isEmptyPayload(x: FetchOk): x is { data: []; message: string } {
  return Array.isArray(x.data) && x.data.length === 0;
}

export default function RightPanel({
  roomId,
  role,
  courseId,
}: {
  roomId: string;
  role: Role;
  courseId: string;
}) {
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [order, setOrder] = useState<OrderSummary | null>(null);
  const [paying, setPaying] = useState(false);
  const [prophetDone, setProphetDone] = useState(false);
  const [customerDone, setCutomerDone] = useState(false);

  const refreshOrder = useCallback(async () => {
    try {
      const res = await api.get(`/api/orders/room/${roomId}`);
      const json: FetchOk = res.data;
      setOrder(isEmptyPayload(json) ? null : json.data[0]);
    } catch (e) {}
  }, [roomId]);

  useEffect(() => {
    let alive = true;
    if (!auth.currentUser) return;
    (async () => {
      setLoading(true);
      setError(null);
      try {
        const res = await api.get(`/api/orders/room/${roomId}`);
        const json: FetchOk = res.data;
        if (!alive) return;
        setOrder(isEmptyPayload(json) ? null : json.data[0]);
      } catch (e: any) {
        setError(e?.message ?? "Failed to fetch");
      } finally {
        if (alive) setLoading(false);
      }
    })();
    return () => {
      alive = false;
    };
  }, [roomId, auth.currentUser]);

  const { connected, messages } = useWebSocket();

  useEffect(() => {
    if (connected) {
      refreshOrder();
    }
  }, [connected, refreshOrder]);

  useEffect(() => {
    if (!messages?.length) return;
    const last = messages[messages.length - 1] as ChatMessage;
    if (last.type !== "notification") return;
    if (last.roomId !== roomId) return;

    switch (last.trigger) {
      case Trigger.OrderPaymentBound: {
        const amt = last.messageDetail.amount;
        if (typeof amt === "number") {
          setOrder((prev) => (prev ? { ...prev, amount: amt } : prev));
        }
        refreshOrder();
        break;
      }
      case Trigger.OrderPaid: {
        const amt = last.messageDetail.amount;
        if (typeof amt === "number") {
          setOrder((prev) => (prev ? { ...prev, amount: amt } : prev));
        }
        refreshOrder();
        break;
      }
      case Trigger.OrderCompleted: {
        refreshOrder();
        break;
      }
      default:
        break;
    }
  }, [messages, roomId, refreshOrder]);

  const canCreate = useMemo(() => role === "customer" && !order, [role, order]);

  const onCreateOrder = useCallback(async () => {
    try {
      setCreating(true);
      setError(null);
      const res = await api.post(`/api/orders`, { courseId, roomId });
      const created: { data: OrderSummary } = res.data;
      setOrder(created.data);
    } catch (e: any) {
      setError(e?.message ?? "Failed to create order");
    } finally {
      setCreating(false);
    }
  }, [courseId, roomId]);

  const onPay = useCallback(async () => {
    if (!order) return;
    try {
      setPaying(true);
      setError(null);
      const { data } = await api.get(`/api/payments/order/${order.order_id}`);
      const paymentId = data.data.payment_id;
      await api.put(`/api/payments/${paymentId}/complete`);
    } catch (e: any) {
      setError(e?.message ?? "Payment failed");
    } finally {
      setPaying(false);
    }
  }, [order]);

  const markCustomerCompleted = useCallback(async () => {
    if (!order) return;
    setError(null);
    try {
      const res = await api.patch(`/api/orders/customer/${order.order_id}`);      
    } catch (e: any) {
      setError(e?.message ?? "Failed to mark customer completed");
    } finally {
      setCutomerDone(true);
    }
  }, [order]);

  const markProphetCompleted = useCallback(async () => {
    if (!order) return;
    setError(null);
    try {
      const res = await api.patch(`/api/orders/prophet/${order.order_id}`);
    } catch (e: any) {
      setError(e?.message ?? "Failed to mark prophet completed");
    } finally {
      setProphetDone(true);
    }
  }, [order]);

  const onProphetDone = markProphetCompleted;
  const onUserDone = markCustomerCompleted;

  const onWriteReview = () => console.log("write");
  const StatusBadge = ({ status }: { status: OrderStatus }) => {
    const map: Record<OrderStatus, string> = {
      PENDING: "bg-yellow-100 text-yellow-700",
      CONFIRMED: "bg-blue-100 text-blue-700",
      PROPHET_DONE: "bg-indigo-100 text-indigo-700",
      CUSTOMER_DONE: "bg-indigo-100 text-indigo-700",
      COMPLETED: "bg-green-100 text-green-700",
    };
    return (
      <span
        className={`px-2.5 py-1 rounded-full text-xs font-semibold ${map[status]}`}
      >
        {status.replaceAll("_", " ")}
      </span>
    );
  };

  const ActionRow = () => {
    if (!order) return null;
    const s = order.status;

    if (s === "PENDING" && role === "customer") {
      return (
        <button
          onClick={onPay}
          disabled={paying}
          className="w-full rounded-xl px-4 py-3 font-semibold bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-60"
        >
          {paying ? "Processing..." : "Pay order"}
        </button>
      );
    }

    if (s === "PENDING" && role === "prophet") {
      return (
        <div>
          <p className="mt-2 text-xs text-gray-500 text-center">
            Waiting for customer to make payment...
          </p>
        </div>
      );
    }

    if (s === "CONFIRMED") {
      if (role === "prophet") {
        return (
          <div>
            <button
              onClick={onProphetDone}
              disabled={order.is_prophet_completed || prophetDone}
              className={`w-full rounded-xl px-4 py-3 font-semibold text-white ${
                order.is_prophet_completed
                  ? "bg-gray-400 cursor-not-allowed"
                  : "bg-indigo-600 hover:bg-indigo-700"
              }`}
            >
              Done Course
            </button>
            {order.is_prophet_completed && !order.is_customer_completed && (
              <p className="mt-2 text-xs text-gray-500 text-center">
                Waiting for customer to done course...
              </p>
            )}
          </div>
        );
      }

      if (role === "customer") {
        return (
          <div>
            <button
              onClick={onUserDone}
              disabled={order.is_customer_completed || customerDone}
              className={`w-full rounded-xl px-4 py-3 font-semibold text-white ${
                order.is_customer_completed
                  ? "bg-gray-400 cursor-not-allowed"
                  : "bg-emerald-600 hover:bg-emerald-700"
              }`}
            >
              Done Course
            </button>
            {order.is_customer_completed && !order.is_prophet_completed && (
              <p className="mt-2 text-xs text-gray-500 text-center">
                Waiting for prophet to done course...
              </p>
            )}
          </div>
        );
      }
    }

    if (s === "COMPLETED" && role === "customer") {
      return <ReviewBox courseId={courseId} onSubmitted={refreshOrder} />;
    }

    if (s === "PROPHET_DONE" && role === "customer") {
      return (
        <button
          onClick={onUserDone}
          className="w-full rounded-xl px-4 py-3 font-semibold bg-emerald-600 text-white hover:bg-emerald-700"
        >
          Done Course
        </button>
      );
    }

    if (s === "CUSTOMER_DONE" && role === "prophet") {
      return (
        <button
          onClick={onProphetDone}
          className="w-full rounded-xl px-4 py-3 font-semibold bg-emerald-600 text-white hover:bg-emerald-700"
        >
          Done Course
        </button>
      );
    }

    return null;
  };

  return (
    <aside className="w-full lg:w-80 xl:w-96 shrink-0 border-l border-gray-200 bg-white">
      <div className="p-4 space-y-4">
        <div className="flex items-center justify-center">
          <h3 className="text-sm font-semibold text-gray-700 text-center">
            Order Status
          </h3>
        </div>

        {loading && (
          <div className="animate-pulse h-24 rounded-xl bg-gray-100" />
        )}
        {error && (
          <div className="rounded-lg border border-rose-200 bg-rose-50 p-3 text-rose-700 text-sm">
            {error}
          </div>
        )}

        {!loading && !order && (
          <div className="rounded-2xl border border-dashed border-gray-300 p-5 bg-gray-50">
            <p className="text-sm text-gray-600 mb-3">
              No order in this room yet.
            </p>
            {canCreate ? (
              <button
                onClick={onCreateOrder}
                disabled={creating}
                className="w-full rounded-xl px-4 py-3 font-semibold bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-60"
              >
                {creating ? "Creating..." : "Create new order"}
              </button>
            ) : (
              <p className="text-xs text-gray-500">
                Waiting for customer to create an order.
              </p>
            )}
          </div>
        )}

        {order && (
          <div className="rounded-2xl border border-gray-200 p-5">
            <div className="flex items-center justify-between">
              <span className="text-sm font-semibold text-gray-800">
                #{order.order_id.slice(0, 8)}
              </span>
              <StatusBadge status={order.status} />
            </div>

            <div className="mt-3 space-y-1 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-500">Ordered</span>
                <span className="text-gray-800">
                  {new Date(order.order_date).toLocaleString()}
                </span>
              </div>
            </div>

            <div className="mt-4">
              <ActionRow />
            </div>
          </div>
        )}
      </div>
    </aside>
  );
}
