"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { doCreateUserWithEmailAndPassword } from "@/firebase/auth";
import { db } from "@/firebase/firebase";
import { doc, serverTimestamp, setDoc } from "firebase/firestore";
import axios from "axios";

export default function RegisterPage() {
  const router = useRouter();
  const [fullname, setFullname] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [role, setRole] = useState("customer");
  const [loading, setLoading] = useState(false);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();

    if (password !== confirm) {
      alert("Passwords do not match");
      return;
    }

    if (!fullname.trim()) {
      alert("Full Name is required");
      return;
    }

    setLoading(true);

    try {
      // Firebase registration
      const userCredential = await doCreateUserWithEmailAndPassword(email, password);
      const user = userCredential.user;
      const tokenId = await user.getIdToken();

      // Save user info in Firestore
      await setDoc(doc(db, "users", user.uid), {
        email: user.email,
        role: role,
        fullname: fullname,
        createdAt: serverTimestamp(),
      });

      // Call API Gateway to register user
      await axios.post(
        `http://localhost:3000/api/users/register`,
        { idToken: tokenId, fullname, role },
        { headers: { Authorization: `Bearer ${tokenId}` } }
      );

      alert("Registration successful!");
      router.push("/signin");
    } catch (error: any) {
      console.error("Error registering:", error.response?.data || error.message);
      alert(error.response?.data?.message || error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center relative">
      <Card className="w-full max-w-md p-6">
        <h1 className="text-3xl font-bold mb-6 text-center">Register</h1>
        <form onSubmit={handleRegister} className="flex flex-col gap-4">
          
          {/* Full Name */}
          <div className="flex flex-col">
            <label className="text-sm text-gray-500 mb-1">Full Name</label>
            <input
              type="text"
              placeholder="Full Name"
              className="border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={fullname}
              onChange={(e) => setFullname(e.target.value)}
              required
            />
          </div>

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

          {/* Confirm Password */}
          <div className="flex flex-col">
            <label className="text-sm text-gray-500 mb-1">Confirm Password</label>
            <input
              type="password"
              placeholder="Confirm Password"
              className="border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              required
            />
          </div>

          {/* Role */}
          <div className="flex flex-col">
            <label className="text-sm text-gray-500 mb-1">Role</label>
            <select
              className="border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={role}
              onChange={(e) => setRole(e.target.value)}
            >
              <option value="customer">Customer</option>
              <option value="prophet">Prophet</option>
            </select>
          </div>

          {/* Submit Button */}
          <Button
            type="submit"
            className="w-full bg-blue-500 hover:bg-blue-600 text-white rounded py-2 mt-2"
            disabled={loading}
          >
            {loading ? "Registering..." : "Register"}
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