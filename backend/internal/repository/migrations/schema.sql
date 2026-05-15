create table quotsint (
    id integer primary key generated always as identity,
    ticker_id smallint not null, 
    time int not null,
    high int not null,
    low int not null,
    open int not null,
    close int not null
);



