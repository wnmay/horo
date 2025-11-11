"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import Image from "next/image";
import { Button } from "@/components/ui/button";
import {
  courseTypesMap,
  courseTypes,
  courseTypeImageMap,
} from "@/types/course-type";

interface CourseSearchBarProps {
  initialSearch?: string;
  initialSort?: "review_score" | "price" | "none";
  initialCourseType?: string | null;
}

export default function CourseSearchBar({
  initialSearch = "",
  initialSort = "none",
  initialCourseType = null,
}: CourseSearchBarProps) {
  const router = useRouter();

  const [tempSearch, setTempSearch] = useState(initialSearch);
  const [tempSort, setTempSort] = useState<"review_score" | "price" | "none">(
    initialSort
  );
  const [selectedCourseType, setSelectedCourseType] = useState<string | null>(
    initialCourseType
  );

  const handleApply = () => {
    const params = new URLSearchParams();

    if (tempSearch) params.append("searchterm", tempSearch);

    if (selectedCourseType) {
      const courseSlug =
        courseTypesMap[selectedCourseType] || selectedCourseType;
      params.append("coursetype", courseSlug);
    }

    if (tempSort !== "none") {
      params.append("sortby", tempSort);
      params.append("order", "asc");
    }

    router.push(`/courses?${params.toString()}`);
  };

  return (
    <div className="space-y-10">
      {/* Tag Filter */}
      <section className="text-center">
        <h2 className="text-2xl font-semibold mb-6 text-zinc-800 dark:text-zinc-100">
          Explore by Category
        </h2>
        <div className="flex flex-wrap justify-center gap-6">
          {courseTypes.map((courseType: string) => (
            <button
              key={courseType}
              onClick={() =>
                setSelectedCourseType(
                  selectedCourseType === courseType ? null : courseType
                )
              }
              className={`flex flex-col items-center border rounded-2xl p-4 w-[110px] text-sm font-medium shadow-sm transition-transform
                ${
                  selectedCourseType === courseType
                    ? "bg-indigo-600 text-white border-indigo-700 scale-105"
                    : "bg-white dark:bg-zinc-800 hover:bg-zinc-100 dark:hover:bg-zinc-700 border-zinc-300 dark:border-zinc-600"
                }`}
            >
              <div className="w-10 h-10 mb-2 relative">
                <Image
                  src={
                    courseTypeImageMap[
                      courseType as keyof typeof courseTypeImageMap
                    ]
                  }
                  alt={courseType}
                  fill
                  sizes="48px"
                  className="object-contain rounded-md"
                  priority
                />
              </div>
              {courseType}
            </button>
          ))}
        </div>
      </section>

      {/* Search + Sort + Apply */}
      <section>
        <div className="flex flex-col md:flex-row justify-center md:justify-between items-center gap-3 bg-white dark:bg-zinc-800 p-4 rounded-xl shadow-md max-w-3xl mx-auto">
          <input
            type="text"
            placeholder="Search courses or prophets..."
            value={tempSearch}
            onChange={(e) => setTempSearch(e.target.value)}
            className="w-full md:w-[50%] border border-zinc-300 dark:border-zinc-600 rounded-lg px-4 py-2 focus:ring-2 focus:ring-indigo-500 focus:outline-none dark:bg-zinc-900"
          />

          <select
            value={tempSort}
            onChange={(e) =>
              setTempSort(e.target.value as "review_score" | "price" | "none")
            }
            className="w-full md:w-[25%] border border-zinc-300 dark:border-zinc-600 rounded-lg px-3 py-2 focus:ring-2 focus:ring-indigo-500 focus:outline-none dark:bg-zinc-900"
          >
            <option value="none">Sort by</option>
            <option value="review_score">Review Score (High → Low)</option>
            <option value="price">Price (Low → High)</option>
          </select>

          <Button
            onClick={handleApply}
            className="bg-blue-600 text-white hover:bg-blue-700 px-6 py-2 rounded-lg"
          >
            Apply
          </Button>
        </div>
      </section>
    </div>
  );
}
