"use client";

import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { useState, useEffect, useRef } from "react";
import CourseCard from "@/components/CourseCard";
import api from "@/lib/api/api-client";
import CourseSearchBar from "@/components/CourseSearchBar";
export default function HomePage() {
  const router = useRouter();
  const menuRef = useRef<HTMLDivElement>(null);
  const [adIndex, setAdIndex] = useState(0);
  const [courses, setCourses] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // --- Ads ---
  const ads = [
    {
      id: 0,
      title: "Welcome to Horo",
      description: "Your personalized horoscope app.",
    },
    {
      id: 1,
      title: "Get a free horoscope reading!",
      description: "Sign up now for your daily horoscope.",
    },
    {
      id: 2,
      title: "Exclusive Tarot Tips",
      description: "Learn top secrets for interpreting tarot.",
    },
  ];
  useEffect(() => {
    const interval = setInterval(
      () => setAdIndex((prev) => (prev + 1) % ads.length),
      3000
    );
    return () => clearInterval(interval);
  }, []);

  // --- Fetch courses ---
  useEffect(() => {
    async function fetchCourses() {
      try {
        const res = await api.get("/api/courses/popular?limit=10");
        setCourses(res.data.data || []);
      } catch (err: any) {
        console.error(err);
        setError("Failed to load courses");
      } finally {
        setLoading(false);
      }
    }
    fetchCourses();
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-b from-white to-zinc-100 dark:from-zinc-900 dark:to-zinc-950">
      <div className="max-w-6xl mx-auto px-4 py-10 space-y-12">
        {/* Hero / Ads */}
        <section className="flex justify-center">
          <Card className="transition-all duration-500 w-full md:w-[650px] p-10 text-center bg-yellow-50 dark:bg-yellow-900 text-zinc-800 dark:text-yellow-100 rounded-2xl shadow-lg">
            <h3 className="font-bold text-3xl md:text-4xl mb-3">
              {ads[adIndex].title}
            </h3>
            <p className="text-base md:text-lg mt-1">
              {ads[adIndex].description}
            </p>
          </Card>
        </section>

        {/* Search + Filter + Sort */}
        <CourseSearchBar />
      </div>

      {/* Featured Courses */}
      <section className="px-6 py-12 bg-white dark:bg-zinc-900 mt-4">
        <h2 className="text-3xl font-semibold text-center text-zinc-800 dark:text-zinc-100 mb-8">
          Featured Courses
        </h2>

        {loading ? (
          <p className="text-center text-zinc-500">Loading courses...</p>
        ) : error ? (
          <p className="text-center text-red-500">{error}</p>
        ) : (
          <>
            <div className="max-w-6xl mx-auto overflow-x-auto scrollbar-thin scrollbar-thumb-zinc-400 dark:scrollbar-thumb-zinc-600 scrollbar-track-transparent">
              <div className="flex gap-6 min-w-max px-2 pb-4">
                {courses.map((course) => (
                  <div key={course.id} className="w-[240px] flex-shrink-0">
                    <CourseCard
                      course={{
                        title: course.coursename,
                        description: course.description,
                        prophet: course.prophetname,
                        price: `$${course.price}`,
                      }}
                    />
                  </div>
                ))}
              </div>
            </div>

            <div className="flex justify-center mt-8">
              <Button
                className="bg-blue-600 text-white hover:bg-blue-700 px-6 py-3 rounded-lg"
                onClick={() => router.push("/courses")}
              >
                See More Courses →
              </Button>
            </div>
          </>
        )}
      </section>
    </div>
  );
}
