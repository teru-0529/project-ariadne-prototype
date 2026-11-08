import { test, expect, vi } from "vitest";
import { render, screen, within, fireEvent } from "@testing-library/svelte";

import App from "../../App.svelte";

// INFO: Task詳細取得のMockデータ
const { initialValue } = vi.hoisted(() => ({
  initialValue: {
    ID: "01M465BRVCYNR0H09A30JE03PW",
    Title: "Task A",
    Description: "詳細取得時の説明",
    Priority: 1,
    Status: "DOING",
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
test("M01 ReferenceからEditへ切り替える - Success", async () => {
  // INFO: GIVEN
  await givenL01ListSelected();
  await givenL02TaskSelected();
  await screen.findByRole("heading", { name: "Task 詳細" });

  // INFO: ACTION
  await fireEvent.click(screen.getByRole("button", { name: "Edit" }));

  // INFO: VERIFY
  // Mode Toggle
  expect(screen.getByRole("button", { name: "Reference" }).getAttribute("aria-pressed")).toBe("false");
  expect(screen.getByRole("button", { name: "Edit" }).getAttribute("aria-pressed")).toBe("true");

  // ID
  expect(screen.getByText(initialValue.ID)).toBeTruthy();

  // title
  const title = screen.getByRole("textbox", { name: "タイトル" });
  expect((title as HTMLInputElement).value).toBe(initialValue.Title);

  // description
  const description = screen.getByRole("textbox", { name: "説明" });
  expect((description as HTMLTextAreaElement).value).toBe(initialValue.Description);

  // priority
  const priority = screen.getByRole("spinbutton", { name: "優先度" });
  expect(priority.getAttribute("aria-valuenow")).toBe(String(initialValue.Priority));

  // status
  const status = screen.getByRole("radio", { name: initialValue.Status });
  expect((status as HTMLInputElement).checked).toBe(true);

  // Validation Error
  expect(screen.queryByText("必須項目です")).toBeNull();
  expect(screen.queryByText("50文字以内で入力してください")).toBeNull();
  expect(screen.queryByText("優先度が5の場合、説明は必須です")).toBeNull();

  // Modify Button
  const modifyButton = screen.getByRole("button", { name: "Modify Task" });
  expect((modifyButton as HTMLButtonElement).disabled).toBe(true);
});
