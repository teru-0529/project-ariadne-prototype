import { test, expect } from "vitest";
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

// TEST:
test("C02 フォーカスアウトする - Single Field Error", async () => {
  // INFO: GIVEN
  await givenC01CreateSelected();

  // title
  const title = screen.getByRole("textbox", { name: "タイトル" });

  // INFO: ACTION
  // title
  await fireEvent.focus(title);
  await fireEvent.focusOut(title);

  // INFO: VERIFY
  // title
  expect(screen.getByText("必須項目です")).toBeTruthy();

  // Create Button
  const createButton = screen.getByRole("button", { name: "Create Task" });
  expect(createButton).toBeTruthy();
  expect((createButton as HTMLButtonElement).disabled).toBe(true);
});

// TEST:
test("C02 フォーカスアウトする - Correlation Errorなし（INITIAL）", async () => {
  // INFO: GIVEN
  await givenC01CreateSelected();

  // title
  await fireEvent.input(screen.getByRole("textbox", { name: "タイトル" }), {
    target: { value: "Test" },
  });

  // priority
  await setPriorityFive();

  // description
  const description = screen.getByRole("textbox", { name: "説明" });

  // INFO: ACTION
  // description
  await fireEvent.focus(description);
  await fireEvent.focusOut(description);

  // INFO: VERIFY
  // description / priority
  expect(screen.queryByText("優先度が5の場合、説明は必須です")).toBeNull();

  // Create Button
  const createButton = screen.getByRole("button", { name: "Create Task" });
  expect(createButton).toBeTruthy();
  expect((createButton as HTMLButtonElement).disabled).toBe(false);
});

// TEST:
test("C02 フォーカスアウトする - Correlation Errorなし（ACTIVE）", async () => {
  // INFO: GIVEN
  await givenC01CreateSelected();

  // title
  await fireEvent.input(screen.getByRole("textbox", { name: "タイトル" }), {
    target: { value: "Test" },
  });

  // priority
  await setPriorityFive();

  // description
  const description = screen.getByRole("textbox", { name: "説明" });
  await fireEvent.input(description, { target: { value: "Details" } });

  // Create Button
  const createButton = screen.getByRole("button", { name: "Create Task" });
  await fireEvent.click(createButton);

  // INFO: ACTION
  // description
  await fireEvent.focus(description);
  await fireEvent.focusOut(description);

  // INFO: VERIFY
  // description / priority
  expect(screen.queryByText("優先度が5の場合、説明は必須です")).toBeNull();

  // Create Button
  expect(createButton).toBeTruthy();
  expect((createButton as HTMLButtonElement).disabled).toBe(false);
});

// TEST:
test("C02 フォーカスアウトする - Correlation Errorあり（ACTIVE）", async () => {
  // INFO: GIVEN
  await givenC01CreateSelected();

  // title
  await fireEvent.input(screen.getByRole("textbox", { name: "タイトル" }), {
    target: { value: "Test" },
  });

  // priority
  await setPriorityFive();

  // description
  const description = screen.getByRole("textbox", { name: "説明" });
  await fireEvent.input(description, { target: { value: "Details" } });

  // Create Button
  const createButton = screen.getByRole("button", { name: "Create Task" });
  await fireEvent.click(createButton);

  // description
  await fireEvent.input(description, { target: { value: "" } });

  // INFO: ACTION
  // description
  await fireEvent.focus(description);
  await fireEvent.focusOut(description);

  // INFO: VERIFY
  // description / priority
  expect(screen.getAllByText("優先度が5の場合、説明は必須です")).toHaveLength(2);

  // Create Button
  expect(createButton).toBeTruthy();
  expect((createButton as HTMLButtonElement).disabled).toBe(true);
});
