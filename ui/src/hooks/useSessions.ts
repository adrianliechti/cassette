import { useEffect, useState } from 'react';

export interface Session {
  id: string;
  created: string;
  origin: string;
  userAgent: string;
}

export default function useSessions() {
  const [sessions, setSessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('/sessions')
      .then((res) => res.json())
      .then((data) => {
        setSessions(data);
        setLoading(false);
      });
  }, []);

  return { data: sessions, loading };
}
