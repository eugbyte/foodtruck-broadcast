/// <reference types="@sveltejs/kit" />
/// <reference lib="webworker" />
/// <reference no-default-lib="true"/>

import { build as jsFiles, files as staticFiles, version } from "$service-worker";
import { onInstall, onFetch, onActivate, onPush } from "./lib/packages/service-workers";
import { nanoid } from "nanoid";

// Create a unique cache name for this deployment
const CACHE_VERSION = `cache-${version}-${nanoid()}`;

const ASSETS = [
	...jsFiles, // the app itself
	...staticFiles // everything in `static`
];

const sw = self as unknown as ServiceWorkerGlobalScope;

sw.addEventListener("install", (event: ExtendableEvent) => onInstall(ASSETS, CACHE_VERSION, event));
sw.addEventListener("activate", (event: ExtendableEvent) => onActivate(CACHE_VERSION, event));
sw.addEventListener("fetch", (event: FetchEvent) => onFetch(CACHE_VERSION, event));
sw.addEventListener("push", (event: PushEvent) => onPush(event, () => undefined));
