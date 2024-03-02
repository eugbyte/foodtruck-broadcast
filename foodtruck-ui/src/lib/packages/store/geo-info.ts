import { Subject } from "rxjs";
import type { GeoInfo, Message } from "$lib/packages/models";

type geoInfoCB = (val: GeoInfo) => void;
// https://stackoverflow.com/a/63934188/6514532
export type GeoInfo$ = Pick<Subject<GeoInfo>, "subscribe" | "unsubscribe">;

/**
 * Create a shared web worker which opens a websocket connection,
 * and emits GeoInfo back to this function running in the main thread.
 * The main thread then transfers the GeoInfo to an observable.
 * @warning do not instantiate the observable in a non browser environment as web worker is utilised.
 * Do not create a singleton instance outside of a component.
 * @returns A hot observable, emitting geoInfo data from the websocket stream.
 */
export async function createGeoInfos(): Promise<GeoInfo$> {
	const subject = new Subject<GeoInfo>();

	// https://vitejs.dev/guide/features.html#import-with-query-suffixes
	const WSWorker = await import("$lib/packages/web-workers/websocket/geo-info?worker");
	const worker: Worker = new WSWorker.default();

	worker.addEventListener("message", (e) => {
		const geoInfo: GeoInfo = e.data;
		subject.next(geoInfo);
	});

	worker.postMessage(<Message>{ action: "ws_init" });

	const handleComplete = () => {
		// complete vs unsubscribe (https://stackoverflow.com/a/38082208/6514532)
		subject.complete();
		worker.postMessage(<Message>{ action: "ws_close" });
	};

	return {
		subscribe: (cb: geoInfoCB) => subject.subscribe(cb),
		unsubscribe: () => handleComplete()
	};
}
