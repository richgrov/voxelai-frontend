<script lang="ts">
	import IconHammer from '~icons/fa6-solid/hammer';

	let textElement: HTMLInputElement;
	export let onSubmit: (query: string) => void;

	function submit() {
		if (!textElement.value) {
			return;
		}

		onSubmit(textElement.value);
	}
</script>

<form on:submit={submit} class="flex relative rounded-full border-black border-2
	shadow-[5px_5px_black] bg-white">
	<input type="text" placeholder="A house made of gold" bind:this={textElement}
		class="w-full p-5 rounded-full outline-none">
	<button class="flex items-center mr-3">
		<div class="bg-pink-300 rounded-full p-3">
			<IconHammer />
		</div>
	</button>
</form>
