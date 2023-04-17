export type BlockFace = '-x' | '+x' | '-y' | '+y' | '-z' | '+z';

export interface Block {
	type: 'solid' | 'cross' | 'custom';
	uv?: (data: number, face: BlockFace) => [number, number],
}

function simple(u: number, v: number): Block {
	return {
		type: 'solid',
		uv: (_data: number, _face: BlockFace) => [u, v],
	}
}

function topSideBottom(topU: number, topV: number, sideU: number, sideV: number, bottomU: number, bottomV: number): Block {
	return {
		type: 'solid',
		uv: (_data: number, face: BlockFace) => {
			switch (face) {
				case '+y':
					return [topU, topV];
				case '-y':
					return [bottomU, bottomV];
				case '-x':
				case '+x':
				case '-z':
				case '+z':
					return [sideU, sideV];
				default:
					throw new Error(`Invalid face ${face}`);
			}
		},
	}
}

const blocks: Array<Block | undefined> = [
	undefined,
	simple(1, 0), // Stone
	topSideBottom(0, 0, 3, 0, 2, 0), // Grass
	simple(2, 0), // Dirt
	simple(0, 1), // Cobblestone
	simple(4, 0), // Planks
];
export default blocks;
