import { Button } from "@/components/ui/button";
import Card from "@/components/ui/card";
import { Pencil } from "lucide-react";

interface Course {
  id: string;
  prophet_id: string;
  prophetname: string;
  coursename: string;
  coursetype: string;
  description: string;
  price: number;
  duration: number;
  created_time: string;
  deleted_at: boolean;
}

interface CourseCardProps {
  course: Course;
  onEdit: (course: Course) => void;
}

export default function CourseCard({ course, onEdit }: CourseCardProps) {
  return (
    <Card
      key={course.coursename}
      className="p-6 flex flex-col justify-between rounded-2xl shadow-lg border border-blue-100 dark:border-zinc-800 bg-white/90 dark:bg-zinc-900/70 backdrop-blur-md hover:shadow-blue-200/60 dark:hover:shadow-blue-900/50"
    >
      <div>
        <h2 className="text-2xl font-semibold text-blue-800 dark:text-blue-300 mb-3">
          {course.coursename}
        </h2>
        <p className="text-zinc-600 dark:text-zinc-400 mb-5 leading-relaxed">
          {course.description}
        </p>

        <div className="space-y-2 text-sm text-zinc-600 dark:text-zinc-400">
          <p>
            <span className="font-medium text-zinc-800 dark:text-zinc-200">
              Course Type:
            </span>{" "}
            {course.coursetype}
          </p>
          <p>
            <span className="font-medium text-zinc-800 dark:text-zinc-200">
              Duration:
            </span>{" "}
            {course.duration} min
          </p>
          <p>
            <span className="font-medium text-zinc-800 dark:text-zinc-200">
              Price:
            </span>{" "}
            ${course.price}
          </p>
        </div>
      </div>

      <Button
        className="mt-6 bg-blue-500 text-white hover:bg-blue-600 hover:shadow-md flex items-center justify-center gap-2 rounded-md transition-all duration-200"
        onClick={() => onEdit(course)}
      >
        <Pencil className="w-4 h-4" />
        Edit
      </Button>
    </Card>
  );
}
