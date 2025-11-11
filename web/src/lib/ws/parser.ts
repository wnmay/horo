import {
  ChatMessage,
  ChatTextMessage,
  OrderCompletedNotification,
  OrderPaymentBoundNotification,
  OrderPaidNotification,
} from "@/types/ws-message";
import { Trigger } from "@/types/contracts";

export function parseChatMessage(raw: string): ChatMessage | null {
  try {
    const obj = JSON.parse(raw);

    if (obj.type === "text") {
      return obj as ChatTextMessage;
    }
    if (obj.type === "notification") {
      switch (obj.trigger) {
        case Trigger.OrderCompleted:
          return obj as OrderCompletedNotification;
        case Trigger.OrderPaymentBound:
          return obj as OrderPaymentBoundNotification;
        case Trigger.OrderPaid:
          return obj as OrderPaidNotification;
        default:
          return null;
      }
    }
    
    return null;
    
} catch (err) {
    return null;
  }
}
