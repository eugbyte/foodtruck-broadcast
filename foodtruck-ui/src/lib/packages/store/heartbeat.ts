import { MAPBOX_TIMEOUT, type MapboxPayload, type Message } from "$lib/packages/models";
import { Subject } from "rxjs";

type VendorID = string;
type heartbeatCB = (val: VendorID) => void;
export type Heartbeat$ = Pick<Subject<VendorID>, "subscribe" | "next" | "unsubscribe">;

/**
 * Create a dedicated web worker, which maintains a series of timeouts for the vendors' GeoInfos.
 * Basically, a heartbeat mechanism to remove stale markers from the map.
 *
 * The timeouts maintains the lifetime of which the Marker, representing the vendor's info, should be displayed on the map.
 * When the timeout expires, the vendorID is emitted back to this function running in the main thread.
 * @warning do not instantiate the observable in a non browser environment as web worker is utilised.
 * Do not create a singleton instance outside of a component.
 * @returns A hot observable, emitting vendorID when the timeout is reached.
 */
export async function createHeartbeats(): Promise<Heartbeat$> {
	const subject = new Subject<VendorID>();
	// https://vitejs.dev/guide/features.html#import-with-query-suffixes
	const TimeoutWorker = await import("$lib/packages/web-workers/timeouts?worker");
	const worker: Worker = new TimeoutWorker.default();

	worker.addEventListener("message", (e) => {
		const msg: Message = e.data;
		if (msg.action != MAPBOX_TIMEOUT.EXPIRED) {
			return;
		}
		const { ID: vendorID } = msg.data;
		console.log("removing action: ", vendorID);
		subject.next(vendorID);
	});

	const handleNext = (ID: string) => {
		const msg: Message = {
			action: MAPBOX_TIMEOUT.BEGIN_TIMEOUT,
			data: <MapboxPayload>{ ID }
		};
		worker.postMessage(msg);
	};

	const handleComplete = () => {
		// complete vs unsubscribe (https://stackoverflow.com/a/38082208/6514532)
		subject.complete();
		worker.postMessage(MAPBOX_TIMEOUT.CLEAR_TIMEOUT);
	};

	return {
		subscribe: (cb: heartbeatCB) => subject.subscribe(cb),
		next: (vendorID: VendorID) => handleNext(vendorID),
		unsubscribe: () => handleComplete()
	};
}
