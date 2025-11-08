"use client";
import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { useState } from "react";
import CourseCard from "@/components/course-card";

export default function HomePage() {
  const router = useRouter();
  const [page, setPage] = useState(0); // page index (0 = first page)
  const COURSES_PER_PAGE = 5;

  // Mock course data (more than 5 now)
  const mockCourses = [
    {
      id: 1,
      title: "Beginner Astrology 101",
      description:
        "Learn the fundamentals of astrology, zodiac signs, and planetary movements.",
      prophet: "Prophet Orion",
      price: "$49",
    },
    {
      id: 2,
      title: "Advanced Horoscope Reading",
      description:
        "Master horoscope interpretations and learn how to create accurate readings.",
      prophet: "Prophet Selene",
      price: "$79",
    },
    {
      id: 3,
      title: "Tarot for Self Discovery",
      description:
        "Explore the meaning of tarot cards and how to use them for personal guidance.",
      prophet: "Prophet Lyra",
      price: "$39",
    },
    {
      id: 4,
      title: "Zodiac Compatibility Secrets",
      description:
        "Understand the relationship dynamics between zodiac signs.",
      prophet: "Prophet Atlas",
      price: "$59",
    },
    {
      id: 5,
      title: "Planetary Transits Explained",
      description:
        "Learn how planetary movements influence your daily horoscope.",
      prophet: "Prophet Vega",
      price: "$69",
    },
    {
      id: 6,
      title: "Astrology and Life Purpose",
      description:
        "Find your path through the stars with this transformative course.",
      prophet: "Prophet Nova",
      price: "$89",
    },
    {
      id: 7,
      title: "Dream Interpretation Basics",
      description:
        "Understand the spiritual meaning behind common dreams.",
      prophet: "Prophet Luna",
      price: "$45",
    },
    {
      id: 8,
      title: "Numerology for Beginners",
      description:
        "Discover how numbers shape your destiny and personality.",
      prophet: "Prophet Orion",
      price: "$55",
    },
  ];

  // Slice data for current page
  const startIndex = page * COURSES_PER_PAGE;
  const visibleCourses = mockCourses.slice(
    startIndex,
    startIndex + COURSES_PER_PAGE
  );

  // Pagination limits
  const hasNext = startIndex + COURSES_PER_PAGE < mockCourses.length;
  const hasPrev = page > 0;

  return (
    <div className="relative min-h-screen bg-gradient-to-b from-white to-zinc-100 dark:from-zinc-900 dark:to-zinc-950">

      {/* Fixed header buttons */}
      <div className="fixed top-4 right-4 p-4 flex gap-4 z-50">
        <Button
          className="bg-blue-500 text-white hover:bg-blue-600"
          onClick={() => router.push("/signin")}
        >
          Sign in
        </Button>
        <Button
          className="bg-blue-500 text-white hover:bg-blue-600"
          onClick={() => router.push("/signup")}
        >
          Sign up
        </Button>
      </div>

      {/* Welcome section */}
      <div className="flex min-h-[60vh] items-center justify-center px-4">
        <Card className="p-8 text-center">
          <h1 className="text-5xl font-bold text-black dark:text-zinc-50 mb-4">
            Welcome to Horo
          </h1>
          <p className="text-lg text-zinc-600 dark:text-zinc-400 max-w-md mx-auto">
            Your personalized horoscope app. Get started by signing in or creating an account.
          </p>
        </Card>
      </div>

      {/* Featured Courses */}
<div className="px-6 py-12 bg-white dark:bg-zinc-900">
  <h2 className="text-3xl font-semibold text-center text-zinc-800 dark:text-zinc-100 mb-8">
    Featured Courses
  </h2>

  {/* Horizontal scroll container */}
  <div className="max-w-6xl mx-auto overflow-x-auto">
    <div className="flex gap-6 min-w-max px-2">
      {visibleCourses.map((course) => (
        <div key={course.id} className="w-[240px] flex-shrink-0">
          <CourseCard course={course} />
        </div>
      ))}
    </div>
  </div>

  {/* Pagination Controls */}
  <div className="flex justify-center gap-4 mt-8">
    <Button
      variant="outline"
      disabled={!hasPrev}
      onClick={() => setPage((p) => p - 1)}
    >
      ← Previous
    </Button>
    <Button
      variant="outline"
      disabled={!hasNext}
      onClick={() => setPage((p) => p + 1)}
    >
      Next →
    </Button>
  </div>

  {/* Optional: page indicator */}
  <p className="text-center text-sm text-zinc-500 mt-4">
    Page {page + 1} of {Math.ceil(mockCourses.length / COURSES_PER_PAGE)}
  </p>
</div>

    </div>
  );
}
