"use client";
import { Button } from "@/components/ui/button";

export default function CourseDetailModal({
  course,
  onClose,
  onChat,
}: {
  course: any;
  onClose: () => void;
  onChat: () => void;
}) {
  return (
    <div className="fixed inset-0 bg-black/50 flex justify-center items-center z-40">
      <div className="bg-white rounded-xl p-8 w-[650px] shadow-xl relative">
        {/* Close Button */}
        <button
          onClick={onClose}
          className="absolute top-3 right-3 text-gray-500 hover:text-gray-700 text-xl"
        >
          ✕
        </button>

        {/* Course Title */}
        <h2 className="text-3xl font-bold mb-3 text-gray-800">{course.title}</h2>
        <p className="text-gray-600 text-sm mb-6">{course.description}</p>

        {/* Divider */}
        <hr className="mb-6 border-gray-300" />

        {/* Prophet Detail Section */}
        <div className="flex items-center gap-5 mb-6">
          {/* Prophet Profile Picture (optional field) */}
          {course.prophetProfileUrl ? (
            <img
              src={course.prophetProfileUrl}
              alt={course.prophet}
              className="w-20 h-20 rounded-full border"
            />
          ) : (
            <div className="w-20 h-20 rounded-full bg-gray-200 flex items-center justify-center text-gray-500 text-xl">
              {course.prophet?.[0]?.toUpperCase() ?? "?"}
            </div>
          )}

          {/* Prophet Info */}
          <div>
            <h3 className="text-xl font-semibold text-gray-800">{course.prophet}</h3>
            <p className="text-gray-500 text-sm">Professional Prophet</p>
            {course.prophetBio && (
              <p className="text-gray-600 text-sm mt-1">{course.prophetBio}</p>
            )}
          </div>
        </div>

        {/* Course Extra Info */}
        <div className="bg-gray-50 p-4 rounded-lg mb-6">
          <p className="text-gray-700 text-sm mb-1">
            <strong>Duration:</strong> {course.duration || "Not specified"}
          </p>
          <p className="text-gray-700 text-sm mb-1">
            <strong>Category:</strong> {course.category || "General"}
          </p>
          <p className="text-gray-700 text-sm">
            <strong>Price:</strong> {course.price}
          </p>
        </div>

        {/* Action Buttons */}
        <div className="flex justify-end gap-3">
          <Button
            onClick={onChat}
            className="bg-green-500 text-white hover:bg-green-600 px-5"
          >
            💬 Chat with Prophet
          </Button>
          <Button
            onClick={onClose}
            className="bg-gray-400 text-white hover:bg-gray-500 px-5"
          >
            Close
          </Button>
        </div>
      </div>
    </div>
  );
}
