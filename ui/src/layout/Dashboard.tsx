import { useState } from "react";
import List from "../List";
import Player from "../Player";
import Loader from "../components/Loader";
import useSessions, { Session } from "../hooks/useSessions";

export default function Dashboard() {
  const { data, loading } = useSessions();

  const [currentSession, setCurrentSession] = useState<Session | undefined>();

  if (loading) {
    return (
      <div className="flex flex-col items-center gap-2 justify-center pt-10">
        <Loader />
        Sessions Loading
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto w-full flex gap-4 flex-grow">
      <List sessions={data} setCurrentSession={setCurrentSession} />
      <Player session={currentSession} />
    </div>
  );
}
