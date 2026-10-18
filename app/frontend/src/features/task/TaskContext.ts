import { TaskListStore } from "./stores/TaskListStore";
import { TaskDetailStore } from "./stores/TaskDetailStore";
import { createProcessedError } from "../../lib/errors/ProcessedError";

// INFO: 画面識別
export const TaskScreen = {
  List: "TaskList",
  Reference: "TaskReference",
  Create: "TaskCreate",
  Edit: "TaskEdit",
} as const;

// INFO: 画面識別一覧
export type TaskScreen = (typeof TaskScreen)[keyof typeof TaskScreen];

// INFO: 画面遷移情報
export type TaskNavigation =
  | { screen: typeof TaskScreen.List }
  | { screen: typeof TaskScreen.Reference; id: string }
  | { screen: typeof TaskScreen.Create };

// INFO: コンテキスト構造
type TaskContextState = {
  screen: TaskScreen;
  stores: {
    taskList?: TaskListStore;
    taskDetail?: TaskDetailStore;
  };
};

// INFO: コンテキスト構造のstoresキー一覧
type StoreKey = keyof TaskContextState["stores"];

// INFO: タスクコンテキスト
export class TaskContext {
  private active: TaskContextState | null = null;
  private prepared: TaskContextState | null = null;

  // FUNCTION: 遷移準備
  async prepare(request: TaskNavigation): Promise<void> {
    if (this.prepared !== null) {
      throw createProcessedError("TaskContext: Prepared already exists");
    }

    if (request.screen === TaskScreen.List) {
      // INFO: ListTask
      const store = await TaskListStore.create();

      this.prepared = {
        screen: request.screen,
        stores: {
          taskList: store,
        },
      };
    } else if (request.screen === TaskScreen.Reference) {
      // INFO: ReferenceTask
      const store = await TaskDetailStore.create(request.id);

      this.prepared = {
        screen: request.screen,
        stores: {
          taskDetail: store,
        },
      };
    } else if (request.screen === TaskScreen.Create) {
      // INFO: CreateTask
      this.prepared = {
        screen: request.screen,
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
  getData<K extends StoreKey>(
    screen: TaskScreen,
    key: K,
  ): ReturnType<NonNullable<TaskContextState["stores"][K]>["getData"]> {
    if (this.active === null || this.active.screen !== screen) {
      throw createProcessedError("TaskContext: Active screen mismatch");
    }

    const store = this.active.stores[key];

    if (store === undefined) {
      throw createProcessedError(`TaskContext: Store does not exist (${key})`);
    }

    return store.getData() as ReturnType<NonNullable<TaskContextState["stores"][K]>["getData"]>;
  }
}
