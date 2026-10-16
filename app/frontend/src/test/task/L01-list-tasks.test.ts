import { test, expect, vi } from "vitest";
import { render, screen, within, fireEvent } from "@testing-library/svelte";

import App from "../../App.svelte";
import { TaskBoundary } from "../../../bindings/ariadne/backend/adapters/wails";

vi.mock("../../../bindings/ariadne/backend/adapters/wails", () => ({
  TaskBoundary: {
    ListTasks: vi.fn(),
  },
}));

// TEST:
test("L01 Task一覧を表示する - Success Tasks Exist", async () => {
  // INFO: GIVEN
  const value = [
    {
      ID: "01M465B5JK59B98Q8AXTJ6HSDR",
      Title: "Grandit API仕様の調査と整理、および既存システムとの連携方法の検討",
      Priority: 5,
      Status: "TODO",
    },
    {
      ID: "01M465BRVCYNR0H09A30JE03PW",
      Title: "Teamcenter連携の設計検討",
      Priority: 1,
      Status: "DOING",
    },
    {
      ID: "01M465BRVCDDWKHCTF84YK2QQD",
      Title: "関係者とのキックオフ",
      Priority: 3,
      Status: "DONE",
    },
  ];
  vi.mocked(TaskBoundary.ListTasks).mockResolvedValue(value);

  render(App);

  // Task Side Navigation
  const listNavigation = screen.getByRole("button", { name: "List" });

  // INFO: ACTION
  await fireEvent.click(listNavigation);

  // INFO: VERIFY
  // Task Side Navigation
  expect(listNavigation.getAttribute("aria-pressed")).toBe("true");

  // main Area - Header
  expect(screen.getByRole("heading", { name: "Task 一覧" })).toBeTruthy();
  expect(screen.getByText("全 3 件")).toBeTruthy();

  // main Area - Table
  const taskList = await screen.findByRole("table");
  const rows = within(taskList).getAllByRole("row");
  expect(rows).toHaveLength(4); // Header + 3 Tasks

  const columnHeaders = within(taskList).getAllByRole("columnheader");
  expect(columnHeaders).toHaveLength(3);
  expect(columnHeaders[0].textContent).toBe("タイトル");
  expect(columnHeaders[1].textContent).toBe("優先度");
  expect(columnHeaders[2].textContent).toBe("ステータス");

  const firstRowCells = within(rows[1]).getAllByRole("cell");
  expect(firstRowCells).toHaveLength(3);
  expect(
    within(firstRowCells[0]).getByText("Grandit API仕様の調査と整理、および既存システムとの連携方法の検討"),
  ).toBeTruthy();
  expect(within(firstRowCells[0]).getByText("01M465B5JK59B98Q8AXTJ6HSDR")).toBeTruthy();
  expect(firstRowCells[1].textContent).toBe("★★★★★");
  expect(firstRowCells[2].textContent).toBe("TODO");

  const secondRowCells = within(rows[2]).getAllByRole("cell");
  expect(secondRowCells).toHaveLength(3);
  expect(within(secondRowCells[0]).getByText("Teamcenter連携の設計検討")).toBeTruthy();
  expect(within(secondRowCells[0]).getByText("01M465BRVCYNR0H09A30JE03PW")).toBeTruthy();
  expect(secondRowCells[1].textContent).toBe("★☆☆☆☆");
  expect(secondRowCells[2].textContent).toBe("DOING");

  const thirdRowCells = within(rows[3]).getAllByRole("cell");
  expect(thirdRowCells).toHaveLength(3);
  expect(within(thirdRowCells[0]).getByText("関係者とのキックオフ")).toBeTruthy();
  expect(within(thirdRowCells[0]).getByText("01M465BRVCDDWKHCTF84YK2QQD")).toBeTruthy();
  expect(thirdRowCells[1].textContent).toBe("★★★☆☆");
  expect(thirdRowCells[2].textContent).toBe("DONE");
  // TODO: L02-⑤ 条件付き行装飾の正式設計時にテスト方法を再検討する。
  // TEMP: DONE行の装飾条件となるステータスバッジの属性を検証する。
  expect(within(rows[3]).getByText("DONE").getAttribute("data-status")).toBe("DONE");
  // expect(rows[3].getAttribute("data-status")).toBe("DONE");
});

