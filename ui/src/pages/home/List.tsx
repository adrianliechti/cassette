import { useSearchParams } from 'react-router-dom';
import { Session } from '../../hooks/useSessions';
import { cn, formatTime } from '../../utils';
import { SessionIcons } from './SessionIcons';

interface ListProps {
  sessions: Session[];
}

export default function List({ sessions }: ListProps) {
  const [searchParams, setSearchParams] = useSearchParams();

  const currentSessionId = searchParams.get('id');

  const preparedList = sessions.reduce(
    (acc, session) => {
      const day = new Date(session.created).toLocaleDateString(undefined, {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      });

      if (!acc[day]) {
        acc[day] = [];
      }

      acc[day].push(session);

      return acc;
    },
    {} as Record<string, Session[]>,
  );

  return (
    <div className="w-80 flex-shrink-0 flex-grow-0 rounded-md border border-gray-200">
      <div className="flex h-16 items-center border-b border-gray-200 p-4 font-bold">Sessions</div>
      <nav className="h-full max-h-[calc(100dvh-12rem)] flex-grow overflow-y-auto">
        {Object.keys(preparedList).map((date) => (
          <div key={date} className="relative">
            <div className="sticky top-0 z-10 flex items-center border-b border-gray-200 bg-gray-100 px-4 py-1 text-sm font-medium">
              <h3>{date}</h3>
            </div>
            <ul role="list" className="divide-low/10 dark:divide-low-dark/20 relative z-0 divide-y">
              {preparedList[date].map((session) => (
                <li
                  key={session.id}
                  className={cn('flex cursor-pointer flex-col gap-2 px-4 py-3 hover:bg-gray-50', {
                    'bg-emerald-500/10 hover:bg-emerald-500/10': session.id === currentSessionId,
                  })}
                  onClick={() => {
                    setSearchParams({ id: session.id });
                  }}
                >
                  <div className="text-xs font-semibold">{session.origin || session.id}</div>
                  <div className="flex items-center justify-between text-xs text-gray-500">
                    <SessionIcons session={session} />
                    <div>{formatTime(session.created)}</div>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </nav>
    </div>
  );
}
