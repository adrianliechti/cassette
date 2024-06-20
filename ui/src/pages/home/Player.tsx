import { useEffect, useRef, useState } from 'react';
import SelectIcon from '../../assets/undraw_select.svg?react';

import { playerMetaData } from '@rrweb/types';
import { Calendar, TimerIcon, Trash2Icon } from 'lucide-react';
import { useSearchParams } from 'react-router-dom';
import rrwebPlayer from 'rrweb-player';
import 'rrweb-player/dist/style.css';
import Loader from '../../components/Loader';
import useDeleteSession from '../../hooks/useDeleteSession';
import useEvents from '../../hooks/useEvents';
import { formatDate, formatTime, getFormatedTimeDiff } from '../../utils';

export default function Player() {
  const [searchParams, setSearchParams] = useSearchParams();

  const playerRef = useRef<rrwebPlayer | null>(null);
  const playerContainerRef = useRef<HTMLDivElement | null>(null);

  const currentSessionId = searchParams.get('id');
  const [metadata, setMetadata] = useState<playerMetaData>();

  const { data: events, loading } = useEvents({ id: currentSessionId });
  const [deleteSession] = useDeleteSession();

  useEffect(() => {
    if (events && events.length > 2 && playerContainerRef.current) {
      const player = playerRef.current?.getReplayer();

      if (player) {
        player.destroy();
        playerContainerRef.current.innerHTML = '';
      }

      playerRef.current = new rrwebPlayer({
        target: document.getElementById('player-container')!,
        props: {
          events: events ?? [],
          width: playerContainerRef.current.clientWidth,
        },
      });

      setMetadata(playerRef.current.getMetaData());
    }
  }, [playerContainerRef.current, events]);

  if (!currentSessionId) {
    return (
      <PlayerWrapper>
        <div className="flex flex-col items-center gap-4 py-32">
          <SelectIcon className="h-40 w-auto" />
          <div className="text-lg font-bold">Select a session to play</div>
        </div>
      </PlayerWrapper>
    );
  }

  if (loading) {
    return (
      <PlayerWrapper>
        <div className="flex flex-col items-center justify-center gap-2 pt-32">
          <Loader />
          Events Loading
        </div>
      </PlayerWrapper>
    );
  }

  return (
    <PlayerWrapper>
      <div className="flex items-center justify-between border-b border-gray-200 p-4">
        <div className="flex flex-col gap-1">
          <div className="text-sm font-bold">{currentSessionId}</div>
          <div className="flex gap-4 text-gray-500">
            <div className="flex items-center gap-1 text-xs">
              <Calendar className="h-4 w-4" />
              {`${formatDate(metadata?.startTime)} - ${formatTime(metadata?.endTime)}`}
            </div>
            <div className="flex items-center gap-1 text-xs">
              <TimerIcon className="h-4 w-4" />
              {getFormatedTimeDiff(metadata?.startTime, metadata?.endTime)}
            </div>
          </div>
        </div>

        <div
          onClick={() => {
            if (confirm(`Delete Session «${currentSessionId}»?`)) {
              deleteSession(currentSessionId).then(() => {
                if (searchParams.has('id')) {
                  searchParams.delete('id');
                  setSearchParams(searchParams);
                }
              });
            }
          }}
          className="cursor-pointer p-2 text-gray-500 hover:text-rose-500"
        >
          <Trash2Icon className="h-5 w-5" />
        </div>
      </div>
      <div id="player-container" ref={playerContainerRef} className="[&_.rr-player]:shadow-none" />
    </PlayerWrapper>
  );
}

function PlayerWrapper({ children }: { children: React.ReactNode }) {
  return <div className="flex-grow rounded-md border border-gray-200">{children}</div>;
}
