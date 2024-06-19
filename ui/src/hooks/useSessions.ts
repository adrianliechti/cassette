import { useEffect, useState } from "react";

export interface Session {
  id: string;
  created: string;
}

export const API_URL = "http://localhost:3000";

export default function useSessions() {
  const [sessions, setSessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(API_URL + "/sessions")
      .then((res) => res.json())
      .then((data) => {
        setSessions(data);
        setLoading(false);
      });
  }, []);

  return { data: sessions, loading };
}
