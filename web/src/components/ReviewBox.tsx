"use client";

import { useState } from "react";
import { auth } from "@/firebase/firebase";
import api from "@/lib/api/api-client";

export default function ReviewBox({
  courseId,
  onSubmitted,
}: {
  courseId: string;
  onSubmitted?: () => void;
}) {
  const [score, setScore] = useState<number>(0);
  const [title, setTitle] = useState("");
  const [desc, setDesc] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [submitted, setSubmitted] = useState(false);

  const Star = ({ i }: { i: number }) => (
    <button
      type="button"
      onClick={() => !submitted && setScore(i)}
      className={`text-2xl leading-none ${
        i <= score ? "text-yellow-500" : "text-gray-300"
      } ${submitted ? "opacity-50 cursor-not-allowed" : ""}`}
      aria-label={`rate ${i}`}
      disabled={submitted}
    >
      ★
    </button>
  );

  const submit = async () => {
    setError(null);
    try {
      const user = auth.currentUser;
      if (!user) {
        setError("Please sign in first.");
        return;
      }
      if (score < 1 || score > 5) {
        setError("Score must be between 1 and 5.");
        return;
      }

      setSubmitting(true);

      const payload = {
        customer_id: user.uid,
        customername: user.displayName || "Anonymous",
        score,
        title,
        description: desc,
      };

      await api.post(`/api/courses/${courseId}/review`, payload);

      setSubmitted(true);
      onSubmitted?.();
    } catch (e: any) {
      setError(e?.message ?? "Failed to submit review");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="rounded-xl border border-violet-200 bg-violet-50/60 p-4">
      <h4 className="text-sm font-semibold text-violet-900">Write a review</h4>

      <div className="mt-3 space-y-3">
        <div className="flex items-center gap-2">
          {[1, 2, 3, 4, 5].map((i) => (
            <Star key={i} i={i} />
          ))}
          <span className="ml-2 text-xs text-gray-600">{score}/5</span>
        </div>

        <input
          disabled={submitted}
          className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-violet-300 disabled:opacity-60"
          placeholder="Title (optional)"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />

        <textarea
          disabled={submitted}
          className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-violet-300 disabled:opacity-60"
          placeholder="Share your experience…"
          rows={4}
          value={desc}
          onChange={(e) => setDesc(e.target.value)}
        />

        {error && (
          <div className="rounded-md border border-rose-200 bg-rose-50 px-3 py-2 text-sm text-rose-700">
            {error}
          </div>
        )}

        {submitted ? (
          <div className="rounded-md border border-emerald-200 bg-emerald-50 px-3 py-2 text-sm text-emerald-700 text-center">
            Review submitted successfully
          </div>
        ) : (
          <button
            onClick={submit}
            disabled={submitting || score < 1 || score > 5}
            className="w-full rounded-xl px-4 py-2.5 font-semibold bg-violet-600 text-white hover:bg-violet-700 disabled:opacity-60"
          >
            {submitting ? "Submitting…" : "Submit review"}
          </button>
        )}
      </div>
    </div>
  );
}
