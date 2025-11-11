"use client";

import { use, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import api from "@/lib/api/api-client";
import { courseTypeImageMap } from "@/types/course-type";
import { toast } from "sonner";

interface Course {
  id: string;
  courseName: string;
  courseType: string;
  description: string;
  duration: number;
  price: number;
  rating: number;
  prophetId: string;
  prophetName: string;
  createdTime: string;
  deletedAt: boolean;
  image: string;
  reviews: {
    id: string;
    courseId: string;
    customerId: string;
    customerName: string;
    title: string;
    comment: string;
    rating: number;
    createdAt: string;
    deletedAt: boolean;
  }[];
}

function getImageForType(type: string): string {
  if (!type) return "/images/default.jpg";
  const match = Object.keys(courseTypeImageMap).find(
    (key) => key.toLowerCase() === type.toLowerCase()
  );
  return match
    ? courseTypeImageMap[match as keyof typeof courseTypeImageMap]
    : "/images/default.jpg";
}

export default function CourseDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const router = useRouter();

  const [course, setCourse] = useState<Course | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchCourse = async () => {
      try {
        setIsLoading(true);
        const response = await api.get(`/api/courses/${id}`);
        const data = response.data.data;

        const mappedCourse: Course = {
          id: data.id,
          courseName: data.coursename,
          courseType: data.coursetype,
          description: data.description,
          duration: data.duration,
          price: data.price,
          rating: data.review_score,
          prophetId: data.prophet_id,
          prophetName: data.prophetname,
          createdTime: data.created_time,
          deletedAt: data.deleted_at,
          image: getImageForType(data.coursetype),
          reviews:
            data.reviews?.map((r: any) => ({
              id: r.id,
              courseId: r.course_id,
              customerId: r.customer_id,
              customerName: r.customername,
              title: r.title,
              comment: r.description,
              rating: r.review_score,
              createdAt: r.created_at,
              deletedAt: r.deleted_at,
            })) || [],
        };

        setCourse(mappedCourse);
      } catch (err) {
        console.error("Error fetching course:", err);
        setError("Failed to load course details");
      } finally {
        setIsLoading(false);
      }
    };

    fetchCourse();
  }, [id]);

  const handleStartChat = async () => {
    if (course) {
      const res = await api.post(`/api/chat/rooms`, { courseId: course.id });
      if (res.status === 200 || res.status === 201) {
        router.push(`/chat`);
      } else {
        throw new Error(res.data.message || "Failed to create chat room");
      }
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-600">Loading course details...</p>
        </div>
      </div>
    );
  }

  if (error || !course) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-600 mb-4">{error || "Course not found"}</p>
          <button
            onClick={() => router.back()}
            className="text-indigo-600 hover:text-indigo-700 font-medium"
          >
            ← Go Back
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-5xl mx-auto">
        {/* Back Button */}
        <button
          onClick={() => router.back()}
          className="mb-6 text-gray-600 hover:text-gray-900 flex items-center gap-2 transition-colors"
        >
          <span>←</span>
          <span>Back to Courses</span>
        </button>

        {/* Main Card */}
        <div className="bg-white shadow-lg rounded-2xl overflow-hidden">
          {/* Hero Image */}
          <div className="relative h-80 w-full bg-gradient-to-br from-indigo-500 to-purple-600">
            <img
              src={course.image}
              alt={course.courseName}
              className="object-cover w-full h-full"
            />
            <div className="absolute inset-0 bg-black bg-opacity-20"></div>
          </div>

          {/* Course Details */}
          <div className="p-8">
            {/* Header */}
            <div className="mb-6">
              <div className="flex items-start justify-between gap-4 mb-3">
                <h1 className="text-4xl font-bold text-gray-900 leading-tight">
                  {course.courseName}
                </h1>
                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-indigo-100 text-indigo-800 whitespace-nowrap">
                  {course.courseType}
                </span>
              </div>

              <div className="flex items-center gap-4 text-sm text-gray-600">
                <span className="flex items-center gap-1">
                  <span className="font-medium">By</span> {course.prophetName}
                </span>
                <span>•</span>
                <span>{course.duration} minutes</span>
                <span>•</span>
                <span className="flex items-center gap-1">
                  <span className="text-yellow-500">⭐</span>
                  <span className="font-medium">
                    {course.rating.toFixed(1)}
                  </span>
                </span>
              </div>
            </div>

            {/* Description */}
            <div className="mb-8">
              <h2 className="text-xl font-semibold text-gray-900 mb-3">
                About This Course
              </h2>
              <p className="text-gray-700 leading-relaxed">
                {course.description}
              </p>
            </div>

            {/* Price and CTA */}
            <div className="flex items-center justify-between p-6 bg-gray-50 rounded-xl border border-gray-200">
              <div>
                <p className="text-sm text-gray-600 mb-1">Course Price</p>
                <p className="text-3xl font-bold text-indigo-600">
                  ฿{course.price.toLocaleString()}
                </p>
              </div>
              <button
                onClick={handleStartChat}
                className="bg-indigo-600 hover:bg-indigo-700 active:bg-indigo-800 text-white px-8 py-3 rounded-lg font-semibold transition-colors shadow-md hover:shadow-lg"
              >
                Start Chat with Prophet
              </button>
            </div>

            {/* Reviews Section */}
            {course.reviews.length > 0 && (
              <div className="mt-10 pt-8 border-t border-gray-200">
                <div className="flex items-center justify-between mb-6">
                  <h2 className="text-2xl font-bold text-gray-900">
                    Customer Reviews
                  </h2>
                  <span className="text-sm text-gray-600">
                    {course.reviews.length}{" "}
                    {course.reviews.length === 1 ? "review" : "reviews"}
                  </span>
                </div>

                <div className="space-y-4">
                  {course.reviews.map((review) => (
                    <div
                      key={review.id}
                      className="p-5 bg-white border border-gray-200 rounded-lg hover:border-indigo-200 hover:shadow-sm transition-all"
                    >
                      <div className="flex items-start justify-between mb-3">
                        <div>
                          <h3 className="font-semibold text-gray-900">
                            {review.customerName}
                          </h3>
                          {review.title && (
                            <p className="text-sm text-gray-600 mt-1">
                              {review.title}
                            </p>
                          )}
                        </div>
                        <div className="flex items-center gap-1 bg-yellow-50 px-2 py-1 rounded">
                          <span className="text-yellow-500">⭐</span>
                          <span className="font-semibold text-gray-900">
                            {review.rating.toFixed(1)}
                          </span>
                        </div>
                      </div>
                      <p className="text-gray-700 leading-relaxed">
                        {review.comment}
                      </p>
                      <p className="text-xs text-gray-500 mt-3">
                        {new Date(review.createdAt).toLocaleDateString(
                          "en-US",
                          {
                            year: "numeric",
                            month: "long",
                            day: "numeric",
                          }
                        )}
                      </p>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
