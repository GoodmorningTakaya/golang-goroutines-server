# golang-goroutines-server

This is a HTTP Server that:
- Communicates on port 12345
- In the background, retrieve the current Unix time based on your IP from worldtimeapi.org
  ever E seconds where E is Euler's Constant from the math library
- Return the last fetched time and start time as an int64 and the uint32 amount of requests
  that have been made to retrieve the time using a GET request to the root
- Each time a request is made, use a channel to send the request IP and last fetched time
  as well as the request time to another go routine which will log this data to a newline
  delimited file named logs in the following format: 
  
    <request-ip>-<current-time>-<request-time>



