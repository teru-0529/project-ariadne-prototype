import { test, expect } from "vitest";
import { render, screen, within, fireEvent } from "@testing-library/svelte";

import App from "../../App.svelte";

// TEST:
test("S02 Primary NavigationでConfigを選択する", async () => {
  render(App);

  // Primary Navigation
  const primaryNavigation = screen.getByRole("navigation", { name: "Primary" });
  const taskNavigation = within(primaryNavigation).getByRole("button", { name: "Task" });
  const configNavigation = within(primaryNavigation).getByRole("button", { name: "Config" });

  // INFO: ACTION
  await fireEvent.click(configNavigation);

  // INFO: VERIFY
  // Primary Navigation
  expect(taskNavigation.getAttribute("aria-current")).toBeNull();
  expect(configNavigation.getAttribute("aria-current")).toBe("page");

  // Task Side Navigation
  expect(screen.queryByRole("navigation", { name: "Task" })).toBeNull();

  // Main Area
  expect(screen.getByText("Configuration is not available in this prototype.")).toBeTruthy();
});
