import React from 'react';
import { cn } from '../utils';

interface CardProps {
  children: React.ReactNode;
  className?: string;
}

export default function Card({ children, className }: CardProps) {
  return (
    <div className={cn('overflow-hidden rounded-md border border-gray-200', className)}>
      <div className="px-4 py-5 sm:p-6">{children}</div>
    </div>
  );
}
