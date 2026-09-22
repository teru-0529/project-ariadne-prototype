-- 受注明細にキャンセル数を登録する。
CREATE OR REPLACE FUNCTION received_order.set_cancel_quantity(
  p_row received_order.cancel_instructions
)
RETURNS received_order.cancel_instructions
LANGUAGE plpgsql
AS $$
DECLARE
BEGIN
  -- キャンセル数を登録する。
  UPDATE received_order.order_details
  SET cancel_quantity = cancel_quantity + p_row.quantity
  WHERE order_no = p_row.order_no AND detail_no = p_row.detail_no;

  RETURN p_row;
END;
$$;
