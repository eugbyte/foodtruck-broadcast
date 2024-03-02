<script lang="ts">
	import { onDestroy, onMount } from "svelte";
	import {
		getPermissionState as getNotifPermissionState,
		requestPermission as _requestPushPermission
	} from "$lib/packages/notification";
	import {
		getPermissionState as getGeoPermissionState,
		geolocate as _requestGeoPermission
	} from "$lib/packages/geolocate";
	import { createGeoWatch, type GeoWatch$ } from "$lib/packages/store";
	import {
		FAQSection,
		GeolocationSection,
		NotificationSection,
		VendorInput
	} from "$lib/components/profile";
	import { _createWSWorker as createWSWorker, type WSWorker } from "./+page";
	import { BroadcastToggle } from "$lib/components/profile";
	import type { GeoInfo } from "$lib/packages/models";
	import type { PageData } from "./$types";
	import { superForm } from "sveltekit-superforms/client";

	let geoPermission: PermissionState = "prompt";
	let notifyPermission: NotificationPermission = "default";

	let geoWatch$: GeoWatch$;
	let wsWorker: WSWorker;

	let canBroadcast = false;
	export let data: PageData;

	onMount(async () => {
		geoWatch$ = createGeoWatch();
		notifyPermission = getNotifPermissionState();
		geoPermission = await getGeoPermissionState();
		wsWorker = await createWSWorker();
	});

	onDestroy(() => {
		wsWorker?.close();
		geoWatch$?.unsubscribe();
	});

	const formHook = superForm(data.form, {
		SPA: true,
		validators: {
			vendorID: (id) => (id.trim() == "" ? "Vendor ID is required to broadcast." : null)
		}
	});

	formHook.form.subscribe(({ vendorID }) => console.log({ vendorID }));

	const requestPushPermission = async () => {
		try {
			notifyPermission = await _requestPushPermission();
		} catch (error) {
			console.log(error);
			notifyPermission = "denied";
		}
	};

	const requestGeoPermission = async () => {
		// geolocation Web API `getCurrentPosition` will request for permission and get position in the same call.
		// not possible to request permission and get geolocation separately.
		try {
			await _requestGeoPermission(1000);
			geoPermission = await getGeoPermissionState();
		} catch (e) {
			console.log(e);
			const error = e as GeolocationPositionError;
			if (error.code == error.TIMEOUT) {
				geoPermission = "granted";
				return;
			}
			geoPermission = "denied";
			canBroadcast = false;
		}
		console.log({ geoPermission });
	};

	$: {
		if (canBroadcast && geoPermission != "granted") {
			requestGeoPermission();
		}
	}

	const { errors } = formHook;
	$: {
		console.log("errors:", $errors.vendorID);
	}

	$: {
		if (!canBroadcast) {
			wsWorker?.close();
			geoWatch$?.stop();
		}

		if (canBroadcast && geoPermission == "granted" && geoWatch$ != null) {
			wsWorker.open();
			geoWatch$.start();
			geoWatch$.subscribe((geoPosition: GeolocationPosition) => {
				const { latitude: lat, longitude: lng, speed, heading } = geoPosition.coords;
				const geoInfo: GeoInfo = {
					vendorID: 1,
					lat,
					lng,
					speed,
					heading
				};
				console.log("device geoInfo:", geoInfo);
				wsWorker.send(geoInfo);
			});
			geoWatch$.error((err) => console.log(err));
		}
	}

	$: isDisabled =
		geoPermission == "denied" || ($errors.vendorID != null && $errors.vendorID?.length > 0);

	const sectionStyle = "w-full sm:w-1/2 my-2 p-2";
</script>

<div class="container h-full mx-auto flex flex-col items-center overflow-y-auto">
	<h2 class="h2 text-left my-1">Broadcast</h2>

	<VendorInput {formHook} className={sectionStyle} />
	<BroadcastToggle bind:canBroadcast {isDisabled} className={sectionStyle} />

	<hr class="!border-t- w-full my-5" />

	<h2 class="h2 text-left">Settings</h2>
	<NotificationSection bind:notifyPermission {requestPushPermission} className={sectionStyle} />
	<GeolocationSection bind:geoPermission {requestGeoPermission} className={sectionStyle} />

	<hr class="!border-t- w-full my-5" />

	<h2 class="h2 text-left">FAQ</h2>
	<FAQSection className={sectionStyle} />
</div>
