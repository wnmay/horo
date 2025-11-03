CREATE USER payment_user WITH ENCRYPTED PASSWORD 'payment_pass';
CREATE DATABASE paymentdb OWNER payment_user;

CREATE USER order_user WITH ENCRYPTED PASSWORD 'order_pass';
CREATE DATABASE orderdb OWNER order_user;