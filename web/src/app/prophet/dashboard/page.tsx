"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";
import { toast } from "sonner";
import { Pencil, Plus } from "lucide-react";

// Mock fetch function
async function fetchProphetCourses(token: string) {
  await new Promise((r) => setTimeout(r, 500)); // simulate network delay
  return [
    {
      courseId: "course-1",
      prophetId: "prophet-1",
      prophetName: "Prophet Name 1",
      courseName: "Course Name 8",
      description: "An in-depth course exploring advanced spiritual practices.",
      price: 750,
      duration: 60,
    },
  ];
}

interface Course {
  courseId: string;
  prophetId: string;
  prophetName: string;
  courseName: string;
  description: string;
  price: number;
  duration: number;
}

export default function ProphetCoursesPage() {
  const [courses, setCourses] = useState<Course[]>([]);
  const [loading, setLoading] = useState(true);
  const [editingCourse, setEditingCourse] = useState<Course | null>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const router = useRouter();

  const token = "auth-token";

  useEffect(() => {
    async function loadCourses() {
      try {
        const data = await fetchProphetCourses(token);
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

      const res = await fetch(`/api/courses/${editingCourse.courseId}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          coursename: editingCourse.courseName,
          description: editingCourse.description,
          price: Number(editingCourse.price),
          duration: Number(editingCourse.duration),
          prophetname: editingCourse.prophetName,
        }),
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || "Failed to update course");
      }

      // Optimistically update the UI
      setCourses((prev) =>
        prev.map((c) =>
          c.prophetId === editingCourse.prophetId ? editingCourse : c
        )
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
            <Card
              key={course.courseName}
              className="p-6 flex flex-col justify-between rounded-2xl shadow-lg border border-blue-100 dark:border-zinc-800 bg-white/90 dark:bg-zinc-900/70 backdrop-blur-md hover:shadow-blue-200/60 dark:hover:shadow-blue-900/50"
            >
              <div>
                <h2 className="text-2xl font-semibold text-blue-800 dark:text-blue-300 mb-3">
                  {course.courseName}
                </h2>
                <p className="text-zinc-600 dark:text-zinc-400 mb-5 leading-relaxed">
                  {course.description}
                </p>

                <div className="space-y-2 text-sm text-zinc-600 dark:text-zinc-400">
                  <p>
                    <span className="font-medium text-zinc-800 dark:text-zinc-200">
                      Prophet:
                    </span>{" "}
                    {course.prophetName}
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
                onClick={() => {
                  setEditingCourse(course);
                  setIsDialogOpen(true);
                }}
              >
                <Pencil className="w-4 h-4" />
                Edit
              </Button>
            </Card>
          ))}
        </div>
      )}

      {/* --- Edit Modal --- */}
      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent
          className="sm:max-w-lg bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl shadow-xl
  data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:slide-in-from-bottom-2 data-[state=open]:zoom-in-95
  data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:slide-out-to-bottom-2 data-[state=closed]:zoom-out-95
  duration-300 ease-[cubic-bezier(0.22,1,0.36,1)]"
        >
          <DialogHeader>
            <DialogTitle className="text-blue-700 dark:text-blue-300 text-xl font-semibold">
              Edit Course
            </DialogTitle>
          </DialogHeader>

          {editingCourse && (
            <div className="space-y-4 mt-4">
              {/* Course Name */}
              <div>
                <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">
                  Course Name
                </label>
                <Input
                  value={editingCourse.courseName}
                  onChange={(e) =>
                    setEditingCourse({
                      ...editingCourse,
                      courseName: e.target.value,
                    })
                  }
                  className="bg-white dark:bg-zinc-900 border border-zinc-300 dark:border-zinc-700 focus:ring-2 focus:ring-blue-400"
                />
              </div>

              {/* Description */}
              <div>
                <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">
                  Description
                </label>
                <Textarea
                  value={editingCourse.description}
                  onChange={(e) =>
                    setEditingCourse({
                      ...editingCourse,
                      description: e.target.value,
                    })
                  }
                  className="bg-white dark:bg-zinc-900 border border-zinc-300 dark:border-zinc-700 focus:ring-2 focus:ring-blue-400"
                />
              </div>

              {/* Price and Duration */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">
                    Price ($)
                  </label>
                  <Input
                    type="number"
                    value={editingCourse.price}
                    onChange={(e) =>
                      setEditingCourse({
                        ...editingCourse,
                        price: Number(e.target.value),
                      })
                    }
                    className="bg-white dark:bg-zinc-900 border border-zinc-300 dark:border-zinc-700 focus:ring-2 focus:ring-blue-400"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">
                    Duration (minutes)
                  </label>
                  <Select
                    value={String(editingCourse.duration)}
                    onValueChange={(value) =>
                      setEditingCourse({
                        ...editingCourse,
                        duration: Number(value),
                      })
                    }
                  >
                    <SelectTrigger className="w-full bg-white dark:bg-zinc-900 border border-zinc-300 dark:border-zinc-700 focus:ring-2 focus:ring-blue-400">
                      <SelectValue placeholder="Select duration" />
                    </SelectTrigger>
                    <SelectContent className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-700 shadow-md">
                      <SelectItem value="15">15</SelectItem>
                      <SelectItem value="30">30</SelectItem>
                      <SelectItem value="45">45</SelectItem>
                      <SelectItem value="60">60</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </div>
          )}

          <DialogFooter className="mt-6">
            <Button
              variant="outline"
              onClick={() => setIsDialogOpen(false)}
              disabled={saving}
            >
              Cancel
            </Button>
            <Button
              className="bg-blue-600 hover:bg-blue-700 text-white"
              onClick={handleSaveEdit}
              disabled={saving}
            >
              {saving ? "Saving..." : "Save Changes"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
