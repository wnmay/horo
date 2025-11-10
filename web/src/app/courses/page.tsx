"use client";

import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import CourseCard from "@/components/course-card";
import { useRouter } from "next/navigation";
import CourseFilter from "@/components/CourseFilter";
import api from "@/lib/api/api-client";

interface Course {
  id: string;
  title: string;
  description: string;
  prophet: string;
  price: string;
  duration: number;
  tags: string[];
}

export default function CoursesPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [sortOption, setSortOption] = useState<"title" | "price" | "none">("none");
  const [selectedTag, setSelectedTag] = useState<string | null>(null);
  const [courses, setCourses] = useState<Course[]>([]); // type-safe
  const router = useRouter();

  const courseTags = ["Love", "Study", "Work", "Health", "Finance", "Personal_Growth"];

  const tagImages: Record<string, string> = {
    Love: "/images/Love.jpg",
    Study: "/images/Study.jpg",
    Work: "/images/Work.jpg",
    Health: "/images/Health.jpg",
    Finance: "/images/Finance.jpg",
    Personal_Growth: "/images/Personal_Growth.jpg",
  };

  // --- Fetch courses ---
  useEffect(() => {
    const fetchCourses = async () => {
      try {
        const res = await api.get("/api/courses");
        if (Array.isArray(res.data.data)) {
          const mapped: Course[] = res.data.data.map((c: any) => ({
            id: c.id || "no-id",
            title: c.coursename || "Untitled",
            description: c.description || "",
            prophet: c.prophet || "Unknown",
            price: c.price != null ? `${c.price}` : "0",
            duration: c.duration != null ? c.duration : 0,
            tags: c.coursetype ? [c.coursetype] : [],
          }));
          setCourses(mapped);
        } else {
          setCourses([]);
        }
      } catch (err) {
        console.error("Error fetching courses:", err);
        setCourses([]);
      }
    };
    fetchCourses();
  }, []);

  // --- Filtering ---
  let filteredCourses = Array.isArray(courses) ? [...courses] : [];
  if (selectedTag) filteredCourses = filteredCourses.filter((c) => c.tags.includes(selectedTag));
  if (searchQuery.trim()) {
    const query = searchQuery.toLowerCase();
    filteredCourses = filteredCourses.filter(
      (c) => c.title.toLowerCase().includes(query) || c.prophet.toLowerCase().includes(query)
    );
  }
  if (sortOption === "title") filteredCourses.sort((a, b) => a.title.localeCompare(b.title));
  if (sortOption === "price") filteredCourses.sort(
    (a, b) => parseFloat(a.price) - parseFloat(b.price)
  );

  return (
    <div className="min-h-screen bg-gradient-to-b from-white to-zinc-100 dark:from-zinc-900 dark:to-zinc-950 py-12">
      <h1 className="text-4xl font-bold text-center text-zinc-800 dark:text-zinc-100 mb-8">
        All Courses
      </h1>

      <CourseFilter
        courseTags={courseTags}
        tagImages={tagImages}
        onApply={({ tag, search, sort }) => {
          setSelectedTag(tag);
          setSearchQuery(search);
          setSortOption(sort);
        }}
      />

      {/* Courses Grid */}
      <div className="max-w-6xl mx-auto mt-8 px-6 grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
        {filteredCourses.map((course) => (
          <CourseCard key={course.id} course={course} />
        ))}
      </div>

      {filteredCourses.length === 0 && (
        <p className="text-center text-zinc-500 mt-12">No courses found.</p>
      )}

      <div className="flex justify-center mt-12">
        <Button
          className="bg-blue-600 text-white hover:bg-blue-700 px-6 py-3 rounded-lg"
          onClick={() => router.push("/")}
        >
          ← Back to Home
        </Button>
      </div>
    </div>
  );
}
