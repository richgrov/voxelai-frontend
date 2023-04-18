<script lang="ts">
	import '../app.css';
	import Search from '../lib/search.svelte';
	import '../firebase';
	import { getFunctions, httpsCallable, connectFunctionsEmulator } from 'firebase/functions';

	const functions = getFunctions();
	const generate = httpsCallable<{ prompt: string }, { jobId: string }>(functions, 'generate');
	if (import.meta.env['MODE'] === 'development') {
		connectFunctionsEmulator(functions, 'localhost', 5001);
	}

	async function submit(prompt: string) {
		try {
			const result = await generate({ prompt });
			const jobId = result.data.jobId;
			window.location.href = `/view?id=${jobId}`;
		} catch (err) {
			console.error(err);
		}
	}
</script>

<div class="flex justify-center items-center flex-col h-4/5">
	<h1 class="text-white text-5xl pb-10">constructify!</h1>
	<div class="w-2/5">
		<Search onSubmit={submit} />
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
	}
</style>
