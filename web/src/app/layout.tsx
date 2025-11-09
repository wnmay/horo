"use client";

import "./globals.css";
import { useEffect, useState, useRef } from "react";
import { Navbar } from "@/components/ui/navbar";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<any>(null);

  // Load user from localStorage
  useEffect(() => {
    const storedUser = localStorage.getItem("user");
    if (storedUser) setUser(JSON.parse(storedUser));
  }, []);

  return (
    <html lang="en">
      <body className="font-sans antialiased bg-white dark:bg-black min-h-screen flex flex-col">
        {/* Header bar */}
        <Navbar role={"customer"} username={"Prim"}/>

        {/* Main content */}
        <main className="flex-1 mt-20">{children}</main>
      </body>
    </html>
  );
}
