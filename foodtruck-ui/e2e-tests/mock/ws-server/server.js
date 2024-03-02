import { faker } from "@faker-js/faker";
import { WebSocketServer } from "ws";

const wss = new WebSocketServer({ port: 3080, path: "/customer" });

console.log("started mock ws server");

wss.on("connection", (ws) => {
	let count = 1;
	setInterval(() => {
		const geoInfo = {
			vendorID: count > 1 ? faker.number.int({ min: 1, max: 2 }).toString() : "1",
			lat: 1.3521 + faker.number.float() / 100,
			lng: 103.8198 + faker.number.float() / 100
		};

		if (count == 1) {
			geoInfo.vendorID = "2";
		}

		const msg = {
			action: "geo_info",
			data: geoInfo
		};
		console.log("sending", msg);
		ws.send(JSON.stringify(msg));

		count += 1;
	}, 1000);
});
