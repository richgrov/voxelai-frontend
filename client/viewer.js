import * as THREE from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';

async function tryRender() {
  const canvas = document.querySelector('canvas');
  if (!canvas) {
    return;
  }

  const url = canvas.getAttribute('data-url');
  const renderer = new THREE.WebGLRenderer({ canvas });
  renderer.setSize(canvas.clientWidth, canvas.clientHeight);

  const scene = new THREE.Scene();
  scene.background = new THREE.Color(0x333333);

  const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
  const controls = new OrbitControls(camera, renderer.domElement);
  controls.update();

  const directional = new THREE.DirectionalLight(0xFFFFFF, 1);
  directional.penumbra = 1;
  directional.position.set(5, 4, 3);
  directional.lookAt(0, 0, 0);
  const ambient = new THREE.AmbientLight(0xFFFFFF, 1);
  scene.add(ambient, directional);

  const loader = new GLTFLoader();
  loader.load(url, gltf => {
    const obj = gltf.scene;

    const size = new THREE.Vector3();
    new THREE.Box3()
      .setFromObject(obj)
      .getSize(size);

    const length = size.length() / 2;
    camera.position.set(length, length * 0.6, length);
    camera.lookAt(0, 0, 0);

    obj.position.set(-size.x/2, -size.y/2, -size.z/2);
    obj.traverse(child => {
      if (child.isMesh) {
        child.material.metalness = 0;
      }
    });
    scene.add(gltf.scene);
    window.requestAnimationFrame(render);
  });

  function render() {
    renderer.render(scene, camera);
  }
  window.requestAnimationFrame(render);
  controls.addEventListener('change', () => window.requestAnimationFrame(render));
}

document.addEventListener("displayMesh", tryRender);

if (document.readyState !== 'complete') {
  document.addEventListener("DOMContentLoaded", tryRender);
} else {
  tryRender();
}
