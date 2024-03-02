<script lang="ts">
	import { Accordion, AccordionItem } from "@skeletonlabs/skeleton";
	import Icon from "svelte-awesome";
	import { info } from "svelte-awesome/icons";

	export let text = "";
	export let header: string;
	const hasText = text.trim().length > 0;
</script>

<Accordion>
	<AccordionItem open={hasText}>
		<Icon slot="lead" data={info} scale={1} />
		<p class="font-bold text-md sm:text-lg" slot="summary">{header}</p>
		<p slot="content" class="whitespace-pre-line text-sm sm:text-base">{text}</p>
	</AccordionItem>
	<br />
</Accordion>
