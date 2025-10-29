import { ReactNode } from "react";

interface CardProps {
  children: ReactNode;
  className?: string; // optional extra classes
}

export default function Card({ children, className = "" }: CardProps) {
  return (
    <div
      className={`
        border-2 border-solid border-gray-400
        rounded-xl
        p-12
        bg-white dark:bg-gray-800
        shadow-sm
        ${className}
      `}
    >
      {children}
    </div>
  );
}
