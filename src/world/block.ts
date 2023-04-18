export type BlockFace = '-x' | '+x' | '-y' | '+y' | '-z' | '+z';
export type UvCoordinates = [number, number];

export interface Block {
	type: 'solid' | 'cross' | 'custom';
	uv?: (data: number, face: BlockFace) => UvCoordinates,
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

function directional(topBottom: UvCoordinates, side: UvCoordinates, front: UvCoordinates): Block {
	return {
		type: 'solid',
		uv: (data, face) => {
			switch (face) {
				case '-x': return data === 4 ? front : side;
				case '+x': return data === 5 ? front : side;
				case '-y': return topBottom;
				case '+y': return topBottom;
				case '-z': return data === 2 ? front : side;
				case '+z': return data === 3 ? front : side;
				default: throw new Error(`Invalid face ${face}`);
			}
		},
	}
}

const stone: UvCoordinates = [1, 0];
const cobbleStone: UvCoordinates = [0, 1];
const planks: UvCoordinates = [4, 0];
const oakLog = topSideBottom(5, 1, 4, 1, 5, 1);
const spruceLog = topSideBottom(5, 1, 4, 7, 5, 1);
const birchLog = topSideBottom(5, 1, 5, 7, 5, 1);
const dispenserFace: UvCoordinates = [14, 2];
const furnaceSide: UvCoordinates = [13, 2];
const furnaceTop: UvCoordinates = [15, 2];
const chestFront: UvCoordinates = [11, 1];
const chestSide: UvCoordinates = [10, 1];
const chestTop: UvCoordinates = [9, 1];
const furnaceFront: UvCoordinates = [12, 2];
const litFurnaceFront: UvCoordinates = [13, 3];
const pumpkinFront: UvCoordinates = [7, 7];
const pumpkinSide: UvCoordinates = [6, 7];
const pumpkinTop: UvCoordinates = [6, 6];
const jackOLanternFront: UvCoordinates = [8, 7];

const oakLeaves = simple(5, 3);
const spruceLeaves = simple(5, 8);
const chest = directional(chestTop, chestSide, chestFront);

const blocks: Array<Block | undefined> = [
	undefined,
	simple(stone[0], stone[1]), // Stone
	topSideBottom(0, 0, 3, 0, 2, 0), // Grass
	simple(2, 0), // Dirt
	simple(cobbleStone[0], cobbleStone[1]), // Cobblestone
	simple(planks[0], planks[1]), // Planks
	simple(15, 0), // Sapling TODO
	simple(1, 1), // Bedrock
	simple(15, 13), // Flowing water TODO
	simple(15, 13), // Water TODO
	simple(15, 15), // Flowing lava TODO
	simple(15, 15), // Lava TODO
	simple(2, 1), // Sand
	simple(3, 1), // Gravel
	simple(0, 2), // Gold ore
	simple(1, 2), // Iron ore
	simple(2, 2), // Coal ore
	{
		type: 'solid',
		uv: (data, face) => {
			switch (data) {
				case 0: return oakLog.uv!(data, face);
				case 1: return spruceLog.uv!(data, face);
				case 2: return birchLog.uv!(data, face);
				default: throw new Error(`Unexpected log type ${data}`);
			}
		},
	}, // Log
	{
		type: 'solid',
		uv: (data, face) => {
			switch (data) {
				case 0:
				case 2:
					return oakLeaves.uv!(data, face);
				case 1:
					return spruceLeaves.uv!(data, face);
				default: throw new Error(`Unexpected leaf type ${data}`);
			}
		},
	}, // Leaves
	simple(0, 3), // Sponge
	simple(1, 3), // Glass TODO
	simple(0, 10), // Lapis ore
	simple(0, 9), // Lapis block
	directional(furnaceTop, furnaceSide, dispenserFace), // Dispenser
	topSideBottom(0, 11, 0, 12, 0, 13), // Sandstone
	simple(10, 4), // Noteblock
	simple(6, 8), // Bed TODO
	simple(10, 3), // Powered rail TODO
	simple(12, 3), // Detector rail TODO
	topSideBottom(10, 6, 12, 6, 13, 6), // Sticky piston TODO
	simple(11, 0), // Web TODO
	simple(8, 3), // Fern TODO
	simple(7, 3), // Dead bush TODO
	topSideBottom(11, 6, 12, 6, 13, 6), // Piston TODO
	simple(11, 6), // Piston head TODO
	{
		type: 'solid',
		uv: (data, _face) => {
			switch (data) {
				case 0: return [4, 0];
				default:
					if (data >= 1 && data <= 7) {
						return [14 - data, 2];
					} else if (data >= 8 && data <= 15) {
						return [14 - (data - 8), 1];
					} else {
						throw new Error(`Invalid wool color ${data}`);
					}
			}
		},
	}, // Wool
	undefined, // Piston extension
	simple(13, 0), // Daisy
	simple(12, 0), // Rose
	simple(13, 1), // Brown mushroom
	simple(12, 1), // Red mushroom
	simple(7, 1), // Gold block
	simple(6, 1), // Iron block
	topSideBottom(6, 0, 5, 0, 6, 0), // Double stone slab
	topSideBottom(6, 0, 5, 0, 6, 0), // Stone slab TODO
	simple(7, 0), // Bricks
	topSideBottom(9, 0, 8, 0, 10, 0), // TNT
	topSideBottom(planks[0], planks[1], 3, 2, planks[0], planks[1]), // Bookshelf
	simple(4, 2), // Moss stone
	simple(5, 2), // Obsidian
	simple(5, 0), // Torch TODO
	simple(13, 15), // Fire TODO
	simple(13, 15), // Mob spawner TODO
	simple(planks[0], planks[1]), // Stairs TODO
	chest, // Chest
	simple(4, 10), // Redstone wire TODO
	simple(2, 3), // Diamond ore
	simple(8, 1), // Diamond block
	{
		type: 'solid',
		uv: (_data, face) => {
			switch (face) {
				case '-x': return [11, 3];
				case '+x': return [11, 3];
				case '-y': return planks;
				case '+y': return [11, 2];
				case '-z': return [12, 3];
				case '+z': return [12, 3];
				default: throw new Error(`Invalid face ${face}`);
			}
		},
	}, // Crafting table
	simple(8, 5), // Seeds TODO
	simple(7, 5), // Farmland TODO
	directional(furnaceTop, furnaceSide, furnaceFront), // Furnace
	directional(furnaceTop, furnaceSide, litFurnaceFront), // Lit furnace
	simple(planks[0], planks[1]), // Wall sign TODO
	simple(5, 1), // Door TODO
	simple(7, 1), // Ladder TODO
	simple(8, 0), // Rail TODO
	simple(cobbleStone[0], cobbleStone[1]), // Cobblestone stairs TODO
	simple(planks[0], planks[1]), // Standing sign TODO
	simple(6, 0), // Lever TODO
	simple(stone[0], stone[1]), // Stone pressure plate TODO
	simple(2, 5), // Iron door TODO
	simple(planks[0], planks[1]), // Wooden pressure plate TODO
	simple(3, 3), // Redstore ore
	simple(3, 3), // Lit redstone ore
	simple(3, 7), // Redstone torch off TODO
	simple(3, 6), // Redstone torch on TODO
	simple(stone[0], stone[1]), // Stone button TODO
	simple(2, 4), // Snow layer TODO
	simple(3, 4), // Ice TODO
	simple(2, 4), // Snow block
	topSideBottom(5, 4, 6, 4, 7, 4), // Cactus TODO
	simple(8, 4), // Clay
	simple(9, 4), // Sugarcane TODO
	simple(11, 4), // Jukebox
	simple(planks[0], planks[1]), // Fence TODO
	directional(pumpkinTop, pumpkinSide, pumpkinFront), // Pumpkin
	simple(7, 6), // Netherrack
	simple(8, 6), // Soulsand
	simple(9, 6), // Glowstone
	simple(13, 15), // Portal TODO
	directional(pumpkinTop, pumpkinSide, jackOLanternFront), // Jack o' lantern
	topSideBottom(9, 7, 10, 7, 12, 7), // Cake TODO
	simple(3, 8), // Redstone repeater off TODO
	simple(3, 9), // Redstone repeater on TODO
	chest, // Locked chest
	simple(4, 5), // Trapdoor TODO
];
export default blocks;
