import { useEffect, useRef, useState, useCallback } from "react";
import { WSClient } from "./client";
import { ChatMessage } from "@/types/ws_message";

export function useWebSocket() {
  const [connected, setConnected] = useState(false);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const wsRef = useRef<WSClient | null>(null);

  // memoize message handler so it’s stable across renders
  const handleMessage = useCallback((msg: ChatMessage) => {
    setMessages((prev) => [...prev, msg]);
  }, []);

  useEffect(() => {
    const client = new WSClient({
      onOpen: () => setConnected(true),
      onClose: () => setConnected(false),
      onMessage: handleMessage,
    });

    client.connect();
    wsRef.current = client;

    return () => {
      client.disconnect();
      wsRef.current = null;
    };
  }, [handleMessage]);

  const send = useCallback((data: object) => {
    wsRef.current?.send(data);
  }, []);

  return { connected, messages, send };
}