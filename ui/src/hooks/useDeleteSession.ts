import { useEffect, useState } from 'react';

export default function useDeleteSession({ id }: { id?: string | null }) {
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) {
      return;
    }

    fetch(`/sessions/${id}`, {
      method: 'DELETE',
    })
      .then((res) => res.json())
      .then(() => {
        setLoading(false);
      });
  }, [id]);

  return { loading };
}
