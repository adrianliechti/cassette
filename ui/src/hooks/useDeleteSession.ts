import { useMutation, useQueryClient } from '@tanstack/react-query';

export default function useDeleteSession() {
  // const [loading, setLoading] = useState<boolean>(false);
  // const [error, setError] = useState<string | null>(null);

  // const deleteFn = async (id: string) => {
  //   setLoading(true);
  //   setError(null);

  //   await fetch(`/sessions/${id}`, {
  //     method: 'DELETE',
  //   });

  //   setLoading(false);
  // };

  // return [deleteFn, loading, error];

  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (id: string) => fetch(`/sessions/${id}`, { method: 'DELETE' }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sessions'] });
    },
  });

  return mutation;
}
