<script lang="ts">
	import { onMount } from 'svelte';
	import * as THREE from 'three';
	import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls';
	import Schematic from '../world/schematic';
	import blocks from '../world/block';
	import * as NBT from 'nbtify';

	export let schematicUrl: string;
	let canvasContainer: Element;

	onMount(async () => {
		const schematic = await fetch(schematicUrl)
			.then(response => response.arrayBuffer())
			.then(buf => NBT.read(buf))
			.then(nbt => nbt.data);

		console.log(schematic);

		const renderer = new THREE.WebGLRenderer();
		renderer.setSize(canvasContainer.clientWidth, canvasContainer.clientHeight);
		canvasContainer.appendChild(renderer.domElement);

		const scene = new THREE.Scene();
		const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
		const controls = new OrbitControls(camera, renderer.domElement);
		camera.position.set(10, 10, 10);
		controls.update();
		controls.addEventListener('change', render);

		const light = new THREE.AmbientLight(0x404040);
		scene.add(light);

		const dirLight = new THREE.DirectionalLight(0xFFFFFF, 0.5);
		dirLight.position.set(1, 2, -2);
		scene.add(dirLight);

		const material = new THREE.MeshLambertMaterial();
		const terrain = new THREE.TextureLoader().load('/terrain.png');
		terrain.wrapS = THREE.RepeatWrapping;
		terrain.wrapT = THREE.RepeatWrapping;
		terrain.magFilter = THREE.NearestFilter;

		material.map = terrain;

		const schem = new Schematic(schematic['Width'], schematic['Height'], schematic['Length']);
		schem.blocks = schematic['Blocks'];
		schem.data = schematic['Data'];

		const mesh = buildMesh(schem, material);
		mesh.position.sub(new THREE.Vector3(schem.xSize / 2, schem.ySize / 2, schem.zSize / 2));
		scene.add(mesh);

		function render() {
			renderer.render(scene, camera);
		}
		// Needs to be delayed. Unsure why
		setTimeout(render, 100);
	});

	function calculateUv(tileUv: [number, number]): [number, number, number, number] {
		const textureIncrement = 1 / 16;
		return [
			textureIncrement * tileUv[0],
			textureIncrement * (tileUv[0] + 1),
			1 - textureIncrement * tileUv[1],
			1 - textureIncrement * (tileUv[1] + 1),
		];
	}

	function buildMesh(schematic: Schematic, material: THREE.Material): THREE.Mesh {
		const geometry = new THREE.BufferGeometry();
		const positions = new Array<number>();
		const uv = new Array<number>();

		for (let x = 0; x < schematic.xSize; x++) {
			for (let y = 0; y < schematic.ySize; y++) {
				for (let z = 0; z < schematic.zSize; z++) {
					const [blockId, data] = schematic.get(x, y, z);
					const blockInfo = blocks[blockId];
					if (!blockInfo) {
						continue;
					}

					if (x === 0 || !schematic.has(x - 1, y, z)) {
						const [startU, endU, startV, endV] = calculateUv(blockInfo.uv(data, '-x'));
						positions.push(x, y, z, x, y + 1, z + 1, x, y + 1, z);
						positions.push(x, y, z + 1, x, y + 1, z + 1, x, y, z);
						uv.push(startU, endV, endU, startV, startU, startV);
						uv.push(endU, endV, endU, startV, startU, endV);
					}
					
					if (x === schematic.xSize - 1 || !schematic.has(x + 1, y, z)) {
						const [startU, endU, startV, endV] = calculateUv(blockInfo.uv(data, '+x'));
						positions.push(x + 1, y, z, x + 1, y + 1, z, x + 1, y + 1, z + 1);
						positions.push(x + 1, y, z, x + 1, y + 1, z + 1, x + 1, y, z + 1);
						uv.push(endU, endV, endU, startV, startU, startV);
						uv.push(endU, endV, startU, startV, startU, endV);
					}
					
					if (y === 0 || !schematic.has(x, y - 1, z)) {
						const [startU, endU, startV, endV] = calculateUv(blockInfo.uv(data, '-y'));
						positions.push(x, y, z, x + 1, y, z, x + 1, y, z + 1);
						positions.push(x, y, z, x + 1, y, z + 1, x, y, z + 1);
						uv.push(startU, startV, endU, startV, endU, endV);
						uv.push(startU, startV, endU, endV, startU, endV);
					}
					
					if (y === schematic.ySize - 1 || !schematic.has(x, y + 1, z)) {
						const [startU, endU, startV, endV] = calculateUv(blockInfo.uv(data, '+y'));
						positions.push(x, y + 1, z, x + 1, y + 1, z + 1, x + 1, y + 1, z);
						positions.push(x, y + 1, z + 1, x + 1, y + 1, z + 1, x, y + 1, z);
						uv.push(startU, startV, endU, endV, endU, startV);
						uv.push(startU, endV, endU, endV, startU, startV);
					}
					
					if (z === 0 || !schematic.has(x, y, z - 1)) {
						const [startU, endU, startV, endV] = calculateUv(blockInfo.uv(data, '-z'));
						positions.push(x, y, z, x + 1, y + 1, z, x + 1, y, z);
						positions.push(x, y + 1, z, x + 1, y + 1, z, x, y, z);
						uv.push(endU, endV, startU, startV, startU, endV);
						uv.push(endU, startV, startU, startV, endU, endV);
					}
					
					if (z === schematic.zSize - 1 || !schematic.has(x, y, z + 1)) {
						const [startU, endU, startV, endV] = calculateUv(blockInfo.uv(data, '+z'));
						positions.push(x, y, z + 1, x + 1, y, z + 1, x + 1, y + 1, z + 1);
						positions.push(x, y, z + 1, x + 1, y + 1, z + 1, x, y + 1, z + 1);
						uv.push(startU, endV, endU, endV, endU, startV);
						uv.push(startU, endV, endU, startV, startU, startV);
					}
				}
			}
		}

		geometry.setAttribute('position', new THREE.Float32BufferAttribute(positions, 3));
		geometry.setAttribute('uv', new THREE.Float32BufferAttribute(uv, 2));
		geometry.computeVertexNormals();
		return new THREE.Mesh(geometry, material);
	}
</script>

<div bind:this={canvasContainer} class="w-full h-full"></div>

