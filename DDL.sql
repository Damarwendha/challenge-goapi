CREATE TABLE mst_customer (
    id INT NOT NULL PRIMARY KEY,
    name VARCHAR(200),
    phone VARCHAR(20)
);

CREATE TABLE mst_service (
    id INT NOT NULL PRIMARY KEY,
    name VARCHAR(200),
    price INT,
	unit_type_id INT
);

CREATE TABLE unit_type (
	id INT NOT NULL PRIMARY KEY,
	name VARCHAR(100)
);

CREATE TABLE trx_laundry (
	id INT NOT NULL PRIMARY KEY,
	customer_id INT,
	laundry_detail_id INT,
	entry_date DATE,
	finish_date DATE
);

CREATE TABLE laundry_detail (
	id INT NOT NULL PRIMARY KEY,
	service_id INT,
	quantity INT,
	total_price INT
);

ALTER TABLE mst_service
ADD CONSTRAINT fk_unit_type_id
FOREIGN KEY (unit_type_id)
REFERENCES unit_type(id);

ALTER TABLE laundry_detail
ADD CONSTRAINT fk_service_id
FOREIGN KEY (service_id)
REFERENCES mst_service(id);

ALTER TABLE trx_laundry
ADD CONSTRAINT fk_laundry_detail_id
FOREIGN KEY (laundry_detail_id)
REFERENCES laundry_detail(id);

ALTER TABLE trx_laundry
ADD CONSTRAINT fk_customer_id
FOREIGN KEY (customer_id)
REFERENCES mst_customer(id);
