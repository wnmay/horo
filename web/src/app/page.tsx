"use client";

import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useState, useEffect } from "react";
import CourseCard from "@/components/course-card";
import CourseFilter from "../components/CourseFilter";
import api from "@/lib/api/api-client";

export default function HomePage() {
  const [page, setPage] = useState(0);
  const [adIndex, setAdIndex] = useState(0);
  const [coursesPerPage, setCoursesPerPage] = useState(6);

  // Filters
  const [selectedTag, setSelectedTag] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [sortOption, setSortOption] = useState<"title" | "price" | "none">("none");
  const [durationFilter, setDurationFilter] = useState<number | null>(null);

  // Courses
  const [courses, setCourses] = useState<any[]>([]);
  const [filteredCourses, setFilteredCourses] = useState<any[]>([]);

  // --- Fetch courses from backend ---
  useEffect(() => {
    const fetchCourses = async () => {
      try {
        const res = await api.get("/api/courses");
        if (Array.isArray(res.data.data)) {
          const mapped = res.data.data.map((c: any) => ({
            id: c.id,
            title: c.coursename,
            description: c.description,
            prophet: c.prophet || "Unknown",
            price: c.price ? c.price : 0,
            duration: c.duration || 0,
            tags: c.coursetype ? [c.coursetype] : [],
          }));
          setCourses(mapped);
          setFilteredCourses(mapped); // initially show all
        } else {
          setCourses([]);
          setFilteredCourses([]);
        }
      } catch (err) {
        console.error("Error fetching courses:", err);
        setCourses([]);
        setFilteredCourses([]);
      }
    };
    fetchCourses();
  }, []);

  // --- Ads ---
  const ads = [
    { id: 0, title: "Welcome to Horo", description: "Your personalized horoscope app." },
    { id: 1, title: "Get a free horoscope reading!", description: "Sign up now for your daily horoscope." },
    { id: 2, title: "Exclusive Tarot Tips", description: "Learn top secrets for interpreting tarot." },
  ];
  useEffect(() => {
    const interval = setInterval(() => setAdIndex((prev) => (prev + 1) % ads.length), 3000);
    return () => clearInterval(interval);
  }, []);

  // --- Apply Handler ---
  const handleApply = () => {
    let filtered = [...courses];

    if (selectedTag) filtered = filtered.filter((c) => c.tags.includes(selectedTag));
    if (searchQuery.trim()) {
      const q = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (c) => c.title.toLowerCase().includes(q) || c.prophet.toLowerCase().includes(q)
      );
    }
    if (durationFilter !== null) filtered = filtered.filter((c) => c.duration === durationFilter);
    if (sortOption === "title") filtered.sort((a, b) => a.title.localeCompare(b.title));
    if (sortOption === "price") filtered.sort((a, b) => a.price - b.price);

    setFilteredCourses(filtered);
    setPage(0); // reset pagination
  };

  // --- Handle Enter key ---
  useEffect(() => {
    const handleKeyPress = (e: KeyboardEvent) => {
      if (e.key === "Enter") {
        e.preventDefault();
        handleApply();
      }
    };
    window.addEventListener("keydown", handleKeyPress);
    return () => window.removeEventListener("keydown", handleKeyPress);
  }, [searchQuery, selectedTag, durationFilter, sortOption, courses]);

  // --- Pagination ---
  const startIndex = page * coursesPerPage;
  const visibleCourses = filteredCourses.slice(startIndex, startIndex + coursesPerPage);

  return (
    <div className="min-h-screen bg-gradient-to-b from-white to-zinc-100 dark:from-zinc-900 dark:to-zinc-950">
      {/* Carousel */}
      <div className="max-w-6xl mx-auto px-4 mb-12">
        <Card className="transition-all duration-500 mx-auto w-full md:w-[600px] p-10 text-center bg-yellow-50 dark:bg-yellow-900 text-zinc-800 dark:text-yellow-100 min-h-[250px] flex flex-col justify-center">
          <h3 className="font-bold text-4xl md:text-5xl mb-4">{ads[adIndex].title}</h3>
          <p className="text-lg md:text-xl mt-2">{ads[adIndex].description}</p>
        </Card>
      </div>

      {/* Filter */}
      <CourseFilter
        courseTags={["Love", "Study", "Work", "Health", "Finance", "Personal_Growth"]}
        tagImages={{
          Love: "/images/Love.jpg",
          Study: "/images/Study.jpg",
          Work: "/images/Work.jpg",
          Health: "/images/Health.jpg",
          Finance: "/images/Finance.jpg",
          Personal_Growth: "/images/Personal_Growth.jpg",
        }}
        selectedTag={selectedTag}
        searchQuery={searchQuery}
        sortOption={sortOption}
        durationFilter={durationFilter}
        onTagChange={setSelectedTag}
        onSearchChange={setSearchQuery}
        onSortChange={setSortOption}
        onDurationChange={setDurationFilter}
        onApply={handleApply}
      />

      {/* Suggested Courses */}
      <div className="px-6 py-12 to-zinc-100 dark:bg-zinc-900">
        <h2 className="text-3xl font-semibold text-center text-zinc-800 dark:text-zinc-100 mb-8">
          Suggested Courses
        </h2>

        <div className="max-w-full mx-auto overflow-x-auto">
          <div className="flex gap-6 min-w-max px-2 justify-center">
            {visibleCourses.map((course) => (
              <CourseCard key={course.id} course={course} />
            ))}
          </div>
        </div>

        {/* Pagination Buttons */}
        <div className="flex items-center justify-between mt-8 max-w-6xl mx-auto px-2">
          <div className="w-1/3"></div>

          <div className="w-1/3 flex justify-end">
            <Button
              className="bg-blue-600 text-white hover:bg-blue-700 px-6 py-3 rounded-lg"
              onClick={handleApply}
            >
              See More Courses →
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
