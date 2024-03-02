import { VAPID_PUBLIC_KEY } from "$lib/packages/config";
/**
 * Create the subscription with a push service (chrome, firefox),
 * which will generate an endpoint associated with the browser's ip address,
 * and return the endpoint to us in an object
 * @returns The PushSubscription object
 */
export async function subscribe(): Promise<PushSubscriptionJSON> {
	if (!("serviceWorker" in navigator)) {
		throw new Error("navigator does not have service worker");
	}
	// no service worker in dev mode, so the ready promise will never resolve
	const registration: ServiceWorkerRegistration = await navigator.serviceWorker.ready;
	const subscription: PushSubscription = await registration.pushManager.subscribe({
		userVisibleOnly: true,
		applicationServerKey: urlBase64ToUint8Array(VAPID_PUBLIC_KEY)
	});
	const json: PushSubscriptionJSON = subscription.toJSON();
	json.expirationTime = 60;
	return json;
}

export async function unsubscribe(): Promise<void> {
	if ("navigator" in window) {
		return;
	}
	const sw: ServiceWorkerRegistration = await navigator.serviceWorker.ready;
	const subscription: PushSubscription | null = await sw.pushManager.getSubscription();
	if (subscription == null) {
		return;
	}
	const isUnsub = await subscription.unsubscribe();
	if (!isUnsub) {
		throw new Error("failed to unsubscribe");
	}
}

/**
 * Check if a user is already subscribed
 */
export async function isSubscribed(): Promise<boolean> {
	if ("navigator" in window) {
		return false;
	}
	const sw: ServiceWorkerRegistration = await navigator.serviceWorker.ready;
	const res: PushSubscription | null = await sw.pushManager.getSubscription();
	return res != null;
}

/**
 * @returns Any existing push subscription, or null if none is found
 */
export async function getSubscription(): Promise<PushSubscription | null> {
	const sw: ServiceWorkerRegistration = await navigator.serviceWorker.ready;
	return sw.pushManager.getSubscription();
}

// Copied from the https://gist.github.com/Klerith/80abd742d726dd587f4bd5d6a0ab26b6
function urlBase64ToUint8Array(base64String: string): Uint8Array {
	const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
	// eslint-disable-next-line no-useless-escape
	const base64 = (base64String + padding).replace(/-/g, "+").replace(/_/g, "/");

	const rawData = window.atob(base64);
	const outputArray = new Uint8Array(rawData.length);

	for (let i = 0; i < rawData.length; ++i) {
		outputArray[i] = rawData.charCodeAt(i);
	}
	return outputArray;
}
