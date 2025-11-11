"use client";

import { useRouter } from "next/navigation";
import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function CourseCard({ course }: any) {
  const router = useRouter();

  const rating = Number(course.reviewScore || course.review_score || 0).toFixed(
    1
  );

  return (
    <div className="cursor-pointer hover:shadow-lg transition">
      <Card className="relative p-4 w-[240px] h-[340px] flex flex-col justify-between shadow-md bg-white dark:bg-zinc-900">
        {/* Top Right Rating */}
        <div className="absolute top-3 right-3 text-yellow-500 text-sm font-semibold flex items-center">
          ⭐{" "}
          <span className="ml-1 text-zinc-800 dark:text-zinc-200">
            {rating}
          </span>
        </div>

        <div className="flex flex-col flex-grow">
          {/* Course Title */}
          <h3 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100 line-clamp-2 pr-8">
            {course.title}
          </h3>

          {/* Description */}
          <p className="text-sm text-zinc-600 dark:text-zinc-400 line-clamp-3 mt-2">
            {course.description}
          </p>

          {/* Type + Duration */}
          <div className="flex justify-between items-center text-xs mt-3 text-zinc-500 dark:text-zinc-400">
            {course.coursetype && (
              <span className="capitalize bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-200 px-2 py-1 rounded-md">
                {course.coursetype.replaceAll("_", " ")}
              </span>
            )}
            {course.duration && <span>{course.duration} min</span>}
          </div>
        </div>

        {/* Prophet */}
        <div className="mt-4">
          <p className="text-sm font-medium text-zinc-700 dark:text-zinc-300">
            by {course.prophet}
          </p>
        </div>

        {/* Price + Detail Button */}
        <div className="mt-3 flex justify-between items-center">
          <p className="text-base text-blue-600 font-semibold">
            {course.price}
          </p>
          <Button
            onClick={() => router.push(`/courses/${course.id}`)}
            className="text-sm px-3 py-1 bg-indigo-600 hover:bg-indigo-700 text-white rounded-md"
          >
            Detail →
          </Button>
        </div>
      </Card>
    </div>
  );
}
