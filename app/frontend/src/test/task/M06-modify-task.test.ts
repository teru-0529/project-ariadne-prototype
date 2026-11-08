import { test, expect, vi } from "vitest";
import { render, screen, within, fireEvent, waitFor } from "@testing-library/svelte";

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
    Title: "Task A",
    Description: "更新後の説明",
    Priority: 3,
    Status: "DOING",
  },
}));

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
    ModifyTask: vi.fn(),
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
  await screen.findByRole("heading", { name: "Task 詳細（参照）" });
  await fireEvent.click(screen.getByRole("button", { name: "Edit" }));
}

// INFO: M04 更新確認Modalを表示する
async function givenM04ModifySubmitted() {
  await fireEvent.input(screen.getByRole("textbox", { name: "説明" }), {
    target: { value: updatedValue.Description },
  });

  const priority = screen.getByRole("spinbutton", { name: "優先度" });
  await fireEvent.keyDown(priority, { key: "ArrowUp" });
  await fireEvent.keyDown(priority, { key: "ArrowUp" });

  await fireEvent.click(screen.getByRole("button", { name: "Modify Task" }));
}

// TEST:
test("M06 Taskを更新する - Success", async () => {
  // INFO: GIVEN
  vi.mocked(TaskBoundary.GetTask).mockResolvedValue(initialValue);
  vi.mocked(TaskBoundary.ModifyTask).mockResolvedValue(undefined);

  await givenL01ListSelected();
  await givenL02TaskSelected();
  await givenM01EditSelected();
  await givenM04ModifySubmitted();

  vi.mocked(TaskBoundary.GetTask).mockResolvedValue(updatedValue);

  // INFO: ACTION
  const dialog = screen.getByRole("dialog");
  await fireEvent.click(within(dialog).getByRole("button", { name: "Confirm" }));

  // INFO: VERIFY
  await waitFor(() => {
    // Confirmation Modal
    expect(screen.queryByRole("dialog")).toBeNull();

    // Success Toast
    const toast = screen.getByRole("alert");
    expect(toast.textContent).toContain("Taskを更新しました。");

    // Reference Mode
    expect(screen.getByRole("button", { name: "Reference" }).getAttribute("aria-pressed")).toBe("true");
    expect(screen.getByRole("button", { name: "Edit" }).getAttribute("aria-pressed")).toBe("false");

    // Updated Values
    expect(screen.getByText(updatedValue.ID)).toBeTruthy();
    expect(screen.getByText(updatedValue.Title)).toBeTruthy();
    expect(screen.getByText(updatedValue.Description)).toBeTruthy();
    expect(screen.getByRole("img", { name: `優先度 ${updatedValue.Priority}` })).toBeTruthy();
    expect(screen.getByText(updatedValue.Status)).toBeTruthy();
  });
});

// TEST:
test("M06 Taskを更新する - Not Found", async () => {
  // INFO: GIVEN
  vi.mocked(TaskBoundary.GetTask).mockResolvedValue(initialValue);

  const error = new Error("指定されたタスクが見つかりませんでした。", {
    cause: {
      Processed: true,
      Message: "指定されたタスクが見つかりませんでした。",
    },
  });
  vi.mocked(TaskBoundary.ModifyTask).mockRejectedValue(error);

  await givenL01ListSelected();
  await givenL02TaskSelected();
  await givenM01EditSelected();
  await givenM04ModifySubmitted();

  // INFO: ACTION
  const dialog = screen.getByRole("dialog");
  await fireEvent.click(within(dialog).getByRole("button", { name: "Confirm" }));

  // INFO: VERIFY
  await waitFor(() => {
    // Confirmation Modal
    expect(screen.queryByRole("dialog")).toBeNull();

    // Error Toast
    const toast = screen.getByRole("alert");
    expect(toast.textContent).toContain("指定されたタスクが見つかりませんでした。");

    // Edit Mode
    expect(screen.getByRole("button", { name: "Edit" }).getAttribute("aria-pressed")).toBe("true");
    expect(screen.getByRole("button", { name: "Reference" }).getAttribute("aria-pressed")).toBe("false");

    // Input Values
    expect(screen.getByText(initialValue.ID)).toBeTruthy();
    expect((screen.getByRole("textbox", { name: "タイトル" }) as HTMLInputElement).value).toBe(initialValue.Title);
    expect((screen.getByRole("textbox", { name: "説明" }) as HTMLTextAreaElement).value).toBe(updatedValue.Description);
    expect(screen.getByRole("spinbutton", { name: "優先度" }).getAttribute("aria-valuenow")).toBe(
      String(updatedValue.Priority),
    );
    expect((screen.getByRole("radio", { name: initialValue.Status }) as HTMLInputElement).checked).toBe(true);
  });
});
