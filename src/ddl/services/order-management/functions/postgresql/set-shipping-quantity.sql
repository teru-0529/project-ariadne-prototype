-- 受注明細に出荷数を登録する。
CREATE OR REPLACE FUNCTION received_order.set_shipping_quantity(
  p_row received_order.shipping_instructions
)
RETURNS received_order.shipping_instructions
LANGUAGE plpgsql
AS $$
DECLARE
BEGIN
  -- 出荷数を登録する。
  UPDATE received_order.order_details
  SET shipping_quantity = shipping_quantity + p_row.quantity
  WHERE order_no = p_row.order_no AND detail_no = p_row.detail_no;

  RETURN p_row;
END;
$$;
