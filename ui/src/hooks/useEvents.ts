import { useEffect, useState } from "react";
import type { eventWithTime } from "@rrweb/types";

export default function useEvents({ id }: { id?: string }) {
  const [sessions, setSessions] = useState<eventWithTime[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) {
      return;
    }

    fetch(`/sessions/${id}/events`)
      .then((res) => res.json())
      .then((data) => {
        setSessions(data);
        setLoading(false);
      });
  }, [id]);

  return { data: sessions, loading };
}
