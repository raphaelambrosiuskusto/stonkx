package repository

const (
	selectMaxTickerID = `select coalesce(max(ticker_id), 0) from quotsint;`

	selectAll = `select * from quotsint;`

	queryTicker = `select time, high, low, open, close from quotsint
	where ticker_id = $1
	order by time ASC;`

	queryTickerUpToTime = `select time, high, low, open, close from quotsint
	where ticker_id = $1 and time < $2
	order by time ASC;`

	insertQuots = `insert into quotsint (ticker_id, time, high, low, open, close) 
    values ($1, $2, $3, $4, $5, $6);`

	deleteQuots = `delete from quotsint
	where "ColumnName" = $1`

	deleteAllQuots = `truncate table quotsint
	restart identity;`

	listTickersAsc = `select distinct ticker_id from quotsint order by ticker_id asc;`

	checkTableExists = `select to_regclass('public.$1') is not null as table_exists;`
)
