import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "./ui/card";
import { useQuery } from "@tanstack/react-query";
import { getApiUrl } from "@/lib/utils";
import type { StockItem } from "@/types/StockItem";
import StockListItem from "./StockListItem";

function AllItemsCard() {
  const { data, isLoading, isError } = useQuery<StockItem[]>({
    queryKey: ['all-items'],
    refetchInterval: 30000,
    queryFn: () => fetch(getApiUrl('/all')).then(res => res.json()),
    initialData: [],
  });

  return (
    <Card className="w-1/2">
      <CardHeader>
        <CardTitle>All Items</CardTitle>
        <CardDescription>All items that are in stock</CardDescription>
      </CardHeader>
      <CardContent>
        {isLoading && !isError && <p>Loading...</p>}
        {isError && <p>Error loading all items</p>}
        {!isError && !isLoading && data.length === 0 && <p>No all items</p>}
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

export default AllItemsCard;