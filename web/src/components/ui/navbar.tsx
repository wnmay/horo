"use client";

import { auth } from "@/firebase/firebase";
import { useEffect } from "react";
import { useState } from "react";
import axios from "axios";
import { useRouter } from "next/navigation";
import { useRef } from "react";
import api from "@/lib/api/api-client";
import { onAuthStateChanged } from "firebase/auth";

interface navbarProps
{
    className?: string;
    username: string;
    role: string;
    onNameUpdate?: () => void;
}

export function Navbar({
    className = "",
    username,
    role,
    onNameUpdate
}: navbarProps)
{
    const [isAuthed, setIsAuthed] = useState<boolean>(false);
    const router = useRouter();
    const [user, setUser] = useState<any>(null);
    const [showMenu, setShowMenu] = useState(false);
    const menuRef = useRef<HTMLDivElement>(null);
    const [changeName, setChangeName] = useState(false);
    const [newUsername, setNewUsername] = useState(username);

    const updateUserName = async () => {
      try{
        await api.patch('/api/users/update-name',{"fullname": newUsername});
        
        // Call the callback to refresh user data from parent
        if (onNameUpdate) {
          await onNameUpdate();
        }
      } catch (error) {
        console.error("Failed to update username:", error);
      }
    }

    // Listen to Firebase auth state changes
    useEffect(() => {
        const unsubscribe = auth.onAuthStateChanged((currentUser) => {
            if (currentUser) {
                setIsAuthed(true);
                setUser(currentUser);
            } else {
                setIsAuthed(false);
                setUser(null);
            }
        });
        return () => unsubscribe();
    }, []);

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
        <header className="fixed w-full h-16 z-50 p-4 flex items-center border-b-2 border-black shadow-md bg-white">
          <h2 className="fixed left-10 text-3xl text-center font-bold">Horo</h2>
          
          {isAuthed ? (
            <div className="fixed right-4 flex items-center" ref={menuRef}>
              {/* Home button */}
              <button 
                onClick={() => router.push("/")}
                className="px-4 py-2 bg-white rounded-lg hover:text-gray-500">Home</button>

              {/* Chat button (visible for all roles) */}
              <button
                onClick={() => router.push("/chat")}
                className="px-4 py-2 bg-white rounded-lg hover:text-gray-400 transition flex items-center gap-2"
                title="Chat"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={1.5}
                  stroke="currentColor"
                  className="w-5 h-5"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M20.25 8.511c.884.284 1.5 1.128 1.5 2.097v4.286c0 1.136-.847 2.1-1.98 2.193-.34.027-.68.052-1.02.072v3.091l-3-3c-1.354 0-2.694-.055-4.02-.163a2.115 2.115 0 01-.825-.242m9.345-8.334a2.126 2.126 0 00-.476-.095 48.64 48.64 0 00-8.048 0c-1.131.094-1.976 1.057-1.976 2.192v4.286c0 .837.46 1.58 1.155 1.951m9.345-8.334V6.637c0-1.621-1.152-3.026-2.76-3.235A48.455 48.455 0 0011.25 3c-2.115 0-4.198.137-6.24.402-1.608.209-2.76 1.614-2.76 3.235v6.226c0 1.621 1.152 3.026 2.76 3.235.577.075 1.157.14 1.74.194V21l4.155-4.155"
                  />
                </svg>
                Chat
              </button>

              {/* Dashboard button (always visible for prophet) */}
              {role === "prophet" && (
                <button
                  onClick={() => router.push("/prophet/dashboard")}
                  className="px-4 py-2 bg-white rounded-lg hover:text-gray-400 transition"
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
                    username || user?.displayName || user?.email || "User"
                  )}&background=random`}
                  alt="Profile"
                  className="w-8 h-8 rounded-full mr-2"
                />
                <span className="text-sm font-medium text-zinc-800 dark:text-zinc-200">
                  {(username || user?.displayName || user?.email || "User")?.split(" ")[0]}
                </span>

                {/* Dropdown menu */}
                {showMenu && (
                  <div className="absolute top-full right-0 mt-2 w-56 flex flex-col bg-white dark:bg-zinc-800 rounded-lg shadow-lg border border-zinc-200 dark:border-zinc-700 z-50 overflow-hidden">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        setChangeName(true)}}
                      className="block w-full text-left px-4 py-2 hover:bg-zinc-100 dark:hover:bg-zinc-700 transition"
                    >
                      Change Full Name
                    </button>
                    {changeName && 
                    <form 
                        onClick={(e) => e.stopPropagation()} 
                        onSubmit={(e) => {
                            e.preventDefault();
                            console.log("New name:", newUsername);
                            setChangeName(false);
                            updateUserName();
                        }}
                        className="w-full flex flex-col items-center gap-2">
                        <input 
                            type="text"
                            placeholder={username || user?.displayName || user?.email || "Enter your name"} 
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
                        auth.signOut();
                        setUser(null);
                        setIsAuthed(false);
                        router.push("/");
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
            <div className="fixed right-4 flex gap-3">
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