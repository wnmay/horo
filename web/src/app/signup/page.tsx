"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { doCreateUserWithEmailAndPassword } from "@/firebase/auth";
import { db } from "@/firebase/firebase"; // adjust path
import { doc, serverTimestamp, setDoc } from "firebase/firestore";

export default function RegisterPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [role, setRole] = useState("customer"); // default role
  const [loading, setLoading] = useState(false);
  

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();

    if (password !== confirm) {
      alert("Passwords do not match");
      return;
    }

    try {
      // TODO: call your Firebase register function here
      const userCredential = await doCreateUserWithEmailAndPassword(email,password)
      const user = userCredential.user

      await setDoc(doc(db,"users", user.uid), {
        email:user.email,
        role : role,
        createAt : serverTimestamp()
      })
      alert("Registration successful!");
      router.push("/signin");
    } catch (error : any) {
      console.error("Error registering:", error.message);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center relative">
      <Card>
        <h1 className="text-3xl font-bold mb-6 text-center">Register</h1>
        <form onSubmit={handleRegister} className="flex flex-col gap-4 w-80">
          <input
            type="email"
            placeholder="Email"
            className="border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />

          <input
            type="password"
            placeholder="Password"
            className="border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />

          <input
            type="password"
            placeholder="Confirm Password"
            className="border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
            required
          />

          {/* Role selection */}
          <select
            className="border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={role}
            onChange={(e) => setRole(e.target.value)}
          >
            <option value="customer">Customer</option>
            <option value="prophet">Prophet</option>
          </select>

          <Button
            type="submit"
            className="w-full bg-blue-500 hover:bg-blue-600 text-white rounded py-2"
          >
            Register
          </Button>
        </form>

        <p className="text-sm text-center text-gray-600 mt-4">
          Already have an account?{" "}
          <button
            onClick={() => router.push("/signin")}
            className="text-green-500 hover:underline"
          >
            Sign In
          </button>
        </p>
      </Card>
    </div>
  );
}
