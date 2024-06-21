import Loader from '../../components/Loader';
import useSessions from '../../hooks/useSessions';
import Help from '../help/Help';
import List from './List';
import Player from './Player';

export default function Home() {
  const { data, isPending } = useSessions();

  if (isPending) {
    return (
      <div className="flex w-full flex-col items-center gap-2 pt-16">
        <Loader />
        Sessions Loading
      </div>
    );
  }

  if (!data || data.length === 0) {
    return (
      <div className="flex w-full flex-col items-center gap-2 pt-16">
        <Help title="No Sessions Found" />
      </div>
    );
  }

  return (
    <div className="mx-auto flex w-full max-w-7xl gap-4">
      <List sessions={data} />
      <Player />
    </div>
  );
}
