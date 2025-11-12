"use client";

import React, { createContext, useContext } from "react";
import { useWebSocket } from "@/lib/ws/useWebSocket";
import type { ChatMessage } from "@/types/ws-message";

// Define what values are exposed to consumers
type WebSocketContextType = {
  connected: boolean;
  messages: ChatMessage[];
  send: (data: object) => void;
  joinRoom: (roomId: string) => void;
  sendMessage: (data: object) => void;
};

// Create a React context
const WebSocketContext = createContext<WebSocketContextType | null>(null);

// Provider component — wraps your app/page
export function WebSocketProvider({ children }: { children: React.ReactNode }) {
  const { connected, messages, send, joinRoom, sendMessage } = useWebSocket();

  return (
    <WebSocketContext.Provider
      value={{ connected, messages, send, joinRoom, sendMessage }}
    >
      {children}
    </WebSocketContext.Provider>
  );
}

// Hook for consuming the context
export function useWebSocketCtx() {
  const ctx = useContext(WebSocketContext);
  if (!ctx) {
    throw new Error("useWebSocketCtx must be used within a WebSocketProvider");
  }
  return ctx;
}
