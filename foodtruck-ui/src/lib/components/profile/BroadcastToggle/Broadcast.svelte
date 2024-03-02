<script lang="ts">
	import { SlideToggle } from "@skeletonlabs/skeleton";

	export let canBroadcast: boolean;
	export let isDisabled: boolean;
	export let className = "";

	const broadcastStates: Record<string, string> = {
		true: "On",
		false: "Off"
	};
</script>

<section class={className}>
	<div class="flex flex-row items-center">
		<span class="font-semibold mr-4">Broadcast your position:</span>
		<SlideToggle
			name="broadcast"
			bind:checked={canBroadcast}
			disabled={isDisabled}
			size="sm"
			class="my-2"
			>{isDisabled ? "Blocked" : broadcastStates[canBroadcast.toString()]}
		</SlideToggle>
	</div>
</section>
