/**
 * Z indexes for the various components so that they overlay each other in the correct order.
 */
export const zIndexes = {
	appbar: "z-40",
	drawer: "z-20",
	map: "z-10"
};

export const mapboxToken =
	"pk.eyJ1IjoiZXVnYnl0ZSIsImEiOiJjbGs2NHVqbXYwd3F1M2Nsbmxka2pqcmtqIn0.bKGwy3dcbmMWufUbe24U4Q";

function wsURL(): string {
	// return "ws://localhost:3080/customer";
	const { MODE } = import.meta.env;

	switch (MODE) {
		case "development":
			return "ws://localhost:8080/customer";
		case "staging":
			return "ws://localhost:3080/customer";
	}

	return "ws://localhost:8080/customer";
}

export const WS_URL = wsURL();
export const subscriberURL =
	"http://115a5e6e.execute-api.localhost.localstack.cloud:4566/api/v1/subscription";
export const VAPID_PUBLIC_KEY =
	"BIvYk7KnTMwTVY9ubq55Gkt0FzMVo2Rsm7Cs43MEhFJjXBwu073JvS6oEH56CMdK7FFcjW6mqZgF_zdsnDgCczA"; // TO DO - change to process.env.VAPID_PUBLIC_KEY depending on stg
