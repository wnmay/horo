"use client";

import { useState } from "react";
import api from "@/lib/api/api-client";

export default function FetchRoomsButton() {
  const [loading, setLoading] = useState(false);

  const fetchRooms = async () => {
    setLoading(true);
    try {
      const res = await api.get(
        "/api/payments/balance"
        // ,{payload} can pass as plain ts object
      );
      console.log(
        "[Chat Rooms] Response JSON:",
        JSON.stringify(res.data, null, 2)
      );
      alert("Fetched! Check console for JSON output.");
    } catch (err: any) {
      console.error("[Chat Rooms] Error:", err.response?.data || err.message);
      alert("Failed to fetch chat rooms. See console for details.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <button
      onClick={fetchRooms}
      disabled={loading}
      className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
    >
      {loading ? "Fetching..." : "Fetch Chat Rooms"}
    </button>
  );
}
