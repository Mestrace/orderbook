# OrderBook
TradeSoft is a well-known company in the crypto market that provides powerful statistics, dashboards and metrics to their customers. As a backend software developer in TradeSoft, you were requested to develop a REST API which integrates with most famous crypto exchanges and exposes aggregated data to be visualized in the application. 

## Project Management

### DOING

### TODO

- [Quality 1] Better error handling, biz_code and err_msg https://github.com/Mestrace/orderbook/issues/5
- [Quality 2] Better service observability and metrics https://github.com/Mestrace/orderbook/issues/6

### DONE

- [Perf 2] Use Redis to cache symbol results https://github.com/Mestrace/orderbook/issues/7
- [Perf 1] Concurrently call to blockchain.com api https://github.com/Mestrace/orderbook/pull/4
- [Prototype 1] https://github.com/Mestrace/orderbook/pull/2

## Project Structure
- `biz` business handler.
  - `handler` core service handler, initializes processor with implementations, execute and parses response.
  - `model` generated api model from thrift definition.
  - `router` generate router and middleware from therift definition.
- `cmd` all main to be build.
- `common` contains common functions to be used across the project.
- `conf` contains configuration definition and conf parser.
- `constant` contains constant definitions.
- `domain` contains domain informations
  - `dao` contains dao interface definition for all data access objects.
  - `dto` contains dto functions that converts between internal and api models.
  - `infra` contains implementation of `dao` that mapped to the actual implementation.
  - `processors` contains that implement the core business logic .
  - `resources` contains initializer and getter for actual resources including connection to db and redis.
- `idl` contains thrift definition of the api.
- `sql` contains the create table sql.
- `thrid_party` contains non-project 3rd party code that used.

### Dependency Injection

`processors` will take the abstracted daos and resources as input, and implement the logic on top of that. Eventually inside the handler, implementation will be picked and intialized and injected into the processor.

## Quick Start

### Setup

- Conf
  - Make a copy of `conf/config.json` as `conf/config_local.json` for your own configuration.
- DB
  - Install and start your db instance. For MacOS, `brew install mysql` and `brew services start mysql`
  - Setup your db and modify `conf/config_local.json` to setup the client and dsn
    - `"master_dsn": "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"`
  - Run SQL under `sql/orderbook.sql` inside your mysql shell to migrate the table structure. 
- Redis
  - Install and start your redis instance. For MacOS, `brew install redis` and `brew services start redis`
  - Setup your redis instance and modify `conf/config_local.json` to setup the redis client.
    - `"addr": "localhost:6379"`
    - `"password": ""`
    - `"db": 0`

### Build and Run

Run `make run_orderbook` the service should start.


## API Design

### GET `exchanges/{exchange-name}/order-books`

given the exchange name, returns the order-books which contains the high level statistics of the order book, for each of the symbols under this exchange.

Assuming that each exchange will expose some form of `l3/book` and `symbol` api that we could fetch the data.

#### Request

| Parameter     | Tag       | Type   | Description                                                                    |
|---------------|-----------|--------|--------------------------------------------------------------------------------|
| exchange_name | api.path  | string | The name of the exchange where data came from.                                 |
| symbol        | api.path  | string | Optional, if not specified, will return all symbols.                           |
| order_type    | api.query | string | Optional, will display ask or bid. 0 = all, 1 = ask, 2 = bid.                  |
| order_by      | api.query | string | Optional, will return the order in specified ordering. 1 = alphabetical order. |

#### Response

| Parameter | Tag      | Type      | Description                                           |
|-----------|----------|-----------|-------------------------------------------------------|
| biz_code  | api.body | uint32    | Business Error Code. 0 = success.                     |
| err_msg   | api.body | string    | Error message.                                        |
| symbols   | api.body | []Symbol  | The list of symbols that contains symbol ask and bid. |

Symbol
| Parameter | Type       | Description        |
|-----------|------------|--------------------|
| symbol    | string     | The string symbol. |
| ask       | SymbolItem | The ask.           |
| bid       | SymbolItem | The bid.           |

SymbolItem
| Parameter | Type   | Description                                               |
|-----------|--------|-----------------------------------------------------------|
| px_avg    | string | The average price, float truncated to 2 digit precision.  |
| qty_total | string | The quantity total, float truncated to 2 digit precision. |

### GET `exchanges/{exchange-name}/metadata`

Given the exchange name, return the key-value metadata.

#### Request

| Parameter | Tag | Type | Description |
| --- | --- | --- | --- |
| exchange-name | api.path | string | The name of the exchange where the data came from |

#### Response

| Parameter | Tag | Type | Description |
| --- | --- | --- | --- |
| biz_code | api.body | uint32 | Business error code |
| err_msg | api.body | string | Error message to display |
| metadata | api.body | Metadata | The list of metadata k-v pairs  |

Metadata 

| Parameter   | Type              | Description                                                   |
|-------------|-------------------|---------------------------------------------------------------|
| description | string            | The string description of this exchange.                      |
| web_site    | string            | The string website of this exchange.                          |
| extinfo     | map[string]string | Any additional info that does not fall into the above fields. |

### POST `exchanges/{exchange-name}/metadata`

#### Request

| Parameter | Tag | Type | Description |
| --- | --- | --- | --- |
| exchange-name | api.path | string | The name of the exchange where the data came from |
| file | api.body | multipart/form | The binary data of the csv file. |

#### Response

| Parameter | Variable | Type | Description |
| --- | --- | --- | --- |
| biz_code | body | uint32 | Business error code |
| err_msg | body | string | Error message to display |


## DB Design

### DB `db_orderbook_main`

#### Table `tab_orderbook_exchange_metadata`

```
CREATE TABLE tab_orderbook_exchange_metadata 
  ( 
     `id`         BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT, 
     `exchange`   VARCHAR(32) NOT NULL, 
     `metadata`   TEXT NOT NULL, 
     `created_at` DATETIME NOT NULL, 
     `updated_at` DATETIME NOT NULL, 
     UNIQUE KEY `uniq_exchange` (`exchange`) 
  ) 
AUTO_INCREMENT = 1000, 
ENGINE=InnoDB, 
DEFAULT CHARSET = utf8mb4;
```

## Further Thoughts

Currently for the demonstration purpose, this prototype is designed to be simple and yet is in complete and cannot ensuer service reliablity. In this section we imagine the overall architecture of the service.

![orderbook](https://user-images.githubusercontent.com/26028388/193398172-c8cab40b-093b-4a43-8ebf-e91a8d6953ba.png)

- Access Layer provides some generic and common ability such as service routing, authentication and rate limiting.
- Service Layer is a user-facing layer that provides user models and some form of logic encapsulation that suits the need of clients.
- Logic Layer contains the core logic for this system
  - Account provides user management and permission
  - Metadata provides metadata and api key management. 
  - Fetcher provides encapsulation of data fetching logic and caching. Internally it would route and use different fetcher implementation to obtain the data from the client. The fetcher implementation could either be something that is implemented inside the fetcher service, or an service that provides generic api to obtain.
- Data Layer contains the data infrastrucutres to be used in this system.
- General abilities includes some of the vital features of modern distributed system.
  - Service Discovery ensures the connectivity between different services, and can provide enhanced abilities such as feature branch deployments.
  - Regional / AZ Deployment provides deployement management across multiple availability zones, in order to provide better service for clients across different geographic regions.
  - Disastor Recovery ensure the service reliability when something happens.
  - Observability provides monitoring and alerting management for all services.
  
