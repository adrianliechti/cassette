import { useEffect, useRef } from "react";
import SelectIcon from "./assets/undraw_select.svg?react";
import useEvents from "./hooks/useEvents";
import { Session } from "./hooks/useSessions";

import rrwebPlayer from "rrweb-player";
import "rrweb-player/dist/style.css";
import Loader from "./components/Loader";

interface PlayerProps {
  session?: Session;
}

export default function Player({ session }: PlayerProps) {
  const playerRef = useRef<rrwebPlayer | null>(null);
  const playerContainerRef = useRef<HTMLDivElement | null>(null);

  const { data: events, loading } = useEvents({ id: session?.id });

  useEffect(() => {
    if (
      events &&
      events.length > 2 &&
      !playerRef.current &&
      playerContainerRef.current
    ) {
      playerRef.current = new rrwebPlayer({
        target: document.getElementById("player-container")!,
        props: {
          events: events ?? [],
        },
      });
    }
  }, [events]);

  if (!session) {
    return (
      <PlayerWrapper>
        <div className="flex flex-col items-center gap-4 pt-32">
          <SelectIcon className="h-40 w-auto" />
          <div className="text-lg font-bold">Select a session to play</div>
        </div>
      </PlayerWrapper>
    );
  }

  if (loading) {
    return (
      <PlayerWrapper>
        <div className="flex flex-col items-center gap-2 justify-center pt-32">
          <Loader />
          Events Loading
        </div>
      </PlayerWrapper>
    );
  }

  return (
    <PlayerWrapper>
      <div
        id="player-container"
        ref={playerContainerRef}
        className="h-full [&>*]:shadow-none"
      />
    </PlayerWrapper>
  );
}

function PlayerWrapper({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex-grow rounded-md border-gray-200 border">
      {children}
    </div>
  );
}
