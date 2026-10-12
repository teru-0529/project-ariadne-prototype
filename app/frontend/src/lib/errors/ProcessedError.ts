// INFO: 共通プロセスエラー
export class ProcessedError extends Error {
  readonly Processed = true;

  constructor(message: string) {
    super(message);
  }
}

// FUNCTION: Processエラーの生成
export function createProcessedError(
  logMessage: string,
  userMessage = "予期しないエラーが発生しました。",
): ProcessedError {
  console.error(logMessage);
  return new ProcessedError(userMessage);
}

// FUNCTION: ARIADNE対処済みエラーの判定
export function isProcessedError(error: unknown): boolean {
  if (error instanceof ProcessedError) {
    return true;
  }

  if (error instanceof Error) {
    const cause = error.cause;

    if (cause !== null && typeof cause === "object" && "Processed" in cause) {
      return cause.Processed === true;
    }
  }

  return false;
}

// FUNCTION: ARIADNE共通エラーメッセージへの変換
export function toErrorMessage(error: unknown): string {
  if (isProcessedError(error)) {
    return (error as Error).message;
  }

  console.error("Unhandled error:", error);
  return "予期しないエラーが発生しました。";
}
