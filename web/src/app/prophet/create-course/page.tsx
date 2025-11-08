"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import Card from "@/components/ui/card";
import { toast } from "sonner";

export default function CreateCoursePage() {
  const [form, setForm] = useState({
    coursename: "",
    description: "",
    price: "",
    duration: "",
    prophetname: "",
  });

  const [submitting, setSubmitting] = useState(false);
  const router = useRouter();
  const token = localStorage.getItem("token"); //TODO: get from auth context once ready

  function handleChange(e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) {
    setForm({ ...form, [e.target.name]: e.target.value });
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    if (!form.coursename || !form.description || !form.price || !form.duration || !form.prophetname) {
      toast("Please fill out all fields.", { position: "top-right" });
      return;
    }

    try {
      setSubmitting(true);

      const body = {
        coursename: form.coursename.trim(),
        description: form.description.trim(),
        price: parseInt(form.price, 10),
        duration: parseInt(form.duration, 10),
        prophetname: form.prophetname.trim(),
      };

      const res = await fetch("/courses", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      });

      if (!res.ok) {
        const errText = await res.text();
        throw new Error(errText || "Failed to create course");
      }

      toast.success("Course created successfully!", { position: "top-right" });
      router.push("/prophet/courses");
    } catch (err) {
      toast.error("Failed to create course", {
        description: (err as Error).message,
        position: "top-right",
      });
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center px-6 py-12">
      <Card className="w-full max-w-xl p-8 rounded-2xl shadow-xl border border-blue-100 dark:border-zinc-800 bg-white/90 dark:bg-zinc-900/80 backdrop-blur-md">
        <h1 className="text-3xl font-extrabold text-center text-blue-900 dark:text-blue-300 mb-8">
          Create New Course
        </h1>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Course Name */}
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">
              Course Name
            </label>
            <Input
              name="coursename"
              placeholder="Enter course name"
              value={form.coursename}
              onChange={handleChange}
              className="focus:ring-2 focus:ring-blue-400"
            />
          </div>

          {/* Description */}
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">
              Description
            </label>
            <Textarea
              name="description"
              placeholder="Describe your course"
              rows={4}
              value={form.description}
              onChange={handleChange}
              className="focus:ring-2 focus:ring-blue-400"
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
                name="price"
                placeholder="Enter price"
                value={form.price}
                onChange={handleChange}
                className="focus:ring-2 focus:ring-blue-400"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">
                Duration (minutes)
              </label>
              <Select
                value={form.duration}
                onValueChange={(value) => setForm({ ...form, duration: value })}
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

          {/* Prophet Name */}
          <div>
            <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1">
              Prophet Name
            </label>
            <Input
              name="prophetname"
              placeholder="Enter prophet name"
              value={form.prophetname}
              onChange={handleChange}
              className="focus:ring-2 focus:ring-blue-400"
            />
          </div>

          {/* Submit Button */}
          <Button
            type="submit"
            className="w-full mt-6 bg-blue-600 hover:bg-blue-700 text-white font-semibold shadow-md hover:shadow-lg rounded-md transition-all duration-200"
            disabled={submitting}
          >
            {submitting ? "Creating..." : "Create Course"}
          </Button>
        </form>
      </Card>
    </div>
  );
}
