import { TaskListStore } from "./stores/TaskListStore";
import { createProcessedError } from "../../lib/errors/ProcessedError";

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
      throw createProcessedError("TaskContext: Prepared already exists");
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
      throw createProcessedError("TaskContext: Prepared does not exist");
    }

    if (this.prepared.screen !== screen) {
      throw createProcessedError("TaskContext: Prepared screen mismatch");
    }

    this.active = this.prepared;
    this.prepared = null;
  }

  // FUNCTION: データ取得
  getData(screen: TaskScreen, key: StoreKey) {
    if (this.active === null || this.active.screen !== screen) {
      throw createProcessedError("TaskContext: Active screen mismatch");
    }

    const store = this.active.stores[key];

    if (store === undefined) {
      throw createProcessedError(`TaskContext: Store does not exist (${key})`);
    }

    return store.getData();
  }
}
