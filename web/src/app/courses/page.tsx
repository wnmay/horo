"use client";
import { useState, useEffect } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import CourseCard from "@/components/CourseCard";
import api from "@/lib/api/api-client";
import { Button } from "@/components/ui/button";
import CourseSearchBar from "@/components/CourseSearchBar";
export default function CoursesPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [courses, setCourses] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchCourses() {
      setLoading(true);
      setError(null);

      try {
        const params = new URLSearchParams();
        const coursetype = searchParams.get("coursetype");
        const sortby = searchParams.get("sortby");
        const order = searchParams.get("order");
        const duration = searchParams.get("duration");
        const searchterm = searchParams.get("searchterm");

        if (duration) params.append("duration", duration);
        if (coursetype) params.append("coursetype", coursetype);
        if (order) params.append("order", order);
        if (sortby) params.append("sortby", sortby);
        if (searchterm) params.append("searchterm", searchterm);

        let url = "/api/courses";

        const queryString = params.toString();
        if (queryString) {
          url += `?${queryString}`;
        }

        const res = await api.get(url);
        setCourses(res.data.data || []);
      } catch (err) {
        console.error(err);
        setError("Failed to load courses");
      } finally {
        setLoading(false);
      }
    }

    fetchCourses();
  }, [searchParams.toString()]);

  if (loading) return <p className="text-center mt-20">Loading...</p>;
  if (error) return <p className="text-center text-red-500 mt-20">{error}</p>;

  return (
    <div className="min-h-screen bg-gradient-to-b from-white to-zinc-100 dark:from-zinc-900 dark:to-zinc-950 py-12">
      <h1 className="text-4xl font-bold text-center text-zinc-800 dark:text-zinc-100 mb-8">
        All Courses
      </h1>

      {/* ✅ Reusable Course Search Bar */}
      <div className="max-w-6xl mx-auto px-4 mb-10">
        <CourseSearchBar />
      </div>

      {/* ✅ Course Grid */}
      <div className="max-w-6xl mx-auto px-6 grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
        {courses.map((course) => (
          <CourseCard
            key={course.id}
            course={{
              id: course.id,
              title: course.coursename,
              description: course.description,
              prophet: course.prophetname,
              price: `$${course.price}`,
              reviewScore: course.review_score,
              reviewCount: course.review_count,
              duration: course.duration,
              coursetype: course.coursetype,
              prophetId: course.prophet_id,
              prophetName: course.prophetname,
              createdTime: course.created_time,
              deletedAt: course.deleted_at,
            }}
          />
        ))}
      </div>

      {courses.length === 0 && (
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
