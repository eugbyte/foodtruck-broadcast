import { goto } from "$app/navigation";

export function goBack(defaultRoute = "/") {
	const ref = document.referrer;
	goto(ref.length > 0 ? ref : defaultRoute);
}

export function routeToPage(route: string, replaceState: boolean) {
	goto(`${route}`, { replaceState });
}
