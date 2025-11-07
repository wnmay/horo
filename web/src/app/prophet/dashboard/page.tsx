"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

interface Course {
  id: string;
  title: string;
  prophet: string;
  status: "active" | "inactive";
}

export default function DashboardPage() {
  const [courses, setCourses] = useState<Course[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [newTitle, setNewTitle] = useState("");

  useEffect(() => {
    async function fetchCourses() {
      const mockCourses: Course[] = [
        { id: "1", title: "How to Read the Stars", prophet: "Flook", status: "active" },
        { id: "2", title: "Dream Interpretation 101", prophet: "Flook", status: "active" },
        { id: "3", title: "Life Path Guidance", prophet: "Flook", status: "inactive" },
      ];
      await new Promise((resolve) => setTimeout(resolve, 300));
      setCourses(mockCourses);
    }
    fetchCourses();
  }, []);

  // handle adding new course
  const handleCreateCourse = () => {
    if (!newTitle.trim()) return;
    const newCourse: Course = {
      id: String(Date.now()),
      title: newTitle,
      prophet: "Flook",
      status: "active",
    };
    setCourses((prev) => [newCourse, ...prev]);
    setNewTitle("");
    setShowModal(false);
  };

  return (
    <div className="p-6">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">My Courses</h1>
        <button
          onClick={() => setShowModal(true)}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 shadow transition"
        >
          + Create Course
        </button>
      </div>

      {/* Course List */}
      <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {courses.map((course) => (
          <li
            key={course.id}
            className={`relative group rounded-xl border shadow-sm overflow-hidden transition-all duration-200 hover:shadow-lg hover:scale-[1.02] ${
              course.status === "inactive" ? "opacity-70" : ""
            }`}
          >
            <Link href={`/prophet/courses/${course.id}`} className="block h-full">
              <div className="absolute inset-0 bg-white group-hover:bg-gray-100 transition-colors duration-300" />
              <div className="relative p-5 flex flex-col justify-between h-full">
                <div>
                  <h2 className="text-xl font-semibold mb-1">{course.title}</h2>
                  <p className="text-gray-600 text-sm">By: {course.prophet}</p>
                </div>

                <div className="flex gap-2 mt-4">
                  <button className="px-3 py-1 text-sm rounded bg-green-100 text-green-700 hover:bg-green-200">
                    Edit
                  </button>
                  <button className="px-3 py-1 text-sm rounded bg-yellow-100 text-yellow-700 hover:bg-yellow-200">
                    {course.status === "active" ? "Deactivate" : "Activate"}
                  </button>
                  <button className="px-3 py-1 text-sm rounded bg-red-100 text-red-700 hover:bg-red-200">
                    Delete
                  </button>
                </div>
              </div>
            </Link>
          </li>
        ))}
      </ul>

      {/* Modal for creating new course */}
      {showModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md shadow-xl relative">
            <h2 className="text-xl font-semibold mb-4">Create New Course</h2>
            <input
              type="text"
              placeholder="Enter course title..."
              value={newTitle}
              onChange={(e) => setNewTitle(e.target.value)}
              className="w-full border rounded-lg p-2 mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <div className="flex justify-end gap-2">
              <button
                onClick={() => setShowModal(false)}
                className="px-4 py-2 bg-gray-200 rounded-lg hover:bg-gray-300"
              >
                Cancel
              </button>
              <button
                onClick={handleCreateCourse}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
              >
                Create
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
