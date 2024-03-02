import type { GeoInfo, Message } from "$lib/packages/models";

(() => {
	const w = self as unknown as DedicatedWorkerGlobalScope;
	let socket: WebSocket;

	// on ws_init, open the ws connection and relay any message received from the ws server to the main thread.
	w.addEventListener("message", (e) => {
		const msg: Message = e.data;
		if (msg.action != "ws_init") {
			return;
		}

		socket = new WebSocket("ws://localhost:5000/vendor");
		socket.addEventListener("open", () => console.log("connection established"));
		socket.addEventListener("error", (e) => console.log(`An error occured: ${e.type}`));
	});

	w.addEventListener("message", (e) => {
		const msg: Message = e.data;
		if (msg.action != "ws_send" || socket == null) {
			return;
		}

		const geoInfo = msg.data as GeoInfo;
		console.log({ geoInfo });
		socket.send(JSON.stringify(geoInfo));
	});

	// on ws_close, close the websocket connection
	w.addEventListener("message", (e) => {
		const msg: Message = e.data;
		if (msg.action == "ws_close") {
			socket?.close();
		}
	});
})();
