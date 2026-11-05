import { test, expect, vi } from "vitest";
import { render, screen, fireEvent, within, waitFor } from "@testing-library/svelte";

import App from "../../App.svelte";
import { TaskBoundary } from "../../../bindings/ariadne/backend/adapters/wails";

vi.mock("../../../bindings/ariadne/backend/adapters/wails", () => ({
  TaskBoundary: {
    CreateTask: vi.fn(),
  },
}));

// INFO: C01 Createを選択する
async function givenC01CreateSelected() {
  render(App);
  await fireEvent.click(screen.getByRole("button", { name: "Create" }));
}

// INFO: C03 Create Buttonを押す（Valid）
async function givenC03CreateSubmitted() {
  await fireEvent.input(screen.getByRole("textbox", { name: "タイトル" }), {
    target: { value: "Test" },
  });
  await fireEvent.input(screen.getByRole("textbox", { name: "説明" }), {
    target: { value: "Test Description" },
  });
  await setPriorityFive();
  await fireEvent.click(screen.getByRole("radio", { name: "DOING" }));

  await fireEvent.click(screen.getByRole("button", { name: "Create Task" }));
}

// INFO: 初期値3からPriorityを5に変更する
async function setPriorityFive() {
  const priority = screen.getByRole("spinbutton", { name: "優先度" });
  await fireEvent.keyDown(priority, { key: "ArrowUp" });
  await fireEvent.keyDown(priority, { key: "ArrowUp" });
}

// TEST:
test("C05 Taskを登録する - Success", async () => {
  // INFO: GIVEN
  vi.mocked(TaskBoundary.CreateTask).mockResolvedValue("01TESTTASKID");
  await givenC01CreateSelected();
  await givenC03CreateSubmitted();

  // INFO: ACTION
  const dialog = screen.getByRole("dialog");
  await fireEvent.click(within(dialog).getByRole("button", { name: "Confirm" }));

  // INFO: VERIFY
  await waitFor(() => {
    // Confirmation Modal
    expect(screen.queryByRole("dialog")).toBeNull();

    // Success Toast
    const toast = screen.getByRole("alert");
    expect(toast.textContent).toMatch(/登録|作成|成功/);

    // Title
    expect((screen.getByRole("textbox", { name: "タイトル" }) as HTMLInputElement).value).toBe("");

    // Description
    expect((screen.getByRole("textbox", { name: "説明" }) as HTMLTextAreaElement).value).toBe("");

    // Priority
    expect(screen.getByRole("spinbutton", { name: "優先度" }).getAttribute("aria-valuenow")).toBe("3");

    // Status
    expect((screen.getByRole("radio", { name: "TODO" }) as HTMLInputElement).checked).toBe(true);

    // Create Button
    expect((screen.getByRole("button", { name: "Create Task" }) as HTMLButtonElement).disabled).toBe(true);
  });
});

// TEST:
test("C05 Taskを登録する - Write Failure", async () => {
  // INFO: GIVEN
  const error = new Error("Taskの登録に失敗しました。", {
    cause: {
      Processed: true,
      Message: "Taskの登録に失敗しました。",
    },
  });
  vi.mocked(TaskBoundary.CreateTask).mockRejectedValue(error);

  await givenC01CreateSelected();
  await givenC03CreateSubmitted();

  // INFO: ACTION
  const dialog = screen.getByRole("dialog");
  await fireEvent.click(within(dialog).getByRole("button", { name: "Confirm" }));

  // INFO: VERIFY
  await waitFor(() => {
    // Confirmation Modal
    expect(screen.queryByRole("dialog")).toBeNull();

    // Error Toast
    const toast = screen.getByRole("alert");
    expect(toast.textContent).toContain("Taskの登録に失敗しました。");

    // Title
    expect((screen.getByRole("textbox", { name: "タイトル" }) as HTMLInputElement).value).toBe("Test");

    // Description
    expect((screen.getByRole("textbox", { name: "説明" }) as HTMLTextAreaElement).value).toBe("Test Description");

    // Priority
    expect(screen.getByRole("spinbutton", { name: "優先度" }).getAttribute("aria-valuenow")).toBe("5");

    // Status
    expect((screen.getByRole("radio", { name: "DOING" }) as HTMLInputElement).checked).toBe(true);

    // Create Button
    expect((screen.getByRole("button", { name: "Create Task" }) as HTMLButtonElement).disabled).toBe(false);
  });
});
