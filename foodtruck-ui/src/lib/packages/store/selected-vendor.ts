import type { GeoInfo } from "$lib/packages/models";
import { Subject } from "rxjs";

export type VendorID = Pick<GeoInfo, "vendorID">;
type selectVendorCB = (val: VendorID) => void;

export type SelectedVendor$ = Pick<Subject<VendorID>, "subscribe" | "next" | "unsubscribe">;

function createSelectedVendor(): SelectedVendor$ {
	const subject = new Subject<VendorID>();

	return {
		subscribe: (cb: selectVendorCB) => subject.subscribe(cb),
		next: (vendorID: VendorID) => subject.next(vendorID),
		unsubscribe: () => subject.complete()
	};
}

export const selectedVendor$ = createSelectedVendor();
