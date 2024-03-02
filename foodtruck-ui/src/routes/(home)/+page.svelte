<script lang="ts">
	import { VendorAccordion, Mapbox, PermissionPrompt, PlayButton } from "$lib/components/(home)";
	import { zIndexes } from "$lib/packages/config";
	import {
		type GeoInfo$,
		type Heartbeat$,
		createGeoInfos,
		createHeartbeats,
		type SelectedVendor$,
		selectedVendor$ as _selectedVendor$,
		geocode$ as _geoCode$,
		type GeoCode$
	} from "$lib/packages/store";
	import { onDestroy, onMount } from "svelte";
	import type { Vendor } from "$lib/packages/models/Vendor";
	import { getVendor } from "$lib/packages/api";
	import {
		_updateGeolocation as updateGeolocation,
		_subscribeGeoPush as subscribeGeoPush
	} from "./+page";
	import { fade } from "svelte/transition";
	import { getPermissionState as getNotifPermission } from "$lib/packages/notification";

	let geoInfos$: GeoInfo$;
	let heartbeats$: Heartbeat$;
	const selectedVendor$: SelectedVendor$ = _selectedVendor$;
	const geocode$: GeoCode$ = _geoCode$;
	let vendor: Vendor | null;
	let intervalId: NodeJS.Timeout;

	$: canAnimate = true;
	const toggleCanAnimate = () => {
		canAnimate = !canAnimate;
	};

	let showPrompt = true;
	const onAllow = async () => {
		let retries = 2;
		while (retries > 0) {
			showPrompt = false;
			retries -= 1;
			try {
				await subscribeGeoPush();
				return;
			} catch (error) {
				console.error(error);
			}
		}
		// if the code reaches here, that means subscription failed
		showPrompt = true;
	};
	const onDismiss = () => {
		showPrompt = false;
	};
	$: console.log(showPrompt);

	onMount(async () => {
		// cannot create web worker during compile time as environment is node.js, not browser.
		geoInfos$ = await createGeoInfos();
		heartbeats$ = await createHeartbeats();
		showPrompt = getNotifPermission() == "default";
		intervalId = setInterval(async () => {
			/**
			 * User has either granted or deny permission for push notification.
			 */
			const hasTouched = !showPrompt;
			if (hasTouched) {
				await updateGeolocation();
			}
		}, 600_000); // 10 min
	});

	selectedVendor$.subscribe(async ({ vendorID }) => {
		try {
			vendor = await getVendor(vendorID);
		} catch (err) {
			console.error(err);
		}
	});

	onDestroy(() => {
		// the unsubscribe does not unsubscribe from the subject, only cleans up
		// https://stackoverflow.com/a/38082208/6514532
		geoInfos$?.unsubscribe();
		heartbeats$?.unsubscribe();
		selectedVendor$?.unsubscribe();
		clearInterval(intervalId);
	});

	$: webworkerReady = geoInfos$ != null && heartbeats$ != null;

	let text = "Click on a map marker to find out more";
	$: {
		if (vendor != null) {
			text = `${vendor.name}: ${vendor.description}`;
		}
	}
</script>

<div class="container h-full mx-auto flex flex-col justify-between items-center overflow-y-auto">
	<!-- Map display -->
	<div class={`w-full lg:w-3/4 px-1 ${zIndexes.map}`}>
		{#if showPrompt}
			<div
				class="my-3 items-start w-full sm:w-fit text-center mx-auto"
				in:fade={{ duration: 2000 }}
			>
				<PermissionPrompt {onAllow} {onDismiss} text="Recieve notifications of vendors near you" />
			</div>
		{/if}

		{#if webworkerReady && selectedVendor$ != null}
			<div class="my-4" />
			<Mapbox className="" {geoInfos$} {heartbeats$} {selectedVendor$} {canAnimate} {geocode$} />
			<div class="flex flex-col items-center space-y-3 mt-2 sm:mt-5">
				<PlayButton {canAnimate} {toggleCanAnimate} />
			</div>
		{/if}
	</div>

	<VendorAccordion {text} header="Vendor Info" />
</div>
