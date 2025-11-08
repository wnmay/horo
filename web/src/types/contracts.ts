// src/types/contract.ts

// Import from contract/amqp.go
export const Trigger = {
  OrderCreated: "order.created",
  OrderCompleted: "order.completed",
  OrderPaymentBound: "order.payment.bound",
  OrderPaid: "order.paid",
  PaymentSuccess: "payment.completed",
  PaymentCreated: "payment.created",
  PaymentSettled: "payment.settled",
  ChatMessageIncoming: "chat.message.incoming",
  ChatMessageOutgoing: "chat.message.outgoing",
} as const;

export type Trigger = (typeof Trigger)[keyof typeof Trigger];
