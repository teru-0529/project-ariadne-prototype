-- order_no の値を新たに採番する。
-- 番号体系は、ORD-YYYYMMDD-999 （処理日付ごとに001からの連番で取得する）
CREATE OR REPLACE FUNCTION received_order.generate_order_no(
  p_row received_order.orders
)
RETURNS varchar
LANGUAGE plpgsql
AS $$
DECLARE
  v_process_date text;
  v_max_no integer;
BEGIN
  -- 採番に使用する処理日付を取得する。
  v_process_date := to_char(current_date, 'YYYYMMDD');

  -- 同一処理日付の採番処理を直列化する。
  -- Transaction終了時にLockは自動的に解放される。
  PERFORM pg_advisory_xact_lock(
    hashtext('received_order.order_no.' || v_process_date)
  );

  -- 同一処理日付ですでに採番されている最大番号を取得する。
  SELECT COALESCE(MAX(RIGHT(order_no, 3)::integer), 0)
    INTO v_max_no
    FROM received_order.orders
   WHERE order_no LIKE 'ORD-' || v_process_date || '-%';

  -- 3桁で採番可能な上限を超えた場合はエラーとする。
  IF v_max_no >= 999 THEN
    RAISE EXCEPTION 'order_no sequence exceeded the daily limit (999): %', v_process_date;
  END IF;

  -- 最大番号に1を加算し、3桁ゼロ埋めして返却する。
  RETURN 'ORD-' || v_process_date || '-' || lpad((v_max_no + 1)::text, 3, '0');
END;
$$;
