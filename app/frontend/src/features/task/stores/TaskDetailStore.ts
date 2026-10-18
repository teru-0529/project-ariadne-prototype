import { TaskBoundary } from "../../../../bindings/ariadne/backend/adapters/wails";
import type { Task } from "../../../../bindings/ariadne/backend/services/publish/models";
import { createProcessedError } from "../../../lib/errors/ProcessedError";

// INFO: タスク詳細
export class TaskDetailStore {
  private readonly data: Task;

  private constructor(data: Task) {
    this.data = data;
  }

  // FUNCTION: Backend連携
  static async create(id: string): Promise<TaskDetailStore> {
    const data = await TaskBoundary.GetTask(id);

    if (data === null) {
      throw createProcessedError("TaskDetailStore: Unexpected null response");
    }

    return new TaskDetailStore(data);
  }

  // FUNCTION: 取得
  getData(): Task {
    return this.data;
  }
}
