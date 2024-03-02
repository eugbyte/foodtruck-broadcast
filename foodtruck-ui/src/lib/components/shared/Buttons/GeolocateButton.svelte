<script lang="ts">
	import { Icon } from "svelte-awesome";
	import { mapMarker } from "svelte-awesome/icons";

	export let onClick: () => void;
</script>

<button class="btn btn-sm variant-filled" on:click={onClick}>
	<Icon slot="lead" data={mapMarker} scale={1} />
	<span>Allow</span>
</button>
