"use client";

import { Button } from "@/components/ui/button";

interface CourseFilterProps {
  courseTags: string[];
  tagImages: Record<string, string>;
  selectedTag?: string | null;
  searchQuery?: string;
  sortOption?: "title" | "price" | "none";
  durationFilter?: number | null;
  onTagChange: (tag: string | null) => void;
  onSearchChange: (query: string) => void;
  onSortChange: (option: "title" | "price" | "none") => void;
  onDurationChange: (duration: number | null) => void;
  onApply: () => void;
}

export default function CourseFilter({
  courseTags,
  tagImages,
  selectedTag = null,
  searchQuery = "",
  sortOption = "none",
  durationFilter = null,
  onTagChange,
  onSearchChange,
  onSortChange,
  onDurationChange,
  onApply,
}: CourseFilterProps) {
  return (
    <div className="flex flex-col items-center gap-6">
      {/* Tag buttons with toggle behavior */}
      <div className="flex flex-wrap justify-center gap-6">
        {courseTags.map((tag) => (
          <button
            key={tag}
            className={`flex flex-col items-center gap-2 p-4 rounded w-32 h-36
              ${selectedTag === tag ? "bg-blue-500 text-white" : "bg-gray-200 text-gray-800"}
              hover:scale-105 transition-transform`}
            onClick={() => {
              // Toggle selection
              if (selectedTag === tag) {
                onTagChange(null); // deselect
              } else {
                onTagChange(tag); // select
              }
            }}
          >
            <img src={tagImages[tag]} alt={tag} className="w-16 h-16 rounded-full object-cover" />
            <span className="text-sm font-medium text-center">{tag}</span>
          </button>
        ))}
      </div>

      {/* Search, Sort, Duration, Apply in same row */}
      <div className="flex flex-wrap justify-center gap-4 items-center">
        <input
          type="text"
          defaultValue={searchQuery}
          onChange={(e) => onSearchChange(e.target.value)}
          placeholder="Search by course or prophet..."
          className="border rounded px-3 py-2 w-64 text-center"
          onKeyDown={(e) => {
            if (e.key === "Enter") onApply();
          }}
        />

        <select
          defaultValue={sortOption}
          onChange={(e) => onSortChange(e.target.value as "title" | "price" | "none")}
          className="border rounded px-3 py-2 text-center"
        >
          <option value="none">Sort: None</option>
          <option value="title">Sort by Title</option>
          <option value="price">Sort by Price</option>
        </select>

        <select
          defaultValue={durationFilter ?? ""}
          onChange={(e) => onDurationChange(e.target.value ? parseInt(e.target.value) : null)}
          className="border rounded px-3 py-2 text-center"
        >
          <option value="">Duration: All</option>
          <option value="1">15 mins</option>
          <option value="2">30 mins</option>
          <option value="3">45 mins</option>
          <option value="4">60 mins</option>
        </select>

        <button
          className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700"
          onClick={onApply}
        >
          Apply
        </button>
      </div>
    </div>
  );
}
