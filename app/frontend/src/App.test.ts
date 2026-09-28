import { vi, test, expect } from "vitest";
import { render, screen } from "@testing-library/svelte";
import App from "./App.svelte";

// モックセットアップ
vi.mock("../bindings/ariadne/backend/services", () => ({
  AppService: {
    GetAppName: vi.fn().mockResolvedValue("Project ARIADNE Prototype"),
  },
}));

// テスト実行
test("[B] Display App Name", async () => {
  render(App);

  const heading = await screen.findByRole("heading", {
    name: "Project ARIADNE Prototype",
  });

  expect(heading).toBeTruthy();
});

// FIXME:
test("Unclassified", () => {});
