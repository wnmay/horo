"use client";

import "./globals.css";
import { useEffect, useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { Toaster } from "sonner";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const [user, setUser] = useState<any>(null);
  const [showMenu, setShowMenu] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  // Load user from localStorage
  useEffect(() => {
    const storedUser = localStorage.getItem("user");
    if (storedUser) setUser(JSON.parse(storedUser));
  }, []);

  // Close menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setShowMenu(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <html lang="en">
      <head />
      <body className="font-sans antialiased bg-white dark:bg-black min-h-screen flex flex-col">
        {/* Header */}
        <header className="fixed top-4 right-4 z-50 flex items-center space-x-3">
          {user ? (
            <div className="flex items-center relative" ref={menuRef}>
              {/* Dashboard button */}
              {user.role === "prophet" && (
                <button
                  onClick={() => router.push("/prophet/dashboard")}
                  className="mr-3 px-4 py-2 bg-white dark:bg-zinc-800 rounded-lg shadow-md hover:shadow-lg transition text-sm font-medium"
                >
                  Dashboard
                </button>
              )}

              {/* Profile Icon */}
              <div
                className="relative flex items-center bg-white dark:bg-zinc-800 rounded-full px-3 py-1 hover:shadow-lg transition cursor-pointer"
                onClick={() => setShowMenu(!showMenu)}
              >
                <img
                  src={`https://ui-avatars.com/api/?name=${encodeURIComponent(
                    user.name || user.email
                  )}&background=random`}
                  alt="Profile"
                  className="w-8 h-8 rounded-full mr-2"
                />
                <span className="text-sm font-medium text-zinc-800 dark:text-zinc-200">
                  {(user.name || user.email)?.split(" ")[0]}
                </span>
              </div>

              {/* Dropdown Menu */}
              {showMenu && (
                <div className="absolute top-full right-0 mt-2 w-56 flex flex-col bg-white dark:bg-zinc-800 rounded-lg shadow-lg border border-zinc-200 dark:border-zinc-700 z-50 overflow-hidden">
                  <button
                    onClick={() => {
                      const newName = prompt(
                        "Enter your full name",
                        user.name || ""
                      );
                      if (newName) {
                        const updatedUser = { ...user, name: newName };
                        setUser(updatedUser);
                        localStorage.setItem(
                          "user",
                          JSON.stringify(updatedUser)
                        );
                      }
                      setShowMenu(false);
                    }}
                    className="block w-full text-left px-4 py-2 hover:bg-zinc-100 dark:hover:bg-zinc-700 transition"
                  >
                    Change Full Name
                  </button>

                  <button
                    onClick={() => {
                      alert("Switch account feature not implemented yet");
                      setShowMenu(false);
                    }}
                    className="block w-full text-left px-4 py-2 hover:bg-zinc-100 dark:hover:bg-zinc-700 transition"
                  >
                    Switch Account
                  </button>

                  <button
                    onClick={() => {
                      localStorage.removeItem("user");
                      setUser(null);
                      router.refresh();
                    }}
                    className="block w-full text-left px-4 py-2 text-red-500 hover:bg-zinc-100 dark:hover:bg-zinc-700 transition"
                  >
                    Logout
                  </button>
                </div>
              )}
            </div>
          ) : (
            <div className="flex gap-3">
              <button
                onClick={() => router.push("/signin")}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
              >
                Sign In
              </button>
              <button
                onClick={() => router.push("/signup")}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
              >
                Sign Up
              </button>
            </div>
          )}
        </header>

        {/* Main Content */}
        <main className="flex-1 mt-20">{children}</main>

        {/* Toast notifications */}
        <Toaster richColors position="top-right" />
      </body>
    </html>
  );
}
