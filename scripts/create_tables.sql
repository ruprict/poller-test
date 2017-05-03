create table bopis_orders (
  id int not null identity(1, 1) primary key,
  brand nvarchar(30) not null,
  shipnode nvarchar(30) not null,
  order_id nvarchar(50) not null,
  customer_name nvarchar(100),
  order_date datetime,
  order_acknowledged bit default '0');
