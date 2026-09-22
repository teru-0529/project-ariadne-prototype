-- 受注ステータスを判定・登録する。
CREATE OR REPLACE FUNCTION received_order.judge_status(
  p_row received_order.order_details
)
RETURNS received_order.order_details
LANGUAGE plpgsql
AS $$
DECLARE
v_preparing_count INT;
v_canceled_count INT;
v_finished_count INT;
v_total_count INT;
v_status received_order.order_status;
BEGIN
  -- 受注明細のステータス別件数を取得する
    SELECT
    COUNT(*) FILTER (WHERE status = 'PREPARING'),
    COUNT(*) FILTER (WHERE status = 'CANCELED'),
    COUNT(*) FILTER (WHERE status IN ('COMPLETED', 'CANCELED')),
    COUNT(*)
  INTO
    v_preparing_count,
    v_canceled_count,
    v_finished_count,
    v_total_count
  FROM received_order.order_details
  WHERE order_no = p_row.order_no;

  -- ステータスを判定する
  IF v_preparing_count = v_total_count THEN
    v_status := 'PREPARING';

  ELSIF v_canceled_count = v_total_count THEN
    v_status := 'CANCELED';

  ELSIF v_finished_count = v_total_count THEN
    v_status := 'COMPLETED';

  ELSE
    v_status := 'IN_PROGRESS';
  END IF;

  -- 受注ステータスを登録する。
  UPDATE received_order.orders
  SET status = v_status
  WHERE order_no = p_row.order_no;

  RETURN p_row;
END;
$$;
