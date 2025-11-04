"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { auth } from "../../firebase/firebase";

export default function ProfilePage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const queryToken = useMemo(
    () => searchParams.get("token") ?? "",
    [searchParams]
  );

  const [idToken, setIdToken] = useState<string>(queryToken);
  const [fullName, setFullName] = useState("");
  const [role, setRole] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    (async () => {
      if (idToken) return;
      const user = auth.currentUser;
      if (user) {
        try {
          const t = await user.getIdToken();
          setIdToken(t);
        } catch (err) {
          console.error("getIdToken error", err);
        }
      }
    })();
  }, [idToken]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!idToken) {
      setError("Missing ID token. Please sign in again.");
      return;
    }

    setError("");
    setIsSubmitting(true);
    try {
      const res = await fetch("http://localhost:3000/api/users/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ idToken, fullName, role }),
      });

      if (!res.ok) throw new Error("Failed to register user");
      router.push("/home");
    } catch (err) {
      console.error(err);
      setError("Something went wrong. Please try again.");
      setIsSubmitting(false);
    }
  };

  return (
    <main className="w-full min-h-screen flex items-center justify-center bg-gray-50">
      <div className="w-96 bg-white shadow-xl p-6 rounded-xl">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">
          Complete Your Profile
        </h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="text-sm font-bold text-gray-600">Full Name</label>
            <input
              type="text"
              value={fullName}
              onChange={(e) => setFullName(e.target.value)}
              className="w-full mt-1 px-3 py-2 border rounded-lg focus:border-indigo-600 outline-none"
              required
            />
          </div>

          <div>
            <label className="text-sm font-bold text-gray-600">Role</label>
            <select
              value={role}
              onChange={(e) => setRole(e.target.value)}
              className="w-full mt-1 px-3 py-2 border rounded-lg focus:border-indigo-600 outline-none"
              required
            >
              <option value="">Select a role</option>
              <option value="customer">ลูกดวง</option>
              <option value="prophet">หมอดู</option>
            </select>
          </div>

          {error && <p className="text-red-600 text-sm font-bold">{error}</p>}

          <button
            type="submit"
            disabled={isSubmitting}
            className={`w-full px-4 py-2 text-white font-medium rounded-lg ${
              isSubmitting
                ? "bg-gray-300 cursor-not-allowed"
                : "bg-indigo-600 hover:bg-indigo-700 transition duration-300"
            }`}
          >
            {isSubmitting ? "Submitting..." : "Submit"}
          </button>
          <p className="text-xs text-gray-500 mt-2">
            {idToken ? "Token ready" : "No token yet — sign in first"}
          </p>
        </form>
      </div>
    </main>
  );
}
