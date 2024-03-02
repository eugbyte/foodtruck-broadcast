export interface Message {
	action: string;
	data: any;
}

export interface MapboxPayload {
	ID: string;
}

export const MAPBOX_BROADCAST_NAME = "mapbox";
export enum MAPBOX_TIMEOUT {
	BEGIN_TIMEOUT = "mapbox_begin_timeout",
	EXPIRED = "mapbox_expired",
	CLEAR_TIMEOUT = "clear_timeout"
}
