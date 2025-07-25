import { useGlobalStore } from "@/stores/global";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "./ui/card";
import { useQuery } from "@tanstack/react-query";
import { getApiUrl } from "@/lib/utils";

function WantedListItem({ item }: { item: string }) {
  const { images: cachedImages } = useGlobalStore();

  return (
    <div
      className="
        flex items-center gap-4 py-2 px-4
        rounded-xl
        bg-white/20 backdrop-blur-md
        border border-white/30
        shadow-lg
        transition-transform duration-200
        hover:scale-105
        hover:shadow-[0_0_16px_4px_rgba(99,102,241,0.4)]
        outline outline-2 outline-transparent
        hover:outline-blue-400/60
      "
    >
      {cachedImages[item] ? (
        <img
          src={cachedImages[item]}
          alt={item}
          className="w-12 h-12 object-cover rounded"
        />
      ) : (
        <div
          className={`
            flex items-center justify-center w-12 h-12 rounded
            text-white text-xl font-bold
          `}
          style={{
            backgroundColor: `hsl(${Math.floor(
              (item.charCodeAt(0) + item.length * 37) % 360
            )}, 70%, 60%)`,
          }}
        >
          {item.trim().split(/\s+/).length > 1
            ? item
                .trim()
                .split(/\s+/)
                .map(word => word[0]?.toUpperCase() || '')
                .join('')
            : item
                .trim()
                .slice(0, 2)
                .toUpperCase()}
        </div>
      )}
      <span className="flex-1 text-center text-base font-medium">{item}</span>
    </div>
  );
}

function WantedItemsCard() {
  const { data, isLoading, isError } = useQuery<string[]>({
    queryKey: ['wanted-items'],
    refetchInterval: 30000,
    queryFn: () => fetch(getApiUrl('/wanted')).then(res => res.json()),
    initialData: [],
  });

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Wanted Items</CardTitle>
        <CardDescription>Items that are in stock</CardDescription>
      </CardHeader>
      <CardContent>
        {isLoading && !isError && <p>Loading...</p>}
        {isError && <p>Error loading wanted items</p>}
        {!isError && !isLoading && data.length === 0 && <p>No wanted items</p>}
        {!isError && !isLoading && data.length > 0 && (
          <ul className="flex flex-col gap-3">
            {data.map((item) => (
              <li key={item} className="px-2 py-1">
                <WantedListItem item={item} />
              </li>
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  );
}

export default WantedItemsCard;