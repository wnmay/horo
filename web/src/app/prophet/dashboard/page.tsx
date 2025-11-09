"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

interface Course {
  id: string;
  title: string;
  description: string;
  price: string;
  status: "active" | "inactive";
}

export default function DashboardPage() {
  const router = useRouter();
  const [courses, setCourses] = useState<Course[]>([]);
  const [editingCourseId, setEditingCourseId] = useState<string | null>(null);
  const [editTitle, setEditTitle] = useState("");
  const [editDescription, setEditDescription] = useState("");
  const [editPrice, setEditPrice] = useState("");

  useEffect(() => {
    async function fetchCourses() {
      const mockCourses: Course[] = [
        { id: "1", title: "How to Read the Stars", description: "Learn the basics of astrology.", price: "$49", status: "active" },
        { id: "2", title: "Dream Interpretation 101", description: "Understand the meaning of dreams.", price: "$59", status: "active" },
        { id: "3", title: "Life Path Guidance", description: "Discover your life path.", price: "$69", status: "inactive" },
      ];
      await new Promise((resolve) => setTimeout(resolve, 300));
      setCourses(mockCourses);
    }
    fetchCourses();
  }, []);

  const startEditing = (course: Course) => {
    setEditingCourseId(course.id);
    setEditTitle(course.title);
    setEditDescription(course.description);
    setEditPrice(course.price);
  };

  const saveEdit = (id: string) => {
    setCourses((prev) =>
      prev.map((c) => (c.id === id ? { ...c, title: editTitle, description: editDescription, price: editPrice } : c))
    );
    setEditingCourseId(null);
  };

  const toggleStatus = (id: string) => {
    setCourses((prev) =>
      prev.map((c) => (c.id === id ? { ...c, status: c.status === "active" ? "inactive" : "active" } : c))
    );
  };

  const deleteCourse = (id: string) => {
    setCourses((prev) => prev.filter((c) => c.id !== id));
  };

  return (
    <div className="p-6">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">My Courses</h1>
        <button
          onClick={() => router.push("/prophet/dashboard/createCourse")}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 shadow transition"
        >
          + Create Course
        </button>
      </div>

      {/* Courses Grid */}
      <ul className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {courses.map((course) => (
          <li
            key={course.id}
            className={`relative rounded-xl border shadow-sm overflow-hidden transition-all duration-200 hover:shadow-lg hover:scale-[1.02] ${
              course.status === "inactive" ? "opacity-70" : ""
            }`}
          >
            <div className="p-5 flex flex-col justify-between h-full">
              {editingCourseId === course.id ? (
                <>
                  <input
                    type="text"
                    value={editTitle}
                    onChange={(e) => setEditTitle(e.target.value)}
                    className="w-full border px-3 py-1 rounded mb-2"
                    placeholder="Course Title"
                  />
                  <textarea
                    value={editDescription}
                    onChange={(e) => setEditDescription(e.target.value)}
                    className="w-full border px-3 py-1 rounded mb-2"
                    placeholder="Course Description"
                  />
                  <input
                    type="text"
                    value={editPrice}
                    onChange={(e) => setEditPrice(e.target.value)}
                    className="w-full border px-3 py-1 rounded mb-2"
                    placeholder="Price"
                  />
                  <div className="flex gap-2 mt-2">
                    <button
                      onClick={() => saveEdit(course.id)}
                      className="px-3 py-1 text-sm rounded bg-green-100 text-green-700 hover:bg-green-200"
                    >
                      Save
                    </button>
                    <button
                      onClick={() => setEditingCourseId(null)}
                      className="px-3 py-1 text-sm rounded bg-gray-100 text-gray-700 hover:bg-gray-200"
                    >
                      Cancel
                    </button>
                  </div>
                </>
              ) : (
                <>
                  <h2 className="text-xl font-semibold mb-1">{course.title}</h2>
                  <p className="text-gray-600 mb-2">{course.description}</p>
                  <p className="text-gray-800 font-medium mb-4">Price: {course.price}</p>

                  {/* Action buttons row */}
                  <div className="flex gap-2 mt-auto">
                    <button
                      onClick={() => startEditing(course)}
                      className="flex-1 px-3 py-1 text-sm rounded bg-green-100 text-green-700 hover:bg-green-200"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => deleteCourse(course.id)}
                      className="flex-1 px-3 py-1 text-sm rounded bg-red-100 text-red-700 hover:bg-red-200"
                    >
                      Delete
                    </button>
                    <button
                      onClick={() => toggleStatus(course.id)}
                      className={`relative w-14 h-7 rounded-full transition-colors duration-300 ${
                        course.status === "active" ? "bg-blue-600" : "bg-gray-300"
                      }`}
                    >
                      <span
                        className={`absolute top-0.5 left-0.5 w-6 h-6 bg-white rounded-full shadow transform transition-transform duration-300 ${
                          course.status === "active" ? "translate-x-7" : "translate-x-0"
                        }`}
                      />
                    </button>
                  </div>
                </>
              )}
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
