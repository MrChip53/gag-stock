import { useGlobalStore } from "@/stores/global";
import type { StockItem } from "@/types/StockItem";
import { useEffect, useState } from "react";

function formatCountdown(ms: number): string {
  if (ms <= 0) return "0:00";
  const totalSeconds = Math.floor(ms / 1000);
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;
  return `${minutes}:${seconds.toString().padStart(2, "0")}`;
}

function StockListItem({ item }: { item: StockItem }) {
  const { images: cachedImages } = useGlobalStore();

  const getRemaining = () => Math.max(item.restockTime - Date.now(), 0);
  const [remaining, setRemaining] = useState(getRemaining());

  const [purchased, setPurchased] = useState(false);

  useEffect(() => {
    if (item.restockTime <= 0) {
      setRemaining(0);
      return;
    }
    setRemaining(getRemaining());
    const interval = setInterval(() => {
      setRemaining(getRemaining());
    }, 250);
    return () => clearInterval(interval);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [item.restockTime, item.name]);

  const handleTogglePurchased = () => {
    setPurchased((prev) => !prev);
  };

  return (
    <div
      className={`
        flex items-center gap-4 py-2 px-4
        rounded-xl
        ${purchased ? "bg-green-300/80" : "bg-white/20"}
        backdrop-blur-md
        border border-white/30
        shadow-lg
        transition-transform duration-200
        hover:scale-105
        hover:shadow-[0_0_16px_4px_rgba(99,102,241,0.4)]
        outline outline-2 outline-transparent
        hover:outline-blue-400/60
        relative
        cursor-pointer
      `}
      onClick={handleTogglePurchased}
      tabIndex={0}
      role="button"
      aria-pressed={purchased}
      title={purchased ? "Mark as not purchased" : "Mark as purchased"}
    >
      {item.restockTime > 0 && (
        remaining > 0 ? (
          <span
            className="
                absolute bottom-2 right-3
                text-xs font-semibold text-blue-700 bg-white/80 px-2 py-0.5 rounded
                shadow
                select-none
                z-10
              "
            title="Restock countdown"
          >
            {formatCountdown(remaining)}
          </span>
        ) : (
          <span
            className="
                absolute bottom-2 right-3
                flex items-center justify-center
                text-xs font-semibold text-white bg-red-600/80 px-2 py-0.5 rounded
                shadow
                select-none
                z-10
              "
            title="Expired"
          >
            EXPIRED
          </span>
        )
      )}
      {cachedImages[item.name] ? (
        <img
          src={cachedImages[item.name]}
          alt={item.name}
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
              (item.name.charCodeAt(0) + item.name.length * 37) % 360
            )}, 70%, 60%)`,
          }}
        >
          {item.name.trim().split(/\s+/).length > 1
            ? item.name
                .trim()
                .split(/\s+/)
                .map(word => word[0]?.toUpperCase() || '')
                .join('')
            : item
                .name.trim()
                .slice(0, 2)
                .toUpperCase()}
        </div>
      )}
      <div className="flex flex-1 flex-col gap-1">
        <span className="text-center text-base font-medium">{item.name}</span>
        <span className="text-center text-sm font-medium">{item.count}</span>
      </div>
    </div>
  );
}

export default StockListItem;