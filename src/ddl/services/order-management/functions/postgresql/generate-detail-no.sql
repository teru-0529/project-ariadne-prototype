-- detail_no の値を新たに採番する。
-- 番号体系は、受注番号単位の連番。
CREATE OR REPLACE FUNCTION received_order.generate_detail_no(
  p_row received_order.order_details
)
RETURNS received_order.order_details
LANGUAGE plpgsql
AS $$
DECLARE
  v_max_no integer;
BEGIN
  -- 同一受注番号の採番処理を直列化する。
  -- Transaction終了時にLockは自動的に解放される。
  PERFORM pg_advisory_xact_lock(
    hashtext('received_order.detail_no.' || p_row.order_no)
  );

  -- 同一受注番号における明細番号の最大値を取得する。
  SELECT COALESCE(MAX(detail_no), 0)
    INTO v_max_no
    FROM received_order.order_details
   WHERE order_no = p_row.order_no;

  -- 最大値に1を加算し、detail_noへ設定する。
  p_row.detail_no := v_max_no + 1;

  RETURN p_row;
END;
$$;
