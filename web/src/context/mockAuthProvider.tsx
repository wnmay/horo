//ref web-firebase\src\contexts\authContext\index.jsx
"use client";

import { createContext, useContext, useState, useEffect, ReactNode } from "react";

// The shape of the auth context
interface AuthContextType {
  token: string | null;
  login: (newToken: string) => void;
  logout: () => void;
}

// Default (empty) context
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Hook for easy access
export function useAuth(): AuthContextType {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used inside an AuthProvider");
  }
  return ctx;
}

// Mock Provider implementation
export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null);

  // mock auto-login behavior (optional)
  useEffect(() => {
    // If you want, pull a token from localStorage to simulate persistence
    const saved = localStorage.getItem("mock_token");
    if (saved) setToken(saved);
  }, []);

  const login = (newToken: string) => {
    setToken(newToken);
    localStorage.setItem("mock_token", newToken);
  };

  const logout = () => {
    setToken(null);
    localStorage.removeItem("mock_token");
  };

  return (
    <AuthContext.Provider value={{ token, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}
