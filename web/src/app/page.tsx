"use client";
import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { useState, useEffect, useRef } from "react";
import CourseCard from "@/components/course-card";
import Image from "next/image";

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
  const [sortOption, setSortOption] = useState<"title" | "price" | "none">("none");

  // --- Temporary (input) filter states before Apply ---
  const [tempSearch, setTempSearch] = useState("");
  const [tempTag, setTempTag] = useState<string | null>(null);
  const [tempSort, setTempSort] = useState<"title" | "price" | "none">("none");

  // --- Dynamic course per page ---
  const [coursesPerPage, setCoursesPerPage] = useState(6);

  useEffect(() => {
    const handleResize = () => {
      const width = window.innerWidth;
      if (width < 640) setCoursesPerPage(1);
      else if (width < 768) setCoursesPerPage(2);
      else if (width < 1024) setCoursesPerPage(3);
      else if (width < 1280) setCoursesPerPage(4);
      else setCoursesPerPage(6);
    };

    handleResize(); // Run initially
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  const courseTags = ["Love", "Study", "Work", "Health", "Finance", "Personal_Growth"];

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

  // --- Image urls ---
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
    { id: 0, title: "Welcome to Horo", description: "Your personalized horoscope app. Get started by signing in or creating an account." },
    { id: 1, title: "Get a free horoscope reading!", description: "Sign up now and receive your personalized daily horoscope." },
    { id: 2, title: "Exclusive Tarot Tips", description: "Learn top secrets for interpreting tarot cards accurately." },
    { id: 3, title: "Numerology Guide", description: "Discover the hidden meaning of numbers in your life." },
    { id: 4, title: "Astrology Workshop", description: "Join our online workshop and deepen your astrology skills." },
    { id: 5, title: "Dream Interpretation", description: "Understand your dreams and their messages every morning." },
  ];

  useEffect(() => {
    const interval = setInterval(() => {
      setAdIndex((prev) => (prev + 1) % ads.length);
    }, 3000);
    return () => clearInterval(interval);
  }, []);

  // --- Mock Courses ---
  const mockCourses = [
    { id: 1, title: "Beginner Astrology 101", description: "Learn astrology basics and understand your birth chart.", prophet: "Prophet Orion", price: "$49", tags: ["Study"] },
    { id: 2, title: "Love Compatibility Reading", description: "Find your perfect match through the stars.", prophet: "Prophet Luna", price: "$59", tags: ["Love"] },
    { id: 3, title: "Work-Life Balance Guidance", description: "Discover how the planets influence your career and health.", prophet: "Prophet Selene", price: "$69", tags: ["Work", "Health"] },
    { id: 4, title: "Financial Fortune Reading", description: "Predict your wealth path and upcoming financial cycles.", prophet: "Prophet Nova", price: "$79", tags: ["Finance"] },
    { id: 5, title: "Healing Energy Workshop", description: "Restore your inner energy through guided meditation.", prophet: "Prophet Vega", price: "$89", tags: ["Health", "Personal_Growth"] },
    { id: 6, title: "Manifest Love & Success", description: "Attract abundance and romance using spiritual techniques.", prophet: "Prophet Atlas", price: "$99", tags: ["Love", "Personal_Growth"] },
    { id: 7, title: "Advanced Astrology Masterclass", description: "Master planetary aspects, houses, and transits.", prophet: "Prophet Orion", price: "$120", tags: ["Study"] },
    { id: 8, title: "Career Path Tarot Reading", description: "Reveal your professional destiny using tarot guidance.", prophet: "Prophet Aria", price: "$65", tags: ["Work"] },
    { id: 9, title: "Chakra Healing Session", description: "Balance your chakras and rejuvenate your spirit.", prophet: "Prophet Vega", price: "$75", tags: ["Health", "Personal_Growth"] },
    { id: 10, title: "Numerology of Wealth", description: "Decode the numbers behind your financial success.", prophet: "Prophet Nova", price: "$85", tags: ["Finance", "Study"] },
    { id: 11, title: "Twin Flame Connection", description: "Discover your twin flame and soul bond purpose.", prophet: "Prophet Luna", price: "$99", tags: ["Love"] },
    { id: 12, title: "Mindful Meditation 101", description: "Train your mind to stay calm and present.", prophet: "Prophet Selene", price: "$49", tags: ["Health", "Personal_Growth"] },
    { id: 13, title: "Dream Symbolism Analysis", description: "Interpret dreams to uncover subconscious truths.", prophet: "Prophet Aria", price: "$55", tags: ["Personal_Growth", "Study"] },
    { id: 14, title: "Entrepreneurial Energy Reading", description: "Find the right time to launch or expand your business.", prophet: "Prophet Nova", price: "$110", tags: ["Work", "Finance"] },
    { id: 15, title: "Emotional Healing Through Stars", description: "Astrological guidance for emotional well-being.", prophet: "Prophet Luna", price: "$70", tags: ["Health", "Love"] },
    { id: 16, title: "Destiny Numbers Reading", description: "Understand your life path and future opportunities.", prophet: "Prophet Atlas", price: "$88", tags: ["Study", "Personal_Growth"] },
    { id: 17, title: "Relationship Renewal Course", description: "Revive and strengthen your relationships.", prophet: "Prophet Luna", price: "$95", tags: ["Love"] },
    { id: 18, title: "Spiritual Success Blueprint", description: "Align your purpose with prosperity and joy.", prophet: "Prophet Orion", price: "$115", tags: ["Personal_Growth", "Finance"] },
    { id: 19, title: "Astro-Health Diagnostic", description: "Use astrology to improve your lifestyle and diet.", prophet: "Prophet Selene", price: "$79", tags: ["Health"] },
    { id: 20, title: "Career Pivot Mentorship", description: "Guided astrological advice for major life transitions.", prophet: "Prophet Aria", price: "$99", tags: ["Work", "Personal_Growth"] },
  ];


  // --- Filtering ---
  let filteredCourses = mockCourses;

  if (selectedTag) filteredCourses = filteredCourses.filter((c) => c.tags.includes(selectedTag));
  if (searchQuery.trim()) {
    const q = searchQuery.toLowerCase();
    filteredCourses = filteredCourses.filter(
      (c) => c.title.toLowerCase().includes(q) || c.prophet.toLowerCase().includes(q)
    );
  }
  if (sortOption === "title") filteredCourses.sort((a, b) => a.title.localeCompare(b.title));
  if (sortOption === "price") filteredCourses.sort((a, b) => parseFloat(a.price.slice(1)) - parseFloat(b.price.slice(1)));

  // --- Pagination ---
  const startIndex = page * coursesPerPage;
  const visibleCourses = filteredCourses.slice(startIndex, startIndex + coursesPerPage);
  const hasNext = startIndex + coursesPerPage < filteredCourses.length;
  const hasPrev = page > 0;

  return (
    <div className="relative min-h-screen bg-gradient-to-b from-white to-zinc-100 dark:from-zinc-900 dark:to-zinc-950">

      {/* Carousel */}
      <div className="max-w-6xl mx-auto px-4 mb-12">
        <Card className="transition-all duration-500 mx-auto w-full md:w-[600px] p-10 text-center bg-yellow-50 dark:bg-yellow-900 text-zinc-800 dark:text-yellow-100 min-h-[250px] flex flex-col justify-center">
          <h3 className="font-bold text-4xl md:text-5xl mb-4">{ads[adIndex].title}</h3>
          <p className="text-lg md:text-xl mt-2">{ads[adIndex].description}</p>
        </Card>
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
            onChange={(e) => setTempSort(e.target.value as "title" | "price" | "none")}
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
      </div>

      {/* Courses */}
      <div className="px-6 py-12 bg-white dark:bg-zinc-900">
        <h2 className="text-3xl font-semibold text-center text-zinc-800 dark:text-zinc-100 mb-8">
          {selectedTag ? `${selectedTag} Courses` : "Featured Courses"}
        </h2>

        <div className="max-w-full mx-auto overflow-x-auto">
          <div className="flex gap-6 min-w-max px-2 justify-center">
            {visibleCourses.map((course) => (
              <div key={course.id} className="w-[240px] flex-shrink-0">
                <CourseCard course={course} />
              </div>
            ))}
          </div>
        </div>

        {(
          <div className="flex items-center justify-between mt-8 max-w-6xl mx-auto px-2">
            <div className="w-1/3"></div>
            <div className="flex gap-4 justify-center w-1/3">
              <Button variant="outline" disabled={!hasPrev} onClick={() => setPage((p) => p - 1)}>
                ← Previous
              </Button>
              <Button variant="outline" disabled={!hasNext} onClick={() => setPage((p) => p + 1)}>
                Next →
              </Button>
            </div>
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
          Page {page + 1} of {Math.ceil(filteredCourses.length / coursesPerPage)}
        </p>
      </div>
    </div>
  );
}
