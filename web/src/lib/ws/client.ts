import { parseChatMessage } from "./parser";
import { ChatMessage } from "@/types/ws_message";

export class WSClient {
  private ws: WebSocket | null = null;
  private readonly url: string;
  private readonly onMessage?: (msg: ChatMessage) => void;
  private readonly onOpen?: () => void;
  private readonly onClose?: () => void;
  private readonly onError?: (e: Event) => void;

  constructor({
    onMessage,
    onOpen,
    onClose,
    onError,
  }: {
    onMessage?: (msg: ChatMessage) => void;
    onOpen?: () => void;
    onClose?: () => void;
    onError?: (e: Event) => void;
  }) {
    const envUrl = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/ws/chat";

    if (!envUrl) {
      throw new Error("Missing NEXT_PUBLIC_WS_URL in your .env file");
    }

    this.url = envUrl;
    this.onMessage = onMessage;
    this.onOpen = onOpen;
    this.onClose = onClose;
    this.onError = onError;
  }

  connect() {
    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => this.onOpen?.();
    this.ws.onmessage = (event) => {
      const msg = parseChatMessage(event.data);
      if (msg) this.onMessage?.(msg);
    };
    this.ws.onerror = (e) => this.onError?.(e);
    this.ws.onclose = () => this.onClose?.();
  }

  send(data: object) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    } else {
      console.warn("WebSocket not connected");
    }
  }

  disconnect() {
    this.ws?.close();
    this.ws = null;
  }
}
