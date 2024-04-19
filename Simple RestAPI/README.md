## Simple RestAPI server
Very simple HTTP server that is receiving `GET` and `POST` requests in parallel and returns HTML responds.

### Arguments:
On start the application takes 2 optional flag arguments:
- `-p`(optional | default: 8000) - defines port number on which port on local address the server will be listening
- `-f`(optional | default: "./threat.html.tmpl")  - defines filepath to HTML template file to parse and returns in `POST` request
    
If other flag arguments are specified, the program ends with error.
### Requests:
- GET - address: `/`(root only), returns simple HTML webpage with textarea and submit button to submit JSON via `POST` request
- POST - address: `/render`, decode posted JSON and returns decoded data placed in parsed HTML template

### Startup
Clone repository, navigate in terminal to folder "Simple RestAPI" and type:   
go run . [-f "TEMPLATE_PATH"] [-p PORT_NUMBER]   
Server is closed upon receiving `SIGTERM` or `SIGINT` signals and shutdown all requests gracefully.
