import { SiApple, SiFirefoxbrowser, SiGooglechrome, SiMicrosoftedge, SiSafari, SiWindows } from '@icons-pack/react-simple-icons';
import { useSearchParams } from 'react-router-dom';
import UAParser from 'ua-parser-js';
import { Session } from '../../hooks/useSessions';
import { cn } from '../../utils';

interface ListProps {
  sessions: Session[];
}

export default function List({ sessions }: ListProps) {
  const [searchParams, setSearchParams] = useSearchParams();

  const currentSessionId = searchParams.get('id');

  return (
    <div className="w-80 flex-shrink-0 flex-grow-0 rounded-md border border-gray-200">
      <div className="border-b border-gray-200 p-4 font-bold">Sessions</div>
      <ul role="list" className="divide-y divide-gray-200">
        {sessions
          .toSorted((a, b) => {
            const dateA = new Date(a.created);
            const dateB = new Date(b.created);

            return dateB.getTime() - dateA.getTime();
          })
          .map((session) => {
            const data = new Date(session.created);

            const date = data.toLocaleDateString(undefined, {
              year: 'numeric',
              month: 'short',
              day: 'numeric',
              hour: '2-digit',
              minute: '2-digit',
            });

            return (
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
                  <ListIcons session={session} />
                  <div>{date}</div>
                </div>
              </li>
            );
          })}
      </ul>
    </div>
  );
}

export function ListIcons({ session }: { session: Session }) {
  let agent = new UAParser(session.userAgent);

  const icons: JSX.Element[] = [];

  const className = 'h-3 w-3 text-gray-400';

  switch (agent.getBrowser().name) {
    case 'Chrome':
      icons.push(<SiGooglechrome className={className} />);
      break;
    case 'Firefox':
      icons.push(<SiFirefoxbrowser className={className} />);
      break;
    case 'Safari':
      icons.push(<SiSafari className={className} />);
      break;
    case 'Edge':
      icons.push(<SiMicrosoftedge className={className} />);
      break;
    default:
      break;
  }

  switch (agent.getOS().name) {
    case 'Mac OS':
      icons.push(<SiApple className={className} />);
      break;
    case 'Windows':
      icons.push(<SiWindows className={className} />);
      break;
    default:
      break;
  }

  return <div className="flex gap-2">{icons}</div>;
}
