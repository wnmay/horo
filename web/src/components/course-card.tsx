"use client";

import { useRouter } from "next/navigation";
import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface Course {
  id: string;
  title: string;
  description: string;
  prophet: string;
  price: string;
  duration: number;
  tags: string[];
}

export default function CourseCard({ course }: { course: Course }) {
  const router = useRouter();

  const goToDetail = (e?: React.MouseEvent) => {
    e?.stopPropagation(); // prevent bubbling
    router.push(`/courses/${course.id}`);
  };

  return (
    <div
      className="cursor-pointer hover:shadow-xl transition"
      onClick={goToDetail}
    >
      <Card className="p-4 w-[240px] h-[280px] flex flex-col justify-between shadow-md bg-white dark:bg-zinc-900 rounded-2xl transition">
        <div className="flex flex-col flex-grow">
          <h3 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100 line-clamp-2">
            {course.title || "Untitled"}
          </h3>
          <p className="text-sm text-zinc-600 dark:text-zinc-400 line-clamp-3 mt-2">
            {course.description || "No description"}
          </p>
        </div>

        {/* Bottom row: prophet, price, duration, see details */}
        <div className="mt-4 flex items-center justify-between gap-2">
          <div className="flex flex-col">
            <p className="text-sm font-medium text-zinc-700 dark:text-zinc-300">
              {course.prophet || "Unknown"}
            </p>
            <p className="text-sm text-blue-600 font-semibold">
              {course.price ?? "0"} ฿
            </p>
            <p className="text-sm text-zinc-500 dark:text-zinc-400">
              {course.duration ?? 0} min
            </p>
          </div>

          <Button
            variant="outline"
            size="sm"
            className="text-xs font-medium border-zinc-300 dark:border-zinc-700 bg-blue-600 hover:bg-blue-700 text-white transition"
            onClick={goToDetail}
          >
            See Details
          </Button>
        </div>
      </Card>
    </div>
  );
}
