<script lang="ts">
	import type { FormHook } from "$lib/packages/form/vendor-id";

	export let className = "";
	export let formHook: FormHook;

	const { enhance, constraints, errors, form, tainted } = formHook;
</script>

<section class={className}>
	<form class="grid grid-cols-3 grid-rows-2 w-full sm:w-1/2" method="POST" use:enhance>
		<div class="self-center">
			<label for="VendorID" class="mr-4 font-semibold whitespace-nowrap">Vendor ID:</label>
		</div>
		<div class="col-span-2">
			<input
				class="input"
				bind:value={$form.vendorID}
				type="text"
				id="VendorID"
				placeholder="Required to broadcast"
				{...$constraints.vendorID}
			/>
		</div>
		<div class="" />
		<div class="col-span-2 self-center justify-self-center">
			{#if $tainted?.vendorID && $errors.vendorID}
				<span class="invalid w-full text-center text-red-500">{$errors.vendorID}</span>
			{/if}
		</div>
	</form>
</section>
