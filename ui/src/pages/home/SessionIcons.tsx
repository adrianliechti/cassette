import { SiApple, SiFirefoxbrowser, SiGooglechrome, SiMicrosoftedge, SiSafari, SiWindows } from '@icons-pack/react-simple-icons';
import UAParser from 'ua-parser-js';
import { Session } from '../../hooks/useSessions';

interface SessionIconsProps {
  session: Session;
  large?: boolean;
}

export function SessionIcons({ session, large = false }: SessionIconsProps) {
  let agent = new UAParser(session.userAgent);

  const icons: JSX.Element[] = [];

  let className = 'text-gray-400 h-3 w-3';

  if (large) {
    className = 'text-gray-500 h-4 w-4';
  }

  switch (agent.getBrowser().name) {
    case 'Chrome':
      icons.push(<SiGooglechrome key="chrome" className={className} />);
      break;
    case 'Firefox':
      icons.push(<SiFirefoxbrowser key="firefox" className={className} />);
      break;
    case 'Safari':
      icons.push(<SiSafari key="safari" className={className} />);
      break;
    case 'Edge':
      icons.push(<SiMicrosoftedge key="edge" className={className} />);
      break;
    default:
      break;
  }

  switch (agent.getOS().name) {
    case 'Mac OS':
      icons.push(<SiApple key="apple" className={className} />);
      break;
    case 'Windows':
      icons.push(<SiWindows key="windows" className={className} />);
      break;
    default:
      break;
  }

  return <div className="flex gap-2">{icons}</div>;
}
