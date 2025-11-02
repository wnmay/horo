"use client"
import { Button } from "@/components/ui/button";
import Card from "@/components/ui/card";

export default function CourseCard({ course }: { course: any }) {
  return (
    <Card className="p-6 flex flex-col justify-between w-[240px] h-[260px] shadow-md">
      <div className="flex flex-col flex-grow">
        <h3 className="text-lg font-bold text-zinc-900 dark:text-zinc-100 mb-2 line-clamp-2">
          {course.title}
        </h3>
        <p className="text-sm text-zinc-600 dark:text-zinc-400 mb-4 line-clamp-3">
          {course.description}
        </p>
      </div>
      <div className="flex items-center justify-between mt-auto">
        <div className="text-left">
          <p className="text-sm font-medium text-zinc-700 dark:text-zinc-300">
            {course.prophet}
          </p>
          <p className="text-sm text-zinc-500">{course.price}</p>
        </div>
        <Button
          className="bg-blue-500 text-white hover:bg-blue-600 text-sm px-3 py-1"
          onClick={() => alert(`Viewing details for: ${course.title}`)}
        >
          View
        </Button>
      </div>
    </Card>
  );
}
