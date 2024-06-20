import { useState } from 'react';

export default function useDeleteSession(): [(id: string) => Promise<any>, boolean, string | null] {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const deleteFn = async (id: string) => {
    setLoading(true);
    setError(null);

    await fetch(`/sessions/${id}`, {
      method: 'DELETE',
    });

    setLoading(false);
  };

  return [deleteFn, loading, error];
}
