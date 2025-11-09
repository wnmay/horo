"use client";

import { useRouter } from "next/navigation";
import Card from "@/components/ui/card";

export default function CourseCard({ course }: any) {
  const router = useRouter();

  return (
    <div
      className="cursor-pointer hover:shadow-lg transition"
      onClick={() => router.push(`/courses/${course.id}`)}
    >
      <Card className="p-4 w-[240px] h-[280px] flex flex-col justify-between shadow-md bg-white dark:bg-zinc-900">
        <div className="flex flex-col flex-grow">
          <h3 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100 line-clamp-2">
            {course.title}
          </h3>
          <p className="text-sm text-zinc-600 dark:text-zinc-400 line-clamp-3 mt-2">
            {course.description}
          </p>
        </div>
        <div className="mt-4">
          <p className="text-sm font-medium text-zinc-700 dark:text-zinc-300">{course.prophet}</p>
          <p className="text-sm text-blue-600 font-semibold">{course.price}</p>
        </div>
      </Card>
    </div>
  );
}
