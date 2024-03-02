import { Subject } from "rxjs";

type BroadcastCB = (g: GeolocationPosition) => void;
type ErrorCB = (e: GeolocationPositionError) => void;

export type GeoWatch$ = Pick<Subject<GeolocationPosition>, "subscribe" | "unsubscribe"> & {
	/**
	 * As a side effect, checks whether the geolocation permission has been granted.
	 * If successful, the coordinates is emitted. If permission is not granted, emits an error.
	 */
	start: () => void;
	/**
	 * stops tracking the user's geolocation. Does NOT close the observable forever, as distinct from `unsubscribe`.
	 */
	stop: () => void;
	error: (cb: ErrorCB) => void;
};

/**
 * Create an observable that emits the device's current geo position.
 *
 * Remeber to call `start()` to begin the broadcast.
 *
 * Note that calling start() after unsubscribe() is called will not work as the observable has been completed.
 * @returns A WatchGeo observable
 */
export function createGeoWatch(): GeoWatch$ {
	let watchID = 0;

	const subject = new Subject<GeolocationPosition>();
	const errors = new Subject<GeolocationPositionError>();

	const onSuccess = (p: GeolocationPosition) => subject.next(p);
	const onError = (e: GeolocationPositionError) => errors.next(e);

	const options: PositionOptions = {
		enableHighAccuracy: true,
		timeout: 5_000,
		maximumAge: 0
	};

	const close = () => {
		subject.complete();
		errors.complete();
		navigator.geolocation.clearWatch(watchID);
	};

	const start = () => {
		navigator.geolocation.clearWatch(watchID);
		console.log("broadcasting position...");
		watchID = navigator.geolocation.watchPosition(onSuccess, onError, options);
	};

	return {
		start: () => start(),
		stop: () => navigator.geolocation.clearWatch(watchID),
		subscribe: (cb: BroadcastCB) => subject.subscribe(cb),
		unsubscribe: () => close(),
		error: (cb: ErrorCB) => errors.subscribe(cb)
	};
}
