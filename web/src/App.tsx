import { useEffect } from 'react'
import { env } from '@/env'
import { useQuery } from '@tanstack/react-query';
import { useGlobalStore } from './stores/global';
import WantedItemsCard from './components/WantedItemsCard';
import { getApiUrl } from './lib/utils';
import AllItemsCard from './components/AllItemsCard';

function App() {
  const { data: images, isLoading: isImagesLoading, isError: isImagesError } = useQuery<Record<string, string>>({
    queryKey: ['images'],
    refetchInterval: 60000,
    queryFn: () =>
      fetch(getApiUrl('/images')).then(res => res.json()),
    initialData: {},
  });

  const { setImages: setCachedImages } = useGlobalStore();

  useEffect(() => {
    if (images) {
      setCachedImages(images);
    }
  }, [images]);

  useEffect(() => {
    if (env.VITE_APP_TITLE) {
      document.title = env.VITE_APP_TITLE
    }
  }, [])

  if (isImagesLoading || isImagesError) {
    return <div>Loading...</div>
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-row justify-center gap-4">
        <WantedItemsCard />
      </div>
      <div className="flex flex-row justify-center gap-4">
        <AllItemsCard />
      </div>
    </div>
  )
}

export default App
