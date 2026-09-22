-- remaining_quantityを計算しstatusを判定する。
CREATE OR REPLACE FUNCTION received_order.calculate_quantity_and_status(
  p_row received_order.order_details
)
RETURNS received_order.order_details
LANGUAGE plpgsql
AS $$
DECLARE
BEGIN
  -- remaining_quantityを計算する。
  p_row.remaining_quantity := p_row.order_quantity - p_row.shipping_quantity - p_row.cancel_quantity;

  -- statusを判定する。
  IF p_row.order_quantity = p_row.remaining_quantity THEN
    p_row.status := 'PREPARING';
  ELSIF p_row.remaining_quantity = 0 AND p_row.shipping_quantity = 0 THEN
    p_row.status := 'CANCELED';
  ELSIF p_row.remaining_quantity = 0 THEN
    p_row.status := 'COMPLETED';
  ELSE
    p_row.status := 'IN_PROGRESS';
  END IF;

  RETURN p_row;
END;
$$;
