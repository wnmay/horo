import { useEffect, useRef, useState, useCallback } from "react";
import { WSClient } from "./client";
import { ChatMessage } from "@/types/ws-message";
import { auth } from "@/firebase/firebase";
import { onAuthStateChanged, onIdTokenChanged } from "firebase/auth";

export function useWebSocket() {
  const [connected, setConnected] = useState(false);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const wsRef = useRef<WSClient | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  // stable message handler
  const handleMessage = useCallback((msg: ChatMessage) => {
    setMessages((prev) => [...prev, msg]);
  }, []);

  useEffect(() => {
    const unsubAuth = onAuthStateChanged(auth, async (user) => {
      if (!user) {
        setToken(null);
        setError("User not logged in");
        return;
      }
      try {
        const t = await user.getIdToken(true);
        setToken(t);
        setError(null);
      } catch (e) {
        setError("Failed to get token");
      }
    });

    const unsubToken = onIdTokenChanged(auth, async (user) => {
      if (!user) return;
      try {
        const t = await user.getIdToken(true);
        setToken(t);
      } catch (e) {
        console.error("[WS] onIdTokenChanged getIdToken failed:", e);
      }
    });

    return () => {
      unsubAuth();
      unsubToken();
    };
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

  const joinRoom = useCallback((roomId: string) => {
    wsRef.current?.send({ action: "join_room", roomId });
  }, []);

  const sendMessage = useCallback((data: object) => {
    wsRef.current?.send(data);
  }, []);

  return { connected, messages, send, joinRoom, sendMessage };
}
