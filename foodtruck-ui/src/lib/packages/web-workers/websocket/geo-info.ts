import type { GeoInfo, Message } from "$lib/packages/models";
import { WS_URL } from "$lib/packages/config";

(() => {
	const w = self as unknown as DedicatedWorkerGlobalScope;
	let socket: WebSocket;

	// on ws_init, open the ws connection and relay any message received from the ws server to the main thread.
	w.addEventListener("message", (pe) => {
		const msg: Message = pe.data;
		if (msg.action != "ws_init") {
			return;
		}

		socket = new WebSocket(WS_URL);

		socket.addEventListener("message", (e) => {
			const msg: Message = JSON.parse(e.data);
			if (msg.action == "geo_info") {
				const geoInfo: GeoInfo = msg.data;
				postMessage(geoInfo);
			}
		});
	});

	// on ws_send, send bounding box of user's geolocation from main thread to ws server.
	w.addEventListener("message", (e) => {
		const msg: Message = e.data;
		if (msg.action != "ws_send" || socket == null) {
			return;
		}

		const boundingBox = msg.data;
		// TO DO
		const box: Record<string, number> = { ...boundingBox };
		const payload: Message = {
			action: "bounding_box",
			data: box
		};

		socket.send(JSON.stringify(payload));
	});

	// on ws_close, close the websocket connection
	w.addEventListener("message", (e) => {
		const msg: Message = e.data;
		if (msg.action == "ws_close") {
			socket?.close();
		}
	});
})();
