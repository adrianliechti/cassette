import { eventWithTime } from '@rrweb/types';
import { useQuery } from '@tanstack/react-query';

export default function useEvents({ id }: { id?: string | null }) {
  const query = useQuery<eventWithTime[]>({
    queryKey: ['events', id],
    queryFn: () => fetch(`/sessions/${id}/events`).then((res) => res.json()),
    enabled: !!id,
  });

  return query;
}
