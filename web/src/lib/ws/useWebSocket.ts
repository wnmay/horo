import { useEffect, useRef, useState, useCallback } from "react";
import { WSClient } from "./client";
import { ChatMessage } from "@/types/ws_message";
import { useAuth } from "@/context/mockAuthProvider"; 

export function useWebSocket() {
  const { token } = useAuth(); // TO DO: replace with the real auth when finished
  const [connected, setConnected] = useState(false);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const wsRef = useRef<WSClient | null>(null);

  // stable message handler
  const handleMessage = useCallback((msg: ChatMessage) => {
    setMessages((prev) => [...prev, msg]);
  }, []);

  useEffect(() => {
    if (!token) {
      return;
    }

    const client = new WSClient({
      onOpen: () => setConnected(true),
      onClose: () => setConnected(false),
      onMessage: handleMessage,
    });

    client.connect(token); 
    wsRef.current = client;

    return () => {
      client.disconnect();
      wsRef.current = null;
    };
  }, [token, handleMessage]);

  const send = useCallback((data: object) => {
    wsRef.current?.send(data);
  }, []);

  return { connected, messages, send };
}
