HTTP Server From Raw TCP, Then net/http
I built an HTTP server in Go two ways. First on raw TCP sockets with no HTTP library. Then again with net/http. The goal was to learn what a web framework does for me by doing it myself first.
Both versions live in this repo and run on their own:
cmd/rawtcp/main.go    the raw TCP version (about 300 lines)
cmd/nethttp/main.go   the net/http version (about 100 lines)
Why I built it
I shipped a full product in Python and FastAPI. It ran in production. But the framework hid the network from me. I could build a route. I could not tell you what happened before my handler ran. This project closes that gap.
The raw TCP version
Accepting connections
The server opens a listener with net.Listen. Each new connection gets its own goroutine. That way one slow client does not block the others.
Reading the request
TCP does not hand you a full request. It hands you bytes, and they can show up in pieces. So I wrote my own read loop. It keeps calling conn.Read and adds each chunk to a buffer until the whole request is in.
The hard part was knowing when the request is done.
An HTTP request has two parts, the headers and the body. A blank line splits them. That line is \r\n\r\n in the raw bytes. Once I found that, the rest fell into place:
Read until the blank line shows up. Everything before it is the request line and headers.
Parse the headers and find Content-Length.
Count how many body bytes came in after the blank line.
Keep reading until the body matches Content-Length.
Without Content-Length the server has no way to know if the body is done or if more bytes are still on the way.
Parsing
I parse the request line myself to get the method, path and version. Then I split each header line on the first colon to get the name and value.
Routing and responses
The server routes each path to its own handler. It supports GET and POST.
Method	Path	What it does
GET	/address	Returns the address of the server
GET	/uptime	Returns how long the server has been running
POST	/submit	Prints "submit" to the server log. This could update a database, but that is outside the scope of this project.
In the net/http version the POST route is /server, and it only accepts application/json.
Every route needs an Authorization: Bearer <key> header. The server checks the token against API_KEY in .env. It uses a constant-time compare, so the check doesn't leak timing hints about the key.
I write every response by hand. That means the status line, the Content-Type and Content-Length headers, the blank line, and then the body. The server returns 200, 400, 401, 404, 405 and 500 where each one fits.
The net/http version
Then I rebuilt it with net/http. The code went from about 300 lines to about 100.
What net/http did for me
Opened the socket and accepted connections.
Ran each request in its own goroutine.
Read the full request, so I did not need my buffer loop.
Found the blank line and used Content-Length for me.
Parsed the request line and headers into r.Method, r.URL and r.Header.
Built the status line and set the headers on the response.
All the work that was hardest in the raw version was gone.
What net/http did not do for me
It did not do auth. It gives me the header with r.Header.Get("Authorization") and a helper with r.BasicAuth(). Deciding who is allowed in is still my job, so I kept my own Bearer token check.
That was the main lesson. net/http replaced the transport layer. It did not replace the app logic on top of it.
Bugs I hit and fixed
A shared struct got its fields overwritten across handler calls.
Off-by-one loop conditions in the read loop.
A missing return after an error response. The code kept going and wrote a second response.
Ignored errors from conn.Read.
Fragile logic for telling header bytes from body bytes.
Each one taught me something the framework would have hidden.
Run it
Create a .env file in the repo root:
API_KEY=your-key-here
Then start either server:
go run ./cmd/rawtcp
go run ./cmd/nethttp
Test it with curl:
curl -i -H "Authorization: Bearer your-key-here" http://localhost:8080/address
curl -i -H "Authorization: Bearer your-key-here" http://localhost:8080/uptime

# raw TCP version
curl -i -X POST -H "Authorization: Bearer your-key-here" -d "hello" http://localhost:8080/submit

# net/http version
curl -i -X POST -H "Authorization: Bearer your-key-here" -H "Content-Type: application/json" -d '{"msg":"hello"}' http://localhost:8080/server
Known limits
The raw version is for learning. It is not a production server.
No chunked transfer encoding. The body must come with a Content-Length.
No keep-alive. Each connection handles one request, then closes.
No HTTPS.
What I took from it
I can now explain what happens from the moment bytes hit the socket to the moment my handler runs. I am glad most languages handle this for us day to day. I am also glad I built it once by hand, because now I know what I am trusting.
