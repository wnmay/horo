"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { signInWithEmailAndPassword, GoogleAuthProvider, signInWithRedirect } from "firebase/auth";
import { auth } from "@/firebase/firebase";
import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function SignInPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleSignIn = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await signInWithEmailAndPassword(auth, email, password);
      alert("login succesful")
      router.push("/");
    } catch (err: any) {
      setError(err.message);
      console.error(err);
    }
  };

  const handleGoogleSignIn = async () => {
    try {
      const provider = new GoogleAuthProvider();
      await signInWithRedirect(auth, provider);
      router.push("/");
    } catch (err: any) {
      setError(err.message);
      console.error(err);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center relative">
      <Card className="w-full max-w-md p-6">
        <h1 className="text-3xl font-bold mb-6 text-center">Sign In</h1>
        <form onSubmit={handleSignIn} className="flex flex-col gap-4">
          
          {/* Email */}
          <div className="flex flex-col">
            <label className="text-sm text-gray-500 mb-1">Email Address</label>
            <input
              type="email"
              placeholder="Email"
              className="border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>

          {/* Password */}
          <div className="flex flex-col">
            <label className="text-sm text-gray-500 mb-1">Password</label>
            <input
              type="password"
              placeholder="Password"
              className="border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          {error && <p className="text-red-500 text-sm text-center">{error}</p>}

          {/* Sign In Button */}
          <Button
            type="submit"
            className="w-full bg-green-500 hover:bg-green-600 text-white rounded py-2 mt-2"
          >
            Sign In
          </Button>

          {/* Google Sign In */}
          <Button
            type="button"
            onClick={handleGoogleSignIn}
            className="w-full bg-red-500 hover:bg-red-600 text-white rounded py-2 mt-2"
          >
            Sign in with Google
          </Button>
        </form>

        <p className="text-sm text-center text-gray-600 mt-4">
          Don’t have an account?{" "}
          <button
            onClick={() => router.push("/signup")}
            className="text-blue-500 hover:underline"
          >
            Register
          </button>
        </p>
      </Card>
    </div>
  );
}
