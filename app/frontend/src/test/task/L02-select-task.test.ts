import { test, expect, vi } from "vitest";
import { render, screen, within, fireEvent } from "@testing-library/svelte";

import App from "../../App.svelte";
import { TaskBoundary } from "../../../bindings/ariadne/backend/adapters/wails";

vi.mock("../../../bindings/ariadne/backend/adapters/wails", () => ({
  TaskBoundary: {
    ListTasks: vi.fn().mockResolvedValue([
      {
        ID: "01M465B5JK59B98Q8AXTJ6HSDR",
        Title: "Task C",
        Priority: 5,
        Status: "TODO",
      },
      {
        ID: "01M465BRVCYNR0H09A30JE03PW",
        Title: "Task A",
        Priority: 1,
        Status: "DOING",
      },
    ]),
    GetTask: vi.fn(),
  },
}));

// TEST:
test("L02 Taskを選択する - Success", async () => {
  // INFO: GIVEN
  const initialValue = {
    ID: "01M465BRVCYNR0H09A30JE03PW",
    Title: "Task A",
    Description: "詳細取得時の説明",
    Priority: 1,
    Status: "DOING",
  };
  vi.mocked(TaskBoundary.GetTask).mockResolvedValue(initialValue);

  render(App);

  // Task Side Navigation
  const listNavigation = screen.getByRole("button", { name: "List" });
  await fireEvent.click(listNavigation);

  // Task List
  const taskList = await screen.findByRole("table");
  const rows = within(taskList).getAllByRole("row");

  // 一覧取得後にGetTaskのMockを変更する
  const changedValue = {
    ID: "01M465BRVCYNR0H09A30JE03PW",
    Title: "Task A - Updated",
    Description: "詳細取得時の説明",
    Priority: 1,
    Status: "DOING",
  };
  vi.mocked(TaskBoundary.GetTask).mockResolvedValue(changedValue);

  // INFO: ACTION
  await fireEvent.click(rows[2]);

  // INFO: VERIFY
  // Task Side Navigation
  expect(listNavigation.getAttribute("aria-pressed")).toBe("true");

  // Main Area - Task詳細画面
  expect(await screen.findByRole("heading", { name: "Task 詳細" })).toBeTruthy();
  expect(screen.queryByRole("heading", { name: "Task 一覧" })).toBeNull();

  // Task詳細画面 - 参照モード
  expect(screen.getByText("01M465BRVCYNR0H09A30JE03PW")).toBeTruthy();
  expect(screen.getByText("Task A - Updated")).toBeTruthy();
  expect(screen.getByText("詳細取得時の説明")).toBeTruthy();
  expect(screen.getByText("★☆☆☆☆")).toBeTruthy();
  expect(screen.getByText("DOING")).toBeTruthy();

  // Error Toast
  expect(screen.queryByRole("alert")).toBeNull();
});

// TEST:
test("L02 Taskを選択する - Not Found", async () => {
  // INFO: GIVEN
  const error = new Error("指定されたタスクが見つかりませんでした。", {
    cause: {
      Processed: true,
      Message: "指定されたタスクが見つかりませんでした。",
    },
  });
  vi.mocked(TaskBoundary.GetTask).mockRejectedValue(error);

  render(App);

  // Task Side Navigation
  const listNavigation = screen.getByRole("button", { name: "List" });
  await fireEvent.click(listNavigation);

  // Task List
  const taskList = await screen.findByRole("table");
  const rows = within(taskList).getAllByRole("row");

  // INFO: ACTION
  await fireEvent.click(rows[2]);

  // INFO: VERIFY
  // Task Side Navigation
  expect(listNavigation.getAttribute("aria-pressed")).toBe("true");

  // Main Area（遷移しない）
  expect(screen.getByRole("heading", { name: "Task 一覧" })).toBeTruthy();
  expect(screen.queryByRole("heading", { name: "Task 詳細" })).toBeNull();

  // Error Toast
  const errorToast = await screen.findByRole("alert");
  expect(errorToast.textContent).toContain("指定されたタスクが見つかりませんでした。");
});

// TEST:
test("L02 Taskを選択する - Read Failure", async () => {
  // INFO: GIVEN
  const error = new Error("タスクを読み込めませんでした。", {
    cause: {
      Processed: true,
      Message: "タスクを読み込めませんでした。",
    },
  });
  vi.mocked(TaskBoundary.GetTask).mockRejectedValue(error);

  render(App);

  // Task Side Navigation
  const listNavigation = screen.getByRole("button", { name: "List" });
  await fireEvent.click(listNavigation);

  // Task List
  const taskList = await screen.findByRole("table");
  const rows = within(taskList).getAllByRole("row");

  // INFO: ACTION
  await fireEvent.click(rows[2]);

  // INFO: VERIFY
  // Task Side Navigation
  expect(listNavigation.getAttribute("aria-pressed")).toBe("true");

  // Main Area（遷移しない）
  expect(screen.getByRole("heading", { name: "Task 一覧" })).toBeTruthy();
  expect(screen.queryByRole("heading", { name: "Task 詳細" })).toBeNull();

  // Error Toast
  const errorToast = await screen.findByRole("alert");
  expect(errorToast.textContent).toContain("タスクを読み込めませんでした。");
});
