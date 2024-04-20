INSERT INTO unit_type (id, name) 
VALUES (1, 'Buah'), (2, 'KG');

INSERT INTO mst_service (id, name, price, unit_type_id) 
VALUES (1, 'Laundry Bedcover', 50000, 1), (2, 'Laundry Shirt', 20000, 1), (3, 'Laundry Jean', 10000, 2), (4, 'Laundry Pack', 70000, 1);

INSERT INTO mst_customer (id, name, phone) 
VALUES (1, 'John Doe', '085856203961');

INSERT INTO trx_laundry (id, customer_id, entry_date, finish_date, laundry_detail_id) 
VALUES (1, 1, '2024-12-22', '2024-12-24', 1);

INSERT INTO laundry_detail (id, service_id, quantity, total_price) 
VALUES (1, 4, 10, 70000);
