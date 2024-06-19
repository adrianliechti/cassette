import { HelpCircle } from "lucide-react";
import CassetteIcon from "../assets/cassette.svg?react";

export default function Header() {
  return (
    <div className="border-b border-gray-200">
      <div className="flex justify-between max-w-7xl py-4 items-center mx-auto">
        <div className="flex items-center gap-4">
          <CassetteIcon className="h-10 w-10" />
          <span className="font-bold text-2xl text-emerald-600">Cassette</span>
        </div>
        <HelpCircle className="h-6 w-6 text-gray-500" />
      </div>
    </div>
  );
}
