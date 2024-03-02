<script lang="ts">
	import type mapboxgl from "mapbox-gl";
	import { Marker, type GeolocateControl } from "mapbox-gl";
	import { onMount } from "svelte";
	import { DefaultMapBox, removeCopyrightText } from "./mapbox";
	import { onDestroy } from "svelte";
	import { animateMarker } from "./animate";
	import type { Subscription } from "rxjs";
	import { differenceInMilliseconds } from "date-fns";
	import type { GeoInfo$, Heartbeat$, SelectedVendor$, GeoCode$ } from "$lib/packages/store";
	import { getPermissionState } from "$lib/packages/geolocate";

	export let className = "";
	/**
	 * Observable stream on geo information of vendors
	 */
	export let geoInfos$: GeoInfo$;
	let geoSub: Subscription; // to unsubscribe when component unmounts
	/**
	 * Observable stream on the expiry of the geo infos.
	 */
	export let heartbeats$: Heartbeat$;
	let timeSub: Subscription; // to unsubscribe when component unmounts
	export let selectedVendor$: SelectedVendor$;
	export let geocode$: GeoCode$;
	let geocodeSub: Subscription;

	let map: mapboxgl.Map;
	let geolocater: GeolocateControl;
	const markers: Record<string, Marker> = {};

	// animations
	const frameIDs: Record<string, [number]> = {};
	const timestamps: Record<string, Date> = {};

	export let canAnimate: boolean;

	$: {
		if (!canAnimate) {
			for (const [frameID] of Object.values(frameIDs)) {
				cancelAnimationFrame(frameID);
			}
		}
	}

	onMount(async () => {
		const mb = new DefaultMapBox("map_foodtruck", 103.8198, 1.3521);
		mb.draw();
		({ map, geolocater } = mb);

		map.on("load", async () => {
			map.resize();
			removeCopyrightText();
			const perm = await getPermissionState();
			if (perm == "granted") {
				geolocater.trigger();
			}
		});
	});

	geoSub = geoInfos$.subscribe((geoInfo) => {
		if (map == null) {
			return;
		}

		const { vendorID, lng: lng, lat: lat } = geoInfo;
		heartbeats$.next(vendorID.toString());
		const prevTime = timestamps[vendorID] ?? new Date();
		timestamps[vendorID] = new Date();

		try {
			if (!(vendorID in markers)) {
				// not possible to instantiate empty marker
				const marker: Marker = new Marker().setLngLat([lng, lat]).addTo(map);
				marker.getElement().addEventListener("click", () => selectedVendor$.next({ vendorID }));
				markers[vendorID] = marker;
				frameIDs[vendorID] = [-1];
			}
		} catch (e) {
			const error = e as Error;
			console.error(`error adding marker to map: ${error.message}`);
			return;
		}

		// code for animation here
		const frameID: [number] = frameIDs[vendorID];
		cancelAnimationFrame(frameID[0]); // clear any existing animation
		if (canAnimate) {
			const animationDuration = differenceInMilliseconds(new Date(), prevTime);
			animateMarker(markers[vendorID], geoInfo, frameID, animationDuration);
		}
	});

	timeSub = heartbeats$.subscribe((vendorID) => {
		if (map == null) {
			return;
		}
		markers[vendorID]?.remove();
		delete markers[vendorID];
		delete frameIDs[vendorID];
		delete timestamps[vendorID];
	});

	geocodeSub = geocode$.subscribe((lnglat) => {
		const { lng, lat } = lnglat;
		map.flyTo({
			center: [lng, lat],
			essential: true // this animation is considered essential with respect to prefers-reduced-motion
		});
	});

	onDestroy(() => {
		map?.remove();
		// unsubscribe from the Subscription rather than the Subject directly (https://stackoverflow.com/a/38082208/6514532)
		geoSub?.unsubscribe();
		timeSub?.unsubscribe();
		geocodeSub?.unsubscribe();
	});
</script>

<div class={className}>
	<main id="map_foodtruck" />
	<p class="text-xs">© Mapbox © OpenStreeMap</p>
</div>

<style>
	@import "mapbox-gl/dist/mapbox-gl.css"; /*  <-- important */
	#map_foodtruck {
		height: 450px;
	}
</style>
