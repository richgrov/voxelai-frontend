<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import * as THREE from 'three';
	import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls';

	import '../../app.css';
	import Schematic from '../../world/schematic';

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

		const material = new THREE.MeshLambertMaterial({ color: 0x00FF00 });

		const schem = new Schematic(10, 10, 10);
		for (let x = 0; x < schem.xSize; x++) {
			for (let y = 0; y < schem.ySize; y++) {
				for (let z = 0; z < schem.zSize; z++) {
					if (Math.random() < 0.15) {
						schem.set(x, y, z, true);
					}
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
		const normals = new Array<number>();

		for (let x = 0; x < schematic.xSize; x++) {
			for (let y = 0; y < schematic.ySize; y++) {
				for (let z = 0; z < schematic.zSize; z++) {
					if (schematic.has(x, y, z)) {
						// -X
						positions.push(x, y, z, x, y + 1, z + 1, x, y + 1, z);
						positions.push(x, y, z + 1, x, y + 1, z + 1, x, y, z);

						// +X
						positions.push(x + 1, y, z, x + 1, y + 1, z, x + 1, y + 1, z + 1);
						positions.push(x + 1, y, z, x + 1, y + 1, z + 1, x + 1, y, z + 1);

						// -Y
						positions.push(x, y, z, x + 1, y, z, x + 1, y, z + 1);
						positions.push(x, y, z, x + 1, y, z + 1, x, y, z + 1);

						// +Y
						positions.push(x, y + 1, z, x + 1, y + 1, z + 1, x + 1, y + 1, z);
						positions.push(x, y + 1, z + 1, x + 1, y + 1, z + 1, x, y + 1, z);

						// -Z
						positions.push(x, y, z, x + 1, y + 1, z, x + 1, y, z);
						positions.push(x, y + 1, z, x + 1, y + 1, z, x, y, z);

						// +Z
						positions.push(x, y, z + 1, x + 1, y, z + 1, x + 1, y + 1, z + 1);
						positions.push(x, y, z + 1, x + 1, y + 1, z + 1, x, y + 1, z + 1);
					}
				}
			}
		}

		geometry.setAttribute('position', new THREE.Float32BufferAttribute(positions, 3));
		geometry.computeVertexNormals();
		return new THREE.Mesh(geometry, material);
	}
</script>
