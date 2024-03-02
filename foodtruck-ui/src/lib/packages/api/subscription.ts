import type { GeoPushSubscription } from "$lib/packages/models/GeoPushSubscription";
import geohasher from "ngeohash";
import axios from "axios";
import { subscriberURL } from "$lib/packages/config";

/**
 * save the user details and subscription to database.
 * @param json The PushSubscriptionJSON object from calling the `navigato` `PushManager` Web API.
 * @param lat latitude of the user.
 * @param lng longitude of the user.
 * @returns The response from saving the subscription to the backend API.
 */
export async function postSubscription(
	json: PushSubscriptionJSON,
	lat: number,
	lng: number
): Promise<Record<string, any>> {
	const keys: Record<string, string> = json.keys ?? {};

	const geoSub: GeoPushSubscription = {
		endpoint: json.endpoint ?? "",
		expiration: json.expirationTime ?? -1,
		p256dh: keys["p256dh"],
		auth: keys["auth"],
		lastSend: 0,
		geohash: geohasher.encode(lat, lng, 6),
		optIn: true
	};
	console.log(geoSub);
	const result = await axios.post(subscriberURL, geoSub);
	return result.data as Record<string, any>;
}
