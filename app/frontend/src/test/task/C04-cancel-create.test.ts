import { test, expect } from "vitest";
import { render, screen, fireEvent, within } from "@testing-library/svelte";
import App from "../../App.svelte";

// INFO: C01 Createを選択する
async function givenC01CreateSelected() {
  render(App);
  await fireEvent.click(screen.getByRole("button", { name: "Create" }));
}

// INFO: C03 Create Buttonを押す（Valid）
async function givenC03CreateSubmitted() {
  await fireEvent.input(screen.getByRole("textbox", { name: "タイトル" }), { target: { value: "Test" } });

  await fireEvent.click(screen.getByRole("button", { name: "Create Task" }));
}

// TEST:
test("C04 登録をCancelする - Cancel", async () => {
  // INFO: GIVEN
  await givenC01CreateSelected();
  await givenC03CreateSubmitted();

  // INFO: ACTION
  const dialog = screen.getByRole("dialog");
  await fireEvent.click(within(dialog).getByRole("button", { name: "Cancel" }));

  // INFO: VERIFY
  // Confirmation Modal
  expect(screen.queryByRole("dialog")).toBeNull();

  // Create Button
  expect(screen.getByRole("button", { name: "Create Task" })).toBeTruthy();

  // title
  expect((screen.getByRole("textbox", { name: "タイトル" }) as HTMLInputElement).value).toBe("Test");
});
