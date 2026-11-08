import { test, expect, vi } from "vitest";
import { render, screen, within, fireEvent } from "@testing-library/svelte";

import App from "../../App.svelte";
import { TaskBoundary } from "../../../bindings/ariadne/backend/adapters/wails";

// INFO: Task詳細取得のMockデータ
const { initialValue, updatedValue } = vi.hoisted(() => ({
  initialValue: {
    ID: "01M465BRVCYNR0H09A30JE03PW",
    Title: "Task A",
    Description: "詳細取得時の説明",
    Priority: 1,
    Status: "DOING",
  },
  updatedValue: {
    ID: "01M465BRVCYNR0H09A30JE03PW",
    Title: "Task A - Updated",
    Description: "更新後の説明",
    Priority: 3,
    Status: "TODO",
  },
}));

// INFO: L01 Task Listを表示する
async function givenL01ListSelected() {
  render(App);
  await fireEvent.click(screen.getByRole("button", { name: "List" }));
}

// INFO: L02 Taskを選択する
async function givenL02TaskSelected() {
  const taskList = await screen.findByRole("table");
  const rows = within(taskList).getAllByRole("row");
  await fireEvent.click(rows[2]);
}

// INFO: M01 Editに切り替える
async function givenM01EditSelected() {
  await screen.findByRole("heading", { name: "Task 詳細（更新）" });
  await fireEvent.click(screen.getByRole("button", { name: "Edit" }));
}

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
    GetTask: vi.fn().mockResolvedValue(initialValue),
  },
}));

// TEST:
test("M02 EditからReferenceへ切り替える - Success", async () => {
  // INFO: GIVEN
  vi.mocked(TaskBoundary.GetTask).mockResolvedValue(initialValue);

  await givenL01ListSelected();
  await givenL02TaskSelected();
  await givenM01EditSelected();

  // 再取得時のMockデータを変更
  vi.mocked(TaskBoundary.GetTask).mockResolvedValue(updatedValue);

  // INFO: ACTION
  await fireEvent.click(screen.getByRole("button", { name: "Reference" }));

  // INFO: VERIFY
  // Mode Toggle
  expect((await screen.findByRole("button", { name: "Reference" })).getAttribute("aria-pressed")).toBe("true");
  expect(screen.getByRole("button", { name: "Edit" }).getAttribute("aria-pressed")).toBe("false");

  // ID
  expect(screen.getByText(updatedValue.ID)).toBeTruthy();

  // title
  expect(screen.getByText(updatedValue.Title)).toBeTruthy();

  // description
  expect(screen.getByText(updatedValue.Description)).toBeTruthy();

  // priority
  expect(screen.getByRole("img", { name: `優先度 ${updatedValue.Priority}` })).toBeTruthy();

  // status
  expect(screen.getByText(updatedValue.Status)).toBeTruthy();

  // Error Toast
  expect(screen.queryByRole("alert")).toBeNull();
});

// TEST:
test("M02 EditからReferenceへ切り替える - Not Found", async () => {
  // INFO: GIVEN
  vi.mocked(TaskBoundary.GetTask).mockResolvedValue(initialValue);

  await givenL01ListSelected();
  await givenL02TaskSelected();
  await givenM01EditSelected();

  // 再取得時にNot Foundを返す
  const error = new Error("指定されたタスクが見つかりませんでした。", {
    cause: {
      Processed: true,
      Message: "指定されたタスクが見つかりませんでした。",
    },
  });
  vi.mocked(TaskBoundary.GetTask).mockRejectedValue(error);

  // INFO: ACTION
  await fireEvent.click(screen.getByRole("button", { name: "Reference" }));

  // INFO: VERIFY
  // Mode Toggle
  expect((await screen.findByRole("button", { name: "Edit" })).getAttribute("aria-pressed")).toBe("true");
  expect(screen.getByRole("button", { name: "Reference" }).getAttribute("aria-pressed")).toBe("false");

  // ID
  expect(screen.getByText(initialValue.ID)).toBeTruthy();

  // title
  expect((screen.getByRole("textbox", { name: "タイトル" }) as HTMLInputElement).value).toBe(initialValue.Title);

  // description
  expect((screen.getByRole("textbox", { name: "説明" }) as HTMLTextAreaElement).value).toBe(initialValue.Description);

  // priority
  expect(screen.getByRole("spinbutton", { name: "優先度" }).getAttribute("aria-valuenow")).toBe(
    String(initialValue.Priority),
  );

  // status
  expect((screen.getByRole("radio", { name: initialValue.Status }) as HTMLInputElement).checked).toBe(true);

  // Error Toast
  const errorToast = await screen.findByRole("alert");
  expect(errorToast.textContent).toContain("指定されたタスクが見つかりませんでした。");
});
