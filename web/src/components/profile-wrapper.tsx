"use client";

import dynamic from "next/dynamic";

const ProfileForm = dynamic(() => import("./profile-form"), {
  ssr: false,
  loading: () => (
    <div className="w-96 bg-white shadow-xl p-6 rounded-xl">
      <div className="h-6 w-40 bg-gray-200 rounded mb-4" />
      <div className="space-y-3">
        <div className="h-10 bg-gray-200 rounded" />
        <div className="h-10 bg-gray-200 rounded" />
        <div className="h-10 bg-gray-200 rounded" />
      </div>
    </div>
  ),
});

export default function ProfileWrapper() {
  return <ProfileForm />;
}
