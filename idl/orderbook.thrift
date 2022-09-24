namespace go tradesoft.exchange.order_book

struct GetExchangeOrderBookReq {
    1: required string ExchangeName ( api.path = "exchange_name" ),
    2: optional string Symbol ( api.body = "symbol" ),
    3: optional i32 OrderType ( api.body = "order_type" ),
    4: optional i32 OrderBy ( api.body = "order_by" ),
}

struct GetExchangeOrderBookResp {
    1: required i32 BizCode ( api.body = "biz_code" ),
    2: required string ErrMsg ( api.body = "err_msg" ),
    3: required map<string, Symbol> Symbols ( api.body = "symbols" ),
}

struct Symbol {
    1: required SymbolItem Bid,
    2: required SymbolItem Ask,
}

struct SymbolItem {
    1: required string PxAvg,    // float number truncated to 2 digits, the price average
    2: required string QtyTotal, // float number truncated to 2 digits, the quantity total
}

struct GetExchangeMetadataReq {
    1: required string ExchangeName ( api.path = "exchange_name" ),
}

struct GetExchangeMetadataResp {
    1: required i32 BizCode ( api.body = "biz_code" ),
    2: required string ErrMsg ( api.body = "err_msg" ),
    3: required map<string, string> Metadata ( api.body = "metadata" ),
}

struct UpdateExchangeMetadataReq {
    1: required string ExchangeName ( api.path = "exchange_name" ),
}

struct UpdateExchangeMetadataResp {
    1: required i32 BizCode ( api.body = "biz_code" ),
    2: required string ErrMsg ( api.body = "err_msg" ),
}

service OrderBookService {
    // GetExchangeOrderBook returns the price average and the total quantities of asks and bids for this exchange_name
    GetExchangeOrderBookResp GetExchangeOrderBook(1: GetExchangeOrderBookReq req) ( api.get = "exchanges/:exchange_name/order-books" ),
    // GetExchangeMetadata returns the list of exchange metadata stored in the database
    GetExchangeMetadataResp GetExchangeMetadata(1: GetExchangeMetadataReq req) ( api.get = "exchanges/:exchange_name/metadata" ),
    // UpdateExchangeMetadata updates the exchange metadata by uploading a csv file
    UpdateExchangeMetadataResp UpdateExchangeMetadata(1: UpdateExchangeMetadataReq req) ( api.post = "exchanges/:exchange_name/metadata" ),
}