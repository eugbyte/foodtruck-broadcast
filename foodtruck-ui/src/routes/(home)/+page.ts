import { sleep } from "$lib/packages/sleep";
import {
	getPermissionState as getNotifPermissionState,
	getSubscription,
	requestPermission as requestPushPermission,
	subscribe as subscribePushNotif
} from "$lib/packages/notification";
import { geolocate, getPermissionState as getGeoPermissionState } from "$lib/packages/geolocate";
import { postSubscription } from "$lib/packages/api";

/**
 * Checks permission, subscribes to the push notification, and then saves the subscription along with the user's geolocation to the DB.
 * @returns
 */
export async function _subscribeGeoPush() {
	// 1. subscribe to push notification
	let notifPermission: NotificationPermission = getNotifPermissionState();

	if (notifPermission == "default") {
		notifPermission = await requestPushPermission();
	}
	if (notifPermission != "granted") {
		return;
	}

	let json: PushSubscriptionJSON | null = await getSubscription();
	if (json == null) {
		json = await subscribePushNotif();
	}
	console.log(json);

	// merely for UX experience, no race conditions involved.
	await sleep(1000);

	// 2. retrieve geolocation
	let [lng, lat] = [103.85, 1.29];
	const geoPermission: PermissionState = await getGeoPermissionState();
	try {
		// geolocation Web API `getCurrentPosition` will request for permission and get position in the same call.
		// not possible to request permission and get geolocation separately.
		if (geoPermission == "prompt" || geoPermission == "granted") {
			const { coords } = await geolocate(10_000);
			console.log("geolocation retrieved", coords);
			({ longitude: lng, latitude: lat } = coords);
		}
	} catch (error) {
		console.error(error);
	}

	// 3. save subscription and user's geolocation to backend DB.
	const response = await postSubscription(json, lat, lng);
	console.log(response);
}

/**
 * Updates the user's geolocation. If the user has not granted permission, the function returns
 * @returns
 */
export async function _updateGeolocation() {
	const json: PushSubscriptionJSON | null = await getSubscription();
	if (json == null) {
		return;
	}

	const geoPermission: PermissionState = await getGeoPermissionState();
	if (geoPermission == "denied") {
		return;
	}

	const { coords } = await geolocate(10_000);
	const { longitude: lng, latitude: lat } = coords;
	const response = await postSubscription(json, lat, lng);
	console.log(response);
}
