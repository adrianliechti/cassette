import Loader from '../../components/Loader';
import useSessions from '../../hooks/useSessions';
import List from './List';
import Player from './Player';

export default function Home() {
  const { data, loading } = useSessions();

  if (loading) {
    return (
      <div className="flex w-full flex-col items-center gap-2 pt-16">
        <Loader />
        Sessions Loading
      </div>
    );
  }

  return (
    <div className="mx-auto flex w-full max-w-7xl flex-grow gap-4">
      <List sessions={data} />
      <Player />
    </div>
  );
}
