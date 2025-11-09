"use client";

import { use, useEffect, useState } from "react";
import { useRouter } from "next/navigation";

interface Course {
  id: string;
  name: string;
  price: number;
  rating: number;
  prophet: string;
  experience: string;
  specialties: string[];
  description: string;
  image: string;
  reviews: { user: string; comment: string; rating: number }[];
}

export default function CourseDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params); // unwrap the Promise using React.use()
  const router = useRouter();

  const [course, setCourse] = useState<Course | null>(null);
  const [showChat, setShowChat] = useState(false);
  const [chatMessages, setChatMessages] = useState<{ sender: string; text: string }[]>([]);
  const [input, setInput] = useState("");

  useEffect(() => {
    // simulate fetching course data
    const mockCourse: Course = {
      id,
      name: "Advanced Astrology Reading",
      price: 1200,
      rating: 4.9,
      prophet: "Master Flook",
      experience: "15 years of experience in astrology and tarot reading",
      specialties: ["Love", "Career", "Destiny", "Future Prediction"],
      description:
        "Discover deep insights into your personal and professional life. This course covers planetary alignments, zodiac compatibility, and real-life case readings to help you master astrology interpretation.",
      image:
        "https://images.unsplash.com/photo-1606112219348-204d7d8b94ee?auto=format&fit=crop&w=900&q=80",
      reviews: [
        { user: "Anna", comment: "Truly life-changing session!", rating: 5 },
        { user: "Mike", comment: "Accurate and insightful — highly recommend!", rating: 5 },
        { user: "Sophie", comment: "Very professional and detailed reading.", rating: 4.5 },
      ],
    };

    setCourse(mockCourse);
  }, [id]);

  const handleSendMessage = () => {
    if (!input.trim()) return;
    setChatMessages((prev) => [...prev, { sender: "You", text: input }]);
    setInput("");
  };

  if (!course)
    return (
      <div className="flex justify-center items-center h-screen text-gray-600">
        Loading course details...
      </div>
    );

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center py-12 px-4">
      <div className="max-w-5xl w-full bg-white shadow-xl rounded-2xl overflow-hidden">
        {/* Top section with image */}
        <div className="relative h-64 w-full">
          <img src={course.image} alt={course.name} className="object-cover w-full h-full" />
        </div>

        {/* Details */}
        <div className="p-8">
          <h1 className="text-4xl font-bold mb-2">{course.name}</h1>
          <p className="text-gray-600 mb-1">By {course.prophet}</p>
          <p className="text-sm text-gray-500 mb-4">{course.experience}</p>

          <p className="text-lg text-gray-700 leading-relaxed mb-6">{course.description}</p>

          <div className="mb-6">
            <h3 className="font-semibold text-lg mb-2">Specialties:</h3>
            <div className="flex flex-wrap gap-2">
              {course.specialties.map((spec, idx) => (
                <span
                  key={idx}
                  className="bg-indigo-100 text-indigo-700 px-3 py-1 rounded-full text-sm"
                >
                  {spec}
                </span>
              ))}
            </div>
          </div>

          <div className="flex justify-between items-center border-t pt-4">
            <div>
              <p className="text-2xl font-semibold text-green-600">฿{course.price}</p>
              <p className="text-yellow-500 mt-1">⭐ {course.rating} / 5</p>
            </div>

            <button
              onClick={() =>
                router.push(`/chat/${course.prophet.toLowerCase().replace(/\s+/g, "-")}`)
              }
              className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-3 rounded-lg font-medium"
            >
              Start Chat
            </button>
          </div>

          {/* Reviews Section */}
          <div className="mt-8 border-t pt-6">
            <h3 className="text-xl font-semibold mb-4">Customer Reviews</h3>
            <div className="space-y-3">
              {course.reviews.map((r, i) => (
                <div key={i} className="border rounded-lg p-4 bg-gray-50 hover:shadow-sm">
                  <div className="flex justify-between">
                    <span className="font-medium text-gray-800">{r.user}</span>
                    <span className="text-yellow-500">⭐ {r.rating}</span>
                  </div>
                  <p className="text-gray-600 mt-2">{r.comment}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Optional Chat Box */}
      {showChat && (
        <div className="fixed bottom-6 right-6 w-80 bg-white rounded-xl shadow-2xl border flex flex-col overflow-hidden">
          <div className="bg-indigo-600 text-white px-4 py-3 flex justify-between items-center">
            <span className="font-semibold">Chat with {course.prophet}</span>
            <button onClick={() => setShowChat(false)} className="text-white">
              ✖
            </button>
          </div>
          <div className="flex-1 p-3 space-y-2 overflow-y-auto max-h-64">
            {chatMessages.length === 0 ? (
              <p className="text-gray-500 text-sm text-center mt-4">
                No messages yet. Start the conversation!
              </p>
            ) : (
              chatMessages.map((msg, idx) => (
                <div
                  key={idx}
                  className={`p-2 rounded-lg text-sm ${
                    msg.sender === "You"
                      ? "bg-indigo-100 text-indigo-900 self-end"
                      : "bg-gray-100 text-gray-800"
                  }`}
                >
                  <strong>{msg.sender}:</strong> {msg.text}
                </div>
              ))
            )}
          </div>
          <div className="border-t flex">
            <input
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Type your message..."
              className="flex-1 px-3 py-2 text-sm outline-none"
            />
            <button
              onClick={handleSendMessage}
              className="bg-indigo-600 text-white px-4 hover:bg-indigo-700"
            >
              Send
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
