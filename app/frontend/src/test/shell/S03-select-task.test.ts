import { test, expect } from "vitest";
import { render, screen, within, fireEvent } from "@testing-library/svelte";

import App from "../../App.svelte";

// TEST:
test("S03 Primary NavigationでTaskを選択する", async () => {
  render(App);

  // Primary Navigation
  const primaryNavigation = screen.getByRole("navigation", { name: "Primary" });
  const taskNavigation = within(primaryNavigation).getByRole("button", { name: "Task" });
  const configNavigation = within(primaryNavigation).getByRole("button", { name: "Config" });

  // INFO: GIVEN
  // S02が保証する機能を利用してconfig ページを表示する
  await fireEvent.click(configNavigation);

  // INFO: ACTION
  await fireEvent.click(taskNavigation);

  // INFO: VERIFY
  // Primary Navigation
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
});
