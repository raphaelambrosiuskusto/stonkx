import { Chart, Tooltip, LinearScale } from "chart.js";
import {
  CandlestickController,
  CandlestickElement,
} from "chartjs-chart-financial";
import { state, Candle } from "./state";
import ZoomPlugin from "chartjs-plugin-zoom";

Chart.register(
  Tooltip,
  LinearScale,
  CandlestickController,
  CandlestickElement,
  ZoomPlugin,
);

let chart: Chart | null = null;

export function initChart(): void {
  const canvas = document.getElementById("chart") as HTMLCanvasElement;
  const ctx = canvas.getContext("2d")!;
  chart = new Chart(ctx, {
    type: "candlestick",
    data: {
      datasets: [{ data: [] }],
    },
    options: {
      animation: false,
      scales: {
        x: {
          type: "linear",
          ticks: {
            callback: (value: number | string) => {
              const totalSeconds =
                typeof value === "string" ? parseInt(value) : value;
              const h = Math.floor(totalSeconds / 3600);
              const m = Math.floor((totalSeconds % 3600) / 60);
              const s = Math.floor(totalSeconds % 60);
              return [h, m, s].map((n) => String(n).padStart(2, "0")).join(":");
            },
          },
        },
        y: {
          type: "linear",
          grace: "10%",
          ticks: {
            callback: (value: number | string) => {
              const n = typeof value === "string" ? parseFloat(value) : value;
              return n.toFixed(2);
            },
          },
        },
      },

      plugins: {
        zoom: {
          zoom: {
            wheel: { enabled: true },
          },
          pan: {
            enabled: true,
            mode: "xy",
          },
        },
      },
    },
  });
}

export function updateChart(): void {
  if (chart === null) return;
  const visible: Candle[] = state.candles.slice(
    0,
    //Math.max(0, state.cursor - 64),
    state.cursor,
  );
  const data = visible.map((c) => ({
    x: c.time,
    o: c.open,
    h: c.high,
    l: c.low,
    c: c.close,
  }));

  chart.data.datasets[0]!.data = data;
  chart.update();
}
