"use client";
import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { useState, useEffect, useRef } from "react";
import CourseCard from "@/components/course-card";

export default function HomePage() {
  const router = useRouter();
  const [page, setPage] = useState(0);
  const [user, setUser] = useState<any>(null);
  const [showMenu, setShowMenu] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const [adIndex, setAdIndex] = useState(0); // For ad carousel
  const COURSES_PER_PAGE = 5;

  // Check if user is logged in
  useEffect(() => {
    const storedUser = localStorage.getItem("user");
    if (storedUser) setUser(JSON.parse(storedUser));
  }, []);

  // Close dropdown when click outside
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setShowMenu(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  // Ads including Welcome as first
  const ads = [
    {
      id: 0,
      title: "Welcome to Horo",
      description: "Your personalized horoscope app. Get started by signing in or creating an account.",
      isWelcome: true,
    },
    { id: 1, title: "Get a free horoscope reading!", description: "Sign up now and receive your personalized daily horoscope." },
    { id: 2, title: "Exclusive Tarot Tips", description: "Learn top secrets for interpreting tarot cards accurately." },
    { id: 3, title: "Numerology Guide", description: "Discover the hidden meaning of numbers in your life." },
    { id: 4, title: "Astrology Workshop", description: "Join our online workshop and deepen your astrology skills." },
    { id: 5, title: "Dream Interpretation", description: "Understand your dreams and their messages every morning." },
  ];

  // Auto-cycle ads
  useEffect(() => {
    const interval = setInterval(() => {
      setAdIndex((prev) => (prev + 1) % ads.length);
    }, 3000);
    return () => clearInterval(interval);
  }, []);

  // Mock course data
  const mockCourses = [
    { id: 1, title: "Beginner Astrology 101", description: "Learn the fundamentals of astrology, zodiac signs, and planetary movements.", prophet: "Prophet Orion", price: "$49" },
    { id: 2, title: "Advanced Horoscope Reading", description: "Master horoscope interpretations and learn how to create accurate readings.", prophet: "Prophet Selene", price: "$79" },
    { id: 3, title: "Tarot for Self Discovery", description: "Explore the meaning of tarot cards and how to use them for personal guidance.", prophet: "Prophet Lyra", price: "$39" },
    { id: 4, title: "Zodiac Compatibility Secrets", description: "Understand the relationship dynamics between zodiac signs.", prophet: "Prophet Atlas", price: "$59" },
    { id: 5, title: "Planetary Transits Explained", description: "Learn how planetary movements influence your daily horoscope.", prophet: "Prophet Vega", price: "$69" },
    { id: 6, title: "Astrology and Life Purpose", description: "Find your path through the stars with this transformative course.", prophet: "Prophet Nova", price: "$89" },
    { id: 7, title: "Dream Interpretation Basics", description: "Understand the spiritual meaning behind common dreams.", prophet: "Prophet Luna", price: "$45" },
    { id: 8, title: "Numerology for Beginners", description: "Discover how numbers shape your destiny and personality.", prophet: "Prophet Orion", price: "$55" },
  ];

  // Pagination logic
  const startIndex = page * COURSES_PER_PAGE;
  const visibleCourses = mockCourses.slice(startIndex, startIndex + COURSES_PER_PAGE);
  const hasNext = startIndex + COURSES_PER_PAGE < mockCourses.length;
  const hasPrev = page > 0;

  // Logout handler
  const handleLogout = () => {
    localStorage.removeItem("user");
    setUser(null);
    router.refresh();
  };

  // Change name handler
  const handleChangeName = () => {
    const newName = prompt("Enter your new full name:", user?.name || "");
    if (newName) {
      const updatedUser = { ...user, name: newName };
      localStorage.setItem("user", JSON.stringify(updatedUser));
      setUser(updatedUser);
    }
  };

  // Switch account
  const handleSwitchAccount = () => {
    localStorage.removeItem("user");
    router.push("/signin");
  };

  return (
    <div className="relative min-h-screen bg-gradient-to-b from-white to-zinc-100 dark:from-zinc-900 dark:to-zinc-950">
      {/* Header */}
      <div className="fixed top-4 right-4 p-4 flex gap-4 items-center z-50">
        {!user ? (
          <>
            <Button className="bg-blue-500 text-white hover:bg-blue-600" onClick={() => router.push("/signin")}>
              Sign in
            </Button>
            <Button className="bg-blue-500 text-white hover:bg-blue-600" onClick={() => router.push("/signup")}>
              Sign up
            </Button>
          </>
        ) : (
          <div className="relative" ref={menuRef}>
            <div
              className="flex items-center bg-white dark:bg-zinc-800 rounded-full shadow px-3 py-1 hover:shadow-lg transition cursor-pointer"
              onClick={() => setShowMenu(!showMenu)}
            >
              <img
                src={`https://ui-avatars.com/api/?name=${encodeURIComponent(user.name || user.email)}&background=random`}
                alt="Profile"
                className="w-8 h-8 rounded-full mr-2"
              />
              <span className="text-sm font-medium text-zinc-800 dark:text-zinc-200">
                {(user.name || user.email)?.split(" ")[0]}
              </span>
            </div>

            {showMenu && (
              <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-zinc-800 rounded-lg shadow-lg border border-zinc-200 dark:border-zinc-700 py-2">
                {user.role === "prophet" ? (
                  <>
                    <button onClick={() => router.push("/prophet/dashboard")} className="block w-full text-left px-4 py-2 hover:bg-zinc-100 dark:hover:bg-zinc-700">
                      Dashboard
                    </button>
                    <button onClick={handleLogout} className="block w-full text-left px-4 py-2 text-red-500 hover:bg-zinc-100 dark:hover:bg-zinc-700">
                      Logout
                    </button>
                  </>
                ) : (
                  <>
                    <button onClick={handleChangeName} className="block w-full text-left px-4 py-2 hover:bg-zinc-100 dark:hover:bg-zinc-700">
                      Change Full Name
                    </button>
                    <button onClick={handleSwitchAccount} className="block w-full text-left px-4 py-2 hover:bg-zinc-100 dark:hover:bg-zinc-700">
                      Switch Account
                    </button>
                    <button onClick={handleLogout} className="block w-full text-left px-4 py-2 text-red-500 hover:bg-zinc-100 dark:hover:bg-zinc-700">
                      Logout
                    </button>
                  </>
                )}
              </div>
            )}
          </div>
        )}
      </div>

      {/* Ad / Welcome Carousel */}
      <div className="max-w-6xl mx-auto px-4 mb-12">
        <Card
          className="transition-all duration-500 mx-auto w-full md:w-[600px] p-10 text-center bg-yellow-50 dark:bg-yellow-900 text-zinc-800 dark:text-yellow-100
                    min-h-[250px] flex flex-col justify-center"
        >
          <h3 className="font-bold text-4xl md:text-5xl mb-4">
            {ads[adIndex].title}
          </h3>
          <p className="text-lg md:text-xl mt-2">
            {ads[adIndex].description}
          </p>
        </Card>
      </div>


      {/* Featured Courses */}
      <div className="px-6 py-12 bg-white dark:bg-zinc-900">
        <h2 className="text-3xl font-semibold text-center text-zinc-800 dark:text-zinc-100 mb-8">Featured Courses</h2>
        <div className="max-w-6xl mx-auto overflow-x-auto">
          <div className="flex gap-6 min-w-max px-2">
            {mockCourses.slice(startIndex, startIndex + COURSES_PER_PAGE).map((course) => (
              <div key={course.id} className="w-[240px] flex-shrink-0">
                <CourseCard course={course} />
              </div>
            ))}
          </div>
        </div>

        <div className="flex justify-center gap-4 mt-8">
          <Button variant="outline" disabled={!hasPrev} onClick={() => setPage((p) => p - 1)}>← Previous</Button>
          <Button variant="outline" disabled={!hasNext} onClick={() => setPage((p) => p + 1)}>Next →</Button>
        </div>

        <p className="text-center text-sm text-zinc-500 mt-4">Page {page + 1} of {Math.ceil(mockCourses.length / COURSES_PER_PAGE)}</p>
      </div>
    </div>
  );
}
