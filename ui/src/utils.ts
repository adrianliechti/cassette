import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const formatDate = (time?: number | string) => {
  if (!time) return 'n/a';
  return new Date(time).toLocaleString();
};

export const formatTime = (time?: number | string) => {
  if (!time) return 'n/a';
  return new Date(time).toLocaleTimeString();
};

export function getFormatedTimeDiff(startDate?: string | number | Date, endDate?: string | number | Date) {
  if (!startDate || !endDate) {
    return '-';
  }

  const start = new Date(startDate);
  const end = new Date(endDate);

  const diffInSeconds = Math.floor((end.getTime() - start.getTime()) / 1000);

  const hours = Math.floor(diffInSeconds / 3600);
  const minutes = Math.floor((diffInSeconds % 3600) / 60);
  const seconds = diffInSeconds % 60;

  const pad = (num: number) => num.toString().padStart(2, '0');

  if (hours === 0 && minutes === 0) {
    return `${seconds}s`;
  }

  if (hours === 0) {
    return `${pad(minutes)}:${pad(seconds)}`;
  }

  return `${pad(hours)}:${pad(minutes)}:${pad(seconds)}`;
}
