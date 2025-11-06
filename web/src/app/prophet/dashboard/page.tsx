"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

interface Course {
  id: string;
  title: string;
  prophet: string;
}

export default function DashboardPage() {
  const [courses, setCourses] = useState<Course[]>([]);

  useEffect(() => {
    async function fetchCourses() {
      const res = await fetch("/api/my-courses");
      const data = await res.json();
      setCourses(data);
    }
    fetchCourses();
  }, []);

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">My Courses</h1>
      <ul className="flex flex-col gap-3">
        {courses.map((course) => (
          <li key={course.id} className="p-3 border rounded hover:bg-gray-50">
            <Link href={`/course/${course.id}`}>
              <div className="font-semibold">{course.title}</div>
              <div className="text-gray-500">By: {course.prophet}</div>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
