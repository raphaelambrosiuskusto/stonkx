import { state } from "./state";
import { initChart, updateChart } from "./chart";
import { fetchCandles } from "./api";
import { buy, sell } from "./trading";

async function main(): Promise<void> {
  const candles = await fetchCandles();
  state.candles = candles;
  console.log("candles loaded:", candles.length);
  initChart();
  state.cursor = 64;
  console.log("cursor:", state.cursor);
  updateChart();

  document.getElementById("btn-buy")!.addEventListener("click", openBuyPopup);
  document.getElementById("btn-sell")!.addEventListener("click", () => {
    sell();
    updateUI();
  });
  document.getElementById("btn-start")!.addEventListener("click", start);
  document
    .getElementById("popup-close")!
    .addEventListener("click", closeBuyPopup);
  document.getElementById("btn-confirm")!.addEventListener("click", confirmBuy);

  document.querySelectorAll(".pct-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      document
        .querySelectorAll(".pct-btn")
        .forEach((b) => b.classList.remove("selected"));
      btn.classList.add("selected");
      document.getElementById("btn-confirm")!.removeAttribute("disabled");
    });
  });
}

function openBuyPopup(): void {
  state.isPaused = true;
  const popup = document.getElementById("buy-popup")!;
  popup.style.display = "block";
  document
    .querySelectorAll(".pct-btn")
    .forEach((b) => b.classList.remove("selected"));
  document.getElementById("btn-confirm")!.setAttribute("disabled", "");
}

function closeBuyPopup(): void {
  state.isPaused = false;
  document.getElementById("buy-popup")!.style.display = "none";
}

function confirmBuy(): void {
  //finds the one button that has both classes — this is how we read which percentage was chosen
  const selected = document.querySelector(".pct-btn.selected");
  if (selected === null) return;

  const pct = parseFloat((selected as HTMLElement).dataset["pct"] ?? "0");
  buy(pct);
  closeBuyPopup();
  updateUI();
}

function start(): void {
  const btn = document.getElementById("btn-start")!;
  if (state.intervalId === null) {
    state.isRunning = true;
    btn.textContent = "Pause";

    state.intervalId = setInterval(() => {
      if (!state.isRunning || state.isPaused) return;
      if (state.cursor >= state.candles.length) return;
      state.cursor++;
      updateChart();
      updateUI();
    }, 1000);
  } else {
    state.isRunning = !state.isRunning;
    btn.textContent = state.isRunning ? "Pause" : "Resume";
  }
}

function updateUI(): void {
  const candle = state.candles[state.cursor - 1];
  const price = candle?.close ?? 0;

  document.getElementById("current-price")!.textContent =
    state.candles[state.cursor - 1]?.close.toFixed(2) ?? "-";

  document.getElementById("balance")!.textContent = state.balance.toFixed(2);

  document.getElementById("position")!.textContent = (
    state.position * price
  ).toFixed(2);
}

main();
