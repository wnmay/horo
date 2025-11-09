"use client";
import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { useState, useEffect, useRef } from "react";
import CourseCard from "@/components/course-card";
import Image from "next/image";
import FetchRoomsButton from "@/components/testFetch";

export default function HomePage() {
  const router = useRouter();
  const [page, setPage] = useState(0);
  const [user, setUser] = useState<any>(null);
  const [showMenu, setShowMenu] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const [adIndex, setAdIndex] = useState(0);

  // --- Filter States ---
  const [selectedTag, setSelectedTag] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [sortOption, setSortOption] = useState<"title" | "price" | "none">(
    "none"
  );

  // --- Temporary (input) filter states before Apply ---
  const [tempSearch, setTempSearch] = useState("");
  const [tempTag, setTempTag] = useState<string | null>(null);
  const [tempSort, setTempSort] = useState<"title" | "price" | "none">("none");

  const COURSES_PER_PAGE = 5;
  const courseTags = [
    "Love",
    "Study",
    "Work",
    "Health",
    "Finance",
    "Personal_Growth",
  ];

  // --- Check login ---
  useEffect(() => {
    const storedUser = localStorage.getItem("user");
    if (storedUser) setUser(JSON.parse(storedUser));
  }, []);

  // --- Close dropdown ---
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setShowMenu(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);
  //--- image url ---
  const tagImages: Record<string, string> = {
    Love: "/images/Love.jpg",
    Study: "/images/Study.jpg",
    Work: "/images/Work.jpg",
    Health: "/images/Health.jpg",
    Finance: "/images/Finance.jpg",
    Personal_Growth: "/images/Personal_Growth.jpg",
  };

  // --- Ads ---
  const ads = [
    {
      id: 0,
      title: "Welcome to Horo",
      description:
        "Your personalized horoscope app. Get started by signing in or creating an account.",
    },
    {
      id: 1,
      title: "Get a free horoscope reading!",
      description: "Sign up now and receive your personalized daily horoscope.",
    },
    {
      id: 2,
      title: "Exclusive Tarot Tips",
      description: "Learn top secrets for interpreting tarot cards accurately.",
    },
    {
      id: 3,
      title: "Numerology Guide",
      description: "Discover the hidden meaning of numbers in your life.",
    },
    {
      id: 4,
      title: "Astrology Workshop",
      description: "Join our online workshop and deepen your astrology skills.",
    },
    {
      id: 5,
      title: "Dream Interpretation",
      description: "Understand your dreams and their messages every morning.",
    },
  ];

  useEffect(() => {
    const interval = setInterval(() => {
      setAdIndex((prev) => (prev + 1) % ads.length);
    }, 3000);
    return () => clearInterval(interval);
  }, []);

  // --- Mock Courses ---
  const mockCourses = [
    {
      id: 1,
      title: "Beginner Astrology 101",
      description:
        "Learn the fundamentals of astrology, zodiac signs, and planetary movements.",
      prophet: "Prophet Orion",
      price: "$49",
    },
    {
      id: 2,
      title: "Advanced Horoscope Reading",
      description:
        "Master horoscope interpretations and learn how to create accurate readings.",
      prophet: "Prophet Selene",
      price: "$79",
    },
    {
      id: 3,
      title: "Tarot for Self Discovery",
      description:
        "Explore the meaning of tarot cards and how to use them for personal guidance.",
      prophet: "Prophet Lyra",
      price: "$39",
    },
    {
      id: 4,
      title: "Zodiac Compatibility Secrets",
      description: "Understand the relationship dynamics between zodiac signs.",
      prophet: "Prophet Atlas",
      price: "$59",
    },
    {
      id: 5,
      title: "Planetary Transits Explained",
      description:
        "Learn how planetary movements influence your daily horoscope.",
      prophet: "Prophet Vega",
      price: "$69",
    },
    {
      id: 6,
      title: "Astrology and Life Purpose",
      description:
        "Find your path through the stars with this transformative course.",
      prophet: "Prophet Nova",
      price: "$89",
    },
    {
      id: 7,
      title: "Dream Interpretation Basics",
      description: "Understand the spiritual meaning behind common dreams.",
      prophet: "Prophet Luna",
      price: "$45",
    },
    {
      id: 8,
      title: "Numerology for Beginners",
      description: "Discover how numbers shape your destiny and personality.",
      prophet: "Prophet Orion",
      price: "$55",
    },
  ];

  let filteredCourses = mockCourses;

  // Slice data for current page
  const startIndex = page * COURSES_PER_PAGE;
  const visibleCourses = filteredCourses.slice(
    startIndex,
    startIndex + COURSES_PER_PAGE
  );
  const hasNext = startIndex + COURSES_PER_PAGE < filteredCourses.length;
  const hasPrev = page > 0;

  // console.log(user)

  return (
    <div className="relative min-h-screen bg-gradient-to-b from-white to-zinc-100 dark:from-zinc-900 dark:to-zinc-950">
      {/* Carousel */}
      <div className="max-w-6xl mx-auto px-4 mb-12">
        <Card className="transition-all duration-500 mx-auto w-full md:w-[600px] p-10 text-center bg-yellow-50 dark:bg-yellow-900 text-zinc-800 dark:text-yellow-100 min-h-[250px] flex flex-col justify-center">
          <h3 className="font-bold text-4xl md:text-5xl mb-4">
            {ads[adIndex].title}
          </h3>
          <p className="text-lg md:text-xl mt-2">{ads[adIndex].description}</p>
        </Card>
      </div>

      {/* Tag Filter */}
      <div className="max-w-6xl mx-auto mb-10 px-4 text-center">
        <h2 className="text-2xl font-semibold mb-6 text-zinc-800 dark:text-zinc-100">
          Explore by Category
        </h2>
        <div className="flex flex-wrap justify-center gap-6">
          {courseTags.map((tag) => (
            <button
              key={tag}
              onClick={() => setTempTag(tempTag === tag ? null : tag)}
              className={`flex flex-col items-center border rounded-2xl p-4 w-[110px] text-sm font-medium shadow-sm transition-transform
                ${
                  tempTag === tag
                    ? "bg-blue-600 text-white border-blue-700 scale-105"
                    : "bg-white dark:bg-zinc-800 hover:bg-zinc-100 dark:hover:bg-zinc-700 border-zinc-300 dark:border-zinc-600"
                }`}
            >
              <div className="w-10 h-10 mb-2 relative">
                <Image
                  src={tagImages[tag]}
                  alt={tag}
                  fill
                  sizes="48px"
                  className="object-contain"
                  priority
                />
              </div>
              {tag.replaceAll("_", " ")}
            </button>
          ))}
        </div>
      </div>

      {/* Search + Sort + Apply */}
      <div className="max-w-6xl mx-auto px-4 mb-10">
        {/* Added max-w-lg to the container to ensure it doesn't stretch infinitely on large screens */}
        <div className="flex flex-col md:flex-row justify-center md:justify-between items-center gap-2 bg-white dark:bg-zinc-800 p-4 rounded-xl shadow max-w-lg mx-auto">
          <input
            type="text"
            placeholder="Search courses or prophets..."
            value={tempSearch}
            onChange={(e) => setTempSearch(e.target.value)}
            // Added md:w-[50%] back for desktop, but kept w-full for mobile flexibility
            className="w-full md:w-[50%] border border-zinc-300 dark:border-zinc-600 rounded-lg px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none dark:bg-zinc-900"
          />

          <select
            value={tempSort}
            onChange={(e) =>
              setTempSort(e.target.value as "title" | "price" | "none")
            }
            // Added md:w-[25%] back for desktop, but kept w-full for mobile flexibility
            className="w-full md:w-[25%] border border-zinc-300 dark:border-zinc-600 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none dark:bg-zinc-900"
          >
            <option value="none">Sort by</option>
            <option value="title">Title (A–Z)</option>
            <option value="price">Price (Low → High)</option>
          </select>

          <Button
            onClick={() => {
              setSearchQuery(tempSearch);
              setSelectedTag(tempTag);
              setSortOption(tempSort);
              setPage(0);
            }}
            className="bg-blue-600 text-white hover:bg-blue-700 px-6 py-2 rounded-lg"
          >
            Apply
          </Button>
        </div>

        {/* Welcome section */}
        <div className="flex min-h-[60vh] items-center justify-center px-4">
          <Card className="p-8 text-center">
            <h1 className="text-5xl font-bold text-black dark:text-zinc-50 mb-4">
              Welcome to Horo
            </h1>
            <p className="text-lg text-zinc-600 dark:text-zinc-400 max-w-md mx-auto">
              Your personalized horoscope app. Get started by signing in or
              creating an account.
            </p>
          </Card>
          <FetchRoomsButton />
        </div>
      </div>

      {/* Courses */}
      <div className="px-6 py-12 bg-white dark:bg-zinc-900">
        <h2 className="text-3xl font-semibold text-center text-zinc-800 dark:text-zinc-100 mb-8">
          {selectedTag ? `${selectedTag} Courses` : "Featured Courses"}
        </h2>
        <div className="max-w-6xl mx-auto overflow-x-auto">
          <div className="flex gap-6 min-w-max px-2">
            {visibleCourses.map((course) => (
              <div key={course.id} className="w-[240px] flex-shrink-0">
                <CourseCard course={course} />
              </div>
            ))}
          </div>
        </div>

        {filteredCourses.length > COURSES_PER_PAGE && (
          <div className="flex items-center justify-between mt-8 max-w-6xl mx-auto px-2">
            {/* Empty div to balance the space on the left */}
            <div className="w-1/3"></div>

            {/* Centered Prev/Next Buttons */}
            <div className="flex gap-4 justify-center w-1/3">
              <Button
                variant="outline"
                disabled={!hasPrev}
                onClick={() => setPage((p) => p - 1)}
              >
                ← Previous
              </Button>
              <Button
                variant="outline"
                disabled={!hasNext}
                onClick={() => setPage((p) => p + 1)}
              >
                Next →
              </Button>
            </div>

            {/* See More Button on the right */}
            <div className="w-1/3 flex justify-end">
              <Button
                className="bg-blue-600 text-white hover:bg-blue-700 px-6 py-3 rounded-lg"
                onClick={() => router.push("/courses")}
              >
                See More Courses →
              </Button>
            </div>
          </div>
        )}

        <p className="text-center text-sm text-zinc-500 mt-4">
          Page {page + 1} of{" "}
          {Math.ceil(filteredCourses.length / COURSES_PER_PAGE)}
        </p>
      </div>
    </div>
  );
}
