ROLLBACK;

---

BEGIN;
SET LOCAL ariadne.audit = '01M31VPPQ2CNQKEHQMSCTWYSWJ::U9985';

WITH new_order AS(
-- orders
  INSERT INTO received_order.orders(order_pic, customer_id, urgent, note, desired_delivery_time)
  VALUES (null, 'C1234567', false, null, '14:30:15')
  RETURNING order_no
)

-- order_details
INSERT INTO received_order.order_details(order_no, product_no, order_quantity, remaining_quantity, selling_price, cost_price, profit_rate)
SELECT order_no, 'P567342', 3, 3, 900, 600, 33.33333
FROM new_order;

COMMIT;

-- TODO: remaining_quantity
-- TODO: profit_rate
---

BEGIN;
SET LOCAL ariadne.audit = '01M31WCW1ECVDRGAPDY76MVH5C::U1532';

-- shipping_instructions
INSERT INTO received_order.shipping_instructions(order_no, detail_no, person_in_charge, quantity, note)
VALUES ('ORD-20260921-001', 1, 'U1532', 2, '緊急出荷');

-- order_details
UPDATE received_order.order_details
SET shipping_quantity = 2, remaining_quantity = 1, status = 'IN_PROGRESS'
WHERE order_no = 'ORD-20260921-001' AND detail_no = 1;

INSERT INTO received_order.order_details(order_no, product_no, order_quantity, remaining_quantity, selling_price, cost_price, profit_rate)
VALUES ('ORD-20260921-001', 'P256870', 1, 1, 2000, 1500, 25);

-- orders
UPDATE received_order.orders
SET order_pic = 'U1532', status = 'IN_PROGRESS'
WHERE order_no = 'ORD-20260921-001';

COMMIT;

-- TODO: shipping_quantity
-- TODO: remaining_quantity
-- TODO: profit_rate
-- TODO: status
---

BEGIN;
SET LOCAL ariadne.audit = '01M31X7E0YS4BWYBCTBPN574PG::U2376';

-- cancel_instructions
INSERT INTO received_order.cancel_instructions(order_no, detail_no, person_in_charge, quantity, note)
VALUES ('ORD-20260921-001', 1, 'U2376', 1, null);

INSERT INTO received_order.cancel_instructions(order_no, detail_no, person_in_charge, quantity, note)
VALUES ('ORD-20260921-001', 2, 'U2376', 1, null);

-- order_details
UPDATE received_order.order_details
SET cancel_quantity = 1, remaining_quantity = 0, status = 'COMPLETED'
WHERE order_no = 'ORD-20260921-001' AND detail_no = 1;

UPDATE received_order.order_details
SET cancel_quantity = 1, remaining_quantity = 0, status = 'CANCELED'
WHERE order_no = 'ORD-20260921-001' AND detail_no = 2;

-- orders
UPDATE received_order.orders
SET status = 'COMPLETED'
WHERE order_no = 'ORD-20260921-001';

COMMIT;

-- TODO: cancel_quantity
-- TODO: remaining_quantity
-- TODO: status
---
