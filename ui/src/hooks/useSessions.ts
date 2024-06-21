import { useQuery } from '@tanstack/react-query';

export interface Session {
  id: string;
  created: string;
  origin: string;
  userAgent: string;
}

export default function useSessions() {
  const query = useQuery<Session[]>({
    queryKey: ['sessions'],
    queryFn: () => fetch(`/sessions`).then((res) => res.json()),
  });

  return query;
}
