<script lang="ts">
	import { GeolocateButton } from "$lib/components/shared";

	export let className = "";
	export let geoPermission: PermissionState;
	export let requestGeoPermission: () => void;
</script>

<section class={className}>
	<div class="flex flex-row items-center">
		<span class="font-semibold mr-4">Geolocation:</span>
		{#if geoPermission == "prompt"}
			<GeolocateButton onClick={requestGeoPermission} />
		{:else}
			<span class="capitalize">{geoPermission}</span>
		{/if}
	</div>
</section>
