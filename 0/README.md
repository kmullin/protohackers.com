0: Smoke Test

View leaderboard

Deep inside Initrode Global's enterprise management framework lies a component that writes data to a server and expects to read the same data back. (Think of it as a kind of distributed system delay-line memory). We need you to write the server to echo the data back.

Accept TCP connections.

Whenever you receive data from a client, send it back unmodified.

Make sure you don't mangle binary data, and that you can handle at least 5 simultaneous clients.

Once the client has finished sending data to you it shuts down its sending side. Once you've reached end-of-file on your receiving side, and sent back all the data you've received, close the socket so that the client knows you've finished.

Your program will implement the TCP Echo Service from RFC 862.