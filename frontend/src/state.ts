export type Candle = {
  time: number;
  high: number;
  low: number;
  open: number;
  close: number;
};

export type State = {
  candles: Candle[];
  cursor: number;
  balance: number;
  position: number;
  isPaused: boolean;
  isRunning: boolean;
  intervalId: number | null;
  buyPrice: number;
};

export const state: State = {
  candles: [],
  cursor: 0,
  balance: 100000,
  position: 0,
  isPaused: false,
  isRunning: false,
  intervalId: null,
  buyPrice: 0,
};
