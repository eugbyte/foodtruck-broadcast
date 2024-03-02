<script lang="ts">
	import { getDrawerStore, initializeStores } from "@skeletonlabs/skeleton";
	import SearchBar from "./SearchBar.svelte";
	import { searchBarFocus$ } from "$lib/packages/store";
	import { computePosition, autoUpdate, offset, shift, flip, arrow } from "@floating-ui/dom";
	import { storePopup } from "@skeletonlabs/skeleton";

	initializeStores();
	storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow });
</script>

<SearchBar value={""} placeholder="Address" {searchBarFocus$} drawerStore$={getDrawerStore()} />
