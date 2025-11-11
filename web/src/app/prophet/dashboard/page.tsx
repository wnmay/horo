"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { Plus } from "lucide-react";
import api from "@/lib/api/api-client";
import { auth } from "@/firebase/firebase";
import CourseCard from "@/components/CourseEditCard";
import CourseEditModal from "@/components/CourseEditModal";

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

export default function ProphetCoursesPage() {
  const [courses, setCourses] = useState<Course[]>([]);
  const [loading, setLoading] = useState(true);
  const [editingCourse, setEditingCourse] = useState<Course | null>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const router = useRouter();

  useEffect(() => {
    async function loadCourses() {
      try {
        if (!auth.currentUser) {
          await new Promise<void>((resolve) => {
            const unsubscribe = auth.onAuthStateChanged((user) => {
              if (user) {
                unsubscribe();
                resolve();
              }
            });
            setTimeout(() => {
              unsubscribe();
              resolve();
            }, 1000);
          });
        }

        const res = await api.get("/api/courses/prophet/courses");
        const data = res.data.data;
        setCourses(data);
      } catch (err) {
        toast("Failed to load prophet courses", {
          description: (err as Error).message,
          duration: 5000,
          position: "top-right",
        });
      } finally {
        setLoading(false);
      }
    }
    loadCourses();
  }, []);

  async function handleSaveEdit() {
    if (!editingCourse) return;

    try {
      setSaving(true);

      const res = await api.patch(
        `/api/courses/${editingCourse.id}`,
        editingCourse
      );

      if (res.status !== 200) {
        throw new Error(res.data.message || "Failed to update course");
      }

      // Optimistically update the UI
      setCourses((prev) =>
        prev.map((c) => (c.id === editingCourse.id ? editingCourse : c))
      );

      toast.success("Course updated successfully!", { position: "top-right" });
      setIsDialogOpen(false);
    } catch (err) {
      toast.error("Failed to update course", {
        description: (err as Error).message,
        position: "top-right",
      });
    } finally {
      setSaving(false);
    }
  }

  function handleEditCourse(course: Course) {
    setEditingCourse(course);
    setIsDialogOpen(true);
  }

  function handleCloseModal() {
    setIsDialogOpen(false);
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-screen text-lg text-blue-600 dark:text-blue-300 animate-pulse">
        Loading courses...
      </div>
    );
  }

  return (
    <div className="min-h-screen px-6 py-12">
      <div className="max-w-6xl mx-auto flex items-center justify-between mb-10">
        <h1 className="text-4xl font-extrabold text-blue-800 dark:text-blue-300 tracking-tight">
          Prophet Courses
        </h1>
        <Button
          onClick={() => router.push("/prophet/create-course")}
          className="flex items-center gap-2 bg-blue-600 text-white hover:bg-blue-700 hover:shadow-lg hover:-translate-y-[1px] transition-all duration-200 rounded-lg px-4 py-2"
        >
          <Plus className="w-4 h-4" />
          Create Course
        </Button>
      </div>

      {courses.length === 0 ? (
        <p className="text-center text-zinc-500 dark:text-zinc-400">
          No courses available.
        </p>
      ) : (
        <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-3 max-w-6xl mx-auto">
          {courses.map((course) => (
            <CourseCard
              key={course.id}
              course={course}
              onEdit={handleEditCourse}
            />
          ))}
        </div>
      )}

      <CourseEditModal
        isOpen={isDialogOpen}
        onClose={handleCloseModal}
        course={editingCourse}
        onCourseChange={setEditingCourse}
        onSave={handleSaveEdit}
        isSaving={saving}
      />
    </div>
  );
}
