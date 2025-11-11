"use client";

import React, { useEffect, useState } from "react";
import { auth } from "@/firebase/firebase";
import ChatRoomList from "@/components/ChatRoomList";
import RightPanel from "@/components/RightChatPanel";

export default function Page() {
  const [role, setRole] = useState<"customer" | "prophet" | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const unsub = auth.onAuthStateChanged(async (user) => {
      if (user) {
        try {
          const tokenResult = await user.getIdTokenResult();
          const claims = tokenResult.claims;
          const userRole =
            (claims.role as "customer" | "prophet") ?? "customer";
          setRole(userRole);
        } catch (err) {
          console.error("getIdTokenResult error:", err);
          setRole("customer");
        }
      } else {
        setRole(null);
      }
      setLoading(false);
    });
    return () => unsub();
  }, []);

  if (loading) return <p className="p-4 text-gray-500">Loading user...</p>;
  if (!role) return <p className="p-4 text-red-500">No user signed in.</p>;

  return (
    <div className="flex w-full">
      <div className="w-[30%]">
        <ChatRoomList />
      </div>
      <div className="flex-1">
        <RightPanel
          roomId="690420c673da5edbbe275321"
          role={role}
          courseId="COURSE-7a230ed4-b2d9-4004-b86e-ffb5da9c4e93"
        />
      </div>
    </div>
  );
}
