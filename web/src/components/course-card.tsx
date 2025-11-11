"use client";

import { useRouter } from "next/navigation";
import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface CourseCardProps {
  course: {
    id: string;
    title: string;
    description: string;
    prophet: string;
    price: string | number;
    duration: number;
    tags: string[];
  };
}

export default function CourseCard({ course }: CourseCardProps) {
  const router = useRouter();

  const goToDetail = (e?: React.MouseEvent) => {
    e?.stopPropagation();
    router.push(`/courses/${course.id}`);
  };

  return (
    <div
      className="cursor-pointer hover:shadow-xl transition"
      onClick={goToDetail}
    >
      <Card className="p-4 w-[240px] h-[280px] flex flex-col justify-between shadow-md bg-white dark:bg-zinc-900 rounded-2xl transition">
        {/* Top: Title + Description */}
        <div className="flex flex-col flex-grow">
          <h3 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100 line-clamp-2">
            {course.title}
          </h3>
          <p className="text-sm text-zinc-600 dark:text-zinc-400 line-clamp-3 mt-2">
            {course.description}
          </p>
        </div>

        {/* Bottom row: Prophet | Price | Duration | See Details */}
        <div className="mt-4 flex items-center justify-between gap-2">
          <div className="flex flex-col">
            <p className="text-sm font-medium text-zinc-700 dark:text-zinc-300">
              {course.prophet}
            </p>
            <p className="text-sm text-blue-600 font-semibold">
              {course.price} ฿
            </p>
            <p className="text-sm text-zinc-500 dark:text-zinc-400">
              {course.duration} min
            </p>
          </div>

          <Button
            variant="outline"
            size="sm"
            className="text-xs font-medium border-zinc-300 dark:border-zinc-700 bg-blue-600 hover:bg-blue-700 text-white transition flex-shrink-0"
            onClick={goToDetail}
          >
            See Details
          </Button>
        </div>
      </Card>
    </div>
  );
}
