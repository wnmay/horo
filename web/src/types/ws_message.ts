import { Trigger } from "@/types/contracts";

interface BaseMessage {
  messageId: string;
  roomId: string;
  senderId: string;
  type: "text" | "notification";
  createdAt: string;
}

/** Text message */
export interface ChatTextMessage extends BaseMessage {
  type: "text";
  content: string;
}

/** Generic notification message shape */
interface BaseNotification<TTrigger extends Trigger, TDetail> extends BaseMessage {
  type: "notification";
  trigger: TTrigger;
  messageDetail: TDetail;
}

/* ===============================
   Notification Detail Interfaces
   =============================== */

export interface OrderCompletedNotificationData {
  orderId: string;
  courseId: string;
  orderStatus: string;
  courseName: string;
}

export interface OrderPaymentBoundNotificationData {
  orderId: string;
  paymentId: string;
  roomId: string;
  customerId: string;
  courseId: string;
  orderStatus: string;
  courseName: string;
  amount: number;
  paymentStatus: string;
}

export interface OrderPaidNotificationData {
  orderId: string;
  paymentId: string;
  roomId: string;
  customerId: string;
  courseId: string;
  orderStatus: string;
  courseName: string;
  amount: number;
  paymentStatus: string;
}

export type OrderCompletedNotification = BaseNotification<
  typeof Trigger.OrderCompleted,
  OrderCompletedNotificationData
>;

export type OrderPaymentBoundNotification = BaseNotification<
  typeof Trigger.OrderPaymentBound,
  OrderPaymentBoundNotificationData
>;

export type OrderPaidNotification = BaseNotification<
  typeof Trigger.OrderPaid,
  OrderPaidNotificationData
>;


export type ChatMessage =
  | ChatTextMessage
  | OrderCompletedNotification
  | OrderPaymentBoundNotification
  | OrderPaidNotification;
