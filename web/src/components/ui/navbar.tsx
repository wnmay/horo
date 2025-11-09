"use client";

import { auth } from "@/firebase/firebase";
import { useEffect } from "react";
import { useState } from "react";
import axios from "axios";
import { useRouter } from "next/navigation";
import { useRef } from "react";

interface navbarProps
{
    className?: string;
    username: string;
    role: string;
}

export function Navbar({
    className = "",
    username,
    role
}: navbarProps)
{
    const [isAuthed, setIsAuthed] = useState<boolean>(false);
    const router = useRouter();
    const [user, setUser] = useState<any>(null);
    const [showMenu, setShowMenu] = useState(false);
    const menuRef = useRef<HTMLDivElement>(null);
    const [changeName, setChangeName] = useState(false);
    const [newUsername, setNewUsername] = useState(username);

    const api = axios.create({
        baseURL: process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080",
    });

    api.interceptors.request.use(async (config) => {
        const token = await auth.currentUser?.getIdToken();
        setIsAuthed(!!token);
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    });

    // Close menu when clicking outside
    useEffect(() => {
        const handleClickOutside = (e: MouseEvent) => {
        if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
            setShowMenu(false);
            setChangeName(false);
        }
        };
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    return (
        <header className="fixed top-4 right-4 z-50 flex items-center space-x-3">
          {true ? (
            <div className="flex items-center relative" ref={menuRef}>
              {/* Dashboard button (always visible for prophet) */}
              {role === "prophet" && (
                <button
                  onClick={() => router.push("/prophet/dashboard")}
                  className="px-4 py-2 bg-white text-blue-600 rounded-lg shadow-md hover:shadow-lg transition"
                >
                  Dashboard
                </button>
              )}

              {/* User icon */}
              <div
                className="relative flex items-center bg-white dark:bg-zinc-800 rounded-full px-3 py-1 hover:shadow-lg transition cursor-pointer ml-2"
                onClick={() => setShowMenu(!showMenu)}
              >
                <img
                  src={`https://ui-avatars.com/api/?name=${encodeURIComponent(
                    // user.name || user.email
                    username
                  )}&background=random`}
                  alt="Profile"
                  className="w-8 h-8 rounded-full mr-2"
                />
                <span className="text-sm font-medium text-zinc-800 dark:text-zinc-200">
                  {(username)?.split(" ")[0]}
                </span>

                {/* Dropdown menu */}
                {showMenu && (
                  <div className="absolute top-full right-0 mt-2 w-56 flex flex-col bg-white dark:bg-zinc-800 rounded-lg shadow-lg border border-zinc-200 dark:border-zinc-700 z-50 overflow-hidden">
                    {role == "customer" && <button
                      onClick={(e) => {
                        e.stopPropagation();
                        setChangeName(true)}}
                      className="block w-full text-left px-4 py-2 hover:bg-zinc-100 dark:hover:bg-zinc-700 transition"
                    >
                      Change Full Name
                    </button>}
                    {role == "customer" && changeName && 
                    <form 
                        onClick={(e) => e.stopPropagation()} 
                        onSubmit={(e) => {
                            e.preventDefault();
                            console.log("New name:", newUsername);
                            // TODO: call API patch/update name if needed
                            setChangeName(false);
                        }}
                        className="w-full flex flex-col items-center gap-2">
                        <input 
                            type="text"
                            placeholder={`${username}`} 
                            className="w-[90%] h-6 rounded border-gray-300 border p-2"
                            value={newUsername}
                            onChange={(e) => setNewUsername(e.target.value)} />
                        <div className="flex w-full items-center justify-center gap-2">
                            <button
                                type="submit"
                                className="w-1/3 bg-blue-300 rounded hover:bg-blue-400">
                                Save
                            </button>
                            <button
                                type="button"
                                className="w-1/3 bg-slate-300 rounded hover:bg-slate-400"
                                onClick={() => setChangeName(false)}>
                                Cancel
                            </button>
                        </div>
                    </form>}
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
    );
}