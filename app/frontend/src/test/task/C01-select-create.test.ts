import { test, expect } from "vitest";
import { render, screen, fireEvent } from "@testing-library/svelte";

import App from "../../App.svelte";

// TEST:
test("C01 Createを選択する", async () => {
  // INFO: GIVEN
  render(App);

  // Task Side Navigation
  const createNavigation = screen.getByRole("button", { name: "Create" });

  // INFO: ACTION
  await fireEvent.click(createNavigation);

  // INFO: VERIFY
  // Task Side Navigation
  expect(createNavigation.getAttribute("aria-pressed")).toBe("true");

  // main Area - Header
  expect(screen.getByRole("heading", { name: "Task 作成" })).toBeTruthy();

  // 入力項目の存在と初期値
  const title = screen.getByRole("textbox", { name: "タイトル" });
  expect(title).toBeTruthy();
  expect((title as HTMLInputElement).value).toBe("");

  const description = screen.getByRole("textbox", { name: "説明" });
  expect(description).toBeTruthy();
  expect((description as HTMLTextAreaElement).value).toBe("");

  const priority = screen.getByRole("spinbutton", { name: "優先度" });
  expect(priority).toBeTruthy();
  expect(priority.getAttribute("aria-valuenow")).toBe("3");

  const status = screen.getByRole("radio", { name: "TODO" });
  expect(status).toBeTruthy();
  expect((status as HTMLInputElement).checked).toBe(true);

  // 各項目のValidation Errorが表示されない
  expect(screen.queryByText(/Title.*Error/i)).toBeNull();
  expect(screen.queryByText(/Description.*Error/i)).toBeNull();
  expect(screen.queryByText(/Priority.*Error/i)).toBeNull();
  expect(screen.queryByText(/Status.*Error/i)).toBeNull();

  // Create Button
  const createButton = screen.getByRole("button", { name: "Create Task" });
  expect(createButton).toBeTruthy();
  // TODO: Validation処理構築後にコメント解除
  // expect((createButton as HTMLButtonElement).disabled).toBe(true);
});
