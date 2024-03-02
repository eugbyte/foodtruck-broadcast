import { mapboxToken } from "$lib/packages/config";
import * as core from "@mapbox/search-js-core"; // cjs vs esm conflict
import axios from "axios";
import debounce from "debounce-promise";
import type { FeatureCollection, Geometry } from "geojson";

export type LngLat = core.LngLat;
export type AddressAutofillSuggestion = core.AddressAutofillSuggestion;

/**
 * Debounced function to suggest addresses based on input text.
 * @param searchText text related to address to prompt address suggestions
 * @returns
 */
export const debounceSuggest = debounce(async (text: string) => {
	const suggestions: AddressAutofillSuggestion[] = await getSuggestions(text);
	return suggestions;
}, 1000);

export async function getSuggestions(searchText: string): Promise<AddressAutofillSuggestion[]> {
	console.log(`getting suggestions for ${searchText}`);
	const autofill = new core.AddressAutofillCore({ accessToken: mapboxToken, country: "sg" });
	const sessionToken = new core.SessionToken();
	const { defaults } = autofill;
	const result = await autofill.suggest(searchText, { sessionToken, ...defaults });
	console.log({ result });
	return result.suggestions;
}

export async function geoCode(address: string): Promise<LngLat> {
	const encode = encodeURIComponent(address);
	const url = `https://api.mapbox.com/geocoding/v5/mapbox.places/${encode}.json?access_token=${mapboxToken}`;
	const json: FeatureCollection = (await axios.get(url)).data;
	console.log(json);
	const features = json.features;
	for (const feature of features) {
		const g = feature.geometry as Geometry;
		if (g.type == "Point") {
			const [lng, lat] = g.coordinates;
			return new core.LngLat(lng, lat);
		}
	}

	throw new Error("could not geocode");
}

export function lngLat({ lng, lat }: { lat: number; lng: number }): LngLat {
	return new core.LngLat(lng, lat);
}
