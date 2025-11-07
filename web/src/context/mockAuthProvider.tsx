"use client";

import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";

// shape of the context
interface AuthContextType {
  token: string | null;
  login: (newToken: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function useAuth(): AuthContextType {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside an AuthProvider");
  return ctx;
}

// your hardcoded mock JWT
const MOCK_TOKEN = 'replace-with-firebase-token'

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    // auto-login immediately with mock token
    setToken(MOCK_TOKEN);
    localStorage.setItem("mock_token", MOCK_TOKEN);
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
