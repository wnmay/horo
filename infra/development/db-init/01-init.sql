CREATE USER payment_user WITH ENCRYPTED PASSWORD 'payment_pass';
CREATE DATABASE paymentdb OWNER payment_user;