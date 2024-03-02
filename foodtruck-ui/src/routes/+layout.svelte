<script lang="ts">
	import "../app.postcss";
	import { initializeStores } from "@skeletonlabs/skeleton";
	import { computePosition, autoUpdate, offset, shift, flip, arrow } from "@floating-ui/dom";
	import { storePopup } from "@skeletonlabs/skeleton";

	initializeStores();
	// needed for the autocomplete to dismiss itself when clicking outside of the select form
	storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow });
</script>

<slot />
