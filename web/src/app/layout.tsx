"use client";

import "./globals.css";
import { useEffect, useState } from "react";
import { Navbar } from "@/components/ui/navbar";
import {auth} from "@/firebase/firebase";
import { onAuthStateChanged } from "firebase/auth";
import api from "@/lib/api/api-client";
import { Toaster } from "sonner";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [user, setUser] = useState<any>(null);
  const [fullname, setFullname] = useState<string>("");
  const [role, setRole] = useState<string>("customer");

  const fetchUserData = async () => {
    try {
      const response = await api.get('/api/users/me');
      const userData = response.data?.data || response.data;
      setFullname(userData?.fullname || "");
      setRole(userData?.role || "customer");
    } catch (error) {
      console.error("Failed to fetch user data:", error);
      setFullname("");
      setRole("customer");
    }
  };

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async (currentUser) => {
      setUser(currentUser);
      
      if (currentUser) {
        await fetchUserData();
      } else {
        setFullname("");
        setRole("customer");
      }
    });

    return () => unsubscribe();
  }, []);

  return (
    <html lang="en">
      <body className="font-sans antialiased bg-white dark:bg-black min-h-screen flex flex-col">
        {/* Header bar */}
        <Navbar role={role} username={fullname} onNameUpdate={fetchUserData} />

        {/* Main content */}
        <main className="flex-1 mt-20">{children}</main>
        <Toaster />
      </body>
    </html>
  );
}
