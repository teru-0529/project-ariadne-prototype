import { test, expect } from "vitest";
import { render, screen, within } from "@testing-library/svelte";

import App from "../../App.svelte";

// TEST:
test("S01 Applicationを起動する", () => {
  // INFO: ACTION
  render(App);

  // INFO: VERIFY
  // Header
  expect(screen.getByRole("heading", { name: "PROJECT ARIADNE" })).toBeTruthy();

  // Primary Navigation
  const primaryNavigation = screen.getByRole("navigation", { name: "Primary" });
  const taskNavigation = within(primaryNavigation).getByRole("button", { name: "Task" });
  const configNavigation = within(primaryNavigation).getByRole("button", { name: "Config" });
  expect(taskNavigation.getAttribute("aria-current")).toBe("page");
  expect(configNavigation.getAttribute("aria-current")).toBeNull();

  // Task Side Navigation
  const taskSideNavigation = screen.getByRole("navigation", { name: "Task" });
  const listNavigation = within(taskSideNavigation).getByRole("button", { name: "List" });
  const createNavigation = within(taskSideNavigation).getByRole("button", { name: "Create" });
  expect(listNavigation.getAttribute("aria-current")).toBeNull();
  expect(createNavigation.getAttribute("aria-current")).toBeNull();

  // Main Area
  expect(screen.getByText("Select an operation.")).toBeTruthy();

  // Footer
  expect(screen.getByRole("contentinfo")).toBeTruthy();
});
