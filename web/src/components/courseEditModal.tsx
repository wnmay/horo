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

interface CourseEditModalProps {
  isOpen: boolean;
  onClose: () => void;
  course: Course | null;
  onCourseChange: (course: Course) => void;
  onSave: () => void;
  isSaving: boolean;
}

export default function CourseEditModal({
  isOpen,
  onClose,
  course,
  onCourseChange,
  onSave,
  isSaving,
}: CourseEditModalProps) {
  if (!course) return null;

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
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

        <div className="space-y-4 mt-4">
          {/* Course Name */}
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">
              Course Name
            </label>
            <Input
              value={course.coursename}
              onChange={(e) =>
                onCourseChange({
                  ...course,
                  coursename: e.target.value,
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
              value={course.description}
              onChange={(e) =>
                onCourseChange({
                  ...course,
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
                value={course.price}
                onChange={(e) =>
                  onCourseChange({
                    ...course,
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
                value={String(course.duration)}
                onValueChange={(value) =>
                  onCourseChange({
                    ...course,
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

        <DialogFooter className="mt-6">
          <Button variant="outline" onClick={onClose} disabled={isSaving}>
            Cancel
          </Button>
          <Button
            className="bg-blue-600 hover:bg-blue-700 text-white"
            onClick={onSave}
            disabled={isSaving}
          >
            {isSaving ? "Saving..." : "Save Changes"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
