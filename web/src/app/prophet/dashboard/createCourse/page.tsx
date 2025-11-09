"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";

export default function CreateCoursePage() {
  const router = useRouter();
  const [userName, setUserName] = useState(""); // automatically get user name
  const [title, setTitle] = useState("");
  const [price, setPrice] = useState("");
  const [description, setDescription] = useState("");
  const [selectedTags, setSelectedTags] = useState<string[]>([]); // multi-select tags

  // Predefined course tags
  const courseTags = ["Love", "Study", "Work", "Health", "Finance", "Personal Growth"];

  useEffect(() => {
    const storedUser = localStorage.getItem("user");
    if (storedUser) {
      const user = JSON.parse(storedUser);
      setUserName(user.name || user.email || "Unknown Prophet");
    }
  }, []);

  const handleTagChange = (tag: string) => {
    setSelectedTags((prev) =>
      prev.includes(tag) ? prev.filter((t) => t !== tag) : [...prev, tag]
    );
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!title.trim() || !price.trim() || !description.trim() || selectedTags.length === 0) {
      alert("Please fill in all fields and select at least one tag.");
      return;
    }

    const newCourse = {
      id: String(Date.now()),
      title,
      prophet: userName,
      price,
      description,
      tags: selectedTags,
      status: "active",
    };

    console.log("Course created:", newCourse);

    // Redirect back to dashboard
    router.push("/dashboard");
  };

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-900 p-6">
      <div className="max-w-lg mx-auto bg-white dark:bg-zinc-800 rounded-lg shadow p-6">
        <h1 className="text-2xl font-bold mb-4">Create New Course</h1>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          {/* Course Title */}
          <div>
            <label className="block mb-1 font-medium">Course Title</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full p-2 rounded border focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Enter course title"
              required
            />
          </div>

          {/* Price */}
          <div>
            <label className="block mb-1 font-medium">Price</label>
            <input
              type="text"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              className="w-full p-2 rounded border focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Enter price (e.g., $49)"
              required
            />
          </div>

          {/* Description */}
          <div>
            <label className="block mb-1 font-medium">Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full p-2 rounded border focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Enter course description"
              rows={4}
              required
            />
          </div>

          {/* Course Tags (Multi-select) */}
          <div>
            <label className="block mb-2 font-medium">Course Tags</label>
            <div className="flex flex-wrap gap-2">
              {courseTags.map((tag) => (
                <label key={tag} className="flex items-center gap-1 bg-gray-100 dark:bg-zinc-700 px-3 py-1 rounded cursor-pointer hover:bg-gray-200 dark:hover:bg-zinc-600">
                  <input
                    type="checkbox"
                    value={tag}
                    checked={selectedTags.includes(tag)}
                    onChange={() => handleTagChange(tag)}
                    className="w-4 h-4"
                  />
                  <span className="text-sm">{tag}</span>
                </label>
              ))}
            </div>
          </div>

          {/* Buttons */}
          <div className="flex gap-2 justify-end mt-4">
            <button
              type="button"
              onClick={() => router.push("/dashboard")}
              className="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Create Course
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
