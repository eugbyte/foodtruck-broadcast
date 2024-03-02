import type { GeoInfo, Message } from "$lib/packages/models";
import { z } from "zod";
import { superValidate } from "sveltekit-superforms/client";
import { schema } from "$lib/packages/form/vendor-id";

export interface WSWorker {
	open: () => void;
	close: () => void;
	send: (geoInfo: GeoInfo) => void;
}

export async function _createWSWorker(): Promise<WSWorker> {
	// https://vitejs.dev/guide/features.html#import-with-query-suffixes
	const WSWorker = await import("$lib/packages/web-workers/websocket/geo-watch?worker");
	const worker: Worker = new WSWorker.default();

	const open = () => worker.postMessage(<Message>{ action: "ws_init" });
	const close = () => worker.postMessage(<Message>{ action: "ws_close" });
	const send = (geoInfo: GeoInfo) =>
		worker.postMessage(<Message>{
			action: "ws_send",
			data: geoInfo
		});

	return {
		open: () => open(),
		close: () => close(),
		send: (geoInfo: GeoInfo) => send(geoInfo)
	};
}

export async function load() {
	const form = await superValidate(
		{
			vendorID: ""
		},
		schema
	);

	return { form };
}
