<script lang="ts">
	import PlayButton from "./PlayButton.svelte";

	let canAnimate = true;
	const toggleCanAnimate = () => {
		canAnimate = !canAnimate;
	};
</script>

<PlayButton {canAnimate} {toggleCanAnimate} />
