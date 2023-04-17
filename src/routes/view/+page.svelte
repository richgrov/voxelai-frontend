<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import * as THREE from 'three';
	import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls';

	import '../../app.css';
	import Schematic from '../../world/schematic';
	import blocks from '../../world/block';

	let canvas: HTMLCanvasElement;

	onMount(() => {
		const renderer = new THREE.WebGLRenderer();
		renderer.setSize(window.innerWidth, window.innerHeight);
		canvas = document.body.appendChild(renderer.domElement);

		const scene = new THREE.Scene();
		const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
		const controls = new OrbitControls(camera, renderer.domElement);
		camera.position.set(10, 10, 10);
		controls.update();

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

		const schem = new Schematic(64, 64, 64);
		for (let x = 0; x < schem.xSize; x++) {
			for (let y = 0; y < schem.ySize; y++) {
				for (let z = 0; z < schem.zSize; z++) {
					const id = Math.floor(Math.random() * 6);
					schem.set(x, y, z, id);
				}
			}
		}

		const mesh = buildMesh(schem, material);
		mesh.position.sub(new THREE.Vector3(schem.xSize / 2, schem.ySize / 2, schem.zSize / 2));
		scene.add(mesh);

		function loop() {
			window.requestAnimationFrame(loop);
			renderer.render(scene, camera);
		}
		window.requestAnimationFrame(loop);
	});

	onDestroy(() => {
		if (canvas) {
			canvas.remove();
		}
	});

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

					const [blockU, blockV] = blockInfo.uv(0, 0);
					const textureIncrement = 1 / 16;
					const startU = textureIncrement * blockU;
					const endU = textureIncrement * (blockU + 1);
					const startV = 1 - textureIncrement * blockV;
					const endV = 1 - textureIncrement * (blockV + 1);

					// -X
					positions.push(x, y, z, x, y + 1, z + 1, x, y + 1, z);
					positions.push(x, y, z + 1, x, y + 1, z + 1, x, y, z);
					uv.push(startU, startV, endU, endV, endU, startV);
					uv.push(startU, endV, endU, endV, startU, startV);

					// +X
					positions.push(x + 1, y, z, x + 1, y + 1, z, x + 1, y + 1, z + 1);
					positions.push(x + 1, y, z, x + 1, y + 1, z + 1, x + 1, y, z + 1);
					uv.push(startU, startV, endU, startV, endU, endV);
					uv.push(startU, startV, endU, endV, startU, endV);

					// -Y
					positions.push(x, y, z, x + 1, y, z, x + 1, y, z + 1);
					positions.push(x, y, z, x + 1, y, z + 1, x, y, z + 1);
					uv.push(startU, startV, endU, startV, endU, endV);
					uv.push(startU, startV, endU, endV, startU, endV);

					// +Y
					positions.push(x, y + 1, z, x + 1, y + 1, z + 1, x + 1, y + 1, z);
					positions.push(x, y + 1, z + 1, x + 1, y + 1, z + 1, x, y + 1, z);
					uv.push(startU, startV, endU, endV, endU, startV);
					uv.push(startU, endV, endU, endV, startU, startV);

					// -Z
					positions.push(x, y, z, x + 1, y + 1, z, x + 1, y, z);
					positions.push(x, y + 1, z, x + 1, y + 1, z, x, y, z);
					uv.push(startU, startV, endU, endV, endU, startV);
					uv.push(startU, endV, endU, endV, startU, startV);

					// +Z
					positions.push(x, y, z + 1, x + 1, y, z + 1, x + 1, y + 1, z + 1);
					positions.push(x, y, z + 1, x + 1, y + 1, z + 1, x, y + 1, z + 1);
					uv.push(startU, startV, endU, startV, endU, endV);
					uv.push(startU, startV, endU, endV, startU, endV);
				}
			}
		}

		geometry.setAttribute('position', new THREE.Float32BufferAttribute(positions, 3));
		geometry.setAttribute('uv', new THREE.Float32BufferAttribute(uv, 2));
		geometry.computeVertexNormals();
		return new THREE.Mesh(geometry, material);
	}
</script>
