<script>
	import PermissionGuide from "./PermissionGuide.svelte";
	import { Accordion, AccordionItem } from "@skeletonlabs/skeleton";
	export let className = "";
</script>

<div class={className}>
	<Accordion>
		<AccordionItem>
			<p class="font-bold" slot="summary">How to modify permission</p>
			<svelte:fragment slot="content">
				<PermissionGuide />
			</svelte:fragment>
		</AccordionItem>
	</Accordion>
</div>
