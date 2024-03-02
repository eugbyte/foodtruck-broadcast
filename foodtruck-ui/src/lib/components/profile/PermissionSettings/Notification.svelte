<script lang="ts">
	import { NotificationButton } from "$lib/components/shared";

	export let notifyPermission: NotificationPermission;
	export let requestPushPermission: () => void;
	export let className = "";
</script>

<section class={className}>
	<div class="flex flex-row items-center">
		<span class="font-semibold mr-4">Notification:</span>
		{#if notifyPermission == "default"}
			<NotificationButton onClick={requestPushPermission} />
		{:else}
			<span class="capitalize">{notifyPermission}</span>
		{/if}
	</div>
</section>
