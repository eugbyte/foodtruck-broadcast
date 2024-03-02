<script lang="ts">
	import { close } from "svelte-awesome/icons";
	import { Icon } from "svelte-awesome";
	import { NotificationButton } from "$lib/components/shared";

	/**
	 * Bind the granted value
	 */
	export let onAllow: () => void;
	export let onDismiss: () => void;
	export let text: string;
	export let className = "";
</script>

<div
	class={`variant-filled-surface flex flex-row rounded-md py-2 text-sm sm:text-base ${className}`}
>
	<aside class="flex flex-row flex-nowrap items-center">
		<!-- Message -->
		<div class="alert-message break-words ml-4">
			<p>{text}</p>
		</div>
		<!-- Actions -->
		<div class="alert-actions mx-4 space-x-4">
			<NotificationButton onClick={onAllow} />
			<button class="btn-icon btn-icon-sm variant-filled" on:click={onDismiss}>
				<Icon slot="lead" data={close} scale={1} />
			</button>
		</div>
	</aside>
</div>
