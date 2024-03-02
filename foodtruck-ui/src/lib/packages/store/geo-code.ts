import { geoCode, type LngLat } from "$lib/packages/geocode";
import { Subject } from "rxjs";
import axios from "axios";

type geocodeCB = (lnglat: LngLat) => void;
type NextFn = (addr: string) => void;
export type GeoCode$ = Pick<Subject<LngLat>, "subscribe" | "unsubscribe"> & { next: NextFn };

function createGeoCode(): GeoCode$ {
	const subject = new Subject<LngLat>();

	const handleNext = async (address: string) => {
		try {
			const lnglat: LngLat = await geoCode(address);
			subject.next(lnglat);
		} catch (err) {
			if (axios.isAxiosError(err)) {
				console.error(err.response?.data);
			} else {
				console.error(err);
			}
		}
	};

	return {
		subscribe: (cb: geocodeCB) => subject.subscribe(cb),
		next: (address: string) => handleNext(address),
		unsubscribe: () => subject.complete()
	};
}

/**
 * GeoCoding store that takes in a string address an emits the LatLng
 */
export const geocode$ = createGeoCode();
