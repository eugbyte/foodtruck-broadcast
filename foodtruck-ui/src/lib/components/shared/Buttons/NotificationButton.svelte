<script lang="ts">
	import { bell } from "svelte-awesome/icons";
	import { Icon } from "svelte-awesome";

	export let onClick: () => void;
</script>

<button class="btn btn-sm variant-filled" on:click={onClick}>
	<Icon slot="lead" data={bell} scale={1} />
	<span>Allow</span>
</button>
