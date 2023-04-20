<script lang="ts">
	import '../../app.css';
	import Scene from '../../lib/scene.svelte';
	import '../../firebase';
	import { connectFirestoreEmulator, getFirestore, doc, onSnapshot } from 'firebase/firestore';
   import { onMount } from 'svelte';

	let prompt = '';
	let status = '';
	let schemUrl = '';

	onMount(() => {
		const jobId = new URLSearchParams(window.location.search).get('id');
		if (!jobId) {
			status = 'Invalid URL';
			return;
		}

		schemUrl = import.meta.env['VITE_STORAGE_BUCKET'] + '/' + jobId + '.schematic';

		const db = getFirestore()
		if (import.meta.env['MODE'] === 'development') {
			connectFirestoreEmulator(db, 'localhost', 8080);
		}

		const unsub = onSnapshot(doc(db, 'jobs', jobId), doc => {
			if (!doc.exists()) {
				status = 'Invalid URL';
				unsub();
				return;
			}

			status = doc.get('status');
			prompt = doc.get('prompt');
		});

		return unsub;
	});
</script>

<h3>{prompt}</h3>
{#if status === 'finished'}
	<Scene schematicUrl={schemUrl} />
{:else}
	<h1>{status}</h1>
{/if}

<style>
	:global(html) {
		height: 100%;
	}

	:global(body) {
		background-color: #222;
		color: white;
		height: 100%;
	}
</style>
