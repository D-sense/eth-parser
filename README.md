## REST API for Ethereum Blockchain Parser
This REST API allows you to interact with an Ethereum blockchain parser that can query transactions for subscribed addresses.

### Endpoints
```azure
GET /current_block
```
Returns the last parsed block number.

**Response**:
```azure
HTTP/1.1 200 OK
Content-Type: application/json

{
"current_block": 123456
}
```

```azure
POST /subscribe/:address
```
Subscribes to incoming/outgoing transactions for the specified Ethereum address.

**Parameters**
address (string, required) - Ethereum address to subscribe to.
**Response**:
```azure
HTTP/1.1 200 OK
Content-Type: application/json

{
    "message": "Successfully subscribed to address 0x1234567890abcdef"
}
```

```azure
GET /transactions/:address
```
Returns a list of inbound or outbound transactions for the specified Ethereum address.

**Parameters**
address (string, required) - Ethereum address to retrieve transactions for.

**Response**:
```azure
HTTP/1.1 200 OK
Content-Type: application/json
{
"transactions": [
        {
            "hash": "0x1234567890abcdef",
            "from": "0x1234567890abcdef",
            "to": "0xabcdef1234567890",
            "value": "1000000000000000000",
            "timestamp": 1645000000
        },
        {
            "hash": "0xabcdef1234567890",
            "from": "0xabcdef1234567890",
            "to": "0x1234567890abcdef",
            "value": "500000000000000000",
            "timestamp": 1644900000
        }
    ]
}
```

### Error Responses
If an error occurs while processing the request, the API will return an error response with a corresponding status code and message.
Example Error Response:
```azure
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
    "status": "400",
    "message": "Address is invalid"
}
```

### Running the Application
- Clone the repository to your local machine.
- Install Go and make sure it is added to your PATH.
- Rename `.env.example` to `.env`, Grab your cloudflare or infura project ID, set `ETHEREUM_GATEWAY_URL` to the project ID value, and source the file (or, you could simply `export` the env variable).
NOTE: I used https://mainnet.infura.io/v3 to test.
- Open a terminal and navigate to `server` directory in the project.
- Build and start the server by using the command `make build-run` from the root directory . 
Alternatively, you can run this program in a container (Dockerfile is provided), by running this command (ensure to have Makefile program installed):
Build the image: `make docker-build`
Run the image: `make docker-run` 
- The server will start listening on localhost:8080. You can use a tool like curl or a web browser to interact with the API.