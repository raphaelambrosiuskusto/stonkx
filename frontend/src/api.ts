import { Candle } from "./state";

type ApiResponse = {
  time: number[];
  high: number[];
  low: number[];
  open: number[];
  close: number[];
};

export async function fetchCandles(): Promise<Candle[]> {
  const response = await fetch("/api/getcharts");
  const raw: ApiResponse = await response.json();
  return raw.time.map((t, i) => ({
    time: t,
    high: raw.high[i]! / 100,
    low: raw.low[i]! / 100,
    open: raw.open[i]! / 100,
    close: raw.close[i]! / 100,
  }));
}
