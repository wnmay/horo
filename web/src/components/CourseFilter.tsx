"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";

interface CourseFilterProps {
  courseTags: string[];
  tagImages: Record<string, string>;
  onApply: (filters: {
    tag: string | null;
    search: string;
    sort: "title" | "price" | "none";
  }) => void;
}

export default function CourseFilter({ courseTags, tagImages, onApply }: CourseFilterProps) {
  const [tempTag, setTempTag] = useState<string | null>(null);
  const [tempSearch, setTempSearch] = useState("");
  const [tempSort, setTempSort] = useState<"title" | "price" | "none">("none");

  return (
    <div className="w-full space-y-10">
      {/* Tag Filter */}
      <div className="text-center">
        <h2 className="text-2xl font-semibold mb-6 text-zinc-800 dark:text-zinc-100">
          Filter by Category
        </h2>
        <div className="flex flex-wrap justify-center gap-6">
          {courseTags.map((tag) => (
            <button
              key={tag}
              onClick={() => setTempTag(tempTag === tag ? null : tag)}
              className={`flex flex-col items-center border rounded-2xl p-4 w-[110px] text-sm font-medium shadow-sm transition-transform
                ${tempTag === tag
                  ? "bg-blue-600 text-white border-blue-700 scale-105"
                  : "bg-white dark:bg-zinc-800 hover:bg-zinc-100 dark:hover:bg-zinc-700 border-zinc-300 dark:border-zinc-600"
                }`}
            >
              <img src={tagImages[tag]} alt={tag} className="w-10 h-10 mb-2 opacity-90" />
              {tag.replaceAll("_", " ")}
            </button>
          ))}
        </div>
      </div>

      {/* Search + Sort + Apply */}
      <div className="max-w-6xl mx-auto px-4">
        <div className="flex flex-col md:flex-row justify-center md:justify-between items-center gap-2 bg-white dark:bg-zinc-800 p-4 rounded-xl shadow max-w-lg mx-auto">
          <input
            type="text"
            placeholder="Search courses or prophets..."
            value={tempSearch}
            onChange={(e) => setTempSearch(e.target.value)}
            className="w-full md:w-[50%] border border-zinc-300 dark:border-zinc-600 rounded-lg px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none dark:bg-zinc-900"
          />

          <select
            value={tempSort}
            onChange={(e) => setTempSort(e.target.value as "title" | "price" | "none")}
            className="w-full md:w-[25%] border border-zinc-300 dark:border-zinc-600 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none dark:bg-zinc-900"
          >
            <option value="none">Sort by</option>
            <option value="title">Title (A–Z)</option>
            <option value="price">Price (Low → High)</option>
          </select>

          <Button
            onClick={() =>
              onApply({
                tag: tempTag,
                search: tempSearch,
                sort: tempSort,
              })
            }
            className="bg-blue-600 text-white hover:bg-blue-700 px-6 py-2 rounded-lg"
          >
            Apply
          </Button>
        </div>
      </div>
    </div>
  );
}
