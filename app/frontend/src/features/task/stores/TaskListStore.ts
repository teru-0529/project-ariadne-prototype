import { TaskBoundary } from "../../../../bindings/ariadne/backend/adapters/wails";
import type { TaskSummary } from "../../../../bindings/ariadne/backend/services/publish/models";
import { createProcessedError } from "../../../lib/errors/ProcessedError";

// INFO: タスクリスト
export class TaskListStore {
  private readonly data: TaskSummary[];

  private constructor(data: TaskSummary[]) {
    this.data = data;
  }

  // FUNCTION: Backend連携
  static async create(): Promise<TaskListStore> {
    const data = await TaskBoundary.ListTasks();

    if (data === null) {
      throw createProcessedError("TaskListStore: Unexpected null response");
    }

    return new TaskListStore(data);
  }

  // FUNCTION: 取得
  getData(): TaskSummary[] {
    return this.data;
  }
}
