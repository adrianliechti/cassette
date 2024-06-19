import { useEffect, useState } from "react";
import type { eventWithTime } from "@rrweb/types";

export const API_URL = "http://localhost:3000";

export default function useEvents({ id }: { id?: string }) {
  const [sessions, setSessions] = useState<eventWithTime[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) {
      return;
    }

    fetch(`${API_URL}/sessions/${id}/events`)
      .then((res) => res.json())
      .then((data) => {
        setSessions(data);
        setLoading(false);
      });
  }, [id]);

  return { data: sessions, loading };
}
