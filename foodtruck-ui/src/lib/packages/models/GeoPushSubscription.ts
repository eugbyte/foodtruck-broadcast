export interface GeoPushSubscription {
	endpoint: string;
	geohash: string;
	lastSend: number;
	expiration: number;
	auth: string;
	p256dh: string;
	optIn: boolean;
}
