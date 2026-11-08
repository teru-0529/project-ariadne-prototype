// INFO: ステータスの種類
export const TASK_STATUSES = ["TODO", "DOING", "DONE"] as const;

// INFO: 入力値集合
export type TaskValues = {
  title: string;
  description: string;
  priority: number;
  status: string;
};
