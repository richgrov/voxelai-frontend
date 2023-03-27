<script lang="ts">
	import '../app.css';
	import Search from '../lib/search.svelte';
</script>

<div class="flex justify-center items-center flex-col h-4/5">
	<h1 class="text-white text-5xl pb-10">constructify!</h1>
	<div class="w-2/5">
		<Search />
	</div>
</div>

<style>
	:global(html) {
		height: 100%;
		font-family: Ubuntu, sans-serif;
	}

	:global(body) {
		background: url('background.png');
		background-size: 50%;
		height: 100%;
		margin: 0;
	}
</style>
