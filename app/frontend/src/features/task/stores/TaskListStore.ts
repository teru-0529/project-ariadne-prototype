import { TaskBoundary } from "../../../../bindings/ariadne/backend/adapters/wails";
import type { TaskSummary } from "../../../../bindings/ariadne/backend/services/publish/models";

// INFO: タスクリスト
export class TaskListStore {
  private readonly data: TaskSummary[] | null;

  private constructor(data: TaskSummary[] | null) {
    this.data = data;
  }

  // FUNCTION: Backend連携
  static async create(): Promise<TaskListStore> {
    const data = await TaskBoundary.ListTasks();
    return new TaskListStore(data);
  }

  // FUNCTION: 取得
  getData(): TaskSummary[] | null {
    return this.data;
  }
}