// TEST:
test("L01 Task一覧を表示する - Success Empty", async () => {
  // INFO: GIVEN
  vi.mocked(TaskBoundary.ListTasks).mockResolvedValue([]);

  render(App);

  // Task Side Navigation
  const listNavigation = screen.getByRole("button", { name: "List" });

  // INFO: ACTION
  await fireEvent.click(listNavigation);

  // INFO: VERIFY
  // Task Side Navigation
  expect(listNavigation.getAttribute("aria-pressed")).toBe("true");

  // Main Area - Header
  expect(screen.getByRole("heading", { name: "Task 一覧" })).toBeTruthy();
  expect(screen.getByText("全 0 件")).toBeTruthy();

  // Main Area - Table
  const taskList = await screen.findByRole("table");
  const rows = within(taskList).getAllByRole("row");
  expect(rows).toHaveLength(2); // Header + No tasks. row

  const columnHeaders = within(taskList).getAllByRole("columnheader");
  expect(columnHeaders).toHaveLength(3);
  expect(columnHeaders[0].textContent).toBe("タイトル");
  expect(columnHeaders[1].textContent).toBe("優先度");
  expect(columnHeaders[2].textContent).toBe("ステータス");

  expect(within(taskList).getByText("No tasks.")).toBeTruthy();
  expect(within(taskList).queryByText("TODO")).toBeNull();
  expect(within(taskList).queryByText("DOING")).toBeNull();
  expect(within(taskList).queryByText("DONE")).toBeNull();
  expect(screen.queryByRole("alert")).toBeNull();
});

// TEST:
test("L01 Task一覧を表示する - Read Failure", async () => {
  // INFO: GIVEN
  const error = new Error("Task Listの読み込みに失敗しました。", {
    cause: {
      Processed: true,
      Message: "Task Listの読み込みに失敗しました。",
    },
  });
  vi.mocked(TaskBoundary.ListTasks).mockRejectedValue(error);

  render(App);

  // Task Side Navigation
  const listNavigation = screen.getByRole("button", { name: "List" });

  // INFO: ACTION
  await fireEvent.click(listNavigation);

  // INFO: VERIFY
  // Task Side Navigation（未選択状態）
  expect(listNavigation.getAttribute("aria-pressed")).toBe("false");

  // Main Area（遷移しない）
  expect(screen.getByText("Select an operation.")).toBeTruthy();
  expect(screen.queryByRole("heading", { name: "Task 一覧" })).toBeNull();

  // Error Toast
  const errorToast = await screen.findByRole("alert");
  expect(errorToast.textContent).toContain("Task Listの読み込みに失敗しました。");
});

// TEST:
test("L01 Task一覧を表示する - Unsupported Major Version", async () => {
  // INFO: GIVEN
  const error = new Error("データのバージョン 3.x には対応していません。対応バージョンは 2.x です。", {
    cause: {
      Processed: true,
      Message: "データのバージョン 3.x には対応していません。対応バージョンは 2.x です。",
    },
  });
  vi.mocked(TaskBoundary.ListTasks).mockRejectedValue(error);

  render(App);

  // Task Side Navigation
  const listNavigation = screen.getByRole("button", { name: "List" });

  // INFO: ACTION
  await fireEvent.click(listNavigation);

  // INFO: VERIFY
  // Task Side Navigation（未選択状態）
  expect(listNavigation.getAttribute("aria-pressed")).toBe("false");

  // Main Area（遷移しない）
  expect(screen.getByText("Select an operation.")).toBeTruthy();
  expect(screen.queryByRole("heading", { name: "Task 一覧" })).toBeNull();

  // Error Toast
  const errorToast = await screen.findByRole("alert");
  expect(errorToast.textContent).toContain("データのバージョン 3.x には対応していません。対応バージョンは 2.x です。");
});
