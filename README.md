# Chat application
Very basic chat application that can handle multiple sessions in a browser. It uses websockets. Only one chat "room" is available where all of the connections are linked to.

### In a browser's console
let socket = new WebSocket("ws://localhost:3000/ws")
socket.onmessage = (event) => { console.log("received from the server: ", event.data) }
socket.send("hello from client")