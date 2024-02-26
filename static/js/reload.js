let socket = new WebSocket(`ws://${window.location.host}/wsreload`);

socket.onclose = function(event) {

    socket = null;

    setTimeout(function() {
        window.location.reload();
    }, 2000)
}

socket.onmessage = function(event) {
    console.debug(event.data);
}
