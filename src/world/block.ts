export interface Block {
	type: 'solid' | 'cross' | 'custom';
	uv?: (data: number, face: number) => [number, number],
}

function simple(u: number, v: number): Block {
	return {
		type: 'solid',
		uv: (_data: number, _face: number) => [u, v],
	}
}

const blocks: Array<Block | undefined> = [
	undefined,
	simple(1, 0), // Stone
	simple(3, 0), // Grass
	simple(2, 0), // Dirt
	simple(0, 1), // Cobblestone
	simple(4, 0), // Planks
];
export default blocks;
