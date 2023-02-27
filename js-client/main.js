const net = require("net");

const PORT = 9876;

function connectToServer(){
    const ws = new net.Socket(URL);
    while(true)
    try{
        return ws.connect(PORT);
    }catch{
        continue;
    }
}

(function(){
    const socket = connectToServer();
    while(!socket.write("login:!jsclient!"));
    socket.on("data", (message) => {
        const parsed = JSON.parse(message.toString());
        console.log("received message: "+ message.toString() +  " " + parsed["id"] + " "+ parsed["message"] + " "+ parsed["status"]);

        socket.write(JSON.stringify({id: parsed["id"], clientName: "jsclient", status: 8}));

    });
})();
