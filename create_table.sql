create table accounts (
	account_id varchar,
	account_name varchar,
	account_email varchar,
	balance float
);

create table transactions (
	transaction_id varchar,
	transaction_type varchar,
	amount float,
	transaction_time timestamp	
);