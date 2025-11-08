"use client";

import { useWebSocket } from "@/lib/ws/useWebSocket";

export default function WSTestPage() {
  const { connected } = useWebSocket();

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full">
        <h1 className="text-2xl font-bold mb-6 text-center">
          WebSocket Connection Test
        </h1>
        
        <div className="flex items-center justify-center space-x-3">
          <div
            className={`w-4 h-4 rounded-full ${
              connected ? "bg-green-500" : "bg-red-500"
            } animate-pulse`}
          />
          <span className="text-lg font-medium">
            {connected ? "Connected ✅" : "Disconnected ❌"}
          </span>
        </div>

        <div className="mt-6 p-4 bg-gray-100 rounded text-sm">
          <p className="font-semibold mb-2">Connection Info:</p>
          <p className="text-gray-600">
            URL: {process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/ws/chat"}
          </p>
          <p className="text-gray-600 mt-1">
            Status: {connected ? "Active" : "Inactive"}
          </p>
        </div>
      </div>
    </div>
  );
}

