<script lang="ts">
	import { Icon } from "svelte-awesome";
	import { pause, play } from "svelte-awesome/icons";
	export let canAnimate: boolean;
	export let toggleCanAnimate: () => void;
</script>

<button type="button" class={`btn btn-sm sm:btn-md variant-filled`} on:click={toggleCanAnimate}>
	{#if canAnimate}
		<span><Icon slot="lead" data={pause} scale={1} /></span>
		<span>Pause</span>
	{:else}
		<span><Icon slot="lead" data={play} scale={1} /></span>
		<span>Resume</span>
	{/if}
</button>
