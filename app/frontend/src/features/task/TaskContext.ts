import { TaskListStore } from "./stores/TaskListStore";

// INFO: 画面識別
export const TaskScreen = {
  List: "TaskList",
  Create: "TaskCreate",
  Edit: "TaskEdit",
} as const;

// INFO: 画面識別一覧
export type TaskScreen = (typeof TaskScreen)[keyof typeof TaskScreen];

// INFO: コンテキスト構造
type TaskContextState = {
  screen: TaskScreen;
  stores: {
    taskList?: TaskListStore;
  };
};

// INFO: コンテキスト構造のstoresキー一覧
type StoreKey = keyof TaskContextState["stores"];

// INFO: タスクコンテキスト
export class TaskContext {
  private active: TaskContextState | null = null;
  private prepared: TaskContextState | null = null;

  // FUNCTION: 遷移準備
  // TODO: Backend引数パラメーター
  async prepare(screen: TaskScreen): Promise<void> {
    if (this.prepared !== null) {
      console.error("TaskContext: Prepared already exists");
      throw new Error("予期しないエラーが発生しました。");
    }

    if (screen === TaskScreen.List) {
      const store = await TaskListStore.create();

      this.prepared = {
        screen: screen,
        stores: {
          taskList: store,
        },
      };
    } else if (screen === TaskScreen.Create) {
      this.prepared = {
        screen: screen,
        stores: {},
      };
    }
  }

  // FUNCTION: 遷移確定
  commit(screen: TaskScreen): void {
    if (this.prepared === null) {
      console.error("TaskContext: Prepared does not exist");
      throw new Error("予期しないエラーが発生しました。");
    }

    if (this.prepared.screen !== screen) {
      console.error("TaskContext: Prepared screen mismatch");
      throw new Error("予期しないエラーが発生しました。");
    }

    this.active = this.prepared;
    this.prepared = null;
  }

  // FUNCTION: データ取得
  getData(screen: TaskScreen, key: StoreKey) {
    if (this.active === null || this.active.screen !== screen) {
      console.error("TaskContext: Active screen mismatch");
      throw new Error("予期しないエラーが発生しました。");
    }

    const store = this.active.stores[key];

    if (store === undefined) {
      console.error("TaskContext: Store does not exist", key);
      throw new Error("予期しないエラーが発生しました。");
    }

    return store.getData();
  }
}
