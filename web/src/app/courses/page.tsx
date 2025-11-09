"use client";
import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import CourseCard from "@/components/course-card";
import { useRouter } from "next/navigation";

export default function CoursesPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [sortOption, setSortOption] = useState<"title" | "price" | "none">("none");
  const [selectedTag, setSelectedTag] = useState<string | null>(null);
  const [tempSearch, setTempSearch] = useState("");
  const [tempTag, setTempTag] = useState<string | null>(null);
  const [tempSort, setTempSort] = useState<"title" | "price" | "none">("none");

  const courseTags = ["Love", "Study", "Work", "Health", "Finance", "Personal_Growth"];
  const router = useRouter()

  // Mock data (same as homepage)
  const mockCourses = [
    { id: 1, title: "Beginner Astrology 101", description: "Learn astrology basics.", prophet: "Prophet Orion", price: "$49", tags: ["Study"] },
    { id: 2, title: "Love Compatibility Reading", description: "Find your perfect match.", prophet: "Prophet Luna", price: "$59", tags: ["Love"] },
    { id: 3, title: "Work-Life Balance Guidance", description: "Career and destiny reading.", prophet: "Prophet Selene", price: "$69", tags: ["Work", "Health"] },
    { id: 4, title: "Financial Fortune Reading", description: "Predict your wealth path.", prophet: "Prophet Nova", price: "$79", tags: ["Finance"] },
    { id: 5, title: "Healing Energy Workshop", description: "Restore your inner energy.", prophet: "Prophet Vega", price: "$89", tags: ["Health", "Personal Growth"] },
    { id: 6, title: "Manifest Love & Success", description: "Use the universe to attract what you desire.", prophet: "Prophet Atlas", price: "$99", tags: ["Love", "Personal Growth"] },
  ];

  const tagImages: Record<string, string> = {
    Love: "/images/Love.jpg",
    Study: "/images/Study.jpg",
    Work: "/images/Work.jpg",
    Health: "/images/Health.jpg",
    Finance: "/images/Finance.jpg",
    Personal_Growth: "/images/Personal_Growth.jpg",
  };

  // --- Filtering logic ---
  let filteredCourses = mockCourses;
  if (selectedTag) {
    filteredCourses = filteredCourses.filter((c) => c.tags.includes(selectedTag));
  }
  if (searchQuery.trim()) {
    const query = searchQuery.toLowerCase();
    filteredCourses = filteredCourses.filter(
      (c) => c.title.toLowerCase().includes(query) || c.prophet.toLowerCase().includes(query)
    );
  }
  if (sortOption === "title") {
    filteredCourses = [...filteredCourses].sort((a, b) => a.title.localeCompare(b.title));
  } else if (sortOption === "price") {
    filteredCourses = [...filteredCourses].sort(
      (a, b) => parseFloat(a.price.slice(1)) - parseFloat(b.price.slice(1))
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-white to-zinc-100 dark:from-zinc-900 dark:to-zinc-950 py-12">
      <h1 className="text-4xl font-bold text-center text-zinc-800 dark:text-zinc-100 mb-8">
        All Courses
      </h1>

      {/* Tag Filter */}
      <div className="max-w-6xl mx-auto mb-10 px-4 text-center">
        <h2 className="text-2xl font-semibold mb-6 text-zinc-800 dark:text-zinc-100">
          Filter by Category
        </h2>
        <div className="flex flex-wrap justify-center gap-6">
          {courseTags.map((tag) => (
            <button
              key={tag}
              onClick={() => setTempTag(tempTag === tag ? null : tag)}
              className={`flex flex-col items-center border rounded-2xl p-4 w-[110px] text-sm font-medium shadow-sm transition-transform
                ${tempTag === tag
                  ? "bg-blue-600 text-white border-blue-700 scale-105"
                  : "bg-white dark:bg-zinc-800 hover:bg-zinc-100 dark:hover:bg-zinc-700 border-zinc-300 dark:border-zinc-600"
                }`}
            >
                <img
                src={tagImages[tag]}
                alt={tag}
                className="w-10 h-10 mb-2 opacity-90"
                />
              {tag.replaceAll("_"," ")}
            </button>
          ))}
        </div>
      </div>

      {/* Search + Sort + Apply */}
      <div className="max-w-6xl mx-auto px-4 mb-10">
        <div className="flex flex-col md:flex-row justify-center md:justify-between items-center gap-2 bg-white dark:bg-zinc-800 p-4 rounded-xl shadow max-w-lg mx-auto">
          <input
            type="text"
            placeholder="Search courses or prophets..."
            value={tempSearch}
            onChange={(e) => setTempSearch(e.target.value)}
            className="w-full md:w-[50%] border border-zinc-300 dark:border-zinc-600 rounded-lg px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none dark:bg-zinc-900"
          />

          <select
            value={tempSort}
            onChange={(e) => setTempSort(e.target.value as "title" | "price" | "none")}
            className="w-full md:w-[25%] border border-zinc-300 dark:border-zinc-600 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none dark:bg-zinc-900"
          >
            <option value="none">Sort by</option>
            <option value="title">Title (A–Z)</option>
            <option value="price">Price (Low → High)</option>
          </select>

          <Button
            onClick={() => {
              setSearchQuery(tempSearch);
              setSelectedTag(tempTag);
              setSortOption(tempSort);
            }}
            className="bg-blue-600 text-white hover:bg-blue-700 px-6 py-2 rounded-lg"
          >
            Apply
          </Button>
        </div>
      </div>

      {/* All Courses */}
      <div className="max-w-6xl mx-auto px-6 grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
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
