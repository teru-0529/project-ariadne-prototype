import { test, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/svelte";
import App from "../../App.svelte";

// INFO: C01 Createを選択する
async function givenC01CreateSelected() {
  render(App);
  await fireEvent.click(screen.getByRole("button", { name: "Create" }));
}

// INFO: 初期値3からPriorityを5に変更する
async function setPriorityFive() {
  const priority = screen.getByRole("spinbutton", { name: "優先度" });
  await fireEvent.keyDown(priority, { key: "ArrowUp" });
  await fireEvent.keyDown(priority, { key: "ArrowUp" });
}

// TODO: A09 Confirmation Modal実装後、暫定alertの無効化を削除
vi.spyOn(window, "alert").mockImplementation(() => {});

// TEST:
test("C03 Createボタンを押す - Invalid Single Field", async () => {
  // INFO: GIVEN
  await givenC01CreateSelected();

  // title
  await fireEvent.input(screen.getByRole("textbox", { name: "タイトル" }), {
    target: { value: "A".repeat(51) },
  });

  // priority
  await setPriorityFive();

  // INFO: ACTION
  // Create Button
  const createButton = screen.getByRole("button", { name: "Create Task" });
  await fireEvent.click(createButton);

  // INFO: VERIFY
  // title
  expect(screen.getByText("タイトルは50文字以内で入力してください。")).toBeTruthy();

  // description / priority
  expect(screen.queryByText("優先度が5の場合、説明は必須です。")).toBeNull();

  // Create Button
  expect((createButton as HTMLButtonElement).disabled).toBe(true);
});

// TEST:
test("C03 Createボタンを押す - Invalid Correlation", async () => {
  // INFO: GIVEN
  await givenC01CreateSelected();

  // title
  await fireEvent.input(screen.getByRole("textbox", { name: "タイトル" }), {
    target: { value: "Test" },
  });

  // priority
  await setPriorityFive();

  // INFO: ACTION
  // Create Button
  const createButton = screen.getByRole("button", { name: "Create Task" });
  await fireEvent.click(createButton);

  // INFO: VERIFY
  expect(screen.getByText("優先度が5の場合、説明は必須です。")).toBeTruthy();

  // Create Button
  expect((createButton as HTMLButtonElement).disabled).toBe(true);
});

// TEST:
test("C03 Createボタンを押す - Valid Descriptionなし", async () => {
  // INFO: GIVEN
  await givenC01CreateSelected();

  // title
  await fireEvent.input(screen.getByRole("textbox", { name: "タイトル" }), {
    target: { value: "Test" },
  });

  // INFO: ACTION
  // Create Button
  await fireEvent.click(screen.getByRole("button", { name: "Create Task" }));

  // INFO: VERIFY
  // title
  expect(screen.queryByText("タイトルは必須です。")).toBeNull();

  // description / priority
  expect(screen.queryByText("優先度が5の場合、説明は必須です。")).toBeNull();

  // TODO: A09 Confirmation Modal実装後、Modal表示の検証を追加
});

// TEST:
test("C03 Createボタンを押す - Valid Descriptionあり", async () => {
  // INFO: GIVEN
  await givenC01CreateSelected();

  // title
  await fireEvent.input(screen.getByRole("textbox", { name: "タイトル" }), {
    target: { value: "Test" },
  });

  // priority
  await setPriorityFive();

  // description
  await fireEvent.input(screen.getByRole("textbox", { name: "説明" }), {
    target: { value: "Details" },
  });

  // INFO: ACTION
  // Create Button
  await fireEvent.click(screen.getByRole("button", { name: "Create Task" }));

  // INFO: VERIFY
  // title
  expect(screen.queryByText("タイトルは必須です。")).toBeNull();

  // description / priority
  expect(screen.queryByText("優先度が5の場合、説明は必須です。")).toBeNull();

  // TODO: A09 Confirmation Modal実装後、Modal表示の検証を追加
});
