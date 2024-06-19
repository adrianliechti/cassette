import React from "react";
import { cn } from "../utils";

interface CardProps {
  children: React.ReactNode;
  className?: string;
}

export default function Card({ children, className }: CardProps) {
  return (
    <div
      className={cn(
        "overflow-hidden rounded-md border-gray-200 border",
        className
      )}
    >
      <div className="px-4 py-5 sm:p-6">{children}</div>
    </div>
  );
}
