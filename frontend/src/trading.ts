import { state } from "./state";

export function buy(percentage: number): void {
  const candle = state.candles[state.cursor - 1];
  if (candle === undefined) return;

  const units = Math.floor((state.balance * percentage) / candle.close);
  if (units === 0) return;
  state.position += units;
  state.buyPrice = candle.close;
  state.balance -= units * candle.close;
}

export function sell(): void {
  if (state.position === 0) return;
  const candle = state.candles[state.cursor - 1];
  if (candle === undefined) return;
  state.balance += state.position * candle.close;
  state.position = 0;
  state.buyPrice = 0;
}
