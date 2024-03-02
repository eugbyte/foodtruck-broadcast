import { MAPBOX_TIMEOUT, type Message, type MapboxPayload } from "$lib/packages/models";
import { differenceInSeconds } from "date-fns";

const timeouts = new Map<string, Date>();
let lock = false;
let intervalID: any;

/**
 * Begin broadcast channel with the main thread to handle multiple setTimeouts (macro task) without affecting performance on main threads
 */

(() => {
	const w = self as unknown as DedicatedWorkerGlobalScope;
	w.addEventListener("message", (e) => {
		const msg: Message = e.data;
		if (msg.action != MAPBOX_TIMEOUT.BEGIN_TIMEOUT) {
			return;
		}
		lock = true;
		const payload: MapboxPayload = msg.data;
		timeouts.set(payload.ID, new Date());
		lock = false;
	});

	w.addEventListener("message", (e) => {
		const msg: Message = e.data;
		if (msg.action == MAPBOX_TIMEOUT.CLEAR_TIMEOUT) {
			clearInterval(intervalID);
		}
	});

	intervalID = setInterval(() => {
		for (const [ID, date] of timeouts) {
			inner: while (lock) {
				continue inner;
			}
			const duration = differenceInSeconds(new Date(), date);
			if (duration >= 30) {
				timeouts.delete(ID);
				postMessage(<Message>{ action: MAPBOX_TIMEOUT.EXPIRED, data: { ID } });
				console.log(`emitted timeout for ${ID}`);
			}
		}
	}, 1000);
})();
