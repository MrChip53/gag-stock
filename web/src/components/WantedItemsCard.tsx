import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "./ui/card";
import { useQuery } from "@tanstack/react-query";
import { getApiUrl } from "@/lib/utils";
import type { StockItem } from "@/types/StockItem";
import StockListItem from "./StockListItem";

function WantedItemsCard() {
  const { data, isLoading, isError } = useQuery<StockItem[]>({
    queryKey: ['wanted-items'],
    refetchInterval: 5000,
    queryFn: () => fetch(getApiUrl('/wanted')).then(res => res.json()),
    initialData: [],
  });

  return (
    <Card className="w-1/2">
      <CardHeader>
        <CardTitle>Wanted Items</CardTitle>
        <CardDescription>Wanted items that are in stock</CardDescription>
      </CardHeader>
      <CardContent>
        {isLoading && !isError && <p>Loading...</p>}
        {isError && <p>Error loading wanted items</p>}
        {!isError && !isLoading && (!data || data.length === 0) && <p>No wanted items</p>}
        {!isError && !isLoading && data.length > 0 && (
          <ul className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
          {data.sort((a, b) => a.name.localeCompare(b.name)).map((item) => (
            <li key={`${item.stockTime}-${item.name}-${item.count}`} className="flex-1">
              <StockListItem item={item} />
            </li>
          ))}
        </ul>
        )}
      </CardContent>
    </Card>
  );
}

export default WantedItemsCard;