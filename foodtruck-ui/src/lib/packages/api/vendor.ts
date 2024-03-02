import type { Vendor } from "$lib/packages/models/Vendor";
import axios from "axios";

export async function getVendor(vendorID: number): Promise<Vendor> {
	const resp = await axios.get(`http://localhost:7000/vendor/${vendorID}`);
	return resp.data as Vendor;
}
