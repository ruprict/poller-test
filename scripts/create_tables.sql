create table bopis_orders (
  id integer not null GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1)
  ,brand varchar(30) not null
  ,store_id varchar(30) not null
  ,order_id varchar(50) not null
  ,customer_name varchar(100)
  ,order_date date not null
  ,order_acknowledged char not null default '0'
  ,PRIMARY KEY (id)
);
