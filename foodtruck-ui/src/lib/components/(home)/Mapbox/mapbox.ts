/// <reference lib="DOM" />

import mapboxgl, { NavigationControl, GeolocateControl } from "mapbox-gl";
import { mapboxToken } from "$lib/packages/config";

/**
 * Create a MapBox with zoom control, rotation control, and geolocation control.
 * 
 * Use it like so:
 * ```
 	const mb = new DefaultMapBox("div_id", lng, lat);
	mb.draw();
	const { map, navControl, geolocater } = mb;	// access the mapboxgl.Map, NavigationControl, GeolocateControl here.
 * ```
 */
export class DefaultMapBox {
	private _map?: mapboxgl.Map;
	/**
	 * mapboxgl.NavigationControl
	 */
	navControl: NavigationControl;
	/**
	 * mapboxgl.GeolocateControl
	 */
	geolocater: GeolocateControl;

	/**
	 * Instantiate the mapbox object.
	 *
	 * Note that there is no side effects. The map will only be rendered to the HTML after calling render().
	 * @param containerID The HTML ID of the container element for the MapBox to attach itself to
	 * @param lng default lattitude
	 * @param lat default longitude
	 */
	constructor(public containerID: string, private lng: number, private lat: number) {
		mapboxgl.accessToken = mapboxToken;

		// 1. Create zoom and rotation controls.
		this.navControl = new NavigationControl();

		// 2. Create geolocation controls.
		this.geolocater = new GeolocateControl({
			positionOptions: {
				enableHighAccuracy: true
			},
			// Draw an arrow next to the location dot to indicate which direction the device is heading.
			showUserHeading: true
		});
	}

	/**
	 * As a side effect, render the Map to the HTML.
	 */
	draw() {
		const { containerID, lng, lat, navControl, geolocater } = this;
		const map = new mapboxgl.Map({
			container: containerID,
			style: "mapbox://styles/mapbox/streets-v12", // style URL
			zoom: 12,
			center: [lng, lat]
		});
		this._map = map;
		// 2. Attach zoom and rotation controls to the rendered map.
		map.addControl(navControl);
		// 3. Attach geolocation to the rendered map.
		map.addControl(geolocater);
	}

	/**
	 * Get the Map object, asserting that it is not null.
	 * The mapbox will be null if it not been called with draw()
	 */
	get map(): mapboxgl.Map {
		if (this._map == null) {
			throw new Error("map is null. Remember to call render() to initialize the map.");
		}
		return this._map;
	}
}

export function removeCopyrightText(): void {
	const copyright = document.querySelector(".mapboxgl-ctrl-attrib-inner");
	if (copyright != null) {
		copyright.remove();
	}
}

export type VendorID = string;
