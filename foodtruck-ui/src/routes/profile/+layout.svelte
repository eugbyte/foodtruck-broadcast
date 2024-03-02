<script lang="ts">
	import { base as _base } from "$app/paths";
	import { goBack } from "$lib/packages/route";
	import { AppBar, AppShell } from "@skeletonlabs/skeleton";
	import { Icon } from "svelte-awesome";
	import { arrowLeft } from "svelte-awesome/icons";

	let base = _base.trim() == "" ? "." : _base;
	// the build folder directory, from `pnpm run build`, places the manifest.webmanifest file just outside the profile folder
	base = `../${base}`;
</script>

<!-- App Shell -->
<AppShell>
	<!-- (header) -->
	<AppBar slot="header">
		<button
			slot="lead"
			type="button"
			class="btn-icon btn-icon-sm variant-filled"
			on:click={() => goBack()}
		>
			<Icon data={arrowLeft} scale={1} />
		</button>
		<h1 class="h3">Profile</h1>
	</AppBar>
	<!-- Page Route Content -->
	<slot />
</AppShell>

<svelte:head>
	<title>Foodtruck - Profile</title>
	<link rel="manifest" href="{base}/manifest.webmanifest" crossorigin="use-credentials" />
</svelte:head>
