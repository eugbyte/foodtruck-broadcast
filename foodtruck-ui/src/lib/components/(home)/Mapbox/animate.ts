import type { GeoInfo } from "$lib/packages/models";
import * as turf from "@turf/turf";
import type { Marker } from "mapbox-gl";

/**
 * Animate the marker moving from its original coordinates to the destination.
 * Mapbox already uses [Web Workers under the hood](https://docs.mapbox.com/mapbox-gl-js/guides/install/#loading-and-transpiling-the-web-worker-separately)
 * to draw the canvas, so there is no need to implement a web worker.
 * @param marker The MapBox marker object
 * @param geoInfo Geo info of the vendor
 * @param frameID The animation frame ID, single element array to store reference to the current animation frame
 * @param animationDuration How long the animation should run in millisecond
 */
export function animateMarker(
	marker: Marker,
	geoInfo: GeoInfo,
	frameID: [number],
	animationDuration: number
): void {
	const { lng: finalLng, lat: finalLat } = geoInfo;

	// code for animation here
	const { lng, lat } = marker.getLngLat();

	const from = turf.point([lng, lat]);
	const to = turf.point([finalLng, finalLat]);

	const bearing = turf.rhumbBearing(from, to);
	const distance = turf.distance(from, to, { units: "meters" });
	const speed = distance / (animationDuration / 1000); // m/s

	if (speed > 0) {
		animate({
			marker,
			speed,
			animationDuration,
			bearing,
			originlng: lng,
			originlat: lat,
			startTime: window.performance.now(),
			timestamp: window.performance.now(),
			frameID
		});
	}
}

interface AnimateProps {
	/**
	 * The mapbox marker
	 */
	marker: Marker;
	/**
	 * m/s
	 */
	speed: number;
	/**
	 * Number of ms for animation to run
	 */
	animationDuration: number;
	/**
	 * Deviance from true north
	 */
	bearing: number;
	/**
	 * lng of the origin
	 */
	originlng: number;
	/**
	 * lat of the origin
	 */
	originlat: number;
	/**
	 * The original [starting time](https://stackoverflow.com/a/21316178) when the animation begins.
	 * Determined with `window.performance.now()`
	 */
	startTime: number;
	/**
	 * The [current time](https://stackoverflow.com/a/21316178) for when requestAnimationFrame starts to fire callbacks.
	 * Determined via the timestamp event via `requestAnimationFrame()`
	 */
	timestamp: number;
	/**
	 * A single element array to store the [ID of the animation frame](https://developer.mozilla.org/en-US/docs/Web/API/window/requestAnimationFrame#return_value).
	 * Can pass this value to `window.cancelAnimationFrame(ID)` to cancel the animation
	 */
	frameID: [number];
}

function animate({
	marker,
	speed,
	bearing,
	originlat,
	originlng,
	startTime,
	animationDuration,
	timestamp,
	frameID
}: AnimateProps): void {
	const runtime = timestamp - startTime;

	if (runtime > animationDuration) {
		return;
	}

	const from = turf.point([originlng, originlat]);
	const newPoint = turf.rhumbDestination(from, speed * (runtime / 1000), bearing, {
		units: "meters"
	});
	const [newLng, newLat] = newPoint.geometry.coordinates;

	try {
		marker.setLngLat([newLng, newLat]);
	} catch (err) {
		console.error("race condition where marker is removed");
		return;
	}

	const id = requestAnimationFrame((timestamp) =>
		animate({
			marker,
			speed,
			bearing,
			originlng,
			originlat,
			startTime,
			timestamp,
			animationDuration,
			frameID
		})
	);
	frameID[0] = id;
}
