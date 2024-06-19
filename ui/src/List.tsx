import { Session } from "./hooks/useSessions";

interface ListProps {
  sessions: Session[];
  setCurrentSession: (session: Session) => void;
}

export default function List({ sessions, setCurrentSession }: ListProps) {
  return (
    <div className="rounded-md border-gray-200 border w-80 flex-grow-0 flex-shrink-0">
      <div className="font-bold p-4 border-b border-gray-200">Sessions</div>
      <ul role="list" className="divide-y divide-gray-200">
        {sessions.map((session) => {
          const data = new Date(session.created);

          const date = data.toLocaleDateString(undefined, {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
          });

          return (
            <li
              key={session.id}
              className="px-4 py-3 flex flex-col gap-1 text-end cursor-pointer hover:bg-gray-50"
              onClick={() => setCurrentSession(session)}
            >
              <div className="text-xs font-semibold">{session.id}</div>
              <div className="text-xs text-gray-500">{date}</div>
            </li>
          );
        })}
      </ul>
    </div>
  );
}
