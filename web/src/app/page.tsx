"use client";
import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import CourseCard from "@/components/course-card";
import CourseFilter from "../components/CourseFilter";
import api from "@/lib/api/api-client";

export default function HomePage() {
  const router = useRouter();
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

  // --- Fetch courses from backend ---
  useEffect(() => {
    const fetchCourses = async () => {
      try {
        const params = new URLSearchParams();
        if (searchQuery.trim()) {
          params.append("coursename", searchQuery);
          params.append("prophetname", searchQuery);
        }
        const res = await api.get(`/api/courses?${params.toString()}`);
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
        } else {
          setCourses([]);
        }
      } catch (err) {
        console.error("Error fetching courses:", err);
        setCourses([]);
      }
    };
    fetchCourses();
  }, [searchQuery]);

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

  // --- Filter courses (frontend) ---
  let filteredCourses = [...courses];
  if (selectedTag) filteredCourses = filteredCourses.filter((c) => c.tags.includes(selectedTag));
  if (durationFilter !== null) filteredCourses = filteredCourses.filter((c) => c.duration === durationFilter);
  if (sortOption === "title") filteredCourses.sort((a, b) => a.title.localeCompare(b.title));
  if (sortOption === "price") filteredCourses.sort((a, b) => a.price - b.price);

  // --- Pagination ---
  const startIndex = page * coursesPerPage;
  const visibleCourses = filteredCourses.slice(startIndex, startIndex + coursesPerPage);
  const hasNext = startIndex + coursesPerPage < filteredCourses.length;
  const hasPrev = page > 0;

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
        onApply={({ tag, search, sort, duration }) => {
          setSelectedTag(tag);
          setSearchQuery(search);
          setSortOption(sort);
          setDurationFilter(duration);
          setPage(0);
        }}
      />

      {/* Courses */}
      <div className="px-6 py-12 to-zinc-100 dark:bg-zinc-900">
        <h2 className="text-3xl font-semibold text-center text-zinc-800 dark:text-zinc-100 mb-8">
          {selectedTag ? `${selectedTag} Courses` : "Featured Courses"}
        </h2>

        <div className="max-w-full mx-auto overflow-x-auto">
          <div className="flex gap-6 min-w-max px-2 justify-center">
            {visibleCourses.map((course) => (
              <CourseCard key={course.id} course={course} />
            ))}
          </div>
        </div>

        {/* Pagination */}
        <div className="flex items-center justify-between mt-8 max-w-6xl mx-auto px-2">
          <div className="w-1/3"></div>
          <div className="flex gap-4 justify-center w-1/3">
            <Button variant="outline" disabled={!hasPrev} onClick={() => setPage(p => p - 1)}>← Previous</Button>
            <Button variant="outline" disabled={!hasNext} onClick={() => setPage(p => p + 1)}>Next →</Button>
          </div>
          <div className="w-1/3 flex justify-end">
            <Button className="bg-blue-600 text-white hover:bg-blue-700 px-6 py-3 rounded-lg" onClick={() => router.push("/courses")}>
              See More Courses →
            </Button>
          </div>
        </div>

        <p className="text-center text-sm text-zinc-500 mt-4">
          Page {page + 1} of {Math.ceil(filteredCourses.length / coursesPerPage)}
        </p>
      </div>
    </div>
  );
}
