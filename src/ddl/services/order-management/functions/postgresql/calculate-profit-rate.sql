-- profit_rate の値を計算する。
CREATE OR REPLACE FUNCTION received_order.calculate_profit_rate(
  p_row received_order.order_details
)
RETURNS received_order.order_details
LANGUAGE plpgsql
AS $$
BEGIN
  -- profit_rateを計算する。
  p_row.profit_rate := (p_row.selling_price - p_row.cost_price) * 100.0 / p_row.selling_price;

  RETURN p_row;
END;
$$;
