# HC-Proxy

Tooling for modification response from any http server

## Use Cases
- Inject custom code on each response from webserver
  - Branding
  - Advertisement
  - Analytics
- Inject custom header on each response from webserver
    - Security
    - Tracking

## UML Diagram

```plantuml
@startuml

title Integrate Custom Elements to 3d party Products
actor User
User->HC-Proxy:User Perform Any Request
HC-Proxy->Backend:Proxy Request to Backend

HC-Proxy<-Backend:Regular HTTP response
User<-HC-Proxy:Modified HTTP response
User->CDN:Get JS script from Modified HTTP Response
User<-CDN:Response with JS Script

@enduml

```
## Configuration
App support configuration via environment variables:
- `HC_PROXY_BACKEND_URL` - the address where the proxy will forward all requests (default: `http://hellcorp.com.ua`)
- `HC_PROXY_INJECTION_SCRIPT_SRC` - the address for custom JS that should be injected in the html response (default: `HC_PROXY_BACKEND_URL + "/js/label.js"`)
- `HC_PROXY_BIND_PORT` - port on which proxy will be listened for requests (default: `8980`)
- `HC_PROXY_BIND_IP` - address on which proxy will be listened for requests (default: `0.0.0.0`)

## HOW TO

### Run Demo

`go run main.go`

### Build

`go build main.go`


## License
- MIT

## Disclaimer
 - This code shared as a POC of fully workable proxy, use with care. 