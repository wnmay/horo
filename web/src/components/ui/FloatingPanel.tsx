"use client";
import { X } from "lucide-react";

export default function FloatingPanel({
  title,
  children,
  onClose,
}: {
  title: string;
  children: React.ReactNode;
  onClose: () => void;
}) {
  return (
    <div className="fixed bottom-4 right-4 w-[320px] bg-white dark:bg-zinc-900 rounded-lg shadow-lg border border-zinc-200 dark:border-zinc-700 z-50">
      <div className="flex justify-between items-center p-3 border-b border-zinc-200 dark:border-zinc-700">
        <h3 className="font-semibold text-zinc-900 dark:text-zinc-100">{title}</h3>
        <button onClick={onClose} className="text-zinc-500 hover:text-zinc-700">
          <X size={18} />
        </button>
      </div>
      <div className="p-4 max-h-[300px] overflow-y-auto">{children}</div>
    </div>
  );
}
